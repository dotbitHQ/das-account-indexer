package txbuilder

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
)

func (d *DasTxBuilder) newTx() error {
	systemScriptCell, err := utils.NewSystemScripts(d.dasCore.Client())
	if err != nil {
		return err
	}
	baseTx := transaction.NewSecp256k1SingleSigTx(systemScriptCell)
	d.Transaction = baseTx
	return nil
}

func (d *DasTxBuilder) equalArgs(src, dst string) bool {
	if common.Has0xPrefix(src) {
		src = src[2:]
	}
	if common.Has0xPrefix(dst) {
		dst = dst[2:]
	}
	return src == dst
}
func (d *DasTxBuilder) addInputsForTx(inputs []*types.CellInput) error {
	if len(inputs) == 0 {
		return fmt.Errorf("inputs is nil")
	}
	startIndex := len(d.Transaction.Inputs)
	_, _, err := transaction.AddInputsForTransaction(d.Transaction, inputs)
	if err != nil {
		return fmt.Errorf("AddInputsForTransaction err: %s", err.Error())
	}

	var cellDepList []*types.CellDep
	for i, v := range inputs {
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return fmt.Errorf("getInputCell err: %s", err.Error())
		}

		if item.Cell.Output.Type != nil {
			if contractName, ok := core.DasContractByTypeIdMap[item.Cell.Output.Type.CodeHash.Hex()]; ok {
				dasContract, err := core.GetDasContractInfo(contractName)
				if err != nil {
					return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
				}
				cellDepList = append(cellDepList, dasContract.ToCellDep())
			}
		}
		if item.Cell.Output.Lock != nil &&
			item.Cell.Output.Lock.CodeHash.Hex() == transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH &&
			d.equalArgs(common.Bytes2Hex(item.Cell.Output.Lock.Args), d.serverArgs) {
			d.ServerSignGroup = append(d.ServerSignGroup, startIndex+i)
		}
		if item.Cell.Output.Lock != nil {
			if contractName, ok := core.DasContractByTypeIdMap[item.Cell.Output.Lock.CodeHash.Hex()]; ok {
				if dasContract, err := core.GetDasContractInfo(contractName); err != nil {
					return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
				} else {
					cellDepList = append(cellDepList, &types.CellDep{OutPoint: dasContract.OutPoint, DepType: types.DepTypeCode})
					if contractName == common.DasContractNameDispatchCellType {
						daf := core.DasAddressFormat{DasNetType: d.dasCore.NetType()}
						ownerHex, managerHex, _ := daf.ArgsToHex(item.Cell.Output.Lock.Args)
						oID, mID := ownerHex.DasAlgorithmId, managerHex.DasAlgorithmId
						oSo, _ := core.GetDasSoScript(oID.ToSoScriptType())
						mSo, _ := core.GetDasSoScript(mID.ToSoScriptType())
						cellDepList = append(cellDepList, oSo.ToCellDep())
						cellDepList = append(cellDepList, mSo.ToCellDep())
					}
				}
			}
		}
	}
	d.addCellDepListIntoMapCellDep(cellDepList)
	return nil
}

func (d *DasTxBuilder) addOutputsForTx(outputs []*types.CellOutput, outputsData [][]byte) error {
	lo := len(outputs)
	lod := len(outputsData)
	if lo == 0 || lod == 0 || lo != lod {
		return fmt.Errorf("outputs[%d] or outputDatas[%d]", lo, lod)
	}
	var cellDepList []*types.CellDep
	for i := 0; i < lo; i++ {
		output := outputs[i]
		outputData := outputsData[i]
		d.Transaction.Outputs = append(d.Transaction.Outputs, output)
		d.Transaction.OutputsData = append(d.Transaction.OutputsData, outputData)

		if output.Type == nil {
			continue
		}
		contractName, ok := core.DasContractByTypeIdMap[output.Type.CodeHash.Hex()]
		if !ok {
			continue
		}
		dasContract, err := core.GetDasContractInfo(contractName)
		if err != nil {
			return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		}
		cellDepList = append(cellDepList, dasContract.ToCellDep())
	}

	d.addCellDepListIntoMapCellDep(cellDepList)
	return nil
}

func (d *DasTxBuilder) checkTxWitnesses() error {
	if len(d.Transaction.Witnesses) == 0 {
		return fmt.Errorf("witness is nil")
	}
	lenI := len(d.Transaction.Inputs)
	lenW := len(d.Transaction.Witnesses)
	if lenW < lenI {
		return fmt.Errorf("len witness[%d]<len inputs[%d]", lenW, lenI)
	} else if lenW > lenI {
		_, err := witness.ActionDataBuilderFromWitness(d.Transaction.Witnesses[lenI])
		//_, err := witness.ActionDataBuilderFromTx(d.Transaction)
		if err != nil {
			return fmt.Errorf("ActionDataBuilderFromTx err: %s", err.Error())
		}
	}
	return nil
}

func (d *DasTxBuilder) addCellDepListIntoMapCellDep(cellDepList []*types.CellDep) {
	for i, v := range cellDepList {
		k := fmt.Sprintf("%s-%d", v.OutPoint.TxHash.Hex(), v.OutPoint.Index)
		d.mapCellDep[k] = cellDepList[i]
	}
}

func (d *DasTxBuilder) addMapCellDepWitnessForBaseTx(cellDepList []*types.CellDep) error {
	configCellMain, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsMain)
	if err != nil {
		return fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	cellDepList = append(cellDepList, configCellMain.ToCellDep())

	tmpMap := make(map[string]bool)
	var tmpCellDeps []*types.CellDep
	for _, v := range cellDepList {
		k := fmt.Sprintf("%s-%d", v.OutPoint.TxHash.Hex(), v.OutPoint.Index)
		if _, ok := tmpMap[k]; ok {
			continue
		}
		tmpMap[k] = true
		tmpCellDeps = append(tmpCellDeps, &types.CellDep{
			OutPoint: v.OutPoint,
			DepType:  v.DepType,
		})
		if _, ok := core.DasConfigCellByTxHashMap.Load(v.OutPoint.TxHash.Hex()); !ok {
			continue
		}
		if res, err := d.dasCore.Client().GetTransaction(d.ctx, v.OutPoint.TxHash); err != nil {
			return fmt.Errorf("GetTransaction err: %s [%s]", err.Error(), k)
		} else {
			d.Transaction.Witnesses = append(d.Transaction.Witnesses, res.Transaction.Witnesses[len(res.Transaction.Witnesses)-1])
		}
	}
	if len(tmpCellDeps) > 0 {
		d.Transaction.CellDeps = append(tmpCellDeps, d.Transaction.CellDeps...)
	}

	for k, v := range d.mapCellDep {
		if _, ok := tmpMap[k]; ok {
			continue
		}
		d.Transaction.CellDeps = append(d.Transaction.CellDeps, &types.CellDep{
			OutPoint: v.OutPoint,
			DepType:  v.DepType,
		})
		if _, ok := core.DasConfigCellByTxHashMap.Load(v.OutPoint.TxHash.Hex()); !ok {
			continue
		}
		if res, err := d.dasCore.Client().GetTransaction(d.ctx, v.OutPoint.TxHash); err != nil {
			return fmt.Errorf("GetTransaction err: %s [%s]", err.Error(), k)
		} else {
			d.Transaction.Witnesses = append(d.Transaction.Witnesses, res.Transaction.Witnesses[len(res.Transaction.Witnesses)-1])
		}
	}
	return nil
}

func (d *DasTxBuilder) SendTransactionWithCheck(needCheck bool) (*types.Hash, error) {
	if needCheck {
		err := d.checkTxBeforeSend()
		if err != nil {
			return nil, fmt.Errorf("checkTxBeforeSend err: %s", err.Error())
		}
	}

	err := d.serverSignTx()
	if err != nil {
		return nil, fmt.Errorf("remoteSignTx err: %s", err.Error())
	}

	log.Info("before sent: ", d.TxString())
	txHash, err := d.dasCore.Client().SendTransactionNoneValidation(d.ctx, d.Transaction)
	if err != nil {
		return nil, fmt.Errorf("SendTransaction err: %v", err)
	}
	log.Info("SendTransaction success:", txHash.Hex())
	return txHash, nil
}

func (d *DasTxBuilder) SendTransaction() (*types.Hash, error) {
	return d.SendTransactionWithCheck(true)
}

func (d *DasTxBuilder) checkTxBeforeSend() error {
	// check total num of inputs and outputs
	if len(d.Transaction.Inputs)+len(d.Transaction.Outputs) > 9000 {
		return fmt.Errorf("checkTxBeforeSend, failed len of inputs: %d, ouputs: %d", len(d.Transaction.Inputs), len(d.Transaction.Outputs))
	}
	// check tx fee < 1 CKB
	totalCapacityFromInputs, err := d.getCapacityFromInputs()
	if err != nil {
		return err
	}
	totalCapacityFromOutputs := d.Transaction.OutputsCapacity()
	txFee := totalCapacityFromInputs - totalCapacityFromOutputs
	if totalCapacityFromInputs <= totalCapacityFromOutputs || txFee >= common.OneCkb {
		return fmt.Errorf("checkTxBeforeSend failed, txFee: %d totalCapacityFromInputs: %d totalCapacityFromOutputs: %d", txFee, totalCapacityFromInputs, totalCapacityFromOutputs)
	}

	// check witness format
	err = d.checkTxWitnesses()
	if err != nil {
		return err
	}
	// check the occupied capacity
	for i, cell := range d.Transaction.Outputs {
		occupied := cell.OccupiedCapacity(d.Transaction.OutputsData[i])
		if cell.Capacity < occupied {
			return fmt.Errorf("checkTxBeforeSend occupied capacity failed, occupied: %d capacity: %d index: %d", occupied, cell.Capacity, i)
		}
	}
	log.Info("check success before sent")
	return nil
}

func (d *DasTxBuilder) getCapacityFromInputs() (uint64, error) {
	total := uint64(0)
	for _, v := range d.Transaction.Inputs {
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return 0, fmt.Errorf("getInputCell err: %s", err.Error())
		}
		total += item.Cell.Output.Capacity
	}
	return total, nil
}
