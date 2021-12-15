package types

import "math/big"

type AlertMessage struct {
	Id          string `json:"id"`
	Message     string `json:"message"`
	NoticeUntil uint64 `json:"notice_until"`
	Priority    string `json:"priority"`
}

type BlockchainInfo struct {
	Alerts                 []*AlertMessage `json:"alerts"`
	Chain                  string          `json:"chain"`
	Difficulty             *big.Int        `json:"difficulty"`
	Epoch                  uint64          `json:"epoch"`
	IsInitialBlockDownload bool            `json:"is_initial_block_download"`
	MedianTime             uint64          `json:"median_time"`
}
