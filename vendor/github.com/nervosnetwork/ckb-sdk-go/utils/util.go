package utils

import (
	"context"
	"errors"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

func ParseSudtAmount(outputData []byte) (*big.Int, error) {
	if len(outputData) == 0 {
		return big.NewInt(0), nil
	}
	tmpData := make([]byte, len(outputData))
	copy(tmpData, outputData)
	if len(tmpData) < 16 {
		return nil, errors.New("invalid sUDT amount")
	}
	b := tmpData[0:16]
	b = reverse(b)

	return big.NewInt(0).SetBytes(b), nil
}

func GenerateSudtAmount(amount *big.Int) []byte {
	b := amount.Bytes()
	b = reverse(b)
	if len(b) < 16 {
		for i := len(b); i < 16; i++ {
			b = append(b, 0)
		}
	}

	return b
}

func reverse(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
	}
	return b
}

func RemoveCellOutput(cellOutputs []*types.CellOutput, index int) []*types.CellOutput {
	ret := make([]*types.CellOutput, 0)
	ret = append(ret, cellOutputs[:index]...)
	return append(ret, cellOutputs[index+1:]...)
}

func RemoveCellOutputData(cellOutputData [][]byte, index int) [][]byte {
	ret := make([][]byte, 0)
	ret = append(ret, cellOutputData[:index]...)
	return append(ret, cellOutputData[index+1:]...)
}

// GetMaxMatureBlockNumber return max mature block number
func GetMaxMatureBlockNumber(client rpc.Client, ctx context.Context) (uint64, error) {
	var cellbaseMaturity *types.EpochParams
	cellbaseMaturity, err := getCellbaseMaturity(client, ctx, cellbaseMaturity)
	if err != nil {
		return 0, err
	}
	tipHeader, err := client.GetTipHeader(ctx)
	if err != nil {
		return 0, err
	}
	tipEpoch := types.ParseEpoch(tipHeader.Epoch)
	if tipEpoch.Number < cellbaseMaturity.Number {
		return 0, errors.New("there are no mature live cells")
	} else {
		maxMatureEpoch, err := client.GetEpochByNumber(ctx, tipEpoch.Number-cellbaseMaturity.Number)
		if err != nil {
			return 0, err
		}
		number, err := calcMaxMatureBlockNumber(tipEpoch, maxMatureEpoch.StartNumber, maxMatureEpoch.Length, cellbaseMaturity)
		if err != nil {
			return 0, err
		}
		return number, nil
	}
}

func getCellbaseMaturity(client rpc.Client, ctx context.Context, cellbaseMaturity *types.EpochParams) (*types.EpochParams, error) {
	nodeInfo, err := client.LocalNodeInfo(ctx)
	if err != nil {
		return nil, err
	}
	major, minor, _, err := parseNodeVersion(nodeInfo.Version)
	if err != nil {
		return nil, err
	}
	if major > 0 || minor >= 39 {
		consensus, err := client.GetConsensus(ctx)
		if err != nil {
			return nil, err
		}
		cellbaseMaturity = types.ParseEpoch(consensus.CellbaseMaturity)
	} else {
		cellbaseMaturity = &types.EpochParams{
			Length: 1,
			Index:  0,
			Number: 4,
		}
	}
	return cellbaseMaturity, nil
}

// startNumber is maxMatureEpoch.StartNumber, length is maxMatureEpoch.Length
func calcMaxMatureBlockNumber(tipEpoch *types.EpochParams, startNumber, length uint64, cellbaseMaturity *types.EpochParams) (uint64, error) {
	tipEpochR := big.NewRat(
		int64(tipEpoch.Number*tipEpoch.Length+tipEpoch.Index),
		int64(tipEpoch.Length),
	)
	cellbaseMaturityR := big.NewRat(
		int64(cellbaseMaturity.Number*cellbaseMaturity.Length+cellbaseMaturity.Index),
		int64(cellbaseMaturity.Length),
	)

	if isTipEpochLessThanCellbaseMaturity(tipEpochR, cellbaseMaturityR) {
		return 0, nil
	} else {
		epochDeltaR := big.NewRat(0, 1).Sub(tipEpochR, cellbaseMaturityR)
		num := new(big.Int).SetInt64(0).Div(epochDeltaR.Num(), epochDeltaR.Denom()).Int64()
		decimalR := big.NewRat(0, 1).Sub(epochDeltaR, big.NewRat(num, 1))
		indexR := big.NewRat(0, 1).Mul(decimalR, big.NewRat(int64(length), 1))
		iNum := new(big.Int).SetInt64(0).Div(indexR.Num(), indexR.Denom()).Uint64()
		blockNumber := iNum + startNumber

		return blockNumber, nil
	}
}

func isTipEpochLessThanCellbaseMaturity(tipEpochR, cellbaseMaturityR *big.Rat) bool {
	if tipEpochR.Cmp(cellbaseMaturityR) < 0 {
		return true
	}
	return false
}

func parseNodeVersion(nodeVersion string) (int, int, int, error) {
	reg, err := regexp.Compile("\\d+(\\.\\d+){0,2}")
	if err != nil {
		return 0, 0, 0, err
	}
	parts := reg.FindString(nodeVersion)
	//parts := strings.Split(nodeVersion, " (")
	versionArr := strings.Split(parts, ".")
	major, err := strconv.Atoi(versionArr[0])
	if err != nil {
		return 0, 0, 0, err
	}
	minor, err := strconv.Atoi(versionArr[1])
	if err != nil {
		return 0, 0, 0, err
	}
	patch, err := strconv.Atoi(versionArr[2])
	if err != nil {
		return 0, 0, 0, err
	}
	return major, minor, patch, nil
}

// IsMature check if a cellbase live cell is mature
func IsMature(cell *indexer.LiveCell, maxMatureBlockNumber uint64) bool {
	return cell.TxIndex > 0 || cell.BlockNumber == 0 || cell.BlockNumber <= maxMatureBlockNumber
}
