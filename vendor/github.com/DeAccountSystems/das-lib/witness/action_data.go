package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type ActionDataBuilder struct {
	ActionData *molecule.ActionData
	Action     common.DasAction
	Params     [][]byte
	ParamsStr  string
}

func (a *ActionDataBuilder) ActionBuyAccountInviterScript() (*molecule.Script, error) {
	if len(a.Params) != 3 {
		return nil, fmt.Errorf("len params err:[%d]", len(a.Params))
	}
	inviterScript, err := molecule.ScriptFromSlice(a.Params[0], false)
	if err != nil {
		return nil, fmt.Errorf("ScriptFromSlice err: %s", err.Error())
	}
	return inviterScript, nil
}

func (a *ActionDataBuilder) ActionBuyAccountChannelScript() (*molecule.Script, error) {
	if len(a.Params) != 3 {
		return nil, fmt.Errorf("len params err:[%d]", len(a.Params))
	}
	channelScript, err := molecule.ScriptFromSlice(a.Params[1], false)
	if err != nil {
		return nil, fmt.Errorf("ScriptFromSlice err: %s", err.Error())
	}
	return channelScript, nil
}

func ActionDataBuilderFromTx(tx *types.Transaction) (*ActionDataBuilder, error) {
	var resp ActionDataBuilder
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeActionData:
			actionData, err := molecule.ActionDataFromSlice(dataBys, false)
			if err != nil {
				return false, fmt.Errorf("ActionDataFromSlice err: %s", err.Error())
			}
			resp.ActionData = actionData
			resp.Action = string(actionData.Action().RawData())
			if resp.Action == common.DasActionBuyAccount {
				raw := actionData.Params().RawData()

				lenRaw := len(raw)
				inviterLockBytesLen, err := molecule.Bytes2GoU32(raw[:4])
				if err != nil {
					return false, fmt.Errorf("Bytes2GoU32 err: %s", err.Error())
				}
				inviterLockRaw := raw[:inviterLockBytesLen]
				channelLockRaw := raw[inviterLockBytesLen : lenRaw-1]

				resp.Params = append(resp.Params, inviterLockRaw)
				resp.Params = append(resp.Params, channelLockRaw)
				resp.Params = append(resp.Params, raw[lenRaw-1:lenRaw])
				resp.ParamsStr = common.GetMaxHashLenParams(common.Bytes2Hex(inviterLockRaw)) + "," + common.GetMaxHashLenParams(common.Bytes2Hex(channelLockRaw)) + "," + common.Bytes2Hex(raw[lenRaw-1:lenRaw])
			} else {
				resp.Params = append(resp.Params, actionData.Params().RawData())
				resp.ParamsStr = common.Bytes2Hex(actionData.Params().RawData())
			}
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if resp.ActionData == nil {
		return nil, fmt.Errorf("not exist action data")
	}
	return &resp, nil
}

func ActionDataBuilderFromWitness(wit []byte) (*ActionDataBuilder, error) {
	if len(wit) <= common.WitnessDasTableTypeEndIndex+1 {
		return nil, fmt.Errorf("action data len is invalid")
	} else if string(wit[0:common.WitnessDasCharLen]) != common.WitnessDas {
		return nil, fmt.Errorf("not a das data")
	}
	actionDataType := common.Bytes2Hex(wit[common.WitnessDasCharLen:common.WitnessDasTableTypeEndIndex])
	dataBys := wit[common.WitnessDasTableTypeEndIndex:]
	if actionDataType != common.ActionDataTypeActionData {
		return nil, fmt.Errorf("not a action data")
	}
	actionData, err := molecule.ActionDataFromSlice(dataBys, false)
	if err != nil {
		return nil, fmt.Errorf("ActionDataFromSlice err: %s", err.Error())
	}
	var resp ActionDataBuilder
	resp.ActionData = actionData
	resp.Action = string(actionData.Action().RawData())
	if resp.Action == common.DasActionBuyAccount {
		raw := actionData.Params().RawData()

		lenRaw := len(raw)
		inviterLockBytesLen, err := molecule.Bytes2GoU32(raw[:4])
		if err != nil {
			return nil, fmt.Errorf("Bytes2GoU32 err: %s", err.Error())
		}
		inviterLockRaw := raw[:inviterLockBytesLen]
		channelLockRaw := raw[inviterLockBytesLen : lenRaw-1]

		resp.Params = append(resp.Params, inviterLockRaw)
		resp.Params = append(resp.Params, channelLockRaw)
		resp.Params = append(resp.Params, raw[lenRaw-1:lenRaw])
		resp.ParamsStr = common.GetMaxHashLenParams(common.Bytes2Hex(inviterLockRaw)) + "," + common.GetMaxHashLenParams(common.Bytes2Hex(channelLockRaw)) + "," + common.Bytes2Hex(raw[lenRaw-1:lenRaw])
	} else {
		resp.Params = append(resp.Params, actionData.Params().RawData())
		resp.ParamsStr = common.Bytes2Hex(actionData.Params().RawData())
	}
	return &resp, nil
}

func GenActionDataWitness(action common.DasAction, params []byte) ([]byte, error) {
	if action == "" {
		return nil, fmt.Errorf("action is nil")
	}
	if params == nil {
		params = []byte{}
	}
	if action == common.DasActionEditRecords {
		params = append(params, common.Hex2Bytes(common.ParamManager)...)
	} else if action == common.DasActionRenewAccount {
		params = []byte{}
	} else {
		params = append(params, common.Hex2Bytes(common.ParamOwner)...)
	}
	actionBytes := molecule.GoString2MoleculeBytes(action)
	paramsBytes := molecule.GoBytes2MoleculeBytes(params)
	actionData := molecule.NewActionDataBuilder().Action(actionBytes).Params(paramsBytes).Build()

	tmp := append([]byte(common.WitnessDas), common.Hex2Bytes(common.ActionDataTypeActionData)...)
	tmp = append(tmp, actionData.AsSlice()...)
	return tmp, nil
}

func GenBuyAccountParams(inviterScript, channelScript *types.Script) []byte {
	iScript := molecule.CkbScript2MoleculeScript(inviterScript)
	paramsInviter := iScript.AsSlice()
	cScript := molecule.CkbScript2MoleculeScript(channelScript)
	paramsChannel := cScript.AsSlice()
	return append(paramsInviter, paramsChannel...)
}
