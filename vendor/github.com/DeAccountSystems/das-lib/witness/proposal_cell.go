package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type ProposalCellDataBuilder struct {
	Index            uint32
	ProposalCellData *molecule.ProposalCellData
	DataEntityOpt    *molecule.DataEntityOpt
}

func ProposalCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*ProposalCellDataBuilder, error) {
	var resp ProposalCellDataBuilder

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeProposalCell:
			dataEntityOpt, dataEntity, err := getDataEntityOpt(dataBys, dataType)
			if err != nil {
				return false, fmt.Errorf("getDataEntityOpt err: %s", err.Error())
			}
			resp.DataEntityOpt = dataEntityOpt

			index, err := molecule.Bytes2GoU32(dataEntity.Index().RawData())
			if err != nil {
				return false, fmt.Errorf("get index err: %s", err.Error())
			}
			resp.Index = index

			proposalCellData, err := molecule.ProposalCellDataFromSlice(dataEntity.Entity().RawData(), false)
			if err != nil {
				return false, fmt.Errorf("ProposalCellDataFromSlice err: %s", err.Error())
			}
			resp.ProposalCellData = proposalCellData
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if resp.ProposalCellData == nil {
		return nil, ErrNotExistWitness
	}
	return &resp, nil
}

type ProposalCellParam struct {
	ProposerLock *molecule.Script
	AccountList  [][]string
	CreateAt     uint64
	Action       common.ActionDataType
	OldIndex     uint32
	NewIndex     uint32
}

func (p *ProposalCellDataBuilder) getOldDataEntityOpt(param *ProposalCellParam) *molecule.DataEntityOpt {
	var oldDataEntity molecule.DataEntity
	oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(p.ProposalCellData.AsSlice())
	oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
		Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(param.OldIndex)).Build()
	oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
	return &oldDataEntityOpt
}
func (a *ProposalCellDataBuilder) GenWitness(p *ProposalCellParam) ([]byte, []byte, error) {
	switch p.Action {
	case common.DasActionPropose:
		proposalSlice := molecule.NewSliceListBuilder()
		for _, l0 := range p.AccountList {
			innerSlice := molecule.NewSLBuilder()
			for j := 0; j < len(l0)-1; j++ {
				// 0: exist, 1: proposed, 2: new
				itemType := molecule.GoU8ToMoleculeU8(uint8(2))
				if j == 0 {
					itemType = molecule.GoU8ToMoleculeU8(uint8(0))
				}
				accountId, err := molecule.AccountIdFromSlice(common.Hex2Bytes(l0[j]), false)
				if err != nil {
					return nil, nil, fmt.Errorf("accountId AccountIdFromSlice err: %s, id: %s", err.Error(), l0[j])
				}
				nextAccountId, err := molecule.AccountIdFromSlice(common.Hex2Bytes(l0[j+1]), false)
				if err != nil {
					return nil, nil, fmt.Errorf("nextAccountId AccountIdFromSlice err: %s, id: %s", err.Error(), l0[j+1])
				}
				item := molecule.NewProposalItemBuilder().
					AccountId(*accountId).
					ItemType(itemType).
					Next(*nextAccountId).
					Build()
				innerSlice.Push(item)
			}
			proposalSlice.Push(innerSlice.Build())
		}
		createAt := molecule.GoU64ToMoleculeU64(p.CreateAt)
		newProposalCellData := molecule.NewProposalCellDataBuilder().
			ProposerLock(*p.ProposerLock).
			CreatedAtHeight(createAt).
			Slices(proposalSlice.Build()).
			Build()
		newProposalCellDataBytes := molecule.GoBytes2MoleculeBytes(newProposalCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newProposalCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(0)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeProposalCell, &tmp)
		return witness, common.Blake2b(newProposalCellData.AsSlice()), nil
	case common.DasActionConfirmProposal, common.DasActionRecycleProposal:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeProposalCell, &tmp)
		return witness, nil, nil
	}
	return nil, nil, fmt.Errorf("not exist action [%s]", p.Action)
}
