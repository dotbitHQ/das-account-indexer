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
		OwnerSubAid:        ownerHex.DasSubAlgorithmId,
		ManagerChainType:   managerHex.ChainType,
		Manager:            managerHex.AddressHex,
		ManagerAlgorithmId: managerHex.DasAlgorithmId,
		ManagerSubAid:      managerHex.DasSubAlgorithmId,
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

	previousBuilder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}
	accountInfo := tables.TableAccountInfo{
		BlockNumber:    req.BlockNumber,
		BlockTimestamp: req.BlockTimestamp,
		Outpoint:       common.OutPoint2String(req.TxHash, 0),
		AccountId:      previousBuilder.AccountId,
		NextAccountId:  previousBuilder.NextAccountId,
	}

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

	if err = b.DbDao.RecycleExpiredAccount(accountInfo, builder.AccountId, builder.EnableSubAccount); err != nil {
		resp.Err = fmt.Errorf("RecycleExpiredAccount err: %s", err.Error())
		return
	}

	return
}

func (b *BlockParser) ActionAccountUpgrade(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version account cross chain tx")
		return
	}
	log.Info("ActionAccountUpgrade:", req.BlockNumber, req.TxHash, req.Action)

	builder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}

	didEntity, err := witness.TxToOneDidEntity(req.Tx, witness.SourceTypeOutputs)
	if err != nil {
		resp.Err = fmt.Errorf("TxToOneDidEntity err: %s", err.Error())
		return
	}
	didCellArgs := common.Bytes2Hex(req.Tx.Outputs[didEntity.Target.Index].Lock.Args)
	accountInfo := tables.TableAccountInfo{
		BlockNumber: req.BlockNumber,
		Outpoint:    common.OutPoint2String(req.TxHash, 0),
		AccountId:   builder.AccountId,
		Status:      tables.AccountStatus(builder.Status),
	}

	didCellInfo := tables.TableDidCellInfo{
		BlockNumber:  req.BlockNumber,
		Outpoint:     common.OutPoint2String(req.TxHash, 0),
		AccountId:    builder.AccountId,
		Args:         didCellArgs,
		LockCodeHash: req.Tx.Outputs[didEntity.Target.Index].Lock.CodeHash.Hex(),
	}

	var recordsInfos []tables.TableRecordsInfo
	recordList := didEntity.DidCellWitnessDataV0.Records
	for _, v := range recordList {
		recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
			AccountId: builder.AccountId,
			Account:   builder.Account,
			Key:       v.Key,
			Type:      v.Type,
			Label:     v.Label,
			Value:     v.Value,
			Ttl:       strconv.FormatUint(uint64(v.TTL), 10),
		})
	}
	if err = b.DbDao.AccountUpgrade(accountInfo, didCellInfo, recordsInfos); err != nil {
		log.Error("AccountUpgrade err:", err.Error(), req.TxHash, req.BlockNumber)
		resp.Err = fmt.Errorf("AccountCrossChain err: %s ", err.Error())
		return
	}
	return
}
