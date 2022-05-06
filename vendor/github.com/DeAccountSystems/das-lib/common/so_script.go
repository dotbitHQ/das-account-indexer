package common

type SoScriptType string

const (
	SoScriptTypeEd25519   SoScriptType = "Ed25519"
	SoScriptTypeEth       SoScriptType = "Eth"
	SoScriptTypeTron      SoScriptType = "Tron"
	SoScriptTypeCkbMulti  SoScriptType = "CkbMulti"
	SoScriptTypeCkbSingle SoScriptType = "CkbSingle"
)
