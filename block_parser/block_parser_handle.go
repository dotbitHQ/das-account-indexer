package block_parser

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type FuncTransactionHandleReq struct {
	Tx             *types.Transaction
	TxHash         string
	BlockNumber    uint64
	BlockTimestamp uint64
	Action         common.DasAction
}

type FuncTransactionHandleResp struct {
	Err error
}

type FuncTransactionHandle func(*FuncTransactionHandleReq) FuncTransactionHandleResp

func (b *BlockParser) registerTransactionHandle() {
	b.MapTransactionHandle = make(map[string]FuncTransactionHandle)
	b.MapTransactionHandle[common.DasActionConfig] = b.ActionConfigCell
	b.MapTransactionHandle[common.DasActionStartAccountSale] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionCancelAccountSale] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionBuyAccount] = b.ActionUpdateAccountInfo

	b.MapTransactionHandle[common.DasActionConfirmProposal] = b.ActionConfirmProposal
	b.MapTransactionHandle[common.DasActionEditRecords] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionEditManager] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionRenewAccount] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionTransferAccount] = b.ActionUpdateAccountInfo

	b.MapTransactionHandle[common.DasActionAcceptOffer] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionLockAccountForCrossChain] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionUnlockAccountForCrossChain] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionForceRecoverAccountStatus] = b.ActionUpdateAccountInfo
	b.MapTransactionHandle[common.DasActionRecycleExpiredAccount] = b.ActionRecycleExpiredAccount

	b.MapTransactionHandle[common.DasActionDeclareReverseRecord] = b.ActionDeclareReverseRecord
	b.MapTransactionHandle[common.DasActionRedeclareReverseRecord] = b.ActionRedeclareReverseRecord
	b.MapTransactionHandle[common.DasActionRetractReverseRecord] = b.ActionRetractReverseRecord
	b.MapTransactionHandle[common.DasActionUpdateReverseRecordRoot] = b.ActionReverseRecordRoot

	b.MapTransactionHandle[common.DasActionEnableSubAccount] = b.ActionEnableSubAccount
	b.MapTransactionHandle[common.DasActionCreateSubAccount] = b.ActionCreateSubAccount
	b.MapTransactionHandle[common.DasActionEditSubAccount] = b.ActionEditSubAccount
	b.MapTransactionHandle[common.DasActionUpdateSubAccount] = b.ActionUpdateSubAccount
	//b.MapTransactionHandle[common.DasActionLockSubAccountForCrossChain] = b.ActionUpdateSubAccountInfo
	//b.MapTransactionHandle[common.DasActionUnlockSubAccountForCrossChain] = b.ActionUpdateSubAccountInfo
	b.MapTransactionHandle[common.DasActionConfigSubAccountCustomScript] = b.ActionConfigSubAccountCreatingScript
	b.MapTransactionHandle[common.DasActionConfigSubAccount] = b.ActionConfigSubAccount
	b.MapTransactionHandle[common.DasActionCreateApproval] = b.ActionCreateApproval
	b.MapTransactionHandle[common.DasActionDelayApproval] = b.ActionDelayApproval
	b.MapTransactionHandle[common.DasActionRevokeApproval] = b.ActionRevokeApproval
	b.MapTransactionHandle[common.DasActionFulfillApproval] = b.ActionFulfillApproval
	b.MapTransactionHandle[common.DasActionBidExpiredAccountAuction] = b.ActionBidExpiredAccountAuction
}

func isCurrentVersionTx(tx *types.Transaction, name common.DasContractName) (bool, error) {
	contract, err := core.GetDasContractInfo(name)
	if err != nil {
		return false, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	isCV := false
	for _, v := range tx.Outputs {
		if v.Type == nil {
			continue
		}
		if contract.IsSameTypeId(v.Type.CodeHash) {
			isCV = true
			break
		}
	}
	return isCV, nil
}
