package types

type LockHashIndexState struct {
	BlockHash   Hash   `json:"block_hash"`
	BlockNumber uint64 `json:"block_number"`
	LockHash    Hash   `json:"lock_hash"`
}

type TransactionPoint struct {
	BlockNumber uint64 `json:"block_number"`
	Index       uint   `json:"index"`
	TxHash      Hash   `json:"tx_hash"`
}

type LiveCell struct {
	CellOutput *CellOutput       `json:"cell_output"`
	CreatedBy  *TransactionPoint `json:"created_by"`
}

type CellTransaction struct {
	ConsumedBy *TransactionPoint `json:"consumed_by"`
	CreatedBy  *TransactionPoint `json:"created_by"`
}
