package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"gorm.io/gorm"
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

func (b *BlockParser) ActionUpdateSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version edit sub account tx")
		return
	}
	log.Info("ActionUpdateSubAccount:", req.BlockNumber, req.TxHash)

	var subAccountNewBuilder witness.SubAccountNewBuilder
	builderMap, err := subAccountNewBuilder.SubAccountNewMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}
	var createBuilderMap = make(map[string]*witness.SubAccountNew)
	var renewBuilderMap = make(map[string]*witness.SubAccountNew)
	var editBuilderMap = make(map[string]*witness.SubAccountNew)
	var recycleBuilderMap = make(map[string]*witness.SubAccountNew)
	var approvalBuilderMap = make(map[string]*witness.SubAccountNew)
	for k, v := range builderMap {
		switch v.Action {
		case common.SubActionCreate:
			createBuilderMap[k] = v
		case common.SubActionRenew:
			renewBuilderMap[k] = v
		case common.SubActionEdit:
			editBuilderMap[k] = v
		case common.SubActionRecycle:
			recycleBuilderMap[k] = v
		case common.SubActionCreateApproval, common.SubActionDelayApproval,
			common.SubActionRevokeApproval, common.SubActionFullfillApproval:
			approvalBuilderMap[k] = v
		default:
			resp.Err = fmt.Errorf("unknow sub-action [%s]", v.Action)
			return
		}
	}
	if err := b.actionUpdateSubAccountForRecycle(recycleBuilderMap); err != nil {
		resp.Err = fmt.Errorf("recycle sub-account err: %s", err.Error())
		return
	}

	if err := b.actionUpdateSubAccountForCreate(req, createBuilderMap); err != nil {
		resp.Err = fmt.Errorf("create sub-account err: %s", err.Error())
		return
	}

	if err := b.actionUpdateSubAccountForRenew(req, renewBuilderMap); err != nil {
		resp.Err = fmt.Errorf("edit err: %s", err.Error())
		return
	}

	if err := b.actionUpdateSubAccountForEdit(req, editBuilderMap); err != nil {
		resp.Err = fmt.Errorf("edit sub-account err: %s", err.Error())
		return
	}

	if err := b.actionUpdateSubAccountForApproval(req, approvalBuilderMap); err != nil {
		resp.Err = fmt.Errorf("edit sub-account err: %s", err.Error())
		return
	}
	return
}

func (b *BlockParser) actionUpdateSubAccountForRecycle(recycleBuilderMap map[string]*witness.SubAccountNew) error {
	if len(recycleBuilderMap) == 0 {
		return nil
	}
	var subAccIds []string
	for _, builder := range recycleBuilderMap {
		subAccIds = append(subAccIds, builder.SubAccountData.AccountId)
	}
	if err := b.DbDao.DelSubAccounts(subAccIds); err != nil {
		return fmt.Errorf("DelSubAccounts err: %s", err.Error())
	}

	return nil
}

func (b *BlockParser) actionUpdateSubAccountForCreate(req *FuncTransactionHandleReq, createBuilderMap map[string]*witness.SubAccountNew) error {
	if len(createBuilderMap) == 0 {
		return nil
	}

	// check sub-account config custom-script-args or not
	contractSub, err := core.GetDasContractInfo(common.DASContractNameSubAccountCellType)
	if err != nil {
		return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	var parentAccountId string
	for _, v := range req.Tx.Outputs {
		if v.Type != nil && contractSub.IsSameTypeId(v.Type.CodeHash) {
			parentAccountId = common.Bytes2Hex(v.Type.Args)
		}
	}

	var parentAccountInfo tables.TableAccountInfo
	var accountInfos []tables.TableAccountInfo
	var subAccountIds []string
	var records []tables.TableRecordsInfo
	for _, v := range createBuilderMap {
		ownerHex, managerHex, err := b.DasCore.Daf().ArgsToHex(v.SubAccountData.Lock.Args)
		if err != nil {
			return fmt.Errorf("ArgsToHex err: %s", err.Error())
		}

		accountInfos = append(accountInfos, tables.TableAccountInfo{
			BlockNumber:          req.BlockNumber,
			BlockTimestamp:       req.BlockTimestamp,
			Outpoint:             common.OutPoint2String(req.TxHash, 0),
			AccountId:            v.SubAccountData.AccountId,
			ParentAccountId:      parentAccountId,
			Account:              v.Account,
			OwnerChainType:       ownerHex.ChainType,
			Owner:                ownerHex.AddressHex,
			OwnerAlgorithmId:     ownerHex.DasAlgorithmId,
			OwnerSubAid:          ownerHex.DasSubAlgorithmId,
			ManagerChainType:     managerHex.ChainType,
			Manager:              managerHex.AddressHex,
			ManagerAlgorithmId:   managerHex.DasAlgorithmId,
			ManagerSubAid:        managerHex.DasSubAlgorithmId,
			Status:               tables.AccountStatus(v.SubAccountData.Status),
			EnableSubAccount:     tables.EnableSubAccount(v.SubAccountData.EnableSubAccount),
			RenewSubAccountPrice: v.SubAccountData.RenewSubAccountPrice,
			Nonce:                v.SubAccountData.Nonce,
			RegisteredAt:         v.SubAccountData.RegisteredAt,
			ExpiredAt:            v.SubAccountData.ExpiredAt,
		})
		subAccountIds = append(subAccountIds, v.SubAccountData.AccountId)

		for _, record := range v.SubAccountData.Records {
			records = append(records, tables.TableRecordsInfo{
				AccountId:       v.SubAccountData.AccountId,
				ParentAccountId: parentAccountId,
				Account:         v.Account,
				Key:             record.Key,
				Type:            record.Type,
				Label:           record.Label,
				Value:           record.Value,
				Ttl:             strconv.FormatUint(uint64(record.TTL), 10),
			})
		}
	}
	if err = b.DbDao.CreateSubAccount(subAccountIds, accountInfos, parentAccountInfo, records); err != nil {
		return fmt.Errorf("CreateSubAccount err: %s", err.Error())
	}

	return nil
}

func (b *BlockParser) actionUpdateSubAccountForRenew(req *FuncTransactionHandleReq, renewBuilderMap map[string]*witness.SubAccountNew) error {
	if len(renewBuilderMap) == 0 {
		return nil
	}

	var accountInfos []tables.TableAccountInfo

	for _, v := range renewBuilderMap {
		subAcc, err := b.DbDao.FindAccountInfoByAccountId(v.CurrentSubAccountData.AccountId)
		if err != nil {
			return err
		}
		if subAcc.Id == 0 {
			return fmt.Errorf("account: [%s] no exist", v.Account)
		}

		accountInfos = append(accountInfos, tables.TableAccountInfo{
			Id:          subAcc.Id,
			BlockNumber: req.BlockNumber,
			Outpoint:    common.OutPoint2String(req.TxHash, 0),
			Nonce:       v.CurrentSubAccountData.Nonce,
			ExpiredAt:   v.CurrentSubAccountData.ExpiredAt,
		})
	}

	if err := b.DbDao.Transaction(func(tx *gorm.DB) error {
		for i := range accountInfos {
			accountInfo := accountInfos[i]
			if err := tx.Where("id=?", accountInfo.Id).Updates(&accountInfo).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("UpdateSubAccountForRenew err: %s", err.Error())
	}
	return nil
}

func (b *BlockParser) actionUpdateSubAccountForEdit(req *FuncTransactionHandleReq, editBuilderMap map[string]*witness.SubAccountNew) error {
	if len(editBuilderMap) == 0 {
		return nil
	}
	for _, builder := range editBuilderMap {
		accountInfo := tables.TableAccountInfo{
			BlockNumber:    req.BlockNumber,
			BlockTimestamp: req.BlockTimestamp,
			Outpoint:       common.OutPoint2String(req.TxHash, 0),
			AccountId:      builder.SubAccountData.AccountId,
			Nonce:          builder.CurrentSubAccountData.Nonce,
		}

		switch builder.EditKey {
		case common.EditKeyOwner:
			ownerHex, managerHex, err := b.DasCore.Daf().ArgsToHex(builder.EditLockArgs)
			if err != nil {
				return fmt.Errorf("ArgsToHex err: %s", err.Error())
			}
			accountInfo.OwnerAlgorithmId = ownerHex.DasAlgorithmId
			accountInfo.OwnerSubAid = ownerHex.DasSubAlgorithmId
			accountInfo.OwnerChainType = ownerHex.ChainType
			accountInfo.Owner = ownerHex.AddressHex
			accountInfo.ManagerAlgorithmId = managerHex.DasAlgorithmId
			accountInfo.ManagerSubAid = managerHex.DasSubAlgorithmId
			accountInfo.ManagerChainType = managerHex.ChainType
			accountInfo.Manager = managerHex.AddressHex
			if err = b.DbDao.EditOwnerSubAccount(accountInfo); err != nil {
				return fmt.Errorf("EditOwnerSubAccount err: %s", err.Error())
			}
		case common.EditKeyManager:
			_, managerHex, err := b.DasCore.Daf().ArgsToHex(builder.EditLockArgs)
			if err != nil {
				return fmt.Errorf("ArgsToHex err: %s", err.Error())
			}
			accountInfo.ManagerAlgorithmId = managerHex.DasAlgorithmId
			accountInfo.ManagerSubAid = managerHex.DasSubAlgorithmId
			accountInfo.ManagerChainType = managerHex.ChainType
			accountInfo.Manager = managerHex.AddressHex
			if err = b.DbDao.EditManagerSubAccount(accountInfo); err != nil {
				return fmt.Errorf("EditManagerSubAccount err: %s", err.Error())
			}
		case common.EditKeyRecords:
			var recordsInfos []tables.TableRecordsInfo
			for _, v := range builder.EditRecords {
				recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
					AccountId:       builder.SubAccountData.AccountId,
					ParentAccountId: common.Bytes2Hex(req.Tx.Outputs[0].Type.Args),
					Account:         builder.Account,
					Key:             v.Key,
					Type:            v.Type,
					Label:           v.Label,
					Value:           v.Value,
					Ttl:             strconv.FormatUint(uint64(v.TTL), 10),
				})
			}
			if err := b.DbDao.EditRecordsSubAccount(accountInfo, recordsInfos); err != nil {
				return fmt.Errorf("EditRecordsSubAccount err: %s", err.Error())
			}
		}
	}
	return nil
}

func (b *BlockParser) actionUpdateSubAccountForApproval(req *FuncTransactionHandleReq, approvalBuilderMap map[string]*witness.SubAccountNew) error {
	if len(approvalBuilderMap) == 0 {
		return nil
	}

	var accounts []map[string]interface{}

	for _, builder := range approvalBuilderMap {
		accountInfo := map[string]interface{}{
			"block_number":    req.BlockNumber,
			"block_timestamp": req.BlockTimestamp,
			"outpoint":        common.OutPoint2String(req.TxHash, 0),
			"account_id":      builder.SubAccountData.AccountId,
			"nonce":           builder.CurrentSubAccountData.Nonce,
			"action":          builder.Action,
		}

		switch builder.Action {
		case common.SubActionCreateApproval:
			accountInfo["status"] = tables.AccountStatusApproval
		case common.SubActionRevokeApproval:
			accountInfo["status"] = tables.AccountStatusNormal
		case common.SubActionFullfillApproval:
			approval := builder.SubAccountData.AccountApproval
			switch approval.Action {
			case witness.AccountApprovalActionTransfer:
				owner, manager, err := b.DasCore.Daf().ScriptToHex(approval.Params.Transfer.ToLock)
				if err != nil {
					return err
				}
				accountInfo["status"] = tables.AccountStatusNormal
				accountInfo["owner"] = owner.AddressHex
				accountInfo["owner_chain_type"] = owner.ChainType
				accountInfo["owner_algorithm_id"] = owner.DasAlgorithmId
				accountInfo["manager"] = manager.AddressHex
				accountInfo["manager_chain_type"] = manager.ChainType
				accountInfo["manager_algorithm_id"] = manager.DasAlgorithmId
			default:
				return fmt.Errorf("unknow approval action [%s]", approval.Action)
			}
		}
		accounts = append(accounts, accountInfo)
	}
	return b.DbDao.UpdateAccounts(accounts)
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

	var subAccountNewBuilder witness.SubAccountNewBuilder
	builderMap, err := subAccountNewBuilder.SubAccountNewMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	var accountInfos []tables.TableAccountInfo
	var subAccountIds []string
	for _, v := range builderMap {
		ownerHex, managerHex, err := b.DasCore.Daf().ArgsToHex(v.SubAccountData.Lock.Args)
		if err != nil {
			resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
			return
		}

		accountInfos = append(accountInfos, tables.TableAccountInfo{
			BlockNumber:          req.BlockNumber,
			BlockTimestamp:       req.BlockTimestamp,
			Outpoint:             common.OutPoint2String(req.TxHash, 0),
			AccountId:            v.SubAccountData.AccountId,
			ParentAccountId:      parentAccountId,
			Account:              v.Account,
			OwnerChainType:       ownerHex.ChainType,
			Owner:                ownerHex.AddressHex,
			OwnerAlgorithmId:     ownerHex.DasAlgorithmId,
			OwnerSubAid:          ownerHex.DasSubAlgorithmId,
			ManagerChainType:     managerHex.ChainType,
			Manager:              managerHex.AddressHex,
			ManagerAlgorithmId:   managerHex.DasAlgorithmId,
			ManagerSubAid:        managerHex.DasSubAlgorithmId,
			Status:               tables.AccountStatus(v.SubAccountData.Status),
			EnableSubAccount:     tables.EnableSubAccount(v.SubAccountData.EnableSubAccount),
			RenewSubAccountPrice: v.SubAccountData.RenewSubAccountPrice,
			Nonce:                v.SubAccountData.Nonce,
			RegisteredAt:         v.SubAccountData.RegisteredAt,
			ExpiredAt:            v.SubAccountData.ExpiredAt,
		})
		subAccountIds = append(subAccountIds, v.SubAccountData.AccountId)
	}

	if err = b.DbDao.CreateSubAccount(subAccountIds, accountInfos, parentAccountInfo, nil); err != nil {
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

	var subAccountNewBuilder witness.SubAccountNewBuilder
	builderMap, err := subAccountNewBuilder.SubAccountNewMapFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("SubAccountBuilderMapFromTx err: %s", err.Error())
		return
	}

	if err := b.actionUpdateSubAccountForEdit(req, builderMap); err != nil {
		resp.Err = fmt.Errorf("edit err: %s", err.Error())
		return
	}

	return
}

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

func (b *BlockParser) ActionConfigSubAccount(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DASContractNameSubAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("ActionConfigSubAccount:", req.BlockNumber, req.TxHash)

	dataBuilder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = err
		return
	}
	resp.Err = b.DbDao.UpdateAccountOutpoint(dataBuilder.AccountId, common.OutPoint2String(req.TxHash, 0))
	return
}
