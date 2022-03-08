package txbuilder

import (
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

func (d *DasTxBuilder) BuildMMJsonObj(evmChainId int64) (*common.MMJsonObj, error) {
	var res common.MMJsonObj
	if err := json.Unmarshal([]byte(common.MMJsonObjStr), &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal err: %s", err.Error())
	}

	res.Domain.ChainID = evmChainId
	if res.Domain.ChainID == 0 {
		res.Domain.ChainID = 1
		if d.dasCore.NetType() != common.DasNetTypeMainNet {
			res.Domain.ChainID = 5
		}
	}

	inputsCapacity, err := d.getCapacityFromInputs()
	if err != nil {
		return nil, fmt.Errorf("getCapacityFromInputs err: %s", err.Error())
	}
	outputsCapacity := d.Transaction.OutputsCapacity()
	feeCapacity := inputsCapacity - outputsCapacity

	res.Message.InputsCapacity = fmt.Sprintf("%s CKB", common.Capacity2Str(inputsCapacity))
	res.Message.OutputsCapacity = fmt.Sprintf("%s CKB", common.Capacity2Str(outputsCapacity))
	res.Message.Fee = fmt.Sprintf("%s CKB", common.Capacity2Str(feeCapacity))

	inputs, err := d.getInputsMMJsonCellInfo()
	if err != nil {
		return nil, fmt.Errorf("getInputsMMJsonCellInfo err: %s", err.Error())
	}
	outputs, err := d.getOutputsMMJsonCellInfo()
	if err != nil {
		return nil, fmt.Errorf("getOutputsMMJsonCellInfo err: %s", err.Error())
	}
	res.Message.Inputs = inputs
	res.Message.Outputs = outputs

	// must be the last one to be executed
	action, dasMessage, err := d.getMMJsonActionAndMessage()
	if err != nil {
		return nil, fmt.Errorf("getMMJsonActionAndMessage err: %s", err.Error())
	}
	res.Message.Action = action
	res.Message.DasMessage = dasMessage

	return &res, nil
}

func (d *DasTxBuilder) getInputsMMJsonCellInfo() ([]common.MMJsonCellInfo, error) {
	var cellList []*types.CellInfo
	for _, v := range d.Transaction.Inputs {
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return nil, fmt.Errorf("getInputCell err: %s", err.Error())
		}
		cellList = append(cellList, item.Cell)
	}
	return d.getMMJsonCellInfo(cellList, common.DataTypeOld)
}

func (d *DasTxBuilder) getOutputsMMJsonCellInfo() ([]common.MMJsonCellInfo, error) {
	var cellList []*types.CellInfo
	for i, output := range d.Transaction.Outputs {
		cell := types.CellInfo{
			Output: output,
			Data: &types.CellData{
				Content: d.Transaction.OutputsData[i],
			},
		}
		cellList = append(cellList, &cell)
	}
	return d.getMMJsonCellInfo(cellList, common.DataTypeNew)
}

func (d *DasTxBuilder) getMMJsonCellInfo(cellList []*types.CellInfo, dataType common.DataType) ([]common.MMJsonCellInfo, error) {
	list := make([]common.MMJsonCellInfo, 0)
	for _, v := range cellList {
		var item common.MMJsonCellInfo
		item.Capacity = fmt.Sprintf("%s CKB", common.Capacity2Str(v.Output.Capacity))
		if v.Data != nil {
			item.Data = common.GetMaxHashLenData(v.Data.Content)
		}
		item.ExtraData = ""
		item.TypeStr = ""
		item.LockStr = ""
		if v.Output.Type == nil {
			continue
		}
		if lockContractName, ok := core.DasContractByTypeIdMap[v.Output.Lock.CodeHash.Hex()]; ok {
			item.LockStr = common.GetMaxHashLenScript(v.Output.Lock, lockContractName)
		}
		if typeContractName, ok := core.DasContractByTypeIdMap[v.Output.Type.CodeHash.Hex()]; ok {
			if typeContractName == common.DasContractNameBalanceCellType {
				continue
			}
			item.TypeStr = common.GetMaxHashLenScript(v.Output.Type, typeContractName)
			switch typeContractName {
			case common.DasContractNameAccountSaleCellType:
				builder, err := witness.AccountSaleCellDataBuilderFromTx(d.Transaction, dataType)
				if err != nil {
					return nil, fmt.Errorf("AccountSaleCellDataBuilderFromTx err: %s", err.Error())
				}
				d.salePrice = builder.Price
			case common.DasContractNameAccountCellType:
				builder, err := witness.AccountCellDataBuilderFromTx(d.Transaction, dataType)
				if err != nil {
					return nil, fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
				}
				d.account = builder.Account
				expiredAt, err := common.GetAccountCellExpiredAtFromOutputData(v.Data.Content)
				if err != nil {
					return nil, fmt.Errorf("GetAccountCellExpiredAtFromOutputData err: %s", err.Error())
				}
				item.Data = fmt.Sprintf("{ account: %s, expired_at: %d }", d.account, expiredAt)
				item.ExtraData = fmt.Sprintf("{ status: %d, records_hash: %s }", builder.Status, common.Bytes2Hex(builder.RecordsHashBys))
			case common.DasContractNameReverseRecordCellType:
				d.account = string(v.Data.Content)
			case common.DASContractNameOfferCellType:
				d.offers++
			}
		}
		list = append(list, item)
	}
	return list, nil
}

func (d *DasTxBuilder) getMMJsonActionAndMessage() (*common.MMJsonAction, string, error) {
	var action common.MMJsonAction
	actionDataBuilder, err := witness.ActionDataBuilderFromTx(d.Transaction)
	if err != nil {
		return nil, "", fmt.Errorf("ActionDataBuilderFromTx err: %s", err.Error())
	}
	action.Action = actionDataBuilder.Action
	action.Params = actionDataBuilder.ParamsStr

	dasMessage := ""
	switch action.Action {
	case common.DasActionEditManager:
		dasMessage = fmt.Sprintf("EDIT MANAGER OF ACCOUNT %s", d.account)
	case common.DasActionEditRecords:
		dasMessage = fmt.Sprintf("EDIT RECORDS OF ACCOUNT %s", d.account)
	case common.DasActionTransferAccount:
		builder, err := witness.AccountCellDataBuilderFromTx(d.Transaction, common.DataTypeNew)
		if err != nil {
			return nil, "", fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		}
		_, _, _, _, oA, _ := core.FormatDasLockToNormalAddress(d.Transaction.Outputs[builder.Index].Lock.Args)
		dasMessage = fmt.Sprintf("TRANSFER THE ACCOUNT %s TO %s", d.account, oA)
	case common.DasActionTransfer, common.DasActionWithdrawFromWallet:
		dasMessage, err = d.getWithdrawDasMessage()
		if err != nil {
			return nil, "", fmt.Errorf("getWithdrawDasMessage err: %s", err.Error())
		}
	case common.DasActionStartAccountSale:
		dasMessage = fmt.Sprintf("SELL %s FOR %s CKB", d.account, common.Capacity2Str(d.salePrice))
	case common.DasActionEditAccountSale:
		dasMessage = fmt.Sprintf("EDIT SALE INFO, CURRENT PRICE IS %s CKB", common.Capacity2Str(d.salePrice))
	case common.DasActionCancelAccountSale:
		dasMessage = fmt.Sprintf("CANCEL SALE OF %s", d.account)
	case common.DasActionBuyAccount:
		dasMessage = fmt.Sprintf("BUY %s WITH %s CKB", d.account, common.Capacity2Str(d.salePrice))
	case common.DasActionDeclareReverseRecord:
		_, _, _, _, oA, _ := core.FormatDasLockToNormalAddress(d.Transaction.Outputs[0].Lock.Args)
		dasMessage = fmt.Sprintf("DECLARE A REVERSE RECORD FROM %s TO %s", oA, d.account)
	case common.DasActionRedeclareReverseRecord:
		_, _, _, _, oA, _ := core.FormatDasLockToNormalAddress(d.Transaction.Outputs[0].Lock.Args)
		dasMessage = fmt.Sprintf("REDECLARE A REVERSE RECORD FROM %s TO %s", oA, d.account)
	case common.DasActionRetractReverseRecord:
		_, _, _, _, oA, _ := core.FormatDasLockToNormalAddress(d.Transaction.Outputs[0].Lock.Args)
		dasMessage = fmt.Sprintf("RETRACT REVERSE RECORDS ON %s", oA)
	case common.DasActionMakeOffer:
		builder, err := witness.OfferCellDataBuilderFromTx(d.Transaction, common.DataTypeNew)
		if err != nil {
			return nil, "", fmt.Errorf("OfferCellDataBuilderFromTx err: %s", err.Error())
		}
		dasMessage = fmt.Sprintf("MAKE AN OFFER ON %s WITH %s CKB", builder.Account, common.Capacity2Str(builder.Price))
	case common.DasActionEditOffer:
		builderOld, err := witness.OfferCellDataBuilderFromTx(d.Transaction, common.DataTypeOld)
		if err != nil {
			return nil, "", fmt.Errorf("OfferCellDataBuilderFromTx err: %s", err.Error())
		}
		builder, err := witness.OfferCellDataBuilderFromTx(d.Transaction, common.DataTypeNew)
		if err != nil {
			return nil, "", fmt.Errorf("OfferCellDataBuilderFromTx err: %s", err.Error())
		}
		dasMessage = fmt.Sprintf("CHANGE THE OFFER ON %s FROM %s CKB TO %s CKB", builder.Account, common.Capacity2Str(builderOld.Price), common.Capacity2Str(builder.Price))
	case common.DasActionCancelOffer:
		dasMessage = fmt.Sprintf("CANCEL %d OFFER(S)", d.offers)
	case common.DasActionAcceptOffer:
		builder, err := witness.OfferCellDataBuilderFromTx(d.Transaction, common.DataTypeOld)
		if err != nil {
			return nil, "", fmt.Errorf("OfferCellDataBuilderFromTx err: %s", err.Error())
		}
		dasMessage = fmt.Sprintf("ACCEPT THE OFFER ON %s WITH %s CKB", builder.Account, common.Capacity2Str(builder.Price))
	case common.DasActionEnableSubAccount:
		dasMessage = fmt.Sprintf("")
	default:
		return nil, "", fmt.Errorf("not support action[%s]", action)
	}

	return &action, dasMessage, nil
}

func (d *DasTxBuilder) getWithdrawDasMessage() (string, error) {
	inputsCapacity, err := d.getCapacityFromInputs()
	if err != nil {
		return "", fmt.Errorf("getCapacityFromInputs err: %s", err.Error())
	}
	item, err := d.getInputCell(d.Transaction.Inputs[0].PreviousOutput)
	if err != nil {
		return "", fmt.Errorf("getInputCell err: %s", err.Error())
	}
	_, _, _, _, oA, _ := core.FormatDasLockToNormalAddress(item.Cell.Output.Lock.Args)
	//dasMessage := fmt.Sprintf("%s:%s(%s CKB) TO ", oCT.String(), oA, common.Capacity2Str(inputsCapacity))
	dasMessage := fmt.Sprintf("%s(%s CKB) TO ", oA, common.Capacity2Str(inputsCapacity))

	// need merge outputs the capacity with the same lock script
	var mapOutputs = make(map[string]uint64)
	var sortList = make([]string, 0)

	dasLock, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return "", fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	mod := address.Testnet
	if d.dasCore.NetType() == common.DasNetTypeMainNet {
		mod = address.Mainnet
	}
	for _, v := range d.Transaction.Outputs {
		chainStrTmp, receiverAddr := "", ""
		if v.Lock.CodeHash.Hex() == dasLock.ContractTypeId.Hex() {
			_, _, oCTTmp, _, oATmp, _ := core.FormatDasLockToNormalAddress(v.Lock.Args)
			chainStrTmp, receiverAddr = oCTTmp.String(), oATmp
		} else {
			addr, _ := common.ConvertScriptToAddress(mod, v.Lock)
			chainStrTmp, receiverAddr = "CKB", addr
		}

		key := fmt.Sprintf("%s-%s", chainStrTmp, receiverAddr)
		if c, ok := mapOutputs[key]; ok {
			mapOutputs[key] = c + v.Capacity
		} else {
			mapOutputs[key] = v.Capacity
			sortList = append(sortList, key)
		}
	}
	for _, v := range sortList {
		capacity := mapOutputs[v]
		tmp := strings.Split(v, "-")
		//dasMessage = dasMessage + fmt.Sprintf("%s:%s(%s CKB), ", tmp[0], tmp[1], common.Capacity2Str(capacity))
		dasMessage = dasMessage + fmt.Sprintf("%s(%s CKB), ", tmp[1], common.Capacity2Str(capacity))
	}

	return fmt.Sprintf("TRANSFER FROM %s", dasMessage[:len(dasMessage)-2]), nil
}
