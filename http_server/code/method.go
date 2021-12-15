package code

type JsonRpcMethod = string

const (
	MethodVersion     JsonRpcMethod = "das_version"
	MethodIndexerInfo JsonRpcMethod = "das_indexerInfo"

	MethodSearchAccount  JsonRpcMethod = "das_searchAccount"
	MethodAddressAccount JsonRpcMethod = "das_getAddressAccount"

	MethodAccountInfo    JsonRpcMethod = "das_accountInfo"
	MethodAccountRecords JsonRpcMethod = "das_accountRecords"
	MethodReverseRecord  JsonRpcMethod = "das_reverseRecord"
)
