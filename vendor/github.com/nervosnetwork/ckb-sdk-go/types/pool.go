package types

type TxPoolInfo struct {
	LastTxsUpdatedAt uint64 `json:"last_txs_updated_at"`
	Orphan           uint64 `json:"orphan"`
	Pending          uint64 `json:"pending"`
	Proposed         uint64 `json:"proposed"`
	TotalTxCycles    uint64 `json:"total_tx_cycles"`
	TotalTxSize      uint64 `json:"total_tx_size"`
}
