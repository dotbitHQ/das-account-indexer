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

func FormatAddressToHex(chainType common.ChainType, addr string) string {
	switch chainType {
	case common.ChainTypeCkb: // todo
		parseAddr, err := address.Parse(addr)
		if err == nil {
			return common.Bytes2Hex(parseAddr.Script.Args)
		}
	case common.ChainTypeEth, common.ChainTypeMixin:
		return addr
	case common.ChainTypeTron:
		if strings.HasPrefix(addr, common.TronBase58PreFix) {
			if addrTron, err := common.TronBase58ToHex(addr); err == nil {
				return addrTron
			}
		}
	}
	return addr
}

func formatAddressToHalfArgs(chainType common.ChainType, addr string) (args []byte) {
	addrHex := FormatAddressToHex(chainType, addr)
	switch chainType {
	case common.ChainTypeCkb:
		args = append(args, common.DasAlgorithmIdCkb.Bytes()...)
		args = append(args, common.Hex2Bytes(addrHex)...)
	case common.ChainTypeEth:
		args = append(args, common.DasAlgorithmIdEth712.Bytes()...)
		args = append(args, common.Hex2Bytes(addrHex)...)
	case common.ChainTypeTron:
		args = append(args, common.DasAlgorithmIdTron.Bytes()...)
		args = append(args, common.Hex2Bytes(strings.TrimPrefix(addrHex, common.TronPreFix))...)
	case common.ChainTypeMixin:
		args = append(args, common.DasAlgorithmIdEd25519.Bytes()...)
		args = append(args, common.Hex2Bytes(addrHex)...)
	}
	return
}

func FormatHexAddressToNormal(chainType common.ChainType, address string) string {
	switch chainType {
	case common.ChainTypeCkb:
		return address // todo
	case common.ChainTypeEth, common.ChainTypeMixin:
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

func FormatOwnerManagerAddressToArgs(oCT, mCT common.ChainType, oA, mA string) (args []byte) {
	ownerArgs := formatAddressToHalfArgs(oCT, oA)
	managerArgs := formatAddressToHalfArgs(mCT, mA)
	args = append(ownerArgs, managerArgs...)
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
	case common.ChainTypeMixin:
		args = common.DasLockEd25519PreFix + strings.TrimPrefix(addr, common.HexPreFix)
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

func FormatDasLockToOwnerAndManager(args []byte) (owner, manager []byte) {
	if len(args) < common.DasLockArgsLen || len(args) > common.DasLockArgsLenMax {
		return
	}
	oID := common.DasAlgorithmId(args[0])
	splitLen := 0
	switch oID {
	case common.DasAlgorithmIdCkb, common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712, common.DasAlgorithmIdTron:
		splitLen = common.DasLockArgsLen / 2
	case common.DasAlgorithmIdEd25519:
		splitLen = common.DasLockArgsLenMax / 2
	default:
		return
	}
	owner = args[:splitLen]
	manager = args[splitLen:]
	return
}

func formatHalfArgsToHexAddress(args []byte) (aId common.DasAlgorithmId, chainType common.ChainType, addr string) {
	aId = common.DasAlgorithmId(args[0])
	switch aId {
	case common.DasAlgorithmIdCkb:
		chainType = common.ChainTypeCkb
		addr = common.HexPreFix + hex.EncodeToString(args[1:])
	case common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712:
		chainType = common.ChainTypeEth
		addr = common.HexPreFix + hex.EncodeToString(args[1:])
	case common.DasAlgorithmIdTron:
		chainType = common.ChainTypeTron
		addr = common.TronPreFix + hex.EncodeToString(args[1:])
	case common.DasAlgorithmIdEd25519:
		chainType = common.ChainTypeMixin
		addr = common.HexPreFix + hex.EncodeToString(args[1:])
	}
	return
}

func FormatDasLockToHexAddress(args []byte) (oID, mID common.DasAlgorithmId, oCT, mCT common.ChainType, oA, mA string) {
	owner, manager := FormatDasLockToOwnerAndManager(args)
	if len(owner) == 0 || len(manager) == 0 {
		return
	}
	oID, oCT, oA = formatHalfArgsToHexAddress(owner)
	mID, mCT, mA = formatHalfArgsToHexAddress(manager)
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
