package sign

import (
	"crypto/ed25519"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Ed25519Signature(privateKey, message []byte) []byte {
	l := len(message)
	if l == 0 {
		return nil
	}
	tmp := append([]byte(fmt.Sprintf(common.Ed25519MessageHeader, l)), message...)
	tmp = crypto.Keccak256(tmp)
	tmp = crypto.Keccak256(tmp)

	return ed25519.Sign(privateKey, tmp)
}

func VerifyEd25519Signature(publicKey, message, sig []byte) bool {
	l := len(message)
	if l == 0 {
		return false
	}
	tmp := append([]byte(fmt.Sprintf(common.Ed25519MessageHeader, l)), message...)
	tmp = crypto.Keccak256(tmp)
	tmp = crypto.Keccak256(tmp)

	return ed25519.Verify(publicKey, tmp, sig)
}
