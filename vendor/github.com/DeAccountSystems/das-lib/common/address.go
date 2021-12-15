package common

import (
	"fmt"
	"github.com/tron-us/go-common/crypto"
)

type ChainType int

const (
	ChainTypeCkb  ChainType = 0
	ChainTypeEth  ChainType = 1
	ChainTypeBtc  ChainType = 2
	ChainTypeTron ChainType = 3

	HexPreFix           = "0x"
	TronPreFix          = "41"
	TronBase58PreFix    = "T"
	DasLockCkbPreFix    = "00"
	DasLockEthPreFix    = "03"
	DasLockTronPreFix   = "04"
	DasLockEth712PreFix = "05"
)

const (
	TronMessageHeader = "\x19TRON Signed Message:\n32"
	EthMessageHeader  = "\x19Ethereum Signed Message:\n32"
)

const (
	DasAccountSuffix = ".bit"
	DasLockArgsLen   = 42
	DasAccountIdLen  = 20
	HashBytesLen     = 32

	ExpireTimeLen    = 8
	NextAccountIdLen = 20

	ExpireTimeEndIndex      = HashBytesLen + DasAccountIdLen + NextAccountIdLen + ExpireTimeLen
	NextAccountIdStartIndex = HashBytesLen + DasAccountIdLen
	NextAccountIdEndIndex   = NextAccountIdStartIndex + NextAccountIdLen
)

func (c ChainType) String() string {
	switch c {
	case ChainTypeCkb:
		return "CKB"
	case ChainTypeBtc:
		return "BTC"
	case ChainTypeEth:
		return "ETH"
	case ChainTypeTron:
		return "TRON"
	}
	return ""
}

func TronHexToBase58(address string) (string, error) {
	tAddr, err := crypto.Encode58Check(&address)
	if err != nil {
		return "", fmt.Errorf("Encode58Check:%v", err)
	}
	return *tAddr, nil
}

func TronBase58ToHex(address string) (string, error) {
	addr, err := crypto.Decode58Check(&address)
	if err != nil {
		return "", fmt.Errorf("Decode58Check:%v", err)
	}
	return *addr, nil
}
