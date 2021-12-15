package indexer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type ScriptType string
type SearchOrder string
type IoType string

const (
	ScriptTypeLock ScriptType = "lock"
	ScriptTypeType ScriptType = "type"

	SearchOrderAsc  SearchOrder = "asc"
	SearchOrderDesc SearchOrder = "desc"

	IOTypeIn  IoType = "input"
	IOTypeOut IoType = "output"
)

type SearchKey struct {
	Script     *types.Script `json:"script"`
	ScriptType ScriptType    `json:"script_type"`
	ArgsLen    uint          `json:"args_len,omitempty"`
	Filter     *CellsFilter  `json:"filter,omitempty"`
}

type CellsFilter struct {
	Script              *types.Script `json:"script"`
	OutputDataLenRange  *[2]uint64    `json:"output_data_len_range"`
	OutputCapacityRange *[2]uint64    `json:"output_capacity_range"`
	BlockRange          *[2]uint64    `json:"block_range"`
}

type LiveCell struct {
	BlockNumber uint64            `json:"block_number"`
	OutPoint    *types.OutPoint   `json:"out_point"`
	Output      *types.CellOutput `json:"output"`
	OutputData  []byte            `json:"output_data"`
	TxIndex     uint              `json:"tx_index"`
}

type LiveCells struct {
	LastCursor string      `json:"last_cursor"`
	Objects    []*LiveCell `json:"objects"`
}

type Transaction struct {
	BlockNumber uint64     `json:"block_number"`
	IoIndex     uint       `json:"io_index"`
	IoType      IoType     `json:"io_type"`
	TxHash      types.Hash `json:"tx_hash"`
	TxIndex     uint       `json:"tx_index"`
}

type Transactions struct {
	LastCursor string         `json:"last_cursor"`
	Objects    []*Transaction `json:"objects"`
}

type TipHeader struct {
	BlockHash   types.Hash `json:"block_hash"`
	BlockNumber uint64     `json:"block_number"`
}

type Capacity struct {
	Capacity    uint64     `json:"capacity"`
	BlockHash   types.Hash `json:"block_hash"`
	BlockNumber uint64     `json:"block_number"`
}

type capacity struct {
	Capacity    hexutil.Uint64 `json:"capacity"`
	BlockHash   types.Hash     `json:"block_hash"`
	BlockNumber hexutil.Uint64 `json:"block_number"`
}

type tipHeader struct {
	BlockHash   types.Hash     `json:"block_hash"`
	BlockNumber hexutil.Uint64 `json:"block_number"`
}

type searchKey struct {
	Script     *script      `json:"script"`
	ScriptType ScriptType   `json:"script_type"`
	ArgsLen    hexutil.Uint `json:"args_len,omitempty"`
	Filter     *cellsFilter `json:"filter,omitempty"`
}

type cellsFilter struct {
	Script              *script            `json:"script"`
	OutputDataLenRange  *[2]hexutil.Uint64 `json:"output_data_len_range"`
	OutputCapacityRange *[2]hexutil.Uint64 `json:"output_capacity_range"`
	BlockRange          *[2]hexutil.Uint64 `json:"block_range"`
}

type outPoint struct {
	TxHash types.Hash   `json:"tx_hash"`
	Index  hexutil.Uint `json:"index"`
}

type script struct {
	CodeHash types.Hash           `json:"code_hash"`
	HashType types.ScriptHashType `json:"hash_type"`
	Args     hexutil.Bytes        `json:"args"`
}

type cellOutput struct {
	Capacity hexutil.Uint64 `json:"capacity"`
	Lock     *script        `json:"lock"`
	Type     *script        `json:"type"`
}

type liveCells struct {
	LastCursor string `json:"last_cursor"`
	Objects    []struct {
		BlockNumber hexutil.Uint64 `json:"block_number"`
		OutPoint    *outPoint      `json:"out_point"`
		Output      *cellOutput    `json:"output"`
		OutputData  hexutil.Bytes  `json:"output_data"`
		TxIndex     hexutil.Uint   `json:"tx_index"`
	} `json:"objects"`
}

type transactions struct {
	LastCursor string `json:"last_cursor"`
	Objects    []struct {
		BlockNumber hexutil.Uint64 `json:"block_number"`
		IoIndex     hexutil.Uint   `json:"io_index"`
		IoType      IoType         `json:"io_type"`
		TxHash      types.Hash     `json:"tx_hash"`
		TxIndex     hexutil.Uint   `json:"tx_index"`
	} `json:"objects"`
}

func toTransactions(transactions transactions) *Transactions {
	result := &Transactions{
		LastCursor: transactions.LastCursor,
	}
	result.Objects = make([]*Transaction, len(transactions.Objects))
	for i := 0; i < len(transactions.Objects); i++ {
		transaction := transactions.Objects[i]
		result.Objects[i] = &Transaction{
			BlockNumber: uint64(transaction.BlockNumber),
			IoIndex:     uint(transaction.IoIndex),
			IoType:      transaction.IoType,
			TxHash:      transaction.TxHash,
			TxIndex:     uint(transaction.TxIndex),
		}
	}
	return result
}

func toLiveCells(cells liveCells) *LiveCells {
	result := &LiveCells{
		LastCursor: cells.LastCursor,
	}
	result.Objects = make([]*LiveCell, len(cells.Objects))
	for i := 0; i < len(cells.Objects); i++ {
		cell := cells.Objects[i]
		result.Objects[i] = &LiveCell{
			BlockNumber: uint64(cell.BlockNumber),
			OutPoint: &types.OutPoint{
				TxHash: cell.OutPoint.TxHash,
				Index:  uint(cell.OutPoint.Index),
			},
			OutputData: cell.OutputData,
			TxIndex:    uint(cell.TxIndex),
		}
		result.Objects[i].Output = &types.CellOutput{
			Capacity: uint64(cell.Output.Capacity),
			Lock: &types.Script{
				CodeHash: cell.Output.Lock.CodeHash,
				HashType: cell.Output.Lock.HashType,
				Args:     cell.Output.Lock.Args,
			},
		}
		if cell.Output.Type != nil {
			result.Objects[i].Output.Type = &types.Script{
				CodeHash: cell.Output.Type.CodeHash,
				HashType: cell.Output.Type.HashType,
				Args:     cell.Output.Type.Args,
			}
		}
	}
	return result
}

func fromSearchKey(key *SearchKey) *searchKey {
	result := &searchKey{
		Script: &script{
			CodeHash: key.Script.CodeHash,
			HashType: key.Script.HashType,
			Args:     key.Script.Args,
		},
		ScriptType: key.ScriptType,
	}

	if key.ArgsLen > 0 {
		result.ArgsLen = hexutil.Uint(key.ArgsLen)
	}

	if key.Filter != nil {
		filter := &cellsFilter{}
		if key.Filter.Script != nil {
			filter.Script = &script{
				CodeHash: key.Filter.Script.CodeHash,
				HashType: key.Filter.Script.HashType,
				Args:     key.Filter.Script.Args,
			}
		}
		if key.Filter.OutputDataLenRange != nil {
			filter.OutputDataLenRange = &[2]hexutil.Uint64{hexutil.Uint64(key.Filter.OutputDataLenRange[0]), hexutil.Uint64(key.Filter.OutputDataLenRange[1])}
		}
		if key.Filter.OutputCapacityRange != nil {
			filter.OutputCapacityRange = &[2]hexutil.Uint64{hexutil.Uint64(key.Filter.OutputCapacityRange[0]), hexutil.Uint64(key.Filter.OutputCapacityRange[1])}
		}
		if key.Filter.BlockRange != nil {
			filter.BlockRange = &[2]hexutil.Uint64{hexutil.Uint64(key.Filter.BlockRange[0]), hexutil.Uint64(key.Filter.BlockRange[1])}
		}
		result.Filter = filter
	}

	return result
}
