package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"strconv"
)

func (b *BlockParser) ActionEnableSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
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
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	// check sub-account config custom-script-args or not
	contractSub, err := core.GetDasContractInfo(common.DASContractNameSubAccountCellType)
	if err != nil {
		resp.Err = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}
	contractAcc, err := core.GetDasContractInfo(common.DasContractNameAccountCellType)
	if err != nil {
		resp.Err = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}
	var parentAccountId, accountCellOutpoint string
	for i, v := range req.Tx.Outputs {
		if v.Type != nil && contractSub.IsSameTypeId(v.Type.CodeHash) {
			parentAccountId = common.Bytes2Hex(v.Type.Args)
		}
		if v.Type != nil && contractAcc.IsSameTypeId(v.Type.CodeHash) {
			accountCellOutpoint = common.OutPoint2String(req.TxHash, uint(i))
		}
	}

	var parentAccountInfo tables.TableAccountInfo
	if accountCellOutpoint != "" {
		parentAccountInfo = tables.TableAccountInfo{
			BlockNumber:    req.BlockNumber,
			BlockTimestamp: req.BlockTimestamp,
			Outpoint:       accountCellOutpoint,
			AccountId:      parentAccountId,
		}
	}

	builderMap, err := witness.SubAccountBuilderMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	var accountInfos []tables.TableAccountInfo
	var subAccountIds []string
	for _, v := range builderMap {
		ownerHex, managerHex, err := b.DasCore.Daf().ArgsToHex(v.SubAccount.Lock.Args)
		if err != nil {
			resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
			return
		}

		accountInfos = append(accountInfos, tables.TableAccountInfo{
			BlockNumber:          req.BlockNumber,
			BlockTimestamp:       req.BlockTimestamp,
			Outpoint:             common.OutPoint2String(req.TxHash, 0),
			AccountId:            v.SubAccount.AccountId,
			ParentAccountId:      parentAccountId,
			Account:              v.Account,
			OwnerChainType:       ownerHex.ChainType,
			Owner:                ownerHex.AddressHex,
			OwnerAlgorithmId:     ownerHex.DasAlgorithmId,
			ManagerChainType:     managerHex.ChainType,
			Manager:              managerHex.AddressHex,
			ManagerAlgorithmId:   managerHex.DasAlgorithmId,
			Status:               tables.AccountStatus(v.SubAccount.Status),
			EnableSubAccount:     tables.EnableSubAccount(v.SubAccount.EnableSubAccount),
			RenewSubAccountPrice: v.SubAccount.RenewSubAccountPrice,
			Nonce:                v.SubAccount.Nonce,
			RegisteredAt:         v.SubAccount.RegisteredAt,
			ExpiredAt:            v.SubAccount.ExpiredAt,
		})
		subAccountIds = append(subAccountIds, v.SubAccount.AccountId)
	}

	if err = b.DbDao.CreateSubAccount(subAccountIds, accountInfos, parentAccountInfo); err != nil {
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
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	builderMap, err := witness.SubAccountBuilderMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	for _, builder := range builderMap {
		accountInfo := tables.TableAccountInfo{
			BlockNumber:    req.BlockNumber,
			BlockTimestamp: req.BlockTimestamp,
			Outpoint:       common.OutPoint2String(req.TxHash, 0),
			AccountId:      builder.SubAccount.AccountId,
			Nonce:          builder.CurrentSubAccount.Nonce,
		}

		subAccount, err := builder.ConvertToEditValue()
		if err != nil {
			resp.Err = fmt.Errorf("ConvertToEditValue err: %s", err.Error())
			return
		}
		switch string(builder.EditKey) {
		case common.EditKeyOwner:
			ownerHex, managerHex, err := b.DasCore.Daf().ArgsToHex(common.Hex2Bytes(subAccount.LockArgs))
			if err != nil {
				resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
				return
			}
			accountInfo.OwnerAlgorithmId = ownerHex.DasAlgorithmId
			accountInfo.OwnerChainType = ownerHex.ChainType
			accountInfo.Owner = ownerHex.AddressHex
			accountInfo.ManagerAlgorithmId = managerHex.DasAlgorithmId
			accountInfo.ManagerChainType = managerHex.ChainType
			accountInfo.Manager = managerHex.AddressHex
			if err = b.DbDao.EditOwnerSubAccount(accountInfo); err != nil {
				resp.Err = fmt.Errorf("EditOwnerSubAccount err: %s", err.Error())
				return
			}
		case common.EditKeyManager:
			_, managerHex, err := b.DasCore.Daf().ArgsToHex(common.Hex2Bytes(subAccount.LockArgs))
			if err != nil {
				resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
				return
			}
			accountInfo.ManagerAlgorithmId = managerHex.DasAlgorithmId
			accountInfo.ManagerChainType = managerHex.ChainType
			accountInfo.Manager = managerHex.AddressHex
			if err = b.DbDao.EditManagerSubAccount(accountInfo); err != nil {
				resp.Err = fmt.Errorf("EditManagerSubAccount err: %s", err.Error())
				return
			}
		case common.EditKeyRecords:
			var recordsInfos []tables.TableRecordsInfo
			for _, v := range subAccount.Records {
				recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
					AccountId:       builder.SubAccount.AccountId,
					ParentAccountId: common.Bytes2Hex(req.Tx.Outputs[0].Type.Args),
					Account:         builder.Account,
					Key:             v.Key,
					Type:            v.Type,
					Label:           v.Label,
					Value:           v.Value,
					Ttl:             strconv.FormatUint(uint64(v.TTL), 10),
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
	log.Info("das tx:", req.Action, req.TxHash)
	return
}

/*func (b *BlockParser) ActionRenewSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	builderMap, err := witness.SubAccountBuilderMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	var accountInfos []tables.TableAccountInfo
	for _, builder := range builderMap {
		accountInfo := tables.TableAccountInfo{
			BlockNumber:    req.BlockNumber,
			BlockTimestamp: req.BlockTimestamp,
			Outpoint:       common.OutPoint2String(req.TxHash, 0),
			AccountId:      builder.SubAccount.AccountId,
			Nonce:          builder.CurrentSubAccount.Nonce,
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
}*/

func (b *BlockParser) ActionRecycleSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	log.Info("das tx:", req.Action, req.TxHash)
	return
}

/*func (b *BlockParser) ActionRecycleSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

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
}*/

func (b *BlockParser) ActionUpdateSubAccountInfo(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	log.Info("das tx:", req.Action, req.TxHash)
	return
}

func (b *BlockParser) ActionConfigSubAccountCreatingScript(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("ActionConfigSubAccountCreatingScript:", req.BlockNumber, req.TxHash)

	// update account cell outpoint
	builder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("witness.AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}
	outpoint := common.OutPoint2String(req.TxHash, uint(builder.Index))
	if err := b.DbDao.UpdateAccountOutpoint(builder.AccountId, outpoint); err != nil {
		resp.Err = fmt.Errorf("UpdateAccountOutpoint err: %s", err.Error())
		return
	}

	return
}
