package types

type DryRunTransactionResult struct {
	Cycles uint64 `json:"cycles"`
}

type EstimateFeeRateResult struct {
	FeeRate uint64 `json:"fee_rate"`
}
