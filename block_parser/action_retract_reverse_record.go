package block_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
)

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
