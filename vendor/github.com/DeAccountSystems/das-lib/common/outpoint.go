package common

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strconv"
	"strings"
)

func String2OutPoint(str string) (txHash string, index uint) {
	list := strings.Split(str, "-")
	if len(list) > 1 {
		txHash = list[0]
		index64, _ := strconv.ParseInt(list[1], 10, 64)
		index = uint(index64)
	}
	return
}

func String2OutPointStruct(str string) *types.OutPoint {
	txHash, index := String2OutPoint(str)
	return &types.OutPoint{
		TxHash: types.HexToHash(txHash),
		Index:  index,
	}
}

func OutPoint2String(txHash string, index uint) string {
	return fmt.Sprintf("%s-%d", txHash, index)
}

func OutPointStruct2String(o *types.OutPoint) string {
	return OutPoint2String(o.TxHash.Hex(), o.Index)
}
