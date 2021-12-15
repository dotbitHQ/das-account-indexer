package bech32

import (
	"errors"
	"fmt"
	"strings"
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var gen = []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

func Decode(bech string) (string, []byte, error) {
	for i := 0; i < len(bech); i++ {
		if bech[i] < 33 || bech[i] > 126 {
			return "", nil, fmt.Errorf("invalid character: '%c'", bech[i])
		}
	}

	lower := strings.ToLower(bech)
	upper := strings.ToUpper(bech)
	if bech != lower && bech != upper {
		return "", nil, errors.New("string not all lowercase or all uppercase")
	}

	bech = lower

	one := strings.LastIndexByte(bech, '1')
	if one < 1 || one+7 > len(bech) {
		return "", nil, fmt.Errorf("invalid index of 1")
	}

	hrp := bech[:one]
	data := bech[one+1:]

	decoded, err := toBytes(data)
	if err != nil {
		return "", nil, errors.New(fmt.Sprintf("failed converting data to bytes: %v", err))
	}

	if !bech32VerifyChecksum(hrp, decoded) {
		moreInfo := ""
		checksum := bech[len(bech)-6:]
		expected, err := toChars(bech32Checksum(hrp,
			decoded[:len(decoded)-6]))
		if err == nil {
			moreInfo = fmt.Sprintf("Expected %v, got %v.", expected, checksum)
		}
		return "", nil, errors.New("checksum failed. " + moreInfo)
	}

	return hrp, decoded[:len(decoded)-6], nil
}

func Encode(hrp string, data []byte) (string, error) {
	checksum := bech32Checksum(hrp, data)
	combined := append(data, checksum...)

	dataChars, err := toChars(combined)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to convert data bytes to chars: %v", err))
	}
	return hrp + "1" + dataChars, nil
}

func toBytes(chars string) ([]byte, error) {
	decoded := make([]byte, 0, len(chars))
	for i := 0; i < len(chars); i++ {
		index := strings.IndexByte(charset, chars[i])
		if index < 0 {
			return nil, errors.New(fmt.Sprintf("invalid character not part of charset: %v", chars[i]))
		}
		decoded = append(decoded, byte(index))
	}
	return decoded, nil
}

func toChars(data []byte) (string, error) {
	result := make([]byte, 0, len(data))
	for _, b := range data {
		if int(b) >= len(charset) {
			return "", errors.New(fmt.Sprintf("invalid data byte: %v", b))
		}
		result = append(result, charset[b])
	}
	return string(result), nil
}

func ConvertBits(data []byte, fromBits, toBits uint8, pad bool) ([]byte, error) {
	if fromBits < 1 || fromBits > 8 || toBits < 1 || toBits > 8 {
		return nil, errors.New("only bit groups between 1 and 8 allowed")
	}

	var regrouped []byte

	nextByte := byte(0)
	filledBits := uint8(0)

	for _, b := range data {
		b = b << (8 - fromBits)

		remFromBits := fromBits
		for remFromBits > 0 {
			remToBits := toBits - filledBits

			toExtract := remFromBits
			if remToBits < toExtract {
				toExtract = remToBits
			}

			nextByte = (nextByte << toExtract) | (b >> (8 - toExtract))

			b = b << toExtract
			remFromBits -= toExtract
			filledBits += toExtract

			if filledBits == toBits {
				regrouped = append(regrouped, nextByte)
				filledBits = 0
				nextByte = 0
			}
		}
	}

	if pad && filledBits > 0 {
		nextByte = nextByte << (toBits - filledBits)
		regrouped = append(regrouped, nextByte)
		filledBits = 0
		nextByte = 0
	}

	if filledBits > 0 && (filledBits > 4 || nextByte != 0) {
		return nil, errors.New("invalid incomplete group")
	}

	return regrouped, nil
}

func bech32Checksum(hrp string, data []byte) []byte {
	integers := make([]int, len(data))
	for i, b := range data {
		integers[i] = int(b)
	}
	values := append(bech32HrpExpand(hrp), integers...)
	values = append(values, []int{0, 0, 0, 0, 0, 0}...)
	polymod := bech32Polymod(values) ^ 1
	var res []byte
	for i := 0; i < 6; i++ {
		res = append(res, byte((polymod>>uint(5*(5-i)))&31))
	}
	return res
}

func bech32Polymod(values []int) int {
	chk := 1
	for _, v := range values {
		b := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ v
		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

func bech32HrpExpand(hrp string) []int {
	v := make([]int, 0, len(hrp)*2+1)
	for i := 0; i < len(hrp); i++ {
		v = append(v, int(hrp[i]>>5))
	}
	v = append(v, 0)
	for i := 0; i < len(hrp); i++ {
		v = append(v, int(hrp[i]&31))
	}
	return v
}

func bech32VerifyChecksum(hrp string, data []byte) bool {
	integers := make([]int, len(data))
	for i, b := range data {
		integers[i] = int(b)
	}
	concat := append(bech32HrpExpand(hrp), integers...)
	return bech32Polymod(concat) == 1
}
