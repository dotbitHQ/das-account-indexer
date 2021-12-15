package types

type SyncState struct {
	Ibd                     bool   `json:"ibd"`
	BestKnownBlockNumber    uint64 `json:"best_known_block_number"`
	BestKnownBlockTimestamp uint64 `json:"best_known_block_timestamp"`
	OrphanBlocksCount       uint64 `json:"orphan_blocks_count"`
	InflightBlocksCount     uint64 `json:"inflight_blocks_count"`
	FastTime                uint64 `json:"fast_time"`
	LowTime                 uint64 `json:"low_time"`
	NormalTime              uint64 `json:"normal_time"`
}
