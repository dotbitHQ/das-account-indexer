package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
)

func (b *BlockParser) ActionDeclareReverseRecord(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameReverseRecordCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersionTx err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	account := string(req.Tx.OutputsData[0])
	ownerHex, _, err := b.DasCore.Daf().ArgsToHex(req.Tx.Outputs[0].Lock.Args)
	if err != nil {
		resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
		return
	}

	reverseInfo := tables.TableReverseInfo{
		BlockNumber:    req.BlockNumber,
		BlockTimestamp: req.BlockTimestamp,
		Outpoint:       common.OutPoint2String(req.TxHash, 0),
		AlgorithmId:    ownerHex.DasAlgorithmId,
		ChainType:      ownerHex.ChainType,
		Address:        ownerHex.AddressHex,
		Account:        account,
		Capacity:       req.Tx.Outputs[0].Capacity,
	}

	if err := b.DbDao.CreateReverseInfo(&reverseInfo); err != nil {
		resp.Err = fmt.Errorf("DeclareReverseRecord err: %s", err.Error())
		return
	}

	return
}

func (b *BlockParser) ActionRedeclareReverseRecord(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameReverseRecordCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersionTx err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	account := string(req.Tx.OutputsData[0])
	ownerHex, _, err := b.DasCore.Daf().ArgsToHex(req.Tx.Outputs[0].Lock.Args)
	if err != nil {
		resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
		return
	}

	reverseInfo := tables.TableReverseInfo{
		BlockNumber:    req.BlockNumber,
		BlockTimestamp: req.BlockTimestamp,
		Outpoint:       common.OutPoint2String(req.TxHash, 0),
		AlgorithmId:    ownerHex.DasAlgorithmId,
		ChainType:      ownerHex.ChainType,
		Address:        ownerHex.AddressHex,
		Account:        account,
		Capacity:       req.Tx.Outputs[0].Capacity,
	}
	lastOutpoint := common.OutPointStruct2String(req.Tx.Inputs[0].PreviousOutput)

	if err := b.DbDao.UpdateReverseInfo(&reverseInfo, lastOutpoint); err != nil {
		resp.Err = fmt.Errorf("UpdateReverseInfo err: %s", err.Error())
		return
	}
	return
}

func (b *BlockParser) ActionRetractReverseRecord(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	res, err := b.DasCore.Client().GetTransaction(b.Ctx, req.Tx.Inputs[0].PreviousOutput.TxHash)
	if err != nil {
		resp.Err = fmt.Errorf("GetTransaction err: %s", err.Error())
		return
	}
	if isCV, err := isCurrentVersionTx(res.Transaction, common.DasContractNameReverseRecordCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersionTx err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	var outpoints []string
	for _, v := range req.Tx.Inputs {
		outpoints = append(outpoints, common.OutPointStruct2String(v.PreviousOutput))
	}

	if err := b.DbDao.DeleteReverseInfo(outpoints); err != nil {
		resp.Err = fmt.Errorf("DeleteReverseInfo err: %s", err.Error())
		return
	}
	return
}
