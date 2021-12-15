package core

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"time"
)

type SoScript struct {
	Name     common.SoScriptType
	OutPoint types.OutPoint
}

func (d *DasCore) InitDasSoScript() error {
	DasSoScriptMap.Store(common.SoScriptTypeEth, &SoScript{Name: common.SoScriptTypeEth})
	DasSoScriptMap.Store(common.SoScriptTypeTron, &SoScript{Name: common.SoScriptTypeTron})
	DasSoScriptMap.Store(common.SoScriptTypeCkb, &SoScript{Name: common.SoScriptTypeCkb})
	DasSoScriptMap.Store(common.SoScriptTypeCkbMulti, &SoScript{Name: common.SoScriptTypeCkbMulti})
	return d.asyncDasSoScript()
}

func (d *DasCore) RunAsyncDasSoScript(t time.Duration) {
	contractTicker := time.NewTicker(t) // update SO
	d.wg.Add(1)
	go func() {
		for {
			select {
			case <-contractTicker.C:
				log.Info("asyncDasSoScript begin ...")
				if err := d.asyncDasSoScript(); err != nil {
					log.Error("asyncDasConfigCell err:", err.Error())
				}
				log.Info("asyncDasSoScript end ...")
			case <-d.ctx.Done():
				d.wg.Done()
				return
			}
		}
	}()
}

func (d *DasCore) asyncDasSoScript() error {
	builder, err := d.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsMain)
	if err != nil {
		return fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	dasLockOutPoint := builder.ConfigCellMain.DasLockOutPointTable()
	DasSoScriptMap.Range(func(key, value interface{}) bool {
		if itemSo, okSo := value.(*SoScript); okSo {
			txHash := types.Hash{}
			switch key {
			case common.SoScriptTypeCkb:
				txHash = types.HexToHash(common.Bytes2Hex(dasLockOutPoint.CkbSignall().TxHash().RawData()))
			case common.SoScriptTypeCkbMulti:
				txHash = types.HexToHash(common.Bytes2Hex(dasLockOutPoint.CkbMultisign().TxHash().RawData()))
			case common.SoScriptTypeEth:
				txHash = types.HexToHash(common.Bytes2Hex(dasLockOutPoint.Eth().TxHash().RawData()))
			case common.SoScriptTypeTron:
				txHash = types.HexToHash(common.Bytes2Hex(dasLockOutPoint.Tron().TxHash().RawData()))
			}
			itemSo.OutPoint = types.OutPoint{
				TxHash: txHash,
				Index:  0,
			}
		}
		return true
	})
	return nil
}

func GetDasSoScript(soScriptName common.SoScriptType) (*SoScript, error) {
	if value, ok := DasSoScriptMap.Load(soScriptName); ok {
		if item, okSo := value.(*SoScript); okSo {
			return item, nil
		}
	}
	return nil, fmt.Errorf("not exist so script: [%s]", soScriptName)
}

func (d *SoScript) ToCellDep() *types.CellDep {
	return &types.CellDep{
		OutPoint: &d.OutPoint,
		DepType:  types.DepTypeCode,
	}
}
