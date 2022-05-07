package common

import (
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type DasContractName string

const (
	DasContractNameAlwaysSuccess  DasContractName = "always-success"
	DasContractNameConfigCellType DasContractName = "config-cell-type"

	DasContractNameDispatchCellType DasContractName = "das-lock"
	DasContractNameAccountCellType  DasContractName = "account-cell-type"
	DasContractNameBalanceCellType  DasContractName = "balance-cell-type"

	DasContractNameApplyRegisterCellType DasContractName = "apply-register-cell-type"
	DasContractNamePreAccountCellType    DasContractName = "pre-account-cell-type"
	DasContractNameProposalCellType      DasContractName = "proposal-cell-type"

	DasContractNameIncomeCellType      DasContractName = "income-cell-type"
	DasContractNameAccountSaleCellType DasContractName = "account-sale-cell-type"

	DasContractNameReverseRecordCellType DasContractName = "reverse-record-cell-type"
	DASContractNameOfferCellType         DasContractName = "offer-cell-type"
	DASContractNameSubAccountCellType    DasContractName = "sub-account-cell-type"
	DASContractNameEip712LibCellType     DasContractName = "eip712-lib"
)

// script to type id
func ScriptToTypeId(script *types.Script) types.Hash {
	bys, _ := script.Serialize()
	bysRet, _ := blake2b.Blake256(bys)
	return types.BytesToHash(bysRet)
}
