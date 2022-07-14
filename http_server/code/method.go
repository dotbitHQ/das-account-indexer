package code

type JsonRpcMethod = string

const (
	MethodVersion    JsonRpcMethod = "das_version"
	MethodServerInfo JsonRpcMethod = "das_serverInfo"

	MethodSearchAccount  JsonRpcMethod = "das_searchAccount"
	MethodAddressAccount JsonRpcMethod = "das_getAddressAccount"

	MethodAccountInfo      JsonRpcMethod = "das_accountInfo"
	MethodAccountList      JsonRpcMethod = "das_accountList"
	MethodAccountRecords   JsonRpcMethod = "das_accountRecords"
	MethodAccountRecordsV2 JsonRpcMethod = "das_accountRecordsV2"
	MethodReverseRecord    JsonRpcMethod = "das_reverseRecord"
)
