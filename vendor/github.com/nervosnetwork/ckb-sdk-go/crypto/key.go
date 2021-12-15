package crypto

import (
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
)

// Key key pair
type Key interface {
	Bytes() []byte
	Sign(data []byte) ([]byte, error)
	Script(systemScripts *utils.SystemScripts) (*types.Script, error)
}

func ZeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
