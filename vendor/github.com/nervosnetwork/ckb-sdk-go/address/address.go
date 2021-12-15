package address

import (
	"encoding/hex"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
	"github.com/pkg/errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/nervosnetwork/ckb-sdk-go/crypto/bech32"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type Mode string
type Type string

const (
	Mainnet Mode = "ckb"
	Testnet Mode = "ckt"

	TypeFull  Type = "Full"
	TypeShort Type = "Short"

	ShortFormat               = "01"
	FullDataFormat            = "02"
	FullTypeFormat            = "04"
	CodeHashIndexSingleSig    = "00"
	CodeHashIndexMultisigSig  = "01"
	CodeHashIndexAnyoneCanPay = "02"
)

var shortPayloadSupportedArgsLens = [2]int{20, 22}

type ParsedAddress struct {
	Mode   Mode
	Type   Type
	Script *types.Script
}

func Generate(mode Mode, script *types.Script) (string, error) {
	if script.HashType == types.HashTypeType && isShortPayloadSupportedArgsLen(len(script.Args)) {
		if transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH == script.CodeHash.String() {
			// generate_short_payload_singleSig_address
			payload := ShortFormat + CodeHashIndexSingleSig + hex.EncodeToString(script.Args)
			data, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
			if err != nil {
				return "", err
			}
			return bech32.Encode((string)(mode), data)
		} else if transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH == script.CodeHash.String() {
			// generate_short_payload_multisig_address
			payload := ShortFormat + CodeHashIndexMultisigSig + hex.EncodeToString(script.Args)
			data, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
			if err != nil {
				return "", err
			}
			return bech32.Encode((string)(mode), data)
		} else if utils.AnyoneCanPayCodeHashOnLina == script.CodeHash.String() || utils.AnyoneCanPayCodeHashOnAggron == script.CodeHash.String() {
			payload := ShortFormat + CodeHashIndexAnyoneCanPay + hex.EncodeToString(script.Args)
			data, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
			if err != nil {
				return "", err
			}
			return bech32.Encode((string)(mode), data)
		}
	}

	hashType := FullTypeFormat
	if script.HashType == types.HashTypeData {
		hashType = FullDataFormat
	}

	return GenerateFullPayloadAddress(hashType, mode, script)
}

func isShortPayloadSupportedArgsLen(argLen int) bool {
	if argLen >= shortPayloadSupportedArgsLens[0] && argLen <= shortPayloadSupportedArgsLens[1] {
		return true
	}
	return false
}

func GenerateFullPayloadAddress(hashType string, mode Mode, script *types.Script) (string, error) {
	payload := hashType + hex.EncodeToString(script.CodeHash.Bytes()) + hex.EncodeToString(script.Args)
	data, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode((string)(mode), data)
}

func Parse(address string) (*ParsedAddress, error) {
	hrp, decoded, err := bech32.Decode(address)
	if err != nil {
		return nil, err
	}
	data, err := bech32.ConvertBits(decoded, 5, 8, false)
	if err != nil {
		return nil, err
	}
	payload := hex.EncodeToString(data)

	var addressType Type
	var script types.Script
	if strings.HasPrefix(payload, "01") {
		addressType = TypeShort
		if CodeHashIndexSingleSig == payload[2:4] {
			script = types.Script{
				CodeHash: types.HexToHash(transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH),
				HashType: types.HashTypeType,
				Args:     common.Hex2Bytes(payload[4:]),
			}
		} else if CodeHashIndexAnyoneCanPay == payload[2:4] {
			script = types.Script{
				HashType: types.HashTypeType,
				Args:     common.Hex2Bytes(payload[4:]),
			}
			if hrp == (string)(Testnet) {
				script.CodeHash = types.HexToHash(utils.AnyoneCanPayCodeHashOnAggron)
			} else {
				script.CodeHash = types.HexToHash(utils.AnyoneCanPayCodeHashOnLina)
			}
		} else {
			script = types.Script{
				CodeHash: types.HexToHash(transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH),
				HashType: types.HashTypeType,
				Args:     common.Hex2Bytes(payload[4:]),
			}
		}
	} else if strings.HasPrefix(payload, "02") {
		addressType = TypeFull
		script = types.Script{
			CodeHash: types.HexToHash(payload[2:66]),
			HashType: types.HashTypeData,
			Args:     common.Hex2Bytes(payload[66:]),
		}
	} else if strings.HasPrefix(payload, "04") {
		addressType = TypeFull
		script = types.Script{
			CodeHash: types.HexToHash(payload[2:66]),
			HashType: types.HashTypeType,
			Args:     common.Hex2Bytes(payload[66:]),
		}
	} else {
		return nil, errors.New("address type error:" + payload[:2])
	}

	result := &ParsedAddress{
		Mode:   Mode(hrp),
		Type:   addressType,
		Script: &script,
	}
	return result, nil
}

func ValidateChequeAddress(addr string, systemScripts *utils.SystemScripts) (*ParsedAddress, error) {
	parsedSenderAddr, err := Parse(addr)
	if err != nil {
		return nil, err
	}
	if isSecp256k1Lock(parsedSenderAddr, systemScripts) {
		return parsedSenderAddr, nil
	}
	return nil, errors.Errorf("address %s is not an SECP256K1 short format address", addr)
}

func isSecp256k1Lock(parsedSenderAddr *ParsedAddress, systemScripts *utils.SystemScripts) bool {
	return parsedSenderAddr.Script.CodeHash == systemScripts.SecpSingleSigCell.CellHash &&
		parsedSenderAddr.Script.HashType == systemScripts.SecpSingleSigCell.HashType &&
		len(parsedSenderAddr.Script.Args) == 20
}
