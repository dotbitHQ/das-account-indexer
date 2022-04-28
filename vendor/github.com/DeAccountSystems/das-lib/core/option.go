package core

import (
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
)

type DasCoreOption func(*DasCore)

func WithClient(client rpc.Client) DasCoreOption {
	return func(dc *DasCore) {
		dc.client = client
	}
}

func WithDasContractCodeHash(dasContractCodeHash string) DasCoreOption {
	return func(dc *DasCore) {
		dc.dasContractCodeHash = dasContractCodeHash
	}
}

func WithDasContractArgs(dasContractArgs string) DasCoreOption {
	return func(dc *DasCore) {
		dc.dasContractArgs = dasContractArgs
	}
}

func WithTHQCodeHash(thqCodeHash string) DasCoreOption {
	return func(dc *DasCore) {
		dc.thqCodeHash = thqCodeHash
	}
}

func WithDasNetType(net common.DasNetType) DasCoreOption {
	return func(dc *DasCore) {
		dc.net = net
		dc.daf = &DasAddressFormat{DasNetType: net}
	}
}
