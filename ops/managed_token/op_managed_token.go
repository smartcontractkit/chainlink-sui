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
	OwnerCapObjectID string
	StateObjectID    string
}

type ManagedTokenInitializeInput struct {
	ManagedTokenPackageID string
	CoinObjectTypeArg     string
	TreasuryCapObjectID   string
	DenyCapObjectID       string // Optional - can be empty
}

var initManagedTokenHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenInitializeInput) (output sui_ops.OpTxResult[ManagedTokenInitializeObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenInitializeObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	var tx *models.SuiTransactionBlockResponse
	if input.DenyCapObjectID != "" {
		tx, err = contract.InitializeWithDenyCap(
			b.GetContext(),
			opts,
			[]string{input.CoinObjectTypeArg},
			bind.Object{Id: input.TreasuryCapObjectID},
			bind.Object{Id: input.DenyCapObjectID},
		)
	} else {
		tx, err = contract.Initialize(
			b.GetContext(),
			opts,
			[]string{input.CoinObjectTypeArg},
			bind.Object{Id: input.TreasuryCapObjectID},
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
		PackageID: input.ManagedTokenPackageID,
		Objects: ManagedTokenInitializeObjects{
			OwnerCapObjectID: obj1,
			StateObjectID:    obj2,
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
	ManagedTokenPackageID string
	CoinObjectTypeArg     string
	StateObjectID         string
	OwnerCapObjectID      string
	MinterAddress         string
	Allowance             uint64
	IsUnlimited           bool
}

var configureNewMinterHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenConfigureNewMinterInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.ConfigureNewMinter(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.OwnerCapObjectID},
		input.MinterAddress,
		input.Allowance,
		input.IsUnlimited,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token configure new minter: %w", err)
	}

	b.Logger.Infow("ConfigureNewMinter on ManagedToken", "ManagedToken PackageID:", input.ManagedTokenPackageID, "Minter:", input.MinterAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageID: input.ManagedTokenPackageID,
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
	ManagedTokenPackageID string
	CoinObjectTypeArg     string
	StateObjectID         string
	OwnerCapObjectID      string
	MintCapObjectID       string
	DenyListObjectID      string
	AllowanceIncrement    uint64
}

var incrementMintAllowanceHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenIncrementMintAllowanceInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.IncrementMintAllowance(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.OwnerCapObjectID},
		bind.Object{Id: input.MintCapObjectID},
		bind.Object{Id: input.DenyListObjectID},
		input.AllowanceIncrement,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token increment mint allowance: %w", err)
	}

	b.Logger.Infow("IncrementMintAllowance on ManagedToken", "ManagedToken PackageID:", input.ManagedTokenPackageID, "Increment:", input.AllowanceIncrement)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageID: input.ManagedTokenPackageID,
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
	ManagedTokenPackageID string
	CoinObjectTypeArg     string
	StateObjectID         string
	MintCapObjectID       string
	DenyListObjectID      string
	Amount                uint64
	Recipient             string
}

var mintHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenMintInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.Mint(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.MintCapObjectID},
		bind.Object{Id: input.DenyListObjectID},
		input.Amount,
		input.Recipient,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token mint: %w", err)
	}

	b.Logger.Infow("Mint on ManagedToken", "ManagedToken PackageID:", input.ManagedTokenPackageID, "Amount:", input.Amount, "Recipient:", input.Recipient)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageID: input.ManagedTokenPackageID,
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
	ManagedTokenPackageID string
	CoinObjectTypeArg     string
	StateObjectID         string
	MintCapObjectID       string
	DenyListObjectID      string
	CoinObjectID          string
	FromAddress           string
}

var burnHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenBurnInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.Burn(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.MintCapObjectID},
		bind.Object{Id: input.DenyListObjectID},
		bind.Object{Id: input.CoinObjectID},
		input.FromAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token burn: %w", err)
	}

	b.Logger.Infow("Burn on ManagedToken", "ManagedToken PackageID:", input.ManagedTokenPackageID, "From:", input.FromAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageID: input.ManagedTokenPackageID,
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
	ManagedTokenPackageID string
	CoinObjectTypeArg     string
	StateObjectID         string
	OwnerCapObjectID      string
	DenyListObjectID      string
}

var pauseHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPauseInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.Pause(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.OwnerCapObjectID},
		bind.Object{Id: input.DenyListObjectID},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token pause: %w", err)
	}

	b.Logger.Infow("Pause on ManagedToken", "ManagedToken PackageID:", input.ManagedTokenPackageID)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageID: input.ManagedTokenPackageID,
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
	ManagedTokenPackageID string
	CoinObjectTypeArg     string
	StateObjectID         string
	OwnerCapObjectID      string
	DenyListObjectID      string
}

var unpauseHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenUnpauseInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	_, err = contract.Unpause(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.OwnerCapObjectID},
		bind.Object{Id: input.DenyListObjectID},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token unpause: %w", err)
	}

	b.Logger.Infow("Unpause on ManagedToken", "ManagedToken PackageID:", input.ManagedTokenPackageID)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    "", // tx.Digest when available
		PackageID: input.ManagedTokenPackageID,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenUnpauseOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "unpause"),
	semver.MustParse("0.1.0"),
	"Unpauses the CCIP Managed Token contract",
	unpauseHandler,
)
