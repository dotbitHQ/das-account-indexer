package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
)

func (b *BlockParser) ActionCreateApproval(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version edit records tx")
		return
	}
	log.Info("DasActionCreateApproval:", req.BlockNumber, req.TxHash)

	accBuilder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}
	resp.Err = b.DbDao.UpdateAccounts([]map[string]interface{}{{
		"action":       common.SubActionCreateApproval,
		"account_id":   accBuilder.AccountId,
		"outpoint":     common.OutPoint2String(req.TxHash, 0),
		"block_number": req.BlockNumber,
		"status":       tables.AccountStatusApproval,
	}})
	return
}

func (b *BlockParser) ActionDelayApproval(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version edit records tx")
		return
	}
	log.Info("DasActionDelayApproval:", req.BlockNumber, req.TxHash)

	accBuilder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}
	resp.Err = b.DbDao.UpdateAccounts([]map[string]interface{}{{
		"action":       common.SubActionDelayApproval,
		"account_id":   accBuilder.AccountId,
		"outpoint":     common.OutPoint2String(req.TxHash, 0),
		"block_number": req.BlockNumber,
	}})
	return
}

func (b *BlockParser) ActionRevokeApproval(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version edit records tx")
		return
	}
	log.Info("DasActionRevokeApproval:", req.BlockNumber, req.TxHash)

	accBuilder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeOld)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}
	resp.Err = b.DbDao.UpdateAccounts([]map[string]interface{}{{
		"action":       common.SubActionRevokeApproval,
		"account_id":   accBuilder.AccountId,
		"outpoint":     common.OutPoint2String(req.TxHash, 0),
		"block_number": req.BlockNumber,
		"status":       tables.AccountStatusNormal,
	}})
	return
}

func (b *BlockParser) ActionFulfillApproval(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version edit records tx")
		return
	}
	log.Info("DasActionFulfillApproval:", req.BlockNumber, req.TxHash)

	accBuilder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeOld)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}

	approval := accBuilder.AccountApproval
	switch approval.Action {
	case witness.AccountApprovalActionTransfer:
		owner, manager, err := b.DasCore.Daf().ScriptToHex(approval.Params.Transfer.ToLock)
		if err != nil {
			resp.Err = fmt.Errorf("ScriptToHex err: %s", err.Error())
			return
		}
		resp.Err = b.DbDao.UpdateAccounts([]map[string]interface{}{{
			"action":               common.SubActionFullfillApproval,
			"outpoint":             common.OutPoint2String(req.TxHash, 0),
			"block_number":         req.BlockNumber,
			"status":               tables.AccountStatusNormal,
			"owner":                owner.AddressHex,
			"owner_chain_type":     owner.ChainType,
			"owner_algorithm_id":   owner.DasAlgorithmId,
			"manager":              manager.AddressHex,
			"manager_chain_type":   manager.ChainType,
			"manager_algorithm_id": manager.DasAlgorithmId,
		}})
	}
	return
}
