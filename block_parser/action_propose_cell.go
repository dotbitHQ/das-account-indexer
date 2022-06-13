package block_parser

import (
	"das-account-indexer/tables"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/witness"
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
			ManagerChainType:   managerHex.ChainType,
			Manager:            managerHex.AddressHex,
			ManagerAlgorithmId: managerHex.DasAlgorithmId,
			Status:             tables.AccountStatus(builder.Status),
			RegisteredAt:       builder.RegisteredAt,
			ExpiredAt:          builder.ExpiredAt,
		})
		if _, ok := mapPreBuilder[builder.Account]; ok {
			accountIdList = append(accountIdList, builder.AccountId)
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
		}
	}

	if err = b.DbDao.UpdateAccountInfoList(accounts, records, accountIdList); err != nil {
		resp.Err = fmt.Errorf("UpdateAccountInfo err: %s", err.Error())
		return
	}

	return
}
