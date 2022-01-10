package address

import (
	"encoding/hex"
	"strings"

	"github.com/nervosnetwork/ckb-sdk-go/utils"
	"github.com/pkg/errors"

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

	TYPE_FULL_WITH_BECH32M    = "00"
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

func ConvertScriptToAddress(mode Mode, script *types.Script) (string, error) {
	return ConvertScriptToBech32mFullAddress(mode, script)
}

// Deprecated: Short address format deprecated because it is limited (only support secp256k1_blake160,
// secp256k1_multisig, anyone_can_pay) and a flaw has been found in its encoding method bech32,
// which could enable attackers to generate valid but unexpected addresses.
// For more please check https://github.com/nervosnetwork/rfcs/blob/master/rfcs/0021-ckb-address-format/0021-ckb-address-format.md
func ConvertScriptToShortAddress(mode Mode, script *types.Script) (string, error) {
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
	return "", errors.New("The given script can not be converted into short address. Unsupported")
}

func isShortPayloadSupportedArgsLen(argLen int) bool {
	if argLen >= shortPayloadSupportedArgsLens[0] && argLen <= shortPayloadSupportedArgsLens[1] {
		return true
	}
	return false
}

// Deprecated: Old full address format is deprecated because a flaw has been found in its encoding method
// bech32, which could enable attackers to generate valid but unexpected addresses.
// For more please check https://github.com/nervosnetwork/rfcs/blob/master/rfcs/0021-ckb-address-format/0021-ckb-address-format.md
func ConvertScriptToFullAddress(hashType string, mode Mode, script *types.Script) (string, error) {
	payload := hashType + hex.EncodeToString(script.CodeHash.Bytes()) + hex.EncodeToString(script.Args)
	data, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode((string)(mode), data)
}

func ConvertScriptToBech32mFullAddress(mode Mode, script *types.Script) (string, error) {
	hashType, err := types.SerializeHashType(script.HashType)
	if err != nil {
		return "", err
	}

	// Payload: type(00) | code hash | hash type | args
	payload := TYPE_FULL_WITH_BECH32M
	payload += script.CodeHash.Hex()[2:]
	payload += hashType

	payload += common.Bytes2Hex(script.Args)

	dataPart, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.EncodeWithBech32m(string(mode), dataPart)
}

func ConvertToBech32mFullAddress(address string) (string, error) {
	parsedAddress, err := Parse(address)
	if err != nil {
		return "", err
	}
	return ConvertScriptToBech32mFullAddress(parsedAddress.Mode, parsedAddress.Script)
}

// Deprecated: Short address format deprecated because it is limited (only support secp256k1_blake160,
// secp256k1_multisig, anyone_can_pay) and a flaw has been found in its encoding method bech32,
// which could enable attackers to generate valid but unexpected addresses.
// For more please check https://github.com/nervosnetwork/rfcs/blob/master/rfcs/0021-ckb-address-format/0021-ckb-address-format.md
func ConvertToShortAddress(address string) (string, error) {
	parsedAddress, err := Parse(address)
	if err != nil {
		return "", err
	}
	return ConvertScriptToShortAddress(parsedAddress.Mode, parsedAddress.Script)
}

// Deprecated: Old full address format is deprecated because a flaw has been found in its encoding method
// bech32, which could enable attackers to generate valid but unexpected addresses.
// For more please check https://github.com/nervosnetwork/rfcs/blob/master/rfcs/0021-ckb-address-format/0021-ckb-address-format.md
func ConvertToBech32FullAddress(address string) (string, error) {
	parsedAddress, err := Parse(address)
	if err != nil {
		return "", err
	}
	return ConvertScriptToFullAddress(FullTypeFormat, parsedAddress.Mode, parsedAddress.Script)
}

func ConvertPublicToAddress(mode Mode, publicKey string) (string, error) {
	script := &types.Script{
		CodeHash: types.HexToHash(transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH),
		HashType: types.HashTypeType,
		Args:     common.FromHex(publicKey),
	}
	return ConvertScriptToBech32mFullAddress(mode, script)
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
	} else if strings.HasPrefix(payload, "00") {
		addressType = TypeFull
		script = types.Script{
			CodeHash: types.HexToHash(payload[2:66]),
			Args:     common.Hex2Bytes(payload[68:]),
		}

		hashType, err := types.DeserializeHashType(payload[66:68])
		if err != nil {
			return nil, err
		}

		script.HashType = hashType

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
