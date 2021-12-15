package utils

import (
	"context"
	"github.com/pkg/errors"

	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

const (
	AnyoneCanPayCodeHashOnLina   = "0xd369597ff47f29fbc0d47d2e3775370d1250b85140c670e4718af712983a2354"
	AnyoneCanPayCodeHashOnAggron = "0x3419a1c09eb2567f6552ee7a8ecffd64155cffe0f1796e6e61ec088d740c1356"
)

type Option func(*SystemScripts)
type SystemScriptCell struct {
	CellHash types.Hash
	OutPoint *types.OutPoint
	HashType types.ScriptHashType
	DepType  types.DepType
}

type SystemScripts struct {
	SecpSingleSigCell *SystemScriptCell
	SecpMultiSigCell  *SystemScriptCell
	DaoCell           *SystemScriptCell
	ACPCell           *SystemScriptCell
	SUDTCell          *SystemScriptCell
	ChequeCell        *SystemScriptCell
}

func secpSingleSigCell(chain string) *SystemScriptCell {
	if chain == "ckb" {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x71a7ba8fc96349fea0ed3a5c47992e3b4084b031a42264a018e0072e8172e46c"),
				Index:  0,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeDepGroup,
		}
	} else {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37"),
				Index:  0,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeDepGroup,
		}
	}
}

func secpMultiSigCell(chain string) *SystemScriptCell {
	if chain == "ckb" {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x5c5069eb0857efc65e1bca0c07df34c31663b3622fd3876c876320fc9634e2a8"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x71a7ba8fc96349fea0ed3a5c47992e3b4084b031a42264a018e0072e8172e46c"),
				Index:  1,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeDepGroup,
		}
	} else {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x5c5069eb0857efc65e1bca0c07df34c31663b3622fd3876c876320fc9634e2a8"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37"),
				Index:  1,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeDepGroup,
		}
	}
}

func daoCell(chain string) *SystemScriptCell {
	if chain == "ckb" {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xe2fb199810d49a4d8beec56718ba2593b665db9d52299a0f9e6e75416d73ff5c"),
				Index:  2,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeCode,
		}
	} else {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x8f8c79eb6671709633fe6a46de93c0fedc9c1b8a6527a18d3983879542635c9f"),
				Index:  2,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeCode,
		}
	}
}

func acpCell(chain string) *SystemScriptCell {
	if chain == "ckb" {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0xd369597ff47f29fbc0d47d2e3775370d1250b85140c670e4718af712983a2354"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x4153a2014952d7cac45f285ce9a7c5c0c0e1b21f2d378b82ac1433cb11c25c4d"),
				Index:  0,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeDepGroup,
		}
	} else {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x3419a1c09eb2567f6552ee7a8ecffd64155cffe0f1796e6e61ec088d740c1356"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xec26b0f85ed839ece5f11c4c4e837ec359f5adc4420410f6453b1f6b60fb96a6"),
				Index:  0,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeDepGroup,
		}
	}
}

func sudtCell(chain string) *SystemScriptCell {
	if chain == "ckb" {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x5e7a36a77e68eecc013dfa2fe6a23f3b6c344b04005808694ae6dd45eea4cfd5"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xc7813f6a415144643970c2e88e0bb6ca6a8edc5dd7c1022746f628284a9936d5"),
				Index:  0,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeCode,
		}
	} else {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0xc5e5dcf215925f7ef4dfaf5f4b4f105bc321c02776d6e7d52a1db3fcd9d011a4"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xe12877ebd2c3c364dc46c5c992bcfaf4fee33fa13eebdf82c591fc9825aab769"),
				Index:  0,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeCode,
		}
	}
}

// mock data
func chequeCell(chain string) *SystemScriptCell {
	if chain == "ckb" {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x5e7a36a77e68eecc013dfa2fe6a23f3b6c344b04005808694ae6dd45eea4cfd5"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xc7813f6a415144643970c2e88e0bb6ca6a8edc5dd7c1022746f628284a9936d5"),
				Index:  0,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeDepGroup,
		}
	} else {
		return &SystemScriptCell{
			CellHash: types.HexToHash("0x9f27f3afc8d26dfa8bc0c8fa21bc033ddcdab6ad83d5e865cdd6d5d0b3b95642"),
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x1dbbeac82db9a330ed07dd33e547facbca14378196f0e2d69ad8e83bce1d5f54"),
				Index:  0,
			},
			HashType: types.HashTypeType,
			DepType:  types.DepTypeDepGroup,
		}
	}
}

// NewSystemScripts returns a SystemScripts object
func NewSystemScripts(client rpc.Client, options ...Option) (*SystemScripts, error) {
	info, err := client.GetBlockchainInfo(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "RPC get_blockchain_info error")
	}
	scripts := &SystemScripts{
		SecpSingleSigCell: secpSingleSigCell(info.Chain),
		SecpMultiSigCell:  secpMultiSigCell(info.Chain),
		DaoCell:           daoCell(info.Chain),
		ACPCell:           acpCell(info.Chain),
		SUDTCell:          sudtCell(info.Chain),
		ChequeCell:        chequeCell(info.Chain),
	}

	for _, option := range options {
		option(scripts)
	}

	return scripts, nil
}

// SecpSingleSigCell set a custom secp single sig cell to SystemScripts
func SecpSingleSigCell(secpSingleSigCell *SystemScriptCell) Option {
	return func(s *SystemScripts) {
		s.SecpSingleSigCell = secpSingleSigCell
	}
}

// SecpMultiSigCell set a custom secp mutisig cell to SystemScripts
func SecpMultiSigCell(secpMultiSigCell *SystemScriptCell) Option {
	return func(s *SystemScripts) {
		s.SecpMultiSigCell = secpMultiSigCell
	}
}

// DaoCell set a custom dao cell to SystemScripts
func DaoCell(daoCell *SystemScriptCell) Option {
	return func(s *SystemScripts) {
		s.DaoCell = daoCell
	}
}

// ACPCell set a custom acp system script cell to SystemScripts
func ACPCell(acpCell *SystemScriptCell) Option {
	return func(s *SystemScripts) {
		s.ACPCell = acpCell
	}
}

// SUDTCell set a custom sudt system script cell to SystemScripts
func SUDTCell(sudtCell *SystemScriptCell) Option {
	return func(s *SystemScripts) {
		s.SUDTCell = sudtCell
	}
}

// ChequeCell set a custom cheque script cell to SystemScripts
func ChequeCell(chequeCell *SystemScriptCell) Option {
	return func(s *SystemScripts) {
		s.ChequeCell = chequeCell
	}
}
