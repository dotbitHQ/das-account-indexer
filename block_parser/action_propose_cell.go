package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"strconv"
)

func (b *BlockParser) ActionConfirmProposal(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersionTx err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	mapPreBuilder, err := witness.PreAccountCellDataBuilderMapFromTx(req.Tx, common.DataTypeOld)
	if err != nil {
		resp.Err = fmt.Errorf("PreAccountCellDataBuilderMapFromTx err: %s", err.Error())
		return
	}

	mapBuilder, err := witness.AccountCellDataBuilderMapFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderMapFromTx err: %s", err.Error())
		return
	}
	var accounts []tables.TableAccountInfo
	var records []tables.TableRecordsInfo
	var accountIdList []string
	for _, builder := range mapBuilder {
		ownerHex, managerHex, err := b.DasCore.Daf().ArgsToHex(req.Tx.Outputs[builder.Index].Lock.Args)
		if err != nil {
			resp.Err = fmt.Errorf("ArgsToHex err: %s", err.Error())
			return
		}
		accounts = append(accounts, tables.TableAccountInfo{
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
		})
		if _, ok := mapPreBuilder[builder.Account]; ok {
			accountIdList = append(accountIdList, builder.AccountId)
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
		}
	}

	// did cell
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
	var didCellList []tables.TableDidCellInfo
	var didCellRecords []tables.TableRecordsInfo
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
				didCellRecords = append(didCellRecords, tables.TableRecordsInfo{
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
	if len(didCellList) > 0 {
		records = didCellRecords
	}

	if err = b.DbDao.UpdateAccountInfoList(accounts, records, accountIdList, didCellList); err != nil {
		resp.Err = fmt.Errorf("UpdateAccountInfo err: %s", err.Error())
		return
	}

	return
}
