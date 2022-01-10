package txbuilder

import (
	"encoding/binary"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"sort"
)

type SignData struct {
	SignType common.DasAlgorithmId `json:"sign_type"` // 签名类型
	SignMsg  string                `json:"sign_msg"`  // 待签名信息
}

func (d *DasTxBuilder) GenerateDigestListFromTx(skipGroups []int) ([]SignData, error) {
	skipGroups = append(skipGroups, d.ServerSignGroup...)
	groups, err := d.getGroupsFromTx()
	if err != nil {
		return nil, fmt.Errorf("getGroupsFromTx err: %s", err.Error())
	}
	var digestList []SignData
	for _, group := range groups {
		if digest, err := d.generateDigestByGroup(group, skipGroups); err != nil {
			return nil, fmt.Errorf("generateDigestByGroup err: %s", err.Error())
		} else {
			digestList = append(digestList, digest)
		}
	}
	return digestList, nil
}

func (d *DasTxBuilder) getGroupsFromTx() ([][]int, error) {
	var tmpMapForGroup = make(map[string][]int)
	var sortList []string
	for i, v := range d.Transaction.Inputs {
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return nil, fmt.Errorf("getInputCell err: %s", err.Error())
		}

		cellHash, err := item.Cell.Output.Lock.Hash()
		if err != nil {
			return nil, fmt.Errorf("inputs lock to hash err: %s", err.Error())
		}
		indexList, okTmp := tmpMapForGroup[cellHash.String()]
		if !okTmp {
			sortList = append(sortList, cellHash.String())
		}
		indexList = append(indexList, i)
		tmpMapForGroup[cellHash.String()] = indexList
	}
	sort.Strings(sortList)
	var list [][]int
	for _, v := range sortList {
		item, _ := tmpMapForGroup[v]
		list = append(list, item)
	}
	return list, nil
}

func (d *DasTxBuilder) generateDigestByGroup(group []int, skipGroups []int) (SignData, error) {
	var signData = SignData{}
	if group == nil || len(group) < 1 {
		return signData, fmt.Errorf("invalid param")
	}

	digest := ""
	data, err := transaction.EmptyWitnessArg.Serialize()
	if err != nil {
		return signData, err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, err := d.Transaction.ComputeHash()
	if err != nil {
		return signData, err
	}

	message := append(hash.Bytes(), length...)
	message = append(message, data...)
	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = d.Transaction.Witnesses[group[i]]
			length := make([]byte, 8)
			binary.LittleEndian.PutUint64(length, uint64(len(data)))
			message = append(message, length...)
			message = append(message, data...)
		}
	}
	// hash witnesses which do not in any input group
	for _, witness := range d.Transaction.Witnesses[len(d.Transaction.Inputs):] {
		length := make([]byte, 8)
		binary.LittleEndian.PutUint64(length, uint64(len(witness)))
		message = append(message, length...)
		message = append(message, witness...)
	}

	message, err = blake2b.Blake256(message)
	if err != nil {
		return signData, err
	}
	digest = common.Bytes2Hex(message)
	item, err := d.getInputCell(d.Transaction.Inputs[group[0]].PreviousOutput)
	if err != nil {
		return signData, fmt.Errorf("getInputCell err: %s", err.Error())
	}

	ownerAlgorithmId, managerAlgorithmId, _, _, _, _ := core.FormatDasLockToHexAddress(item.Cell.Output.Lock.Args)
	signData.SignMsg = digest
	signData.SignType = ownerAlgorithmId
	if actionBuilder, err := witness.ActionDataBuilderFromTx(d.Transaction); err == nil && actionBuilder.Action == common.DasActionEditRecords {
		signData.SignType = managerAlgorithmId
	}

	// skip useless signature
	if len(skipGroups) != 0 {
		skip := false
		for i := range group {
			for j := range skipGroups {
				if group[i] == skipGroups[j] {
					skip = true
					break
				}
			}
			if skip {
				break
			}
		}
		if skip {
			signData.SignMsg = ""
		}
	}
	return signData, nil
}

func (d *DasTxBuilder) getInputCell(o *types.OutPoint) (*types.CellWithStatus, error) {
	if o == nil {
		return nil, fmt.Errorf("OutPoint is nil")
	}
	key := fmt.Sprintf("%s-%d", o.TxHash.Hex(), o.Index)
	if item, ok := d.MapInputsCell[key]; ok {
		if item.Cell != nil && item.Cell.Output != nil && item.Cell.Output.Lock != nil {
			return item, nil
		}
	}
	if cell, err := d.dasCore.Client().GetLiveCell(d.ctx, o, true); err != nil {
		return nil, fmt.Errorf("GetLiveCell err: %s", err.Error())
	} else if cell.Cell.Output.Lock != nil {
		d.MapInputsCell[key] = cell
		return cell, nil
	} else {
		log.Warn("GetLiveCell:", key, cell.Status)
		return nil, fmt.Errorf("cell [%s] not live", key)
	}
}
