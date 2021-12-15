package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

var (
	DataEntityVersion1 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion1)
	DataEntityVersion2 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion2)
)

func AccountSaleCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*AccountSaleCellDataBuilder, error) {
	var resp AccountSaleCellDataBuilder
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeAccountSaleCell:
			dataEntityOpt, dataEntity, err := getDataEntityOpt(dataBys, dataType)
			if err != nil {
				return false, fmt.Errorf("getDataEntityOpt err: %s", err.Error())
			}
			resp.DataEntityOpt = dataEntityOpt
			index, err := molecule.Bytes2GoU32(dataEntity.Index().RawData())
			if err != nil {
				return false, fmt.Errorf("get index err")
			}
			resp.Index = index

			accountSaleData, err := molecule.AccountSaleCellDataFromSlice(dataEntity.Entity().RawData(), false)
			if err != nil {
				return false, fmt.Errorf("AccountSaleCellDataFromSlice err: %s", err.Error())
			}
			resp.AccountSaleCellData = accountSaleData
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if resp.AccountSaleCellData == nil {
		return nil, fmt.Errorf("not exist account sale cell")
	}
	return &resp, nil
}

type AccountSaleCellDataBuilder struct {
	Index               uint32
	AccountSaleCellData *molecule.AccountSaleCellData
	DataEntityOpt       *molecule.DataEntityOpt
}

func (a *AccountSaleCellDataBuilder) Account() string {
	return string(a.AccountSaleCellData.Account().RawData())
}

func (a *AccountSaleCellDataBuilder) Description() string {
	return string(a.AccountSaleCellData.Description().RawData())
}

func (a *AccountSaleCellDataBuilder) Price() (uint64, error) {
	return molecule.Bytes2GoU64(a.AccountSaleCellData.Price().RawData())
}

func (a *AccountSaleCellDataBuilder) StartedAt() (uint64, error) {
	return molecule.Bytes2GoU64(a.AccountSaleCellData.StartedAt().RawData())
}

type AccountSaleCellParam struct {
	Price       uint64
	Description string
	Account     string
	StartAt     uint64
	Action      string
}

func (a *AccountSaleCellDataBuilder) GenWitness(p *AccountSaleCellParam) ([]byte, []byte, error) {
	switch p.Action {
	case common.DasActionBuyAccount:
		oldAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountSaleCellData.AsSlice())
		oldDataEntity := molecule.NewDataEntityBuilder().Entity(oldAccountSaleCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(1)).Build()
		oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountSaleCell, &tmp)
		return witness, nil, nil
	case common.DasActionEditAccountSale:
		oldAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountSaleCellData.AsSlice())
		oldDataEntity := molecule.NewDataEntityBuilder().Entity(oldAccountSaleCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(0)).Build()
		oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()

		newBuilder := a.AccountSaleCellData.AsBuilder()
		newAccountSaleCellData := newBuilder.Price(molecule.GoU64ToMoleculeU64(p.Price)).
			Description(molecule.GoString2MoleculeBytes(p.Description)).Build()
		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(0)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().Old(oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountSaleCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionStartAccountSale:
		accountId, err := molecule.AccountIdFromSlice(common.GetAccountIdByAccount(p.Account), false)
		if err != nil {
			return nil, nil, fmt.Errorf("AccountIdFromSlice err: %s", err.Error())
		}
		startAt := molecule.GoU64ToMoleculeU64(p.StartAt)
		price := molecule.GoU64ToMoleculeU64(p.Price)

		newAccountSaleCellData := molecule.NewAccountSaleCellDataBuilder().
			Account(molecule.GoString2MoleculeBytes(p.Account)).
			AccountId(*accountId).
			Description(molecule.GoString2MoleculeBytes(p.Description)).
			StartedAt(startAt).
			Price(price).
			Build()

		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(1)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountSaleCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionCancelAccountSale:
		oldAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountSaleCellData.AsSlice())
		oldDataEntity := molecule.NewDataEntityBuilder().Entity(oldAccountSaleCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(1)).Build()
		oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountSaleCell, &tmp)
		return witness, nil, nil
	}
	return nil, nil, fmt.Errorf("not exist action [%s]", p.Action)
}
