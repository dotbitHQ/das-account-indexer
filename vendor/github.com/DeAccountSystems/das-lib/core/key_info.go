package core

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
)

type ChainTypeAddress struct {
	Type    string  `json:"type"` // blockchain
	KeyInfo KeyInfo `json:"key_info"`
}

type KeyInfo struct {
	CoinType common.CoinType `json:"coin_type"`
	ChainId  common.ChainId  `json:"chain_id"`
	Key      string          `json:"key"`
}

func (c *ChainTypeAddress) FormatChainTypeAddress(net common.DasNetType, is712 bool) (*DasAddressHex, error) {
	if c.Type != "blockchain" {
		return nil, fmt.Errorf("not support type[%s]", c.Type)
	}
	dasChainType := common.FormatCoinTypeToDasChainType(c.KeyInfo.CoinType)
	if dasChainType == -1 {
		dasChainType = common.FormatChainIdToDasChainType(net, c.KeyInfo.ChainId)
	}
	if dasChainType == -1 {
		return nil, fmt.Errorf("not support coin type[%s]-chain id[%s]", c.KeyInfo.CoinType, c.KeyInfo.ChainId)
	}

	daf := DasAddressFormat{DasNetType: net}
	addrHex, err := daf.NormalToHex(DasAddressNormal{
		ChainType:     dasChainType,
		AddressNormal: c.KeyInfo.Key,
		Is712:         is712,
	})
	if err != nil {
		return nil, fmt.Errorf("address NormalToHex err: %s", err.Error())
	}

	return &addrHex, nil
}
