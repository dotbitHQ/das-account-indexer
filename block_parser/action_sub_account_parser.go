package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"strconv"
	"time"
)

func (b *BlockParser) ActionEnableSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version enable sub account tx")
		return
	}

	log.Info("ActionEnableSubAccount:", req.BlockNumber, req.TxHash)

	builder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}

	accountInfo := tables.TableAccountInfo{
		BlockNumber:          req.BlockNumber,
		BlockTimestamp:       req.BlockTimestamp,
		Outpoint:             common.OutPoint2String(req.TxHash, 0),
		AccountId:            builder.AccountId,
		EnableSubAccount:     tables.AccountEnableStatusOn,
		RenewSubAccountPrice: builder.RenewSubAccountPrice,
	}

	if err = b.DbDao.EnableSubAccount(accountInfo); err != nil {
		resp.Err = fmt.Errorf("EnableSubAccount err: %s", err.Error())
		return
	}

	return
}

func (b *BlockParser) ActionCreateSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version create sub account tx")
		return
	}

	log.Info("ActionCreateSubAccount:", req.BlockNumber, req.TxHash)

	builder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}

	builderMap, err := witness.SubAccountBuilderMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	var accountInfos []tables.TableAccountInfo
	for _, v := range builderMap {
		oID, mID, oCT, mCT, oA, mA := core.FormatDasLockToHexAddress(v.SubAccount.Lock.Args)

		accountInfos = append(accountInfos, tables.TableAccountInfo{
			BlockNumber:          req.BlockNumber,
			BlockTimestamp:       req.BlockTimestamp,
			Outpoint:             common.OutPoint2String(req.TxHash, 0),
			AccountId:            v.SubAccount.AccountId,
			ParentAccountId:      builder.AccountId,
			Account:              v.Account,
			OwnerChainType:       oCT,
			Owner:                oA,
			OwnerAlgorithmId:     oID,
			ManagerChainType:     mCT,
			Manager:              mA,
			ManagerAlgorithmId:   mID,
			Status:               tables.AccountStatus(v.SubAccount.Status),
			EnableSubAccount:     tables.EnableSubAccount(v.SubAccount.EnableSubAccount),
			RenewSubAccountPrice: v.SubAccount.RenewSubAccountPrice,
			Nonce:                v.SubAccount.Nonce,
			RegisteredAt:         v.SubAccount.RegisteredAt,
			ExpiredAt:            v.SubAccount.ExpiredAt,
		})
	}
	accountInfo := tables.TableAccountInfo{
		BlockNumber:    req.BlockNumber,
		BlockTimestamp: req.BlockTimestamp,
		Outpoint:       common.OutPoint2String(req.TxHash, 0),
		AccountId:      builder.AccountId,
	}
	if err = b.DbDao.CreateSubAccount(accountInfos, accountInfo); err != nil {
		resp.Err = fmt.Errorf("CreateSubAccount err: %s", err.Error())
		return
	}

	return
}

func (b *BlockParser) ActionEditSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version edit sub account tx")
		return
	}

	log.Info("ActionEditSubAccount:", req.BlockNumber, req.TxHash)

	builderMap, err := witness.SubAccountBuilderMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	for _, builder := range builderMap {
		builder.SubAccount.Nonce++
		accountInfo := tables.TableAccountInfo{
			BlockNumber:    req.BlockNumber,
			BlockTimestamp: req.BlockTimestamp,
			Outpoint:       common.OutPoint2String(req.TxHash, 0),
			AccountId:      builder.SubAccount.AccountId,
			Nonce:          builder.SubAccount.Nonce,
		}

		subAccount, err := builder.ConvertToEditValue()
		if err != nil {
			resp.Err = fmt.Errorf("ConvertToEditValue err: %s", err.Error())
			return
		}
		switch string(builder.EditKey) {
		case common.EditKeyOwner:
			oID, _, oCT, _, oA, _ := core.FormatDasLockToHexAddress(common.Hex2Bytes(subAccount.LockArgs))
			accountInfo.OwnerAlgorithmId = oID
			accountInfo.OwnerChainType = oCT
			accountInfo.Owner = oA
			if err = b.DbDao.EditOwnerSubAccount(accountInfo); err != nil {
				resp.Err = fmt.Errorf("EditOwnerSubAccount err: %s", err.Error())
				return
			}
		case common.EditKeyManager:
			_, mID, _, mCT, _, mA := core.FormatDasLockToHexAddress(common.Hex2Bytes(subAccount.LockArgs))
			accountInfo.ManagerAlgorithmId = mID
			accountInfo.ManagerChainType = mCT
			accountInfo.Manager = mA
			if err = b.DbDao.EditManagerSubAccount(accountInfo); err != nil {
				resp.Err = fmt.Errorf("EditManagerSubAccount err: %s", err.Error())
				return
			}
		case common.EditKeyRecords:
			var recordsInfos []tables.TableRecordsInfo
			for _, v := range subAccount.Records {
				recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
					AccountId: builder.SubAccount.AccountId,
					Account:   builder.Account,
					Key:       v.Key,
					Type:      v.Type,
					Label:     v.Label,
					Value:     v.Value,
					Ttl:       strconv.FormatUint(uint64(v.TTL), 10),
				})
			}
			if err = b.DbDao.EditRecordsSubAccount(accountInfo, recordsInfos); err != nil {
				resp.Err = fmt.Errorf("EditRecordsSubAccount err: %s", err.Error())
				return
			}
		}
	}

	return
}

func (b *BlockParser) ActionRenewSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version renew sub account tx")
		return
	}

	log.Info("ActionRenewSubAccount:", req.BlockNumber, req.TxHash)

	builderMap, err := witness.SubAccountBuilderMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	var accountInfos []tables.TableAccountInfo
	for _, builder := range builderMap {
		builder.SubAccount.Nonce++
		accountInfo := tables.TableAccountInfo{
			BlockNumber:    req.BlockNumber,
			BlockTimestamp: req.BlockTimestamp,
			Outpoint:       common.OutPoint2String(req.TxHash, 0),
			AccountId:      builder.SubAccount.AccountId,
			Nonce:          builder.SubAccount.Nonce,
		}

		subAccount, err := builder.ConvertToEditValue()
		if err != nil {
			resp.Err = fmt.Errorf("ConvertToEditValue err: %s", err.Error())
			return
		}
		switch string(builder.EditKey) {
		case common.EditKeyExpiredAt:
			accountInfo.ExpiredAt = subAccount.ExpiredAt
			accountInfos = append(accountInfos, accountInfo)
		}
	}
	if err = b.DbDao.RenewSubAccount(accountInfos); err != nil {
		resp.Err = fmt.Errorf("RenewSubAccount err: %s", err.Error())
		return
	}

	return
}

func (b *BlockParser) ActionRecycleSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version recycle sub account tx")
		return
	}

	log.Info("ActionRecycleSubAccount:", req.BlockNumber, req.TxHash)

	builderMap, err := witness.SubAccountBuilderMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	var accountIds []string
	for _, builder := range builderMap {
		// if expired time greater than three months ago, then reject the recycle of sub_account.
		if builder.SubAccount.ExpiredAt > uint64(time.Now().Add(-time.Hour*24*90).Unix()) {
			resp.Err = fmt.Errorf("not yet arrived expired time: %d", builder.SubAccount.ExpiredAt)
			return
		}

		accountIds = append(accountIds, builder.SubAccount.AccountId)
	}
	if err = b.DbDao.RecycleSubAccount(accountIds); err != nil {
		resp.Err = fmt.Errorf("RecycleSubAccount err: %s", err.Error())
		return
	}

	return
}
