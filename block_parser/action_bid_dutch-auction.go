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

	// did cell
	var didCellList []tables.TableDidCellInfo
	if refundLock := builder.GetRefundLock(); refundLock != nil {
		_, req.TxDidCellMap, err = b.DasCore.TxToDidCellEntityAndAction(req.Tx)
		if err != nil {
			resp.Err = fmt.Errorf("TxToDidCellEntityAndAction err: %s", err.Error())
			return
		}
		txDidEntityWitness, err := witness.GetDidEntityFromTx(req.Tx)
		if err != nil {
			resp.Err = fmt.Errorf("witness.GetDidEntityFromTx err: %s", err.Error())
			return
		}

		for k, v := range req.TxDidCellMap.Outputs {
			_, cellDataNew, err := v.GetDataInfo()
			if err != nil {
				resp.Err = fmt.Errorf("GetDataInfo new err: %s[%s]", err.Error(), k)
				return
			}
			acc := cellDataNew.Account
			accId := common.Bytes2Hex(common.GetAccountIdByAccount(acc))
			tmp := tables.TableDidCellInfo{
				BlockNumber:  req.BlockNumber,
				Outpoint:     common.OutPointStruct2String(v.OutPoint),
				AccountId:    accId,
				Account:      acc,
				Args:         common.Bytes2Hex(v.Lock.Args),
				LockCodeHash: v.Lock.CodeHash.Hex(),
				ExpiredAt:    cellDataNew.ExpireAt,
			}
			didCellList = append(didCellList, tmp)
			if w, yes := txDidEntityWitness.Outputs[v.Index]; yes {
				for _, r := range w.DidCellWitnessDataV0.Records {
					recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
						AccountId: accId,
						Account:   acc,
						Key:       r.Key,
						Type:      r.Type,
						Label:     r.Label,
						Value:     r.Value,
						Ttl:       strconv.FormatUint(uint64(r.TTL), 10),
					})
				}
			}
		}
	} else {
		for _, v := range builder.Records {
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
	}

	if err := b.DbDao.BidExpiredAccountAuction(accountInfo, recordsInfos, didCellList); err != nil {
		log.Error("ActionBidExpiredAccountAuction err:", err.Error(), toolib.JsonString(accountInfo))
		resp.Err = fmt.Errorf("ActionBidExpiredAccountAuction err: %s", err.Error())
	}
	return
}
