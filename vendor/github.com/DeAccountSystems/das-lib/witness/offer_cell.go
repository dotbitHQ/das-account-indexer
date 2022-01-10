package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type OfferCellBuilder struct {
	Index         uint32
	Version       uint32
	OfferCellData *molecule.OfferCellData
	DataEntityOpt *molecule.DataEntityOpt
	Account       string
	Price         uint64
	Message       string
	InviterLock   *molecule.Script
	ChannelLock   *molecule.Script
}

func OfferCellDataBuilderMapFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*OfferCellBuilder, error) {
	var respMap = make(map[string]*OfferCellBuilder)
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeOfferCell:
			var resp OfferCellBuilder

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

			version, err := molecule.Bytes2GoU32(dataEntity.Version().RawData())
			if err != nil {
				return false, fmt.Errorf("get version err: %s", err.Error())
			}
			resp.Version = version

			offerCellData, err := molecule.OfferCellDataFromSlice(dataEntity.Entity().RawData(), false)
			if err != nil {
				return false, fmt.Errorf("AccountSaleCellDataFromSlice err: %s", err.Error())
			}
			resp.OfferCellData = offerCellData
			resp.Account = string(offerCellData.Account().RawData())
			resp.Price, _ = molecule.Bytes2GoU64(offerCellData.Price().RawData())
			resp.Message = string(offerCellData.Message().RawData())
			resp.InviterLock, _ = molecule.ScriptFromSlice(offerCellData.InviterLock().AsSlice(), false)
			resp.ChannelLock, _ = molecule.ScriptFromSlice(offerCellData.ChannelLock().AsSlice(), false)
			if resp.InviterLock == nil {
				tmp := molecule.ScriptDefault()
				resp.InviterLock = &tmp
			}
			if resp.ChannelLock == nil {
				tmp := molecule.ScriptDefault()
				resp.ChannelLock = &tmp
			}
			key := fmt.Sprintf("%s-%d", tx.Hash.Hex(), index)
			respMap[key] = &resp
			return true, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(respMap) == 0 {
		return nil, fmt.Errorf("not exist offer cell")
	}
	return respMap, nil
}

func OfferCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*OfferCellBuilder, error) {
	respMap, err := OfferCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}
	for k, _ := range respMap {
		return respMap[k], nil
	}
	return nil, fmt.Errorf("not exist offer cell")
}

type OfferCellParam struct {
	Action        common.DasAction
	Account       string
	Price         uint64
	Message       string
	InviterScript *types.Script
	ChannelScript *types.Script
	OldIndex      uint32
}

func (o *OfferCellBuilder) getOldDataEntityOpt(p *OfferCellParam) *molecule.DataEntityOpt {
	var oldDataEntity molecule.DataEntity

	oldOfferCellDataBytes := molecule.GoBytes2MoleculeBytes(o.OfferCellData.AsSlice())
	oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldOfferCellDataBytes).
		Version(DataEntityVersion1).
		Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).
		Build()
	oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()

	return &oldDataEntityOpt
}

func (o *OfferCellBuilder) getNewAccountCellDataBuilder() *molecule.OfferCellDataBuilder {
	var newBuilder molecule.OfferCellDataBuilder
	newBuilder = o.OfferCellData.AsBuilder()
	return &newBuilder
}

func (o *OfferCellBuilder) GenWitness(p *OfferCellParam) ([]byte, []byte, error) {
	switch p.Action {
	case common.DasActionMakeOffer:
		iScript := molecule.CkbScript2MoleculeScript(p.InviterScript)
		cScript := molecule.CkbScript2MoleculeScript(p.ChannelScript)
		offerCellData := molecule.NewOfferCellDataBuilder().
			Account(molecule.GoString2MoleculeBytes(p.Account)).
			Price(molecule.GoU64ToMoleculeU64(p.Price)).
			Message(molecule.GoString2MoleculeBytes(p.Message)).
			InviterLock(iScript).ChannelLock(cScript).Build()

		offerCellDataBytes := molecule.GoBytes2MoleculeBytes(offerCellData.AsSlice())
		newDataEntity := molecule.NewDataEntityBuilder().Entity(offerCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(0)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeOfferCell, &tmp)
		return witness, common.Blake2b(offerCellData.AsSlice()), nil
	case common.DasActionEditOffer:
		oldDataEntityOpt := o.getOldDataEntityOpt(p)

		newBuilder := o.getNewAccountCellDataBuilder()
		newOfferCellData := newBuilder.Price(molecule.GoU64ToMoleculeU64(p.Price)).
			Message(molecule.GoString2MoleculeBytes(p.Message)).
			Build()

		newOfferCellDataBytes := molecule.GoBytes2MoleculeBytes(newOfferCellData.AsSlice())
		newDataEntity := molecule.NewDataEntityBuilder().Entity(newOfferCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(0)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeOfferCell, &tmp)
		return witness, common.Blake2b(newOfferCellData.AsSlice()), nil
	case common.DasActionCancelOffer, common.DasActionAcceptOffer:
		oldDataEntityOpt := o.getOldDataEntityOpt(p)
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeOfferCell, &tmp)
		return witness, nil, nil
	}

	return nil, nil, fmt.Errorf("not exist action [%s]", p.Action)
}
