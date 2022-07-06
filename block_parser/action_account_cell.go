package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"strconv"
)

func (b *BlockParser) ActionUpdateAccountInfo(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersionTx err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	builder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}
	ownerHex, managerHex, err := b.DasCore.Daf().ArgsToHex(req.Tx.Outputs[builder.Index].Lock.Args)
	if err != nil {
		resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
		return
	}

	accountInfo := tables.TableAccountInfo{
		BlockNumber:        req.BlockNumber,
		BlockTimestamp:     req.BlockTimestamp,
		Outpoint:           common.OutPoint2String(req.TxHash, uint(builder.Index)),
		AccountId:          builder.AccountId,
		NextAccountId:      builder.NextAccountId,
		Account:            builder.Account,
		OwnerChainType:     ownerHex.ChainType,
		Owner:              ownerHex.AddressHex,
		OwnerAlgorithmId:   ownerHex.DasAlgorithmId,
		ManagerChainType:   managerHex.ChainType,
		Manager:            managerHex.AddressHex,
		ManagerAlgorithmId: managerHex.DasAlgorithmId,
		Status:             tables.AccountStatus(builder.Status),
		RegisteredAt:       builder.RegisteredAt,
		ExpiredAt:          builder.ExpiredAt,
	}

	var records []tables.TableRecordsInfo
	list := builder.Records
	for _, v := range list {
		records = append(records, tables.TableRecordsInfo{
			Account:   builder.Account,
			AccountId: builder.AccountId,
			Key:       v.Key,
			Type:      v.Type,
			Label:     v.Label,
			Value:     v.Value,
			Ttl:       strconv.FormatUint(uint64(v.TTL), 10),
		})
	}

	if err = b.DbDao.UpdateAccountInfo(&accountInfo, records); err != nil {
		resp.Err = fmt.Errorf("UpdateAccountInfo err: %s", err.Error())
		return
	}

	return
}

func (b *BlockParser) ActionRecycleExpiredAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	var builder *witness.AccountCellDataBuilder
	builderMap, err := witness.AccountCellDataBuilderMapFromTx(req.Tx, common.DataTypeOld)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderMapFromTx err: %s", err.Error())
		return
	}
	for _, v := range builderMap {
		if v.Index == 1 {
			builder = v
		}
	}
	if builder == nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilder is nil")
		return
	}

	accountId, err := common.OutputDataToAccountId(req.Tx.OutputsData[0])
	if err != nil {
		resp.Err = fmt.Errorf("OutputDataToAccountId err: %s", err.Error())
		return
	}
	nextAccountId, err := common.GetAccountCellNextAccountIdFromOutputData(req.Tx.OutputsData[0])
	if err != nil {
		resp.Err = fmt.Errorf("GetAccountCellNextAccountIdFromOutputData err: %s", err.Error())
		return
	}
	accountInfo := tables.TableAccountInfo{
		BlockNumber:    req.BlockNumber,
		BlockTimestamp: req.BlockTimestamp,
		Outpoint:       common.OutPoint2String(req.TxHash, 0),
		AccountId:      common.Bytes2Hex(accountId),
		NextAccountId:  common.Bytes2Hex(nextAccountId),
	}

	if err = b.DbDao.RecycleExpiredAccount(accountInfo, builder.AccountId, builder.EnableSubAccount); err != nil {
		resp.Err = fmt.Errorf("RecycleExpiredAccount err: %s", err.Error())
		return
	}

	return
}
