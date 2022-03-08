package common

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/tron-us/go-common/crypto"
)

type ChainType int

const (
	ChainTypeCkb   ChainType = 0
	ChainTypeEth   ChainType = 1
	ChainTypeBtc   ChainType = 2
	ChainTypeTron  ChainType = 3
	ChainTypeMixin ChainType = 4

	HexPreFix            = "0x"
	TronPreFix           = "41"
	TronBase58PreFix     = "T"
	DasLockCkbPreFix     = "00"
	DasLockEthPreFix     = "03"
	DasLockTronPreFix    = "04"
	DasLockEth712PreFix  = "05"
	DasLockEd25519PreFix = "06"
)

const (
	TronMessageHeader    = "\x19TRON Signed Message:\n%d"
	EthMessageHeader     = "\x19Ethereum Signed Message:\n%d"
	Ed25519MessageHeader = "\x18Ed25519 Signed Message:\n%d"
)

const (
	DasAccountSuffix  = ".bit"
	DasLockArgsLen    = 42
	DasLockArgsLenMax = 66
	DasAccountIdLen   = 20
	HashBytesLen      = 32

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
	case ChainTypeMixin:
		return "MIXIN"
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

func ConvertScriptToAddress(mode address.Mode, script *types.Script) (string, error) {
	if transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH == script.CodeHash.String() ||
		transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH == script.CodeHash.String() {
		return address.ConvertScriptToShortAddress(mode, script)
	}
	return address.ConvertScriptToAddress(mode, script)

	//if script.HashType == types.HashTypeType && len(script.Args) >= 20 && len(script.Args) <= 22 {
	//	return address.ConvertScriptToShortAddress(mode, script)
	//}
	//
	//hashType := address.FullTypeFormat
	//if script.HashType == types.HashTypeData {
	//	hashType = address.FullDataFormat
	//}
	//return address.ConvertScriptToFullAddress(hashType, mode, script)
}
