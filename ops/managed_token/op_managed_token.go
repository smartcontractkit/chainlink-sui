package managedtokenops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/block-vision/sui-go-sdk/models"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_managed_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/managed_token/managed_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

// MANAGED_TOKEN -- INITIALIZE
type ManagedTokenInitializeObjects struct {
	OwnerCapObjectId string
	StateObjectId    string
}

type ManagedTokenInitializeInput struct {
	ManagedTokenPackageId string
	CoinObjectTypeArg     string
	TreasuryCapObjectId   string
	DenyCapObjectId       string // Optional - can be empty
}

var initManagedTokenHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenInitializeInput) (output sui_ops.OpTxResult[ManagedTokenInitializeObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenInitializeObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	var tx *models.SuiTransactionBlockResponse
	if input.DenyCapObjectId != "" {
		tx, err = contract.InitializeWithDenyCap(
			b.GetContext(),
			opts,
			[]string{input.CoinObjectTypeArg},
			bind.Object{Id: input.TreasuryCapObjectId},
			bind.Object{Id: input.DenyCapObjectId},
		)
	} else {
		tx, err = contract.Initialize(
			b.GetContext(),
			opts,
			[]string{input.CoinObjectTypeArg},
			bind.Object{Id: input.TreasuryCapObjectId},
		)
	}

	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenInitializeObjects]{}, fmt.Errorf("failed to execute managed token initialization: %w", err)
	}

	// Find the object IDs from the transaction
	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "managed_token", "TokenState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[ManagedTokenInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: err1=%w, err2=%w", err1, err2)
	}

	return sui_ops.OpTxResult[ManagedTokenInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPackageId,
		Objects: ManagedTokenInitializeObjects{
			OwnerCapObjectId: obj1,
			StateObjectId:    obj2,
		},
	}, nil
}

var ManagedTokenInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Managed Token contract",
	initManagedTokenHandler,
)

// MANAGED_TOKEN -- configure_new_minter
type NoObjects struct {
}

type ManagedTokenConfigureNewMinterInput struct {
	ManagedTokenPackageId string
	CoinObjectTypeArg     string
	StateObjectId         string
	OwnerCapObjectId      string
	MinterAddress         string
	Allowance             uint64
	IsUnlimited           bool
}

var configureNewMinterHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenConfigureNewMinterInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.ConfigureNewMinter(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.MinterAddress,
		input.Allowance,
		input.IsUnlimited,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token configure new minter: %w", err)
	}

	b.Logger.Infow("ConfigureNewMinter on ManagedToken", "ManagedToken PackageId:", input.ManagedTokenPackageId, "Minter:", input.MinterAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageId: input.ManagedTokenPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenConfigureNewMinterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "configure_new_minter"),
	semver.MustParse("0.1.0"),
	"Configures a new minter for the CCIP Managed Token contract",
	configureNewMinterHandler,
)

// MANAGED_TOKEN -- increment_mint_allowance
type ManagedTokenIncrementMintAllowanceInput struct {
	ManagedTokenPackageId string
	CoinObjectTypeArg     string
	StateObjectId         string
	OwnerCapObjectId      string
	MintCapObjectId       string
	DenyListObjectId      string
	AllowanceIncrement    uint64
}

var incrementMintAllowanceHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenIncrementMintAllowanceInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.IncrementMintAllowance(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.MintCapObjectId},
		bind.Object{Id: input.DenyListObjectId},
		input.AllowanceIncrement,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token increment mint allowance: %w", err)
	}

	b.Logger.Infow("IncrementMintAllowance on ManagedToken", "ManagedToken PackageId:", input.ManagedTokenPackageId, "Increment:", input.AllowanceIncrement)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageId: input.ManagedTokenPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenIncrementMintAllowanceOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "increment_mint_allowance"),
	semver.MustParse("0.1.0"),
	"Increments mint allowance for the CCIP Managed Token contract",
	incrementMintAllowanceHandler,
)

// MANAGED_TOKEN -- mint
type ManagedTokenMintInput struct {
	ManagedTokenPackageId string
	CoinObjectTypeArg     string
	StateObjectId         string
	MintCapObjectId       string
	DenyListObjectId      string
	Amount                uint64
	Recipient             string
}

var mintHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenMintInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.Mint(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.MintCapObjectId},
		bind.Object{Id: input.DenyListObjectId},
		input.Amount,
		input.Recipient,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token mint: %w", err)
	}

	b.Logger.Infow("Mint on ManagedToken", "ManagedToken PackageId:", input.ManagedTokenPackageId, "Amount:", input.Amount, "Recipient:", input.Recipient)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageId: input.ManagedTokenPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenMintOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "mint"),
	semver.MustParse("0.1.0"),
	"Mints tokens in the CCIP Managed Token contract",
	mintHandler,
)

// MANAGED_TOKEN -- burn
type ManagedTokenBurnInput struct {
	ManagedTokenPackageId string
	CoinObjectTypeArg     string
	StateObjectId         string
	MintCapObjectId       string
	DenyListObjectId      string
	CoinObjectId          string
	FromAddress           string
}

var burnHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenBurnInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.Burn(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.MintCapObjectId},
		bind.Object{Id: input.DenyListObjectId},
		bind.Object{Id: input.CoinObjectId},
		input.FromAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token burn: %w", err)
	}

	b.Logger.Infow("Burn on ManagedToken", "ManagedToken PackageId:", input.ManagedTokenPackageId, "From:", input.FromAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageId: input.ManagedTokenPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenBurnOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "burn"),
	semver.MustParse("0.1.0"),
	"Burns tokens in the CCIP Managed Token contract",
	burnHandler,
)

// MANAGED_TOKEN -- pause
type ManagedTokenPauseInput struct {
	ManagedTokenPackageId string
	CoinObjectTypeArg     string
	StateObjectId         string
	OwnerCapObjectId      string
	DenyListObjectId      string
}

var pauseHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPauseInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.Pause(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.DenyListObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token pause: %w", err)
	}

	b.Logger.Infow("Pause on ManagedToken", "ManagedToken PackageId:", input.ManagedTokenPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageId: input.ManagedTokenPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenPauseOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "pause"),
	semver.MustParse("0.1.0"),
	"Pauses the CCIP Managed Token contract",
	pauseHandler,
)

// MANAGED_TOKEN -- unpause
type ManagedTokenUnpauseInput struct {
	ManagedTokenPackageId string
	CoinObjectTypeArg     string
	StateObjectId         string
	OwnerCapObjectId      string
	DenyListObjectId      string
}

var unpauseHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenUnpauseInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.Unpause(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.DenyListObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token unpause: %w", err)
	}

	b.Logger.Infow("Unpause on ManagedToken", "ManagedToken PackageId:", input.ManagedTokenPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageId: input.ManagedTokenPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenUnpauseOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "unpause"),
	semver.MustParse("0.1.0"),
	"Unpauses the CCIP Managed Token contract",
	unpauseHandler,
)
