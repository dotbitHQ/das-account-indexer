package sign

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TronSignature(signType bool, data []byte, hexPrivateKey string) ([]byte, error) {
	l := len(data)
	if l == 0 {
		return nil, errors.New("invalid raw data")
	}

	if signType {
		l = 32 // fix tron sign
		data = append([]byte(fmt.Sprintf(common.TronMessageHeader, l)), data...)
	}

	tmpHash := crypto.Keccak256(data)

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}

	signData, err := crypto.Sign(tmpHash[:], privateKey)
	if err == nil && len(signData) == 65 && signData[64] < 27 {
		signData[64] = signData[64] + 27
	}
	return signData, err
}

func TronVerifySignature(signType bool, sign []byte, rawByte []byte, base58Addr string) bool {
	l := len(rawByte)
	if len(sign) != 65 || l == 0 { // sign check
		return false
	}

	if sign[64] >= 27 {
		sign[64] = sign[64] - 27
	}

	if signType {
		l = 32 // fix tron sign
		rawByte = append([]byte(fmt.Sprintf(common.TronMessageHeader, l)), rawByte...)
	}

	pubKey, err := GetSignedPubKey(rawByte, sign)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	tronAddr := common.TronPreFix + recoveredAddr.String()[2:]

	base58address, err := common.TronHexToBase58(tronAddr)
	if err != nil {
		return false
	}
	return base58address == base58Addr
}

func GetSignedPubKey(rawByte []byte, sign []byte) (*ecdsa.PublicKey, error) {
	if len(sign) != 65 { // sign check
		return nil, errors.New("invalid transaction signature, should be 65 length bytes")
	}
	hash := crypto.Keccak256(rawByte)

	pub, err := crypto.Ecrecover(hash[:], sign)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPubkey(pub)
}
