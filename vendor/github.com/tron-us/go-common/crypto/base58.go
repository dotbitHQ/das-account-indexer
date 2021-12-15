package crypto

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcutil/base58"
)

var (
	ErrDecodeLength = errors.New("base58 decode length error")
	ErrDecodeCheck  = errors.New("base58 check failed")
	ErrEncodeLength = errors.New("base58 encode length error")
)

// Decode58Check converts *base58-string to hex and check.
func Decode58Check(input *string) (*string, error) {
	if input == nil {
		return nil, nil
	}
	decodeCheck := base58.Decode(*input)
	if len(decodeCheck) <= 4 {
		return nil, ErrDecodeLength
	}
	decodeData := decodeCheck[:len(decodeCheck)-4]
	hash0 := Hash(decodeData)
	hash1 := Hash(hash0)
	if hash1[0] == decodeCheck[len(decodeData)] && hash1[1] == decodeCheck[len(decodeData)+1] &&
		hash1[2] == decodeCheck[len(decodeData)+2] && hash1[3] == decodeCheck[len(decodeData)+3] {
		s := hex.EncodeToString(decodeData)
		return &s, nil
	}
	return nil, ErrDecodeCheck
}

// Encode58Check converts *hex-string to base58 and check.
func Encode58Check(input *string) (*string, error) {
	if input == nil {
		return nil, nil
	}
	b, err := hex.DecodeString(*input)
	if err != nil {
		return nil, err
	}
	hash0 := Hash(b)
	hash1 := Hash(hash0)
	// Since hash (sha256) never fails, the hash should always have length >4
	inputCheck := append(b, hash1[:4]...)
	result := base58.Encode(inputCheck)

	return &result, nil
}

// Encode58CheckLen is Encode58Check + length requirement for non-nil string
func Encode58CheckLen(input *string, length int) (*string, error) {
	s, err := Encode58Check(input)
	if err != nil {
		return nil, err
	}
	if s == nil || len(*s) == length {
		return s, nil
	}
	return nil, ErrEncodeLength
}
