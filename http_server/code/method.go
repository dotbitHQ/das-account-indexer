package code

type JsonRpcMethod = string

const (
	MethodVersion    JsonRpcMethod = "das_version"
	MethodDidNumber  JsonRpcMethod = "das_didNumber"
	MethodServerInfo JsonRpcMethod = "das_serverInfo"

	MethodSearchAccount  JsonRpcMethod = "das_searchAccount"
	MethodAddressAccount JsonRpcMethod = "das_getAddressAccount"

	MethodAccountInfo           JsonRpcMethod = "das_accountInfo"
	MethodAccountList           JsonRpcMethod = "das_accountList"
	MethodAccountRecords        JsonRpcMethod = "das_accountRecords"
	MethodBatchAccountRecords   JsonRpcMethod = "das_batchAccountRecords"
	MethodAccountRecordsV2      JsonRpcMethod = "das_accountRecordsV2"
	MethodReverseRecord         JsonRpcMethod = "das_reverseRecord"
	MethodBatchReverseRecord    JsonRpcMethod = "das_batchReverseRecord"
	MethodBatchRegisterInfo     JsonRpcMethod = "das_batchRegisterInfo"
	MethodAccountReverseAddress JsonRpcMethod = "das_accountReverseAddress"

	MethodSubAccountList   JsonRpcMethod = "das_subAccountList"
	MethodSubAccountVerify JsonRpcMethod = "das_subAccountVerify"
	MethodDidCellList      JsonRpcMethod = "das_didCellList"
)
