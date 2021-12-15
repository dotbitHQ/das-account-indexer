package types

type TransactionProof struct {
	Proof         *Proof `json:"proof"`
	BlockHash     Hash   `json:"block_hash"`
	WitnessesRoot Hash   `json:"witnesses_root"`
}

type Proof struct {
	Indices []uint `json:"indices"`
	Lemmas  []Hash `json:"lemmas"`
}
