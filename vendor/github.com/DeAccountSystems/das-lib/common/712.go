package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"math/big"
	"strings"
)

type MMJsonObj struct {
	Types struct {
		EIP712Domain []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"EIP712Domain"`
		Action []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"Action"`
		Cell []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"Cell"`
		Transaction []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"Transaction"`
	} `json:"types"`
	PrimaryType string `json:"primaryType"`
	Domain      struct {
		ChainID           int64  `json:"chainId"`
		Name              string `json:"name"`
		VerifyingContract string `json:"verifyingContract"`
		Version           string `json:"version"`
	} `json:"domain"`
	Message struct {
		DasMessage      string           `json:"DAS_MESSAGE"`
		InputsCapacity  string           `json:"inputsCapacity"`
		OutputsCapacity string           `json:"outputsCapacity"`
		Fee             string           `json:"fee"`
		Digest          string           `json:"digest"`
		Action          *MMJsonAction    `json:"action"`
		Inputs          []MMJsonCellInfo `json:"inputs"`
		Outputs         []MMJsonCellInfo `json:"outputs"`
	} `json:"message"`
}

type MMJsonAction struct {
	Action string `json:"action"`
	Params string `json:"params"`
}

type MMJsonCellInfo struct {
	Capacity  string `json:"capacity"`
	LockStr   string `json:"lock"`
	TypeStr   string `json:"type"`
	Data      string `json:"data"`
	ExtraData string `json:"extraData"`
}

func (m *MMJsonObj) String() string {
	bys, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(bys)
}

const (
	MaxHashLen   = 20
	MMJsonObjStr = `{
  "types": {
    "EIP712Domain": [
      {"name": "chainId", "type": "uint256"},
      {"name": "name", "type": "string"},
      {"name": "verifyingContract", "type": "address"},
      {"name": "version", "type": "string"}
    ],
    "Action": [
      {"name": "action", "type": "string"},
      {"name": "params", "type": "string"}
    ],
    "Cell": [
      {"name": "capacity", "type": "string"},
      {"name": "lock", "type": "string"},
      {"name": "type", "type": "string"},
      {"name": "data", "type": "string"},
      {"name": "extraData", "type": "string"}
    ],
    "Transaction": [
      {"name": "DAS_MESSAGE", "type": "string"},
      {"name": "inputsCapacity", "type": "string"},
      {"name": "outputsCapacity", "type": "string"},
      {"name": "fee", "type": "string"},
      {"name": "action", "type": "Action"},
      {"name": "inputs", "type": "Cell[]"},
      {"name": "outputs", "type": "Cell[]"},
      {"name": "digest", "type": "bytes32"}
    ]
  },
  "primaryType": "Transaction",
  "domain": {
    "chainId": 1,
    "name": "da.systems",
    "verifyingContract": "0x0000000000000000000000000000000020210722",
    "version": "1"
  },
  "message": {
    "DAS_MESSAGE": "",
    "inputsCapacity": "",
    "outputsCapacity": "",
    "fee": "",
    "action": {},
    "inputs": [],
    "outputs": []
  }
}`
)

func GetMaxHashLenParams(s string) string {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s) > MaxHashLen {
		s = s[:MaxHashLen] + "..."
	}
	return "0x" + s
}

func GetMaxHashLenData(data []byte) string {
	if len(data) > MaxHashLen {
		return "0x" + hex.EncodeToString(data[:MaxHashLen]) + "..."
	} else {
		if len(data) == 0 {
			return ""
		}
		return "0x" + hex.EncodeToString(data)
	}
}

func GetMaxHashLenScript(script *types.Script, dasContractName DasContractName) string {
	if script == nil || dasContractName == "" {
		return ""
	}
	tmp := ""
	if len(script.Args) > MaxHashLen {
		tmp = "0x" + hex.EncodeToString(script.Args[:MaxHashLen]) + "..."
	} else {
		tmp = "0x" + hex.EncodeToString(script.Args)
	}
	return fmt.Sprintf("%s,0x01,%s", dasContractName, tmp)
}

func GetAccountCellExpiredAtFromOutputData(data []byte) (uint64, error) {
	if size := len(data); size < ExpireTimeEndIndex {
		return 0, fmt.Errorf("invalid data, len not enough, your: %d, want: %d", size, ExpireTimeEndIndex)
	}
	expireTime, err := molecule.Bytes2GoU64(data[ExpireTimeEndIndex-8 : ExpireTimeEndIndex])
	if err != nil {
		return 0, fmt.Errorf("BytesToGoU64 err: %s", err)
	}
	return expireTime, nil
}

func GetAccountCellNextAccountIdFromOutputData(data []byte) ([]byte, error) {
	if size := len(data); size < NextAccountIdEndIndex {
		return nil, fmt.Errorf("invalid data, len not enough, your: %d, want: %d", size, NextAccountIdEndIndex)
	}
	return data[NextAccountIdStartIndex:NextAccountIdEndIndex], nil
}

func Capacity2Str(capacity uint64) string {
	capacityRat := new(big.Rat).SetInt(new(big.Int).SetUint64(capacity))
	oneCkbRat := new(big.Rat).SetInt(new(big.Int).SetUint64(OneCkb))
	capacityStr := new(big.Rat).Quo(capacityRat, oneCkbRat).FloatString(8)
	if !strings.Contains(capacityStr, ".") {
		return capacityStr
	}
	for i := len(capacityStr) - 1; i >= 0; i-- {
		if capacityStr[i] == '.' {
			return capacityStr[:i]
		}
		if capacityStr[i] != '0' {
			return capacityStr[:i+1]
		}

	}
	return capacityStr
}
