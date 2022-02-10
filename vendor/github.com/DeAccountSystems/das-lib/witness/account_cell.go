package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type AccountCellDataBuilder struct {
	Index             uint32
	Version           uint32
	AccountId         string
	NextAccountId     string
	Account           string
	Status            uint8
	RegisteredAt      uint64
	ExpiredAt         uint64
	RecordsHashBys    []byte
	Records           *molecule.Records
	AccountCellDataV1 *molecule.AccountCellDataV1
	AccountCellData   *molecule.AccountCellData
	DataEntityOpt     *molecule.DataEntityOpt
}

type AccountCellParam struct {
	OldIndex              uint32
	NewIndex              uint32
	Status                uint8
	Action                string
	AccountId             string
	RegisterAt            uint64
	SubAction             string
	AccountChars          *molecule.AccountChars
	LastEditRecordsAt     int64
	LastEditManagerAt     int64
	LastTransferAccountAt int64
	Records               []AccountCellRecord
}

func AccountCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*AccountCellDataBuilder, error) {
	respMap, err := AccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}
	for k, _ := range respMap {
		return respMap[k], nil
	}
	return nil, fmt.Errorf("not exist account cell")
}

func AccountCellDataBuilderMapFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*AccountCellDataBuilder, error) {
	var respMap = make(map[string]*AccountCellDataBuilder)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeAccountCell:
			var resp AccountCellDataBuilder
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
				return false, fmt.Errorf("get index err: %s", err.Error())
			}
			resp.Index = index
			if dataType == common.DataTypeNew {
				expiredAt, err := common.GetAccountCellExpiredAtFromOutputData(tx.OutputsData[index])
				if err != nil {
					return false, fmt.Errorf("GetAccountCellExpiredAtFromOutputData err: %s", err.Error())
				}
				resp.ExpiredAt = expiredAt
				nextAccountId, err := common.GetAccountCellNextAccountIdFromOutputData(tx.OutputsData[index])
				if err != nil {
					return false, fmt.Errorf("GetAccountCellNextAccountIdFromOutputData err: %s", err.Error())
				}
				resp.NextAccountId = common.Bytes2Hex(nextAccountId)
			}

			switch version {
			case common.GoDataEntityVersion1:
				accountData, err := molecule.AccountCellDataV1FromSlice(dataEntity.Entity().RawData(), false)
				if err != nil {
					return false, fmt.Errorf("AccountCellDataV1FromSlice err: %s", err.Error())
				}
				resp.AccountCellDataV1 = accountData
				resp.Account = common.AccountCharsToAccount(accountData.Account())
				resp.AccountId = common.Bytes2Hex(accountData.Id().RawData())
				resp.Status, _ = molecule.Bytes2GoU8(accountData.Status().RawData())
				resp.RegisteredAt, _ = molecule.Bytes2GoU64(accountData.RegisteredAt().RawData())
				resp.Records = accountData.Records()
				resp.RecordsHashBys, _ = blake2b.Blake256(accountData.Records().AsSlice())
				respMap[resp.Account] = &resp
			case common.GoDataEntityVersion2:
				accountData, err := molecule.AccountCellDataFromSlice(dataEntity.Entity().RawData(), false)
				if err != nil {
					return false, fmt.Errorf("AccountSaleCellDataFromSlice err: %s", err.Error())
				}
				resp.AccountCellData = accountData
				resp.Account = common.AccountCharsToAccount(accountData.Account())
				resp.AccountId = common.Bytes2Hex(accountData.Id().RawData())
				resp.Status, _ = molecule.Bytes2GoU8(accountData.Status().RawData())
				resp.RegisteredAt, _ = molecule.Bytes2GoU64(accountData.RegisteredAt().RawData())
				resp.Records = accountData.Records()
				resp.RecordsHashBys, _ = blake2b.Blake256(accountData.Records().AsSlice())
				respMap[resp.Account] = &resp
			default:
				return false, fmt.Errorf("account version: %d", version)
			}
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(respMap) == 0 {
		return nil, fmt.Errorf("not exist account cell")
	}
	return respMap, nil
}

func AccountIdCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*AccountCellDataBuilder, error) {
	respMap, err := AccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*AccountCellDataBuilder)
	for k, v := range respMap {
		k1 := v.AccountId
		retMap[k1] = respMap[k]
	}
	return retMap, nil
}
func (a *AccountCellDataBuilder) getOldDataEntityOpt(p *AccountCellParam) *molecule.DataEntityOpt {
	var oldDataEntity molecule.DataEntity
	switch a.Version {
	case common.GoDataEntityVersion1:
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountCellDataV1.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	case common.GoDataEntityVersion2:
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountCellData.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	}
	oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
	return &oldDataEntityOpt
}

func (a *AccountCellDataBuilder) getNewAccountCellDataBuilder() *molecule.AccountCellDataBuilder {
	var newBuilder molecule.AccountCellDataBuilder
	switch a.Version {
	case common.GoDataEntityVersion1:
		temNewBuilder := molecule.NewAccountCellDataBuilder()
		temNewBuilder.Records(*a.AccountCellDataV1.Records()).Id(*a.AccountCellDataV1.Id()).
			Status(*a.AccountCellDataV1.Status()).Account(*a.AccountCellDataV1.Account()).
			RegisteredAt(*a.AccountCellDataV1.RegisteredAt()).
			LastTransferAccountAt(molecule.TimestampDefault()).
			LastEditRecordsAt(molecule.TimestampDefault()).
			LastEditManagerAt(molecule.TimestampDefault()).Build()
		newBuilder = *temNewBuilder
	case common.GoDataEntityVersion2:
		newBuilder = a.AccountCellData.AsBuilder()
	}
	return &newBuilder
}

func (a *AccountCellDataBuilder) GenWitness(p *AccountCellParam) ([]byte, []byte, error) {

	switch p.Action {
	case common.DasActionRenewAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()
		newAccountSaleCellData := newBuilder.Build()
		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionEditRecords:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		lastEditRecordsAt := molecule.NewTimestampBuilder().Set(molecule.GoTimeUnixToMoleculeBytes(p.LastEditRecordsAt)).Build()
		newBuilder.LastEditRecordsAt(lastEditRecordsAt)
		if len(p.Records) == 0 {
			newBuilder.Records(molecule.RecordsDefault())
		} else {
			records := molecule.RecordsDefault()
			recordsBuilder := records.AsBuilder()
			for _, v := range p.Records {
				record := molecule.RecordDefault()
				recordBuilder := record.AsBuilder()
				recordBuilder.RecordKey(molecule.GoString2MoleculeBytes(v.Key)).
					RecordType(molecule.GoString2MoleculeBytes(v.Type)).
					RecordLabel(molecule.GoString2MoleculeBytes(v.Label)).
					RecordValue(molecule.GoString2MoleculeBytes(v.Value)).
					RecordTtl(molecule.GoU32ToMoleculeU32(v.TTL))
				recordsBuilder.Push(recordBuilder.Build())
			}
			newBuilder.Records(recordsBuilder.Build())
		}
		newAccountSaleCellData := newBuilder.Build()
		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionEditManager:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		lastEditManagerAt := molecule.NewTimestampBuilder().Set(molecule.GoTimeUnixToMoleculeBytes(p.LastEditManagerAt)).Build()
		newBuilder.LastEditManagerAt(lastEditManagerAt)
		newAccountSaleCellData := newBuilder.Build()
		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionTransferAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		newBuilder.Records(molecule.RecordsDefault())
		lastTransferAccountAt := molecule.NewTimestampBuilder().Set(molecule.GoTimeUnixToMoleculeBytes(p.LastTransferAccountAt)).Build()
		newBuilder.LastTransferAccountAt(lastTransferAccountAt)
		newAccountSaleCellData := newBuilder.Build()
		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionBuyAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		newBuilder.Records(molecule.RecordsDefault())
		newAccountSaleCellData := newBuilder.Status(molecule.GoU8ToMoleculeU8(p.Status)).Build()
		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionCancelAccountSale, common.DasActionStartAccountSale, common.DasActionAcceptOffer:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		newAccountSaleCellData := newBuilder.Status(molecule.GoU8ToMoleculeU8(p.Status)).Build()
		newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
	case common.DasActionPropose, common.DasActionDeclareReverseRecord, common.DasActionRedeclareReverseRecord:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		tmp := molecule.NewDataBuilder().Dep(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, nil, nil
	case common.DasActionConfirmProposal:
		if p.SubAction == "exist" {
			oldDataEntityOpt := a.getOldDataEntityOpt(p)

			newBuilder := a.getNewAccountCellDataBuilder()
			newAccountSaleCellData := newBuilder.Build()
			newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

			newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
				Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
			newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

			tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
			witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
			return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
		} else if p.SubAction == "new" {
			accountId, err := molecule.AccountIdFromSlice(common.Hex2Bytes(p.AccountId), false)
			if err != nil {
				return nil, nil, fmt.Errorf("AccountIdFromSlice err: %s", err.Error())
			}
			newAccountSaleCellData := molecule.NewAccountCellDataBuilder().
				Status(molecule.GoU8ToMoleculeU8(uint8(0))).
				Records(molecule.RecordsDefault()).
				LastTransferAccountAt(molecule.TimestampDefault()).
				LastEditRecordsAt(molecule.TimestampDefault()).
				LastEditManagerAt(molecule.TimestampDefault()).
				RegisteredAt(molecule.GoU64ToMoleculeU64(p.RegisterAt)).
				Id(*accountId).
				Account(*p.AccountChars).
				Build()
			newAccountSaleCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountSaleCellData.AsSlice())

			newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountSaleCellDataBytes).
				Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
			newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
			tmp := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
			witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
			return witness, common.Blake2b(newAccountSaleCellData.AsSlice()), nil
		} else {
			return nil, nil, fmt.Errorf("not exist sub action [%s]", p.SubAction)
		}
	}
	return nil, nil, fmt.Errorf("not exist action [%s]", p.Action)
}

type AccountCellRecord struct {
	Key   string
	Type  string
	Label string
	Value string
	TTL   uint32
}

func (a *AccountCellDataBuilder) RecordList() []AccountCellRecord {
	var list []AccountCellRecord
	for index, lenRecords := uint(0), a.Records.Len(); index < lenRecords; index++ {
		record := a.Records.Get(index)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		list = append(list, AccountCellRecord{
			Key:   string(record.RecordKey().RawData()),
			Type:  string(record.RecordType().RawData()),
			Label: string(record.RecordLabel().RawData()),
			Value: string(record.RecordValue().RawData()),
			TTL:   ttl,
		})
	}
	return list
}
