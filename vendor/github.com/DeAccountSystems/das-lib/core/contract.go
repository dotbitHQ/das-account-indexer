package core

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"time"
)

type DasContractInfo struct {
	ContractName   common.DasContractName
	OutPoint       *types.OutPoint   // contract outpoint
	OutPut         *types.CellOutput // contract script
	ContractTypeId types.Hash        // contract type id
}

func (d *DasCore) InitDasContract(mapDasContractTypeArgs map[common.DasContractName]string) {
	outPutLock := common.GetNormalLockScript(d.dasContractArgs)
	if d.net == common.DasNetTypeMainNet {
		outPutLock = common.GetNormalLockScriptByMultiSig(d.dasContractArgs)
	}
	for k, v := range mapDasContractTypeArgs {
		if v == "" {
			continue
		}
		dasContract := DasContractInfo{
			ContractName: k,
			OutPoint: &types.OutPoint{
				TxHash: types.Hash{},
				Index:  0,
			},
			OutPut: &types.CellOutput{
				Capacity: 0,
				Lock:     outPutLock,
				Type: &types.Script{
					CodeHash: types.HexToHash(d.dasContractCodeHash),
					HashType: types.HashTypeType,
					Args:     common.Hex2Bytes(v),
				},
			},
			ContractTypeId: types.Hash{},
		}
		// contract type id
		dasContract.ContractTypeId = common.ScriptToTypeId(dasContract.OutPut.Type)
		DasContractMap.Store(k, &dasContract)
		DasContractByTypeIdMap[dasContract.ContractTypeId.Hex()] = k
	}
	d.asyncDasContract()
}

// update contract info
func (d *DasCore) RunAsyncDasContract(t time.Duration) {
	contractTicker := time.NewTicker(t)
	d.wg.Add(1)
	go func() {
		for {
			select {
			case <-contractTicker.C:
				log.Info("asyncDasContracts begin ...")
				d.asyncDasContract()
				log.Info("asyncDasContracts end ...")
			case <-d.ctx.Done():
				d.wg.Done()
				return
			}
		}
	}()
}

func (d *DasCore) asyncDasContract() {
	DasContractMap.Range(func(key, value interface{}) bool {
		item, ok := value.(*DasContractInfo)
		if !ok {
			return true
		}
		searchKey := &indexer.SearchKey{
			Script:     item.OutPut.Lock,
			ScriptType: indexer.ScriptTypeLock,
			Filter: &indexer.CellsFilter{
				Script: item.OutPut.Type,
			},
		}
		now := time.Now()
		res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 1, "")
		if err != nil {
			log.Error("GetCells err:", key, err.Error())
			return true
		}
		if len(res.Objects) > 0 {
			item.OutPoint.Index = res.Objects[0].OutPoint.Index
			item.OutPoint.TxHash = res.Objects[0].OutPoint.TxHash
			log.Info("contract:", key, item.ContractTypeId, item.OutPoint.TxHash, item.OutPoint.Index, time.Since(now).Seconds())
		}
		return true
	})
}

func GetDasContractInfo(contractName common.DasContractName) (*DasContractInfo, error) {
	if value, ok := DasContractMap.Load(contractName); ok {
		if item, ok := value.(*DasContractInfo); ok {
			return item, nil
		}
	}
	return nil, fmt.Errorf("not exist contract name: [%s]", contractName)
}

func (d *DasContractInfo) ToCellDep() *types.CellDep {
	return &types.CellDep{
		OutPoint: d.OutPoint,
		DepType:  types.DepTypeCode,
	}
}

func (d *DasContractInfo) ToScript(args []byte) *types.Script {
	return &types.Script{
		CodeHash: d.ContractTypeId,
		HashType: types.HashTypeType,
		Args:     args,
	}
}

func (d *DasContractInfo) IsSameTypeId(codeHash types.Hash) bool {
	return d.ContractTypeId.Hex() == codeHash.Hex()
}
