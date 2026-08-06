package evm

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	router "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
)

// erc20ABI is the minimal ABI for ERC-20 approve, allowance, and balanceOf.
var erc20ABI = func() abi.ABI {
	a, err := abi.JSON(strings.NewReader(`[
		{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
		{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},
		{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}
	]`))
	if err != nil {
		panic(fmt.Sprintf("failed to parse ERC-20 ABI: %v", err))
	}
	return a
}()

// ApproveRouterForTokens approves the Router to spend the total token amount once at the start of the run.
//
// totalAmount = tokenAmount * messageCount. If the existing allowance is already >= totalAmount,
// the approval transaction is skipped.
func ApproveRouterForTokens(
	ctx context.Context,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	tokenAddress common.Address,
	routerAddress common.Address,
	totalAmount *big.Int,
) error {
	if totalAmount == nil || totalAmount.Sign() <= 0 {
		return fmt.Errorf("total approval amount must be > 0")
	}

	tokenContract := bind.NewBoundContract(tokenAddress, erc20ABI, client, client, client)

	var allowance *big.Int
	results := []any{&allowance}
	if err := tokenContract.Call(&bind.CallOpts{Context: ctx}, &results, "allowance", auth.From, routerAddress); err != nil {
		return fmt.Errorf("failed to query token allowance: %w", err)
	}

	if allowance != nil && allowance.Cmp(totalAmount) >= 0 {
		slog.Info("Existing ERC-20 allowance is sufficient, skipping approval",
			"token", tokenAddress.Hex(),
			"allowance", allowance.String(),
			"required", totalAmount.String(),
		)
		return nil
	}

	tx, err := tokenContract.Transact(auth, "approve", routerAddress, totalAmount)
	if err != nil {
		return fmt.Errorf("failed to approve Router for token %s: %w", tokenAddress.Hex(), err)
	}

	slog.Info("Approved Router for ERC-20 tokens",
		"token", tokenAddress.Hex(),
		"router", routerAddress.Hex(),
		"amount", totalAmount.String(),
		"txHash", tx.Hash().Hex(),
	)

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for token approval tx %s: %w", tx.Hash().Hex(), err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("token approval tx %s failed", tx.Hash().Hex())
	}

	return nil
}

// CheckTokenBalance returns the ERC-20 balance of the signer.
func CheckTokenBalance(ctx context.Context, client *ethclient.Client, tokenAddress common.Address, owner common.Address) (*big.Int, error) {
	tokenContract := bind.NewBoundContract(tokenAddress, erc20ABI, client, client, client)

	var balance *big.Int
	results := []any{&balance}
	if err := tokenContract.Call(&bind.CallOpts{Context: ctx}, &results, "balanceOf", owner); err != nil {
		return nil, fmt.Errorf("failed to query token balance: %w", err)
	}
	if balance == nil {
		return big.NewInt(0), nil
	}
	return balance, nil
}

// SendTokenMessage sends a CCIP token transfer from EVM to Sui.
//
// tokenAmount is the per-message token amount in base units.
// receiver is the 32-byte Sui receiver package ID (zero for token-only EOA transfers).
// data is the message payload (empty for token-only).
// extraArgs is the serialized SuiExtraArgsV1.
//
// Returns the message ID and transaction hash.
func SendTokenMessage(
	ctx context.Context,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	routerAddress common.Address,
	destChainSelector uint64,
	receiver []byte,
	data []byte,
	tokenAddress common.Address,
	tokenAmount *big.Int,
	extraArgs []byte,
) (messageID string, txHash string, err error) {
	routerContract, err := router.NewRouter(routerAddress, client)
	if err != nil {
		return "", "", fmt.Errorf("failed to instantiate Router: %w", err)
	}

	var tokenAmounts []router.ClientEVMTokenAmount
	if tokenAddress != (common.Address{}) && tokenAmount != nil && tokenAmount.Sign() > 0 {
		tokenAmounts = []router.ClientEVMTokenAmount{{Token: tokenAddress, Amount: tokenAmount}}
	}

	msg := router.ClientEVM2AnyMessage{
		Receiver:     receiver,
		Data:         data,
		TokenAmounts: tokenAmounts,
		FeeToken:     common.Address{}, // zero address = native ETH
		ExtraArgs:    extraArgs,
	}

	fee, err := routerContract.GetFee(&bind.CallOpts{Context: ctx}, destChainSelector, msg)
	if err != nil {
		return "", "", fmt.Errorf("failed to get fee: %w", err)
	}

	feeWithBuffer := new(big.Int).Add(fee, new(big.Int).Div(fee, big.NewInt(5)))
	auth.Value = feeWithBuffer

	slog.Info("Sending EVM→Sui token message",
		"fee", fee.String(),
		"feeWithBuffer", feeWithBuffer.String(),
		"token", tokenAddress.Hex(),
		"tokenAmount", tokenAmount.String(),
	)

	tx, err := routerContract.CcipSend(auth, destChainSelector, msg)
	if err != nil {
		return "", "", fmt.Errorf("failed to send CCIP token message: %w", err)
	}

	txHash = tx.Hash().Hex()

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		slog.Warn("Could not wait for transaction to be mined, using tx hash as message ID",
			"txHash", txHash,
			"error", err,
		)
		return txHash, txHash, nil
	}

	messageID, _, err = ExtractMessageIDFromReceipt(receipt, destChainSelector)
	if err != nil {
		slog.Warn("Could not extract message ID from receipt, using tx hash",
			"txHash", txHash,
			"error", err,
		)
		messageID = txHash
	}

	return messageID, txHash, nil
}
