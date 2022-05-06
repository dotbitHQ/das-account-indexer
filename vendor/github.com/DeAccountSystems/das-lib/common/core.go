package common

type DasNetType = int

const (
	DasNetTypeMainNet  DasNetType = 1
	DasNetTypeTestnet2 DasNetType = 2
	DasNetTypeTestnet3 DasNetType = 3
)

type DasAlgorithmId int

const (
	DasAlgorithmIdCkb       DasAlgorithmId = 0
	DasAlgorithmIdCkbMulti  DasAlgorithmId = 1
	DasAlgorithmIdCkbSingle DasAlgorithmId = 2
	DasAlgorithmIdEth       DasAlgorithmId = 3
	DasAlgorithmIdTron      DasAlgorithmId = 4
	DasAlgorithmIdEth712    DasAlgorithmId = 5
	DasAlgorithmIdEd25519   DasAlgorithmId = 6
)

func (d DasAlgorithmId) Bytes() []byte {
	return []byte{uint8(d)}
}

func (d DasAlgorithmId) ToSoScriptType() SoScriptType {
	switch d {
	case DasAlgorithmIdCkbSingle:
		return SoScriptTypeCkbSingle
	case DasAlgorithmIdCkbMulti:
		return SoScriptTypeCkbMulti
	case DasAlgorithmIdEth, DasAlgorithmIdEth712:
		return SoScriptTypeEth
	case DasAlgorithmIdTron:
		return SoScriptTypeTron
	case DasAlgorithmIdEd25519:
		return SoScriptTypeEd25519
	default:
		return SoScriptTypeCkbSingle
	}
}

func (d DasAlgorithmId) ToChainType() ChainType {
	switch d {
	case DasAlgorithmIdCkbSingle:
		return ChainTypeCkbSingle
	case DasAlgorithmIdCkbMulti:
		return ChainTypeCkbMulti
	case DasAlgorithmIdEth, DasAlgorithmIdEth712:
		return ChainTypeEth
	case DasAlgorithmIdTron:
		return ChainTypeTron
	case DasAlgorithmIdEd25519:
		return ChainTypeMixin
	default:
		return ChainTypeCkb
	}
}
