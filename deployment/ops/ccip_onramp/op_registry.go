package onrampops

var AllOperationsOnramp = []any{
	*TransferOwnershipOnRampOp,
	*AcceptOwnershipOnRampOp,
	*ExecuteOwnershipTransferToMcmsOnRampOp,
	*ApplyAllowListUpdateOp,
	*ApplyDestChainConfigUpdateOp,
	*SetDynamicConfigOp,
	*WithdrawFeeTokensOp,
}
