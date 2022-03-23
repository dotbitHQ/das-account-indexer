package witness

import (
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
)

var (
	log                   = mylog.NewLogger("witness", mylog.LevelDebug)
	ErrDataEntityOptIsNil = errors.New("DataEntityOpt is nil")
	ErrNotExistWitness    = errors.New("the witness wanted not exist")

	DataEntityVersion1 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion1)
	DataEntityVersion2 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion2)
	DataEntityVersion3 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion3)
)

func GetWitnessDataFromTx(tx *types.Transaction, handle FuncParseWitness) error {
	inputsSize := len(tx.Inputs)
	witnessesSize := len(tx.Witnesses)
	for i := inputsSize; i < witnessesSize; i++ {
		dataBys := tx.Witnesses[i]
		if len(dataBys) <= common.WitnessDasTableTypeEndIndex+1 {
			continue
		} else if string(dataBys[0:common.WitnessDasCharLen]) != common.WitnessDas {
			continue
		} else {
			actionDataType := common.Bytes2Hex(dataBys[common.WitnessDasCharLen:common.WitnessDasTableTypeEndIndex])
			if goON, err := handle(actionDataType, dataBys[common.WitnessDasTableTypeEndIndex:]); err != nil {
				return err
			} else if !goON {
				return nil
			}
		}
	}
	return nil
}

type FuncParseWitness func(actionDataType common.ActionDataType, dataBys []byte) (bool, error)

func getDataEntityOpt(dataBys []byte, dataType common.DataType) (*molecule.DataEntityOpt, *molecule.DataEntity, error) {
	data, err := molecule.DataFromSlice(dataBys, false)
	if err != nil {
		return nil, nil, fmt.Errorf("DataFromSlice err: %s", err.Error())
	}
	var dataEntityOpt *molecule.DataEntityOpt
	switch dataType {
	case common.DataTypeNew:
		dataEntityOpt = data.New()
	case common.DataTypeOld:
		dataEntityOpt = data.Old()
	case common.DataTypeDep:
		dataEntityOpt = data.Dep()
	}
	if dataEntityOpt == nil || dataEntityOpt.IsNone() {
		return nil, nil, ErrDataEntityOptIsNil
	}
	dataEntity, err := molecule.DataEntityFromSlice(dataEntityOpt.AsSlice(), false)
	if err != nil {
		return nil, nil, fmt.Errorf("DataEntityFromSlice err: %s", err.Error())
	}

	return dataEntityOpt, dataEntity, nil
}

func GenDasDataWitness(action common.ActionDataType, data *molecule.Data) []byte {
	tmp := append([]byte(common.WitnessDas), common.Hex2Bytes(action)...)
	tmp = append(tmp, data.AsSlice()...)
	return tmp
}

func GenDasDataWitnessWithByte(action common.ActionDataType, data []byte) []byte {
	tmp := append([]byte(common.WitnessDas), common.Hex2Bytes(action)...)
	tmp = append(tmp, data...)
	return tmp
}
