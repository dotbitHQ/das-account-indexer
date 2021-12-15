package types

type RawTxPool struct {
	Pending  []*Hash `json:"pending"`
	Proposed []*Hash `json:"proposed"`
}
