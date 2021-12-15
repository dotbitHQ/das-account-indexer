package types

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	HashLength = 32
)

var (
	hashT = reflect.TypeOf(Hash{})
)

type Hash [HashLength]byte

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func HexToHash(s string) Hash {
	return BytesToHash(common.FromHex(s))
}

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) Hex() string {
	return hexutil.Encode(h[:])
}

func (h Hash) String() string {
	return h.Hex()
}

func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}
