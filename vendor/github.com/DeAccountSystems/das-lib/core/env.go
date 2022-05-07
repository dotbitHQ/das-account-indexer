package core

import "github.com/DeAccountSystems/das-lib/common"

type Env struct {
	THQCodeHash      string
	ContractArgs     string
	ContractCodeHash string
	MapContract      map[common.DasContractName]string
}

var EnvMainNet = Env{
	THQCodeHash:      "0x9e537bf5b8ec044ca3f53355e879f3fd8832217e4a9b41d9994cf0c547241a79",
	ContractArgs:     "0xc126635ece567c71c50f7482c5db80603852c306",
	ContractCodeHash: "0x00000000000000000000000000000000000000000000000000545950455f4944",
	MapContract: map[common.DasContractName]string{
		common.DasContractNameConfigCellType:        "0x3775c65aabe8b79980c4933dd2f4347fa5ef03611cef64328685618aa7535794",
		common.DasContractNameAccountCellType:       "0x96dc231bbbee6aa474076468640f9e0ad27cf13b1343716a7ce04b116ea18ba8",
		common.DasContractNameBalanceCellType:       "0xbdc8f42643ccad23e8df3d2e8dbdea9201812cd1b7f84c46e69b020529629822",
		common.DasContractNameDispatchCellType:      "0xda22fd296682488687a6035b5fc97c269b72d7de128034389bd03041b40309c0",
		common.DasContractNameIncomeCellType:        "0x108fba6a9b9f2898b4cdf11383ba2a6ed3da951b458c48e5f5de0353bbca2d46",
		common.DasContractNameAccountSaleCellType:   "0xb782d5f4e24603340997494871ba8f7d175e6920e63ead6b8137170b2e370469",
		common.DasContractNameAlwaysSuccess:         "0xca5016f232830f8a73e6827b5e1108aca68e7cf8baea4847ac40ef1da43c4c50",
		common.DasContractNameApplyRegisterCellType: "0xf18c3eab9fd28adbb793c38be9a59864989c1f739c22d2b6dc3f4284f047a69d",
		common.DasContractNamePreAccountCellType:    "0xf6574955079797010689a22cd172ce55b52d2c34d1e9bc6596e97babc2906f7e",
		common.DasContractNameProposalCellType:      "0xd7b779b1b30f86a77db6b292c9492906f2437b7d88a8a5994e722619bb1d41c8",
		common.DasContractNameReverseRecordCellType: "0x000f3e1a89d85d268ed6d36578d474ecf91d8809f4f696dd2e5b97fe67b84a2e",
		common.DASContractNameOfferCellType:         "0x3ffc0f8b0ce4bc09f700ca84355a092447d79fc5224a6fbd64af95af840af91b",
		common.DASContractNameSubAccountCellType:    "0x97b19f14184f24d55b1247596a5d7637f133c7bb7735f0ae962dc709c5fc1e2e",
		common.DASContractNameEip712LibCellType:     "",
	},
}

var EnvTestnet2 = Env{
	THQCodeHash:      "0x96248cdefb09eed910018a847cfb51ad044c2d7db650112931760e3ef34a7e9a",
	ContractArgs:     "0xbc502a34a430e3e167c82a24db6f9237b15ebf35",
	ContractCodeHash: "0x00000000000000000000000000000000000000000000000000545950455f4944",
	MapContract: map[common.DasContractName]string{
		common.DasContractNameConfigCellType:        "0x34363fad2018db0b3b6919c26870f302da74c3c4ef4456e5665b82c4118eda51",
		common.DasContractNameAccountCellType:       "0x6f0b8328b703617508d62d1f017b0d91fab2056de320a7b7faed4c777a356b7b",
		common.DasContractNameBalanceCellType:       "0x27560fe2daa6150b771621300d1d4ea127832b7b326f2d70eed63f5333b4a8a9",
		common.DasContractNameDispatchCellType:      "0xeedd10c7d8fee85c119daf2077fea9cf76b9a92ddca546f1f8e0031682e65aee",
		common.DasContractNameIncomeCellType:        "0xd7b9d8213671aec03f3a3ab95171e0e79481db2c084586b9ea99914c00ff3716",
		common.DasContractNameAccountSaleCellType:   "0xed5d7fc00a3f8605bfe3f6717747bb0ed529fa064c2b8ce56e9677a0c46c2c1c",
		common.DasContractNameAlwaysSuccess:         "0x7821c662b7efd50e7f6cf2b036efe53e07eccaf2e3447a2a470ee07ae455ab92",
		common.DasContractNameApplyRegisterCellType: "0xc78fa9066af1624e600ccfb21df9546f900b2afe5d7940d91aefc115653f90d9",
		common.DasContractNamePreAccountCellType:    "0xd3f7ad59632a2ebdc2fe9d41aa69708ed1069b074cd8b297b205f835335d3a6b",
		common.DasContractNameProposalCellType:      "0x03d0bb128bd10e666975d9a07c148f6abebe811f511e9574048b30600b065b9a",
		common.DasContractNameReverseRecordCellType: "0x334d7841eb156b8aa5abd7b09277e91e782d840140905496bb5bff0ea6ce9d75",
		common.DASContractNameOfferCellType:         "0x443b2d1b3b00ffab1a2287d84c47b2c31a11aad24b183d732c213a69e3d6d390",
		common.DASContractNameSubAccountCellType:    "0x63ca3e26cc69809f06735c6d9139ec2d84f2a277f13509a54060d6ee19423b5b",
		common.DASContractNameEip712LibCellType:     "0x16549cab7e92afb5f157141bc9da7781ce692a3144e47e2b8879a8d5a57b87c6",
	},
}

var EnvTestnet3 = Env{
	THQCodeHash:      "0x96248cdefb09eed910018a847cfb51ad044c2d7db650112931760e3ef34a7e9a",
	ContractArgs:     "0xbc502a34a430e3e167c82a24db6f9237b15ebf35",
	ContractCodeHash: "0x00000000000000000000000000000000000000000000000000545950455f4944",
	MapContract: map[common.DasContractName]string{
		common.DasContractNameConfigCellType:        "0xe9fa679290f63ba1debc74f04e96f299d1b61a03a24e7a1a51c7ccad416ec16a",
		common.DasContractNameAccountCellType:       "0x1c7d62798d412351ec9a4262aacc2c2837712b780d90e0720051ffa6d6304e32",
		common.DasContractNameBalanceCellType:       "0x301127f501c2620174b60945f305d78a4f5cdf67dd44009001265ca133d0088d",
		common.DasContractNameDispatchCellType:      "0xeedd10c7d8fee85c119daf2077fea9cf76b9a92ddca546f1f8e0031682e65aee",
		common.DasContractNameIncomeCellType:        "0x0af1f2332b1f742136ef0ffec849d415c4b5b8c4426de5d245f2e7ac0a9f5773",
		common.DasContractNameAccountSaleCellType:   "0xaf446d756fca5b51fad915a3f0526221527f1dd5f32898408ee86e73dc9d9814",
		common.DasContractNameAlwaysSuccess:         "0x7821c662b7efd50e7f6cf2b036efe53e07eccaf2e3447a2a470ee07ae455ab92",
		common.DasContractNameApplyRegisterCellType: "0xb12e454a3896b910030e669e8a81bf99c351cb58be8d1529e1c7023265db6084",
		common.DasContractNamePreAccountCellType:    "0x919e8036ded6670b923958db207527b5d1d3b95955e446e3475bfd7f949898c1",
		common.DasContractNameProposalCellType:      "0x2eecc2c70936609d14872acdedc9ba8b265c26240d49d66b6db027c5cfdb4256",
		common.DasContractNameReverseRecordCellType: "0x80963278cbdc61cdafd5250555984b71ad016798b8879adc0e6b1ee7e01b7912",
		common.DASContractNameOfferCellType:         "0xc69186c17e41fead0f87eb1f94829778e98a398be202655ac59fdb9567d05bae",
		common.DASContractNameSubAccountCellType:    "0x57498a2df0c0137146ced681fa1854599e404da5804c1a5ff45d954c3cc89bfd",
		common.DASContractNameEip712LibCellType:     "",
	},
}

func InitEnv(net common.DasNetType) Env {
	switch net {
	case common.DasNetTypeMainNet:
		return EnvMainNet
	case common.DasNetTypeTestnet2:
		return EnvTestnet2
	case common.DasNetTypeTestnet3:
		return EnvTestnet3
	default:
		return EnvMainNet
	}
}

func InitEnvOpt(net common.DasNetType, names ...common.DasContractName) Env {
	switch net {
	case common.DasNetTypeMainNet:
		return initEnvOpt(EnvMainNet, names...)
	case common.DasNetTypeTestnet2:
		return initEnvOpt(EnvTestnet2, names...)
	case common.DasNetTypeTestnet3:
		return initEnvOpt(EnvTestnet3, names...)
	default:
		return initEnvOpt(EnvMainNet, names...)
	}
}

func initEnvOpt(envNet Env, names ...common.DasContractName) Env {
	env := Env{
		THQCodeHash:      envNet.THQCodeHash,
		ContractArgs:     envNet.ContractArgs,
		ContractCodeHash: envNet.ContractCodeHash,
		MapContract:      map[common.DasContractName]string{},
	}
	for _, v := range names {
		if contract, ok := envNet.MapContract[v]; ok {
			env.MapContract[v] = contract
		}
	}
	return env
}
