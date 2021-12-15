package core

import (
	"bytes"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/collector"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

func getDasAlgorithmId(chainType common.ChainType) (list []common.DasAlgorithmId) {
	switch chainType {
	case common.ChainTypeCkb:
		list = append(list, common.DasAlgorithmIdCkb)
	case common.ChainTypeEth:
		list = append(list, common.DasAlgorithmIdEth712, common.DasAlgorithmIdEth)
	case common.ChainTypeTron:
		list = append(list, common.DasAlgorithmIdTron)
	}
	return
}

func dasAlgorithmIdAddressToArgs(dasAlgorithmId common.DasAlgorithmId, addr string) ([]byte, error) {
	switch dasAlgorithmId {
	case common.DasAlgorithmIdCkb:
		if parseAddr, err := address.Parse(addr); err != nil {
			return nil, err
		} else {
			return parseAddr.Script.Args, nil
		}
	case common.DasAlgorithmIdEth712, common.DasAlgorithmIdEth:
		return append(dasAlgorithmId.Bytes(), common.Hex2Bytes(addr)...), nil
	case common.DasAlgorithmIdTron:
		if strings.HasPrefix(addr, common.TronBase58PreFix) {
			if addrHex, err := common.TronBase58ToHex(addr); err != nil {
				return nil, err
			} else {
				return append(dasAlgorithmId.Bytes(), common.Hex2Bytes(strings.TrimPrefix(addrHex, common.TronPreFix))...), nil
			}
		}
		return append(dasAlgorithmId.Bytes(), common.Hex2Bytes(strings.TrimPrefix(addr, common.TronPreFix))...), nil
	}
	return nil, fmt.Errorf("not exist DasAlgorithmId: %d", dasAlgorithmId)
}

func getDasAccountArgs(ownerDasAlgorithmId, managerDasAlgorithmId common.DasAlgorithmId, ownerAddress, managerAddress string) ([]byte, error) {
	ownerArgs, err := dasAlgorithmIdAddressToArgs(ownerDasAlgorithmId, ownerAddress)
	if err != nil {
		return nil, err
	}
	managerArgs, err := dasAlgorithmIdAddressToArgs(managerDasAlgorithmId, managerAddress)
	if err != nil {
		return nil, err
	}
	return append(ownerArgs, managerArgs...), nil
}

func (d *DasCore) GetAccountCellOnChain(ownerChainType, managerChainType common.ChainType, ownerAddress, managerAddress, account string) (*indexer.LiveCell, error) {
	accountId := common.GetAccountIdByAccount(account)
	// account cell code hash
	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	contractAcc, err := GetDasContractInfo(common.DasContractNameAccountCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	// 0303 0305 0503 0505
	ownerDasAlgorithmIds := getDasAlgorithmId(ownerChainType)
	managerDasAlgorithmIds := getDasAlgorithmId(managerChainType)
	var lockScriptInSearchKey []*types.Script
	for _, ov := range ownerDasAlgorithmIds {
		for _, mv := range managerDasAlgorithmIds {
			if args, err := getDasAccountArgs(ov, mv, ownerAddress, managerAddress); err != nil {
				return nil, fmt.Errorf("getDasAccountArgs err: %s", err.Error())
			} else {
				dasLockScript := contractDispatch.ToScript(args)
				log.Info("search account cell:", account, common.Bytes2Hex(args))
				lockScriptInSearchKey = append(lockScriptInSearchKey, dasLockScript)
			}
		}
	}
	// search
	accountCellDataLenMin := uint64(common.ExpireTimeEndIndex + 5)   // .bit for 4 and at least 1 length
	accountCellDataLenMax := uint64(common.ExpireTimeEndIndex + 100) // consider for emoji
	for _, dasLockScript := range lockScriptInSearchKey {
		searchKey := &indexer.SearchKey{
			Script:     contractAcc.ToScript(nil),
			ScriptType: indexer.ScriptTypeType,
			Filter: &indexer.CellsFilter{
				Script:             dasLockScript,
				OutputDataLenRange: &[2]uint64{accountCellDataLenMin, accountCellDataLenMax},
			},
		}
		co := collector.NewLiveCellCollector(d.client, searchKey, indexer.SearchOrderAsc, indexer.SearchLimit, "")
		co.TypeScript = searchKey.Script
		iterator, err := co.Iterator()
		if err != nil {
			return nil, fmt.Errorf("get cell err: %s", err.Error())
		}
		for iterator.HasNext() {
			liveCell, err := iterator.CurrentItem()
			if err != nil {
				return nil, fmt.Errorf("CurrentItem err:%s", err.Error())
			}
			searchAccountId, err := common.OutputDataToAccountId(liveCell.OutputData)
			if err != nil {
				continue
			}
			if bytes.Compare(searchAccountId, accountId) == 0 {
				log.Info("get account:", account, liveCell.OutPoint.TxHash, liveCell.OutPoint.Index)
				return liveCell, nil
			}
			if err = iterator.Next(); err != nil {
				return nil, fmt.Errorf("next err:%s", err.Error())
			}
		}
	}
	return nil, fmt.Errorf("not exist acc: %s", account)
}

func (d *DasCore) GetAccountCellOnChainByAlgorithmId(oID, mID common.DasAlgorithmId, oA, mA, account string) (*indexer.LiveCell, error) {
	accountId := common.GetAccountIdByAccount(account)
	// account cell code hash
	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	contractAcc, err := GetDasContractInfo(common.DasContractNameAccountCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	args := append(oID.Bytes(), common.Hex2Bytes(oA)...)
	args = append(args, mID.Bytes()...)
	args = append(args, common.Hex2Bytes(mA)...)
	log.Info("GetAccountCellOnChainByAlgorithmId:", common.Bytes2Hex(args))
	dasLockScript := contractDispatch.ToScript(args)

	// search
	accountCellDataLenMin := uint64(common.ExpireTimeEndIndex + 5)
	accountCellDataLenMax := uint64(common.ExpireTimeEndIndex + 100)
	searchKey := &indexer.SearchKey{
		Script:     contractAcc.ToScript(nil),
		ScriptType: indexer.ScriptTypeType,
		Filter: &indexer.CellsFilter{
			Script:             dasLockScript,
			OutputDataLenRange: &[2]uint64{accountCellDataLenMin, accountCellDataLenMax},
		},
	}
	co := collector.NewLiveCellCollector(d.client, searchKey, indexer.SearchOrderAsc, indexer.SearchLimit, "")
	co.TypeScript = searchKey.Script
	iterator, err := co.Iterator()
	if err != nil {
		return nil, fmt.Errorf("get cell err: %s", err.Error())
	}
	for iterator.HasNext() {
		liveCell, err := iterator.CurrentItem()
		if err != nil {
			return nil, fmt.Errorf("CurrentItem err:%s", err.Error())
		}
		searchAccountId, err := common.OutputDataToAccountId(liveCell.OutputData)
		if err != nil {
			continue
		}
		if bytes.Compare(searchAccountId, accountId) == 0 {
			log.Info("get account:", account, liveCell.OutPoint.TxHash, liveCell.OutPoint.Index)
			return liveCell, nil
		}
		if err = iterator.Next(); err != nil {
			return nil, fmt.Errorf("next err:%s", err.Error())
		}
	}
	return nil, fmt.Errorf("not exist acc: %s", account)
}

func (d *DasCore) UpdateAccountCellDasLockToEip712(cell *indexer.LiveCell) {
	if cell == nil || cell.Output == nil || cell.Output.Type == nil {
		return
	}
	accContract, err := GetDasContractInfo(common.DasContractNameAccountCellType)
	if err != nil {
		log.Error("GetDasContractInfo err:", err.Error())
		return
	}
	dasLockContract, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		log.Error("GetDasContractInfo err:", err.Error())
		return
	}
	dasLock := cell.Output.Lock.CodeHash.Hex()
	dasType := cell.Output.Type.CodeHash.Hex()
	if dasLock != dasLockContract.ContractTypeId.Hex() || dasType != accContract.ContractTypeId.Hex() {
		return
	}

	if common.DasAlgorithmId(cell.Output.Lock.Args[0]) == common.DasAlgorithmIdEth {
		cell.Output.Lock.Args[0] = byte(common.DasAlgorithmIdEth712)
	}
	if common.DasAlgorithmId(cell.Output.Lock.Args[common.DasLockArgsLen/2]) == common.DasAlgorithmIdEth {
		cell.Output.Lock.Args[common.DasLockArgsLen/2] = byte(common.DasAlgorithmIdEth712)
	}
}

func UpdateAccountCellDasLockToEip712(cell *indexer.LiveCell) error {
	if cell == nil || cell.Output == nil || cell.Output.Type == nil {
		return nil
	}
	accContract, err := GetDasContractInfo(common.DasContractNameAccountCellType)
	if err != nil {
		return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	dasLockContract, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	dasLock := cell.Output.Lock.CodeHash.Hex()
	dasType := cell.Output.Type.CodeHash.Hex()
	if dasLock != dasLockContract.ContractTypeId.Hex() || dasType != accContract.ContractTypeId.Hex() {
		return nil
	}

	if common.DasAlgorithmId(cell.Output.Lock.Args[0]) == common.DasAlgorithmIdEth {
		cell.Output.Lock.Args[0] = byte(common.DasAlgorithmIdEth712)
	}
	if common.DasAlgorithmId(cell.Output.Lock.Args[common.DasLockArgsLen/2]) == common.DasAlgorithmIdEth {
		cell.Output.Lock.Args[common.DasLockArgsLen/2] = byte(common.DasAlgorithmIdEth712)
	}
	return nil
}
