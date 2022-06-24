package code

import "github.com/dotbitHQ/das-lib/common"

type CoinType string // EIP-155

const (
	CoinTypeCKB   = "309"
	CoinTypeEth   = "60"
	CoinTypeTrx   = "195"
	CoinTypeBNB   = "714"
	CoinTypeMatic = "966"
)

func FormatCoinTypeToDasChainType(coinType CoinType) common.ChainType {
	switch coinType {
	case CoinTypeCKB:
		return common.ChainTypeCkbMulti
	case CoinTypeEth, CoinTypeBNB, CoinTypeMatic:
		return common.ChainTypeEth
	case CoinTypeTrx:
		return common.ChainTypeTron
	}
	return -1
}
