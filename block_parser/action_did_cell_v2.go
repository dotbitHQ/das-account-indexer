package block_parser

import (
	"bytes"
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"strconv"
)

func (b *BlockParser) DidCellActionUpdate(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	log.Info("DidCellActionUpdate:", req.BlockNumber, req.TxHash, req.Action)

	if len(req.TxDidCellMap.Inputs) != len(req.TxDidCellMap.Outputs) {
		resp.Err = fmt.Errorf("len(req.TxDidCellMap.Inputs)!=len(req.TxDidCellMap.Outputs)")
		return
	}
	txDidEntityWitness, err := witness.GetDidEntityFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("witness.GetDidEntityFromTx err: %s", err.Error())
		return
	}

	var oldOutpointList []string
	var list []tables.TableDidCellInfo
	var accountIds []string
	var records []tables.TableRecordsInfo

	for k, v := range req.TxDidCellMap.Inputs {
		_, cellDataOld, err := v.GetDataInfo()
		if err != nil {
			resp.Err = fmt.Errorf("GetDataInfo old err: %s[%s]", err.Error(), k)
			return
		}
		n, ok := req.TxDidCellMap.Outputs[k]
		if !ok {
			resp.Err = fmt.Errorf("TxDidCellMap diff err: %s[%s]", err.Error(), k)
			return
		}
		_, cellDataNew, err := n.GetDataInfo()
		if err != nil {
			resp.Err = fmt.Errorf("GetDataInfo new err: %s[%s]", err.Error(), k)
			return
		}
		account := cellDataOld.Account
		accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))

		oldOutpoint := common.OutPointStruct2String(v.OutPoint)
		oldOutpointList = append(oldOutpointList, oldOutpoint)

		tmp := tables.TableDidCellInfo{
			BlockNumber:  req.BlockNumber,
			Outpoint:     common.OutPointStruct2String(n.OutPoint),
			AccountId:    accountId,
			Account:      account,
			Args:         common.Bytes2Hex(n.Lock.Args),
			LockCodeHash: n.Lock.CodeHash.Hex(),
			ExpiredAt:    cellDataNew.ExpireAt,
		}
		list = append(list, tmp)
		if bytes.Compare(cellDataOld.WitnessHash, cellDataNew.WitnessHash) != 0 {
			accountIds = append(accountIds, accountId)
			if w, yes := txDidEntityWitness.Outputs[v.Index]; yes {
				for _, r := range w.DidCellWitnessDataV0.Records {
					records = append(records, tables.TableRecordsInfo{
						AccountId: accountId,
						Account:   account,
						Key:       r.Key,
						Type:      r.Type,
						Label:     r.Label,
						Value:     r.Value,
						Ttl:       strconv.FormatUint(uint64(r.TTL), 10),
					})
				}
			}
		}
	}

	if err := b.DbDao.DidCellUpdateList(oldOutpointList, list, accountIds, records); err != nil {
		resp.Err = fmt.Errorf("DidCellUpdateList err: %s", err.Error())
		return
	}

	return
}

func (b *BlockParser) DidCellActionRecycle(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	log.Info("DidCellActionRecycle:", req.BlockNumber, req.TxHash, req.Action)

	var oldOutpointList []string
	var accountIds []string
	for k, v := range req.TxDidCellMap.Inputs {
		oldOutpoint := common.OutPointStruct2String(v.OutPoint)
		oldOutpointList = append(oldOutpointList, oldOutpoint)

		_, cellData, err := v.GetDataInfo()
		if err != nil {
			resp.Err = fmt.Errorf("GetDataInfo err: %s[%s]", err.Error(), k)
			return
		}
		account := cellData.Account
		accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
		accountIds = append(accountIds, accountId)
	}

	if err := b.DbDao.DidCellRecycleList(oldOutpointList, accountIds); err != nil {
		resp.Err = fmt.Errorf("DidCellRecycleList err: %s", err.Error())
		return
	}
	return
}
