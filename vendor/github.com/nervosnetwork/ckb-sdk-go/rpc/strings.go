package rpc

import (
	"encoding/json"

	"github.com/nervosnetwork/ckb-sdk-go/types"
)

func TransactionString(tx *types.Transaction) (string, error) {
	itx := fromTransaction(tx)
	bytes, err := json.Marshal(itx)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func TransactionFromString(tx string) (*types.Transaction, error) {
	var itx transaction
	err := json.Unmarshal([]byte(tx), &itx)
	if err != nil {
		return nil, err
	}

	return toTransaction(itx), nil
}
