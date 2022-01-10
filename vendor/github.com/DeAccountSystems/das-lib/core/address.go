package core

import (
	"encoding/hex"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

// Deprecated: format normal ckb lock to address will delete in future
func FormatNormalCkbLockToAddress(net common.DasNetType, args []byte) (addr string, err error) {
	lockScript := common.GetNormalLockScript(common.Bytes2Hex(args))
	netMode := address.Mainnet
	switch net {
	case common.DasNetTypeMainNet:
		netMode = address.Mainnet
	case common.DasNetTypeTestnet2, common.DasNetTypeTestnet3:
		netMode = address.Testnet
	}

	addr, err = common.ConvertScriptToAddress(netMode, lockScript)
	return
}

func FormatDasLockToNormalAddress(args []byte) (oID, mID common.DasAlgorithmId, oCT, mCT common.ChainType, oA, mA string) {
	oID, mID, oCT, mCT, oA, mA = FormatDasLockToHexAddress(args)
	if oCT == common.ChainTypeTron {
		oA, _ = common.TronHexToBase58(oA)
	}
	if mCT == common.ChainTypeTron {
		mA, _ = common.TronHexToBase58(mA)
	}
	return
}

func FormatDasLockToHexAddress(args []byte) (oID, mID common.DasAlgorithmId, oCT, mCT common.ChainType, oA, mA string) {
	if len(args) < common.DasLockArgsLen {
		return
	}
	ownerBytes := args[:common.DasLockArgsLen/2]
	oID = common.DasAlgorithmId(ownerBytes[0])
	switch oID {
	case common.DasAlgorithmIdCkb:
		oCT = common.ChainTypeCkb
		oA = common.HexPreFix + hex.EncodeToString(ownerBytes[1:])
	case common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712:
		oCT = common.ChainTypeEth
		oA = common.HexPreFix + hex.EncodeToString(ownerBytes[1:])
	case common.DasAlgorithmIdTron:
		oCT = common.ChainTypeTron
		oA = common.TronPreFix + hex.EncodeToString(ownerBytes[1:])
	}

	managerBytes := args[common.DasLockArgsLen/2:]
	mID = common.DasAlgorithmId(managerBytes[0])
	switch mID {
	case common.DasAlgorithmIdCkb:
		mCT = common.ChainTypeCkb
		mA = common.HexPreFix + hex.EncodeToString(managerBytes[1:])
	case common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712:
		mCT = common.ChainTypeEth
		mA = common.HexPreFix + hex.EncodeToString(managerBytes[1:])
	case common.DasAlgorithmIdTron:
		mCT = common.ChainTypeTron
		mA = common.TronPreFix + hex.EncodeToString(managerBytes[1:])
	}
	return
}

func (d *DasCore) FormatAddressToDasLockScript(chainType common.ChainType, addr string, is712 bool) (lockScript, typeScript *types.Script, e error) {
	addr = FormatAddressToHex(chainType, addr)
	args := ""
	switch chainType {
	case common.ChainTypeCkb:
		if parseAddr, err := address.Parse(addr); err != nil {
			e = err
			return
		} else {
			args = common.DasLockCkbPreFix + hex.EncodeToString(parseAddr.Script.Args)
		}
	case common.ChainTypeEth:
		if is712 {
			args = common.DasLockEth712PreFix + strings.TrimPrefix(addr, common.HexPreFix)
			if contractInfo, err := GetDasContractInfo(common.DasContractNameBalanceCellType); err != nil {
				e = err
				return
			} else {
				typeScript = contractInfo.ToScript(nil)
			}
		} else {
			args = common.DasLockEthPreFix + strings.TrimPrefix(addr, common.HexPreFix)
		}
	case common.ChainTypeTron:
		args = common.DasLockTronPreFix + strings.TrimPrefix(addr, common.TronPreFix)
	default:
		e = fmt.Errorf("unknow chain type [%d]", chainType)
		return
	}
	args = common.HexPreFix + args + args

	if contractInfo, err := GetDasContractInfo(common.DasContractNameDispatchCellType); err != nil {
		e = err
		return
	} else {
		lockScript = contractInfo.ToScript(common.Hex2Bytes(args))
		return
	}
}

func FormatAddressToHex(chainType common.ChainType, addr string) string {
	switch chainType {
	case common.ChainTypeCkb:
		parseAddr, err := address.Parse(addr)
		if err == nil {
			return common.Bytes2Hex(parseAddr.Script.Args)
		}
	case common.ChainTypeBtc, common.ChainTypeEth:
		return addr
	case common.ChainTypeTron:
		if strings.HasPrefix(addr, common.TronBase58PreFix) {
			if addr, err := common.TronBase58ToHex(addr); err == nil {
				return addr
			}
		}
	}
	return addr
}

func FormatHexAddressToNormal(chainType common.ChainType, address string) string {
	switch chainType {
	case common.ChainTypeCkb, common.ChainTypeBtc, common.ChainTypeEth:
		return address
	case common.ChainTypeTron:
		if strings.HasPrefix(address, common.TronPreFix) {
			if addr, err := common.TronHexToBase58(address); err == nil {
				return addr
			}
		}
	}
	return address
}

func FormatOwnerManagerAddressToArgs(oCT, mCT common.ChainType, oA, mA string) []byte {
	oA = FormatAddressToHex(oCT, oA)
	mA = FormatAddressToHex(mCT, mA)
	var args []byte
	switch oCT {
	case common.ChainTypeCkb:
		args = append(args, common.DasAlgorithmIdCkb.Bytes()...)
		args = append(args, common.Hex2Bytes(oA)...)
	case common.ChainTypeEth:
		args = append(args, common.DasAlgorithmIdEth712.Bytes()...)
		args = append(args, common.Hex2Bytes(oA)...)
	case common.ChainTypeTron:
		args = append(args, common.DasAlgorithmIdTron.Bytes()...)
		args = append(args, common.Hex2Bytes(strings.TrimPrefix(oA, common.TronPreFix))...)
	}
	switch mCT {
	case common.ChainTypeCkb:
		args = append(args, common.DasAlgorithmIdCkb.Bytes()...)
		args = append(args, common.Hex2Bytes(mA)...)
	case common.ChainTypeEth:
		args = append(args, common.DasAlgorithmIdEth712.Bytes()...)
		args = append(args, common.Hex2Bytes(mA)...)
	case common.ChainTypeTron:
		args = append(args, common.DasAlgorithmIdTron.Bytes()...)
		args = append(args, common.Hex2Bytes(strings.TrimPrefix(mA, common.TronPreFix))...)
	}
	return args
}
