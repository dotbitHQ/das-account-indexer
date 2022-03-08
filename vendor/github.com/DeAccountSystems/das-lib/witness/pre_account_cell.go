package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type PreAccountCellDataBuilder struct {
	Index              uint32
	Account            string
	PreAccountCellData *molecule.PreAccountCellData
	DataEntityOpt      *molecule.DataEntityOpt
}

func PreAccountCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*PreAccountCellDataBuilder, error) {
	respMap, err := PreAccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}
	for k, _ := range respMap {
		return respMap[k], nil
	}
	return nil, fmt.Errorf("not exist pre account cell")
}
func PreAccountIdCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*PreAccountCellDataBuilder, error) {
	respMap, err := PreAccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*PreAccountCellDataBuilder)
	for k, v := range respMap {
		k1 := common.Bytes2Hex(common.GetAccountIdByAccount(k))
		retMap[k1] = v
	}
	return retMap, nil
}
func PreAccountCellDataBuilderMapFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*PreAccountCellDataBuilder, error) {
	var respMap = make(map[string]*PreAccountCellDataBuilder)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypePreAccountCell:
			var resp PreAccountCellDataBuilder
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

			preAccountCellData, err := molecule.PreAccountCellDataFromSlice(dataEntity.Entity().RawData(), false)
			if err != nil {
				return false, fmt.Errorf("PreAccountCellDataFromSlice err: %s", err.Error())
			}
			resp.PreAccountCellData = preAccountCellData
			resp.Account = common.AccountCharsToAccount(preAccountCellData.Account())
			respMap[resp.Account] = &resp
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(respMap) == 0 {
		return nil, fmt.Errorf("not exist pre account cell")
	}
	return respMap, nil
}

func (p *PreAccountCellDataBuilder) AccountName() (string, error) {
	if p.PreAccountCellData != nil {
		return common.AccountCharsToAccount(p.PreAccountCellData.Account()), nil
	}
	return "", fmt.Errorf("AccountName is nil")
}

func (p *PreAccountCellDataBuilder) InviterId() (string, error) {
	if p.PreAccountCellData != nil {
		return common.Bytes2Hex(p.PreAccountCellData.InviterId().RawData()), nil
	}
	return "", fmt.Errorf("PreAccountCellData is nil")
}

func (p *PreAccountCellDataBuilder) InviterLock() (*molecule.Script, error) {
	if p.PreAccountCellData != nil {
		if len(p.PreAccountCellData.InviterLock().AsSlice()) == 0 {
			return nil, nil
		}
		return molecule.ScriptFromSlice(p.PreAccountCellData.InviterLock().AsSlice(), false)
	}
	return nil, fmt.Errorf("PreAccountCellData is nil")
}

func (p *PreAccountCellDataBuilder) ChannelLock() (*molecule.Script, error) {
	if p.PreAccountCellData != nil {
		if len(p.PreAccountCellData.ChannelLock().AsSlice()) == 0 {
			return nil, nil
		}
		return molecule.ScriptFromSlice(p.PreAccountCellData.ChannelLock().AsSlice(), false)
	}
	return nil, fmt.Errorf("PreAccountCellData is nil")
}

func (p *PreAccountCellDataBuilder) RefundLock() (*molecule.Script, error) {
	if p.PreAccountCellData != nil {
		return p.PreAccountCellData.RefundLock(), nil
	}
	return nil, fmt.Errorf("PreAccountCellData is nil")
}

func (p *PreAccountCellDataBuilder) OwnerLockArgsStr() (string, error) {
	if p.PreAccountCellData != nil {
		return common.Bytes2Hex(p.PreAccountCellData.OwnerLockArgs().RawData()), nil
	}
	return "", fmt.Errorf("PreAccountCellData is nil")
}

type PreAccountCellParam struct {
	OldIndex uint32
	NewIndex uint32
	Status   uint8
	Action   string

	CreatedAt       int64
	InvitedDiscount uint32
	Quote           uint64
	InviterScript   *types.Script
	ChannelScript   *types.Script
	InviterId       []byte
	OwnerLockArgs   []byte
	RefundLock      *types.Script
	Price           molecule.PriceConfig
	AccountChars    molecule.AccountChars
}

func (p *PreAccountCellDataBuilder) getOldDataEntityOpt(param *PreAccountCellParam) *molecule.DataEntityOpt {
	var oldDataEntity molecule.DataEntity
	oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(p.PreAccountCellData.AsSlice())
	oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
		Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(param.OldIndex)).Build()
	oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
	return &oldDataEntityOpt
}
func (p *PreAccountCellDataBuilder) GenWitness(param *PreAccountCellParam) ([]byte, []byte, error) {

	switch param.Action {
	case common.DasActionPreRegister:
		createdAt := molecule.NewUint64Builder().Set(molecule.GoTimeUnixToMoleculeBytes(param.CreatedAt)).Build()
		invitedDiscount := molecule.GoU32ToMoleculeU32(param.InvitedDiscount)
		quote := molecule.GoU64ToMoleculeU64(param.Quote)
		var iScript, cScript molecule.ScriptOpt
		if param.InviterScript != nil {
			iScript = molecule.NewScriptOptBuilder().Set(molecule.CkbScript2MoleculeScript(param.InviterScript)).Build()
		} else {
			iScript = *molecule.ScriptOptFromSliceUnchecked([]byte{})
		}
		if param.ChannelScript != nil {
			cScript = molecule.NewScriptOptBuilder().Set(molecule.CkbScript2MoleculeScript(param.ChannelScript)).Build()
		} else {
			cScript = *molecule.ScriptOptFromSliceUnchecked([]byte{})
		}
		inviterId := molecule.GoBytes2MoleculeBytes(param.InviterId)
		ownerLockArgs := molecule.GoBytes2MoleculeBytes(param.OwnerLockArgs)
		refundLock := molecule.CkbScript2MoleculeScript(param.RefundLock)

		preAccountCellData := molecule.NewPreAccountCellDataBuilder().
			Account(param.AccountChars).
			RefundLock(refundLock).
			OwnerLockArgs(ownerLockArgs).
			InviterId(inviterId).
			InviterLock(iScript).
			ChannelLock(cScript).
			Price(param.Price).
			Quote(quote).
			InvitedDiscount(invitedDiscount).
			CreatedAt(createdAt).Build()
		newDataBytes := molecule.GoBytes2MoleculeBytes(preAccountCellData.AsSlice())
		newDataEntity := molecule.NewDataEntityBuilder().Entity(newDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(param.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypePreAccountCell, &tmp)
		return witness, common.Blake2b(preAccountCellData.AsSlice()), nil
	case common.DasActionPropose:
		oldDataEntityOpt := p.getOldDataEntityOpt(param)
		tmp := molecule.NewDataBuilder().Dep(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypePreAccountCell, &tmp)
		return witness, nil, nil
	case common.DasActionConfirmProposal:
		oldDataEntityOpt := p.getOldDataEntityOpt(param)
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypePreAccountCell, &tmp)
		return witness, nil, nil
	}
	return nil, nil, fmt.Errorf("not exist action [%s]", param.Action)
}
