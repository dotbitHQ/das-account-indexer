package code

import "github.com/DeAccountSystems/das-lib/common"

type ChainId string //BIP-44

const (
	ChainIdEthMainNet     = "1"
	ChainIdBscMainNet     = "56"
	ChainIdPolygonMainNet = "137"

	ChainIdEthTestNet     = "5" // Goerli
	ChainIdBscTestNet     = "97"
	ChainIdPolygonTestNet = "80001"
)

func FormatChainIdToDasChainType(netType common.DasNetType, chainId ChainId) common.ChainType {
	if netType == common.DasNetTypeMainNet {
		switch chainId {
		case ChainIdEthMainNet, ChainIdBscMainNet, ChainIdPolygonMainNet:
			return common.ChainTypeEth
		}
	} else {
		switch chainId {
		case ChainIdEthTestNet, ChainIdBscTestNet, ChainIdPolygonTestNet:
			return common.ChainTypeEth
		}
	}
	return -1
}
