package core

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"time"
)

// config cell info
type DasConfigCellInfo struct {
	Name        string
	OutPoint    types.OutPoint
	BlockNumber uint64
}

func (d *DasCore) InitDasConfigCell() error {
	DasConfigCellMap.Store(common.ConfigCellTypeArgsAccount, &DasConfigCellInfo{Name: "ConfigCellAccount"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsApply, &DasConfigCellInfo{Name: "ConfigCellApply"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsIncome, &DasConfigCellInfo{Name: "ConfigCellIncome"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsMain, &DasConfigCellInfo{Name: "ConfigCellMain"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPrice, &DasConfigCellInfo{Name: "ConfigCellPrice"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsProposal, &DasConfigCellInfo{Name: "ConfigCellProposal"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsProfitRate, &DasConfigCellInfo{Name: "ConfigCellProfitRate"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsRecordNamespace, &DasConfigCellInfo{Name: "ConfigCellRecordNamespace"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsRelease, &DasConfigCellInfo{Name: "ConfigCellRelease"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsUnavailable, &DasConfigCellInfo{Name: "ConfigCellUnavailable"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsSecondaryMarket, &DasConfigCellInfo{Name: "ConfigCellSecondaryMarket"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsReverseRecord, &DasConfigCellInfo{Name: "ConfigCellReverseRecord"})

	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount00, &DasConfigCellInfo{Name: "PreservedAccount00"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount01, &DasConfigCellInfo{Name: "PreservedAccount01"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount02, &DasConfigCellInfo{Name: "PreservedAccount02"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount03, &DasConfigCellInfo{Name: "PreservedAccount03"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount04, &DasConfigCellInfo{Name: "PreservedAccount04"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount05, &DasConfigCellInfo{Name: "PreservedAccount05"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount06, &DasConfigCellInfo{Name: "PreservedAccount06"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount07, &DasConfigCellInfo{Name: "PreservedAccount07"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount08, &DasConfigCellInfo{Name: "PreservedAccount08"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount09, &DasConfigCellInfo{Name: "PreservedAccount09"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount10, &DasConfigCellInfo{Name: "PreservedAccount10"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount11, &DasConfigCellInfo{Name: "PreservedAccount11"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount12, &DasConfigCellInfo{Name: "PreservedAccount12"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount13, &DasConfigCellInfo{Name: "PreservedAccount13"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount14, &DasConfigCellInfo{Name: "PreservedAccount14"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount15, &DasConfigCellInfo{Name: "PreservedAccount15"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount16, &DasConfigCellInfo{Name: "PreservedAccount16"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount17, &DasConfigCellInfo{Name: "PreservedAccount17"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount18, &DasConfigCellInfo{Name: "PreservedAccount18"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsPreservedAccount19, &DasConfigCellInfo{Name: "PreservedAccount19"})

	DasConfigCellMap.Store(common.ConfigCellTypeArgsCharSetEmoji, &DasConfigCellInfo{Name: "CharSetEmoji"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsCharSetDigit, &DasConfigCellInfo{Name: "CharSetDigit"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsCharSetEn, &DasConfigCellInfo{Name: "CharSetEn"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsCharSetHanS, &DasConfigCellInfo{Name: "CharSetHanS"})
	DasConfigCellMap.Store(common.ConfigCellTypeArgsCharSetHanT, &DasConfigCellInfo{Name: "CharSetHanT"})

	return d.AsyncDasConfigCell()
}

func (d *DasCore) RunAsyncDasConfigCell(t time.Duration) {
	contractTicker := time.NewTicker(t) // update config cell
	d.wg.Add(1)
	go func() {
		for {
			select {
			case <-contractTicker.C:
				log.Info("asyncDasConfigCell begin ...")
				if err := d.AsyncDasConfigCell(); err != nil {
					log.Error("asyncDasConfigCell err:", err.Error())
				}
				log.Info("asyncDasConfigCell end ...")
			case <-d.ctx.Done():
				d.wg.Done()
				return
			}
		}
	}()
}

func (d *DasCore) AsyncDasConfigCell() error {
	configCellContract, err := GetDasContractInfo(common.DasContractNameConfigCellType)
	if err != nil {
		return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	// search
	searchKey := &indexer.SearchKey{
		Script:     configCellContract.OutPut.Lock,
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			Script: configCellContract.ToScript(nil),
		},
	}
	res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 200, "")
	if err != nil {
		return fmt.Errorf("GetCells err: %s", err.Error())
	}
	//fmt.Println(len(res.Objects))
	// list
	for _, v := range res.Objects {
		configCellArgs := common.Bytes2Hex(v.Output.Type.Args)
		if value, ok := DasConfigCellMap.Load(configCellArgs); ok {
			if item, ok1 := value.(*DasConfigCellInfo); ok1 {
				if v.BlockNumber > item.BlockNumber && item.OutPoint.TxHash != v.OutPoint.TxHash {
					if _, ok2 := DasConfigCellByTxHashMap.Load(item.OutPoint.TxHash.Hex()); ok2 {
						DasConfigCellByTxHashMap.Delete(item.OutPoint.TxHash.Hex())
					}
					item.BlockNumber = v.BlockNumber
					item.OutPoint.TxHash = v.OutPoint.TxHash
					item.OutPoint.Index = v.OutPoint.Index

					DasConfigCellByTxHashMap.Store(item.OutPoint.TxHash.Hex(), true)
				}
			}
		}
	}
	return nil
}

func (d *DasConfigCellInfo) ToCellDep() *types.CellDep {
	return &types.CellDep{
		OutPoint: &d.OutPoint,
		DepType:  types.DepTypeCode,
	}
}

func GetDasConfigCellInfo(configCellTypeArgs common.ConfigCellTypeArgs) (*DasConfigCellInfo, error) {
	if value, ok := DasConfigCellMap.Load(configCellTypeArgs); ok {
		if item, okC := value.(*DasConfigCellInfo); okC {
			return item, nil
		}
	}
	return nil, fmt.Errorf("not exits ConfigCellInfo: [%s]", configCellTypeArgs)
}

func (d *DasCore) ConfigCellDataBuilderByTypeArgs(configCellTypeArgs common.ConfigCellTypeArgs) (*witness.ConfigCellDataBuilder, error) {
	configCell, err := GetDasConfigCellInfo(configCellTypeArgs)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	if res, err := d.client.GetTransaction(d.ctx, configCell.OutPoint.TxHash); err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	} else {
		return witness.ConfigCellDataBuilderByTypeArgs(res.Transaction, configCellTypeArgs)
	}
}

func (d *DasCore) ConfigCellDataBuilderByTypeArgsList(list ...common.ConfigCellTypeArgs) (*witness.ConfigCellDataBuilder, error) {
	var builder witness.ConfigCellDataBuilder
	for _, v := range list {
		configCell, err := GetDasConfigCellInfo(v)
		if err != nil {
			return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
		}
		if res, err := d.client.GetTransaction(d.ctx, configCell.OutPoint.TxHash); err != nil {
			return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
		} else if err = witness.ConfigCellDataBuilderRefByTypeArgs(&builder, res.Transaction, v); err != nil {
			return nil, fmt.Errorf("ConfigCellDataBuilderRefByTypeArgs err: %s", err.Error())
		}
	}
	return &builder, nil
}
