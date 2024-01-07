package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/scorpiotzh/toolib"
	"strconv"
)

func (b *BlockParser) ActionBidExpiredAccountAuction(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version transfer account tx")
		return
	}
	log.Info("BidExpiredAccountAuction:", req.BlockNumber, req.TxHash)

	builder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}
	account := builder.Account
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
	oHex, mHex, err := b.DasCore.Daf().ArgsToHex(req.Tx.Outputs[builder.Index].Lock.Args)
	if err != nil {
		resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
		return
	}
	accountInfo := tables.TableAccountInfo{
		BlockNumber:        req.BlockNumber,
		Outpoint:           common.OutPoint2String(req.TxHash, uint(builder.Index)),
		AccountId:          accountId,
		Account:            account,
		OwnerChainType:     oHex.ChainType,
		Owner:              oHex.AddressHex,
		OwnerAlgorithmId:   oHex.DasAlgorithmId,
		OwnerSubAid:        oHex.DasSubAlgorithmId,
		ManagerChainType:   mHex.ChainType,
		Manager:            mHex.AddressHex,
		ManagerAlgorithmId: mHex.DasAlgorithmId,
		ManagerSubAid:      mHex.DasSubAlgorithmId,
		ExpiredAt:          builder.ExpiredAt,
		RegisteredAt:       builder.RegisteredAt,
		Status:             tables.AccountStatus(builder.Status),
	}
	log.Info("ActionBidExpiredAccountAuction:", accountInfo)

	var recordsInfos []tables.TableRecordsInfo

	recordList := builder.Records
	for _, v := range recordList {
		recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
			AccountId: accountId,
			Account:   account,
			Key:       v.Key,
			Type:      v.Type,
			Label:     v.Label,
			Value:     v.Value,
			Ttl:       strconv.FormatUint(uint64(v.TTL), 10),
		})
	}

	if err := b.DbDao.BidExpiredAccountAuction(accountInfo, recordsInfos); err != nil {
		log.Error("ActionBidExpiredAccountAuction err:", err.Error(), toolib.JsonString(accountInfo))
		resp.Err = fmt.Errorf("ActionBidExpiredAccountAuction err: %s", err.Error())
	}
	return
}
