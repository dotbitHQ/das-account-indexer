package common

import (
	"encoding/hex"
	"fmt"
	"github.com/Andrew-M-C/go.emoji/official"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

const (
	OneCkb                            = uint64(1e8)
	DasLockWithBalanceTypeOccupiedCkb = uint64(116 * 1e8)
	ProposalCellOccupiedCkb           = uint64(106 * 1e8)
	MinCellOccupiedCkb                = uint64(61 * 1e8)
	PercentRateBase                   = 1e4
	UsdRateBase                       = 1e6

	AccountStatusNormal    uint8 = 0
	AccountStatusOnSale    uint8 = 1
	AccountStatusOnAuction uint8 = 2
	AccountStatusOnCross   uint8 = 3

	OneYearSec = int64(3600 * 24 * 365)
)

func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func Hex2Bytes(s string) []byte {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	h, _ := hex.DecodeString(s)
	return h
}

func Bytes2Hex(b []byte) string {
	h := hex.EncodeToString(b)
	if len(h) == 0 {
		h = "0"
	}
	return "0x" + h
}

func GetScript(codeHash, args string) *types.Script {
	return &types.Script{
		CodeHash: types.HexToHash(codeHash),
		HashType: types.HashTypeType,
		Args:     Hex2Bytes(args),
	}
}

// GetNormalLockScript normal script
func GetNormalLockScript(args string) *types.Script {
	return GetScript(transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH, args)
}

// GetNormalLockScriptByMultiSig multi sig
func GetNormalLockScriptByMultiSig(args string) *types.Script {
	return GetScript(transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH, args)
}

func Blake2b(acc []byte) []byte {
	bys, _ := blake2b.Blake256(acc)
	return bys
}

func GetAccountIdByAccount(acc string) []byte {
	if acc != "" && !strings.HasSuffix(acc, DasAccountSuffix) {
		acc = acc + DasAccountSuffix
	}
	bys, _ := blake2b.Blake160([]byte(acc))
	return bys
}

func OutputDataToAccountId(data []byte) ([]byte, error) {
	if size := len(data); size < HashBytesLen+DasAccountIdLen {
		return nil, fmt.Errorf("len not enough: %d", size)
	}
	return data[HashBytesLen : HashBytesLen+DasAccountIdLen], nil
}

func OutputDataToSMTRoot(data []byte) ([]byte, error) {
	if size := len(data); size < HashBytesLen {
		return nil, fmt.Errorf("len not enough: %d", size)
	}
	return data[0:HashBytesLen], nil
}

func GetAccountLength(account string) uint8 {
	account = strings.TrimSuffix(account, DasAccountSuffix)
	nextIndex := 0
	accLen := uint8(0)
	for i, _ := range account {
		if i < nextIndex {
			continue
		}
		match, length := official.AllSequences.HasEmojiPrefix(account[i:])
		if match {
			nextIndex = i + length
		}
		accLen++
	}
	return accLen
}
