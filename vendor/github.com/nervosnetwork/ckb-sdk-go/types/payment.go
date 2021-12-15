package types

import "math/big"

type ReceiverInfo struct {
	Receiver *Script
	Amount   *big.Int
}
