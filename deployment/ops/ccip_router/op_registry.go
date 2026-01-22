package routerops

var AllOperationsRouter = []any{
	*TransferOwnershipOp,
	*AcceptOwnershipOp,
	*ExecuteOwnershipTransferOp,
	*SetOnRampsOp,
}
