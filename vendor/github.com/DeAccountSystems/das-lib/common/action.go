package common

type DasAction = string

const (
	DasActionConfig                 DasAction = "config"
	DasActionApplyRegister          DasAction = "apply_register"
	DasActionRefundApply            DasAction = "refund_apply"
	DasActionPreRegister            DasAction = "pre_register"
	DasActionPropose                DasAction = "propose"
	DasActionTransferAccount        DasAction = "transfer_account"
	DasActionRenewAccount           DasAction = "renew_account"
	DasActionExtendPropose          DasAction = "extend_proposal"
	DasActionConfirmProposal        DasAction = "confirm_proposal"
	DasActionRecycleProposal        DasAction = "recycle_proposal"
	DasActionWithdrawFromWallet     DasAction = "withdraw_from_wallet"
	DasActionEditManager            DasAction = "edit_manager"
	DasActionEditRecords            DasAction = "edit_records"
	DasActionStartAccountSale       DasAction = "start_account_sale"
	DasActionEditAccountSale        DasAction = "edit_account_sale"
	DasActionCancelAccountSale      DasAction = "cancel_account_sale"
	DasActionBuyAccount             DasAction = "buy_account"
	DasActionSellAccount            DasAction = "sell_account"
	DasActionCreateIncome           DasAction = "create_income"
	DasActionConsolidateIncome      DasAction = "consolidate_income"
	DasActionDeclareReverseRecord   DasAction = "declare_reverse_record"
	DasActionRedeclareReverseRecord DasAction = "redeclare_reverse_record"
	DasActionRetractReverseRecord   DasAction = "retract_reverse_record"
	DasActionTransfer               DasAction = "transfer"

	DasActionMakeOffer   DasAction = "make_offer"
	DasActionEditOffer   DasAction = "edit_offer"
	DasActionCancelOffer DasAction = "cancel_offer"
	DasActionAcceptOffer DasAction = "accept_offer"

	DasActionEnableSubAccount  DasAction = "enable_sub_account"
	DasActionCreateSubAccount  DasAction = "create_sub_account"
	DasActionEditSubAccount    DasAction = "edit_sub_account"
	DasActionRenewSubAccount   DasAction = "renew_sub_account"
	DasActionRecycleSubAccount DasAction = "recycle_sub_account"
)

const (
	ParamOwner   = "0x00"
	ParamManager = "0x01"
)
