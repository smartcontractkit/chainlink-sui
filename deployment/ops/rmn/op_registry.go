package rmn

var AllOperationsRMN = []any{
	*CurseChainOp,
	*UncurseChainOp,
	*CreateCurserCapOp,
	*McmsMintAndRegisterCurserCapOp,
	*McmsCreateCurserCapAndTransferOp,
	*McmsRegisterCurserCapOp,
	*McmsInitializeAllowedCurserCapsOp,
	*McmsRegisterCurserCapIdsOp,
	*McmsSetCurserCapAllowlistEnabledOp,
	*McmsDeregisterCurserCapIdsOp,
	*CurseWithCurserCapOp,
}
