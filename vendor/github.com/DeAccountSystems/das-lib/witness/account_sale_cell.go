package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
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
			version, err := molecule.Bytes2GoU32(dataEntity.Version().RawData())
			if err != nil {
				return false, fmt.Errorf("get version err: %s", err.Error())
			}
			resp.Version = version
			index, err := molecule.Bytes2GoU32(dataEntity.Index().RawData())
			if err != nil {
				return false, fmt.Errorf("get index err")
			}
			resp.Index = index

			switch version {
			case common.GoDataEntityVersion1:
				accountSaleData, err := molecule.AccountSaleCellDataV1FromSlice(dataEntity.Entity().RawData(), false)
				if err != nil {
					return false, fmt.Errorf("AccountSaleCellDataV1FromSlice err: %s", err.Error())
				}
				resp.AccountSaleCellDataV1 = accountSaleData
				resp.Account = string(accountSaleData.Account().RawData())
				resp.Description = string(accountSaleData.Description().RawData())
				resp.Price, _ = molecule.Bytes2GoU64(accountSaleData.Price().RawData())
				resp.StartedAt, _ = molecule.Bytes2GoU64(accountSaleData.StartedAt().RawData())
				resp.BuyerInviterProfitRate = 100
			case common.GoDataEntityVersion2:
				accountSaleData, err := molecule.AccountSaleCellDataFromSlice(dataEntity.Entity().RawData(), false)
				if err != nil {
					return false, fmt.Errorf("AccountSaleCellDataFromSlice err: %s", err.Error())
				}
				resp.AccountSaleCellData = accountSaleData
				resp.Account = string(accountSaleData.Account().RawData())
				resp.Description = string(accountSaleData.Description().RawData())
				resp.Price, _ = molecule.Bytes2GoU64(accountSaleData.Price().RawData())
				resp.StartedAt, _ = molecule.Bytes2GoU64(accountSaleData.StartedAt().RawData())
				resp.BuyerInviterProfitRate, _ = molecule.Bytes2GoU32(accountSaleData.BuyerInviterProfitRate().RawData())
			default:
				return false, fmt.Errorf("account sale version: %d", version)
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if resp.AccountSaleCellData == nil && resp.AccountSaleCellDataV1 == nil {
		return nil, fmt.Errorf("not exist account sale cell")
	}
	return &resp, nil
}

type AccountSaleCellDataBuilder struct {
	Index                  uint32
	Version                uint32
	AccountSaleCellDataV1  *molecule.AccountSaleCellDataV1
	AccountSaleCellData    *molecule.AccountSaleCellData
	DataEntityOpt          *molecule.DataEntityOpt
	Account                string
	Description            string
	Price                  uint64
	StartedAt              uint64
	BuyerInviterProfitRate uint32
}

type AccountSaleCellParam struct {
	OldIndex               uint32
	NewIndex               uint32
	Account                string
	Description            string
	Price                  uint64
	StartAt                uint64
	BuyerInviterProfitRate uint32
	Action                 string
}

func (a *AccountSaleCellDataBuilder) GenWitness(p *AccountSaleCellParam) ([]byte, []byte, error) {
	switch p.Action {
	case common.DasActionEditAccountSale:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountSaleCellDataBuilder()

		newAccountSaleCellData := newBuilder.Price(molecule.GoU64ToMoleculeU64(p.Price)).
			Description(molecule.GoString2MoleculeBytes(p.Description)).
			BuyerInviterProfitRate(molecule.GoU32ToMoleculeU32(p.BuyerInviterProfitRate)).Build()
		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountSaleCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionStartAccountSale:
		accountId, err := molecule.AccountIdFromSlice(common.GetAccountIdByAccount(p.Account), false)
		if err != nil {
			return nil, nil, fmt.Errorf("AccountIdFromSlice err: %s", err.Error())
		}
		newAccountSaleCellData := molecule.NewAccountSaleCellDataBuilder().
			Account(molecule.GoString2MoleculeBytes(p.Account)).
			AccountId(*accountId).
			Description(molecule.GoString2MoleculeBytes(p.Description)).
			StartedAt(molecule.GoU64ToMoleculeU64(p.StartAt)).
			Price(molecule.GoU64ToMoleculeU64(p.Price)).
			BuyerInviterProfitRate(molecule.GoU32ToMoleculeU32(p.BuyerInviterProfitRate)).
			Build()

		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountSaleCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionCancelAccountSale, common.DasActionBuyAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)

		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountSaleCell, &tmp)
		return witness, nil, nil
	}
	return nil, nil, fmt.Errorf("not exist action [%s]", p.Action)
}

func (a *AccountSaleCellDataBuilder) getOldDataEntityOpt(p *AccountSaleCellParam) *molecule.DataEntityOpt {
	var oldDataEntity molecule.DataEntity
	switch a.Version {
	case common.GoDataEntityVersion1:
		oldAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountSaleCellDataV1.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountSaleCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	case common.GoDataEntityVersion2:
		oldAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountSaleCellData.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	}
	oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
	return &oldDataEntityOpt
}

func (a *AccountSaleCellDataBuilder) getNewAccountSaleCellDataBuilder() *molecule.AccountSaleCellDataBuilder {
	var newBuilder molecule.AccountSaleCellDataBuilder
	switch a.Version {
	case common.GoDataEntityVersion1:
		temNewBuilder := molecule.NewAccountSaleCellDataBuilder()
		temNewBuilder.Account(*a.AccountSaleCellDataV1.Account()).
			AccountId(*a.AccountSaleCellDataV1.AccountId()).
			Price(*a.AccountSaleCellDataV1.Price()).
			Description(*a.AccountSaleCellDataV1.Description()).
			StartedAt(*a.AccountSaleCellDataV1.StartedAt()).Build()
		newBuilder = *temNewBuilder
	case common.GoDataEntityVersion2:
		newBuilder = a.AccountSaleCellData.AsBuilder()
	}
	return &newBuilder
}
