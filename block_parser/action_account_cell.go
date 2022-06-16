package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/witness"
	"strconv"
	"time"
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
	list := builder.RecordList()
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
	res, err := b.DasCore.Client().GetTransaction(b.Ctx, req.Tx.Inputs[1].PreviousOutput.TxHash)
	if err != nil {
		resp.Err = fmt.Errorf("GetTransaction err: %s", err.Error())
		return
	}
	if isCV, err := isCurrentVersionTx(res.Transaction, common.DasContractNameAccountCellType); err != nil {
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
	builder, err := witness.AccountCellDataBuilderFromTx(res.Transaction, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}
	builderConfig, err := b.DasCore.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsAccount)
	if err != nil {
		resp.Err = fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
		return
	}
	gracePeriod, err := builderConfig.ExpirationGracePeriod()
	if err != nil {
		resp.Err = fmt.Errorf("ExpirationGracePeriod err: %s", err.Error())
		return
	}

	if builder.Status != 0 {
		resp.Err = fmt.Errorf("ActionRecycleExpiredAccount: account is not normal status")
		return
	}
	if builder.ExpiredAt+uint64(gracePeriod) > uint64(time.Now().Unix()) {
		resp.Err = fmt.Errorf("ActionRecycleExpiredAccount: account has not expired yet")
		return
	}
	oHex, _, err := b.DasCore.Daf().ArgsToHex(res.Transaction.Outputs[req.Tx.Inputs[1].PreviousOutput.Index].Lock.Args)
	if err != nil {
		resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
		return
	}

	var subAccountIds []string
	if builder.EnableSubAccount == 1 {
		accountInfos, err := b.DbDao.GetAccountInfoByParentAccountId(builder.AccountId)
		if err != nil {
			resp.Err = fmt.Errorf("GetAccountInfoByParentAccountId err: %s", err.Error())
			return
		}
		for _, accountInfo := range accountInfos {
			subAccountIds = append(subAccountIds, accountInfo.AccountId)
		}
	}

	log.Info("ActionRecycleExpiredAccount:", builder.Account, oHex.DasAlgorithmId, oHex.ChainType, oHex.AddressHex, len(subAccountIds))

	if err = b.DbDao.RecycleExpiredAccount(previousBuilder.AccountId, previousBuilder.NextAccountId, builder.AccountId, subAccountIds); err != nil {
		resp.Err = fmt.Errorf("RecycleExpiredAccount err: %s", err.Error())
		return
	}

	return
}
