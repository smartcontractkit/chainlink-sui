package sui

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	managedtokenpool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// AppendManagedTokenPoolLockOrBurn appends a managed_token_pool.lock_or_burn call to the PTB.
//
// Parameters:
//   - tokenPoolPkgID: latest token pool package ID
//   - tokenCoinID: the token coin object to lock/burn (must be distinct from gas/fee coins)
//   - tokenTransferParams: PTB argument result from create_token_transfer_params
//   - destChainSelector: destination chain selector
//   - ccipObjectRefID: shared CCIPObjectRef object ID
//   - poolStateObjectID: shared pool state object ID
//   - denyListObjectID: shared deny list object ID
//   - tokenStateObjectID: shared token state object ID
//
// The lock_or_burn function signature for managed_token_pool is:
// lock_or_burn<T>(ref: &CCIPObjectRef, token_transfer_params: &mut onramp_sh::TokenTransferParams,
//   c: Coin<T>, remote_chain_selector: u64, clock: &Clock, deny_list: &DenyList,
//   token_state: &mut TokenState<T>, state: &mut ManagedTokenPoolState<T>)
func AppendManagedTokenPoolLockOrBurn(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	ptb *transaction.Transaction,
	callOpts *bind.CallOpts,
	tokenPoolPkgID string,
	coinType string,
	tokenCoinID string,
	tokenTransferParams transaction.Argument,
	destChainSelector uint64,
	ccipObjectRefID string,
	poolStateObjectID string,
	denyListObjectID string,
	tokenStateObjectID string,
) error {
	poolContract, err := managedtokenpool.NewManagedTokenPool(tokenPoolPkgID, ptbClient)
	if err != nil {
		return fmt.Errorf("failed to create managed token pool contract: %w", err)
	}

	// Use LockOrBurnWithArgs to pass the PTB argument (transaction.Argument) for tokenTransferParams.
	// LockOrBurn expects bind.Object values; LockOrBurnWithArgs accepts any values including transaction.Argument.
	encoded, err := poolContract.Encoder().LockOrBurnWithArgs(
		[]string{coinType},
		bind.Object{Id: ccipObjectRefID},
		tokenTransferParams,
		bind.Object{Id: tokenCoinID},
		destChainSelector,
		bind.Object{Id: "0x6"},
		bind.Object{Id: denyListObjectID},
		bind.Object{Id: tokenStateObjectID},
		bind.Object{Id: poolStateObjectID},
	)
	if err != nil {
		return fmt.Errorf("failed to encode lock_or_burn: %w", err)
	}

	_, err = poolContract.Bound().AppendPTB(ctx, callOpts, ptb, encoded)
	if err != nil {
		return fmt.Errorf("failed to append lock_or_burn to PTB: %w", err)
	}

	return nil
}
