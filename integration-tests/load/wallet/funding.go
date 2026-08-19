package wallet

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"

	"github.com/block-vision/sui-go-sdk/transaction"
	ethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	suibind "github.com/smartcontractkit/chainlink-sui/bindings/bind"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const (
	suiCoinObjectType      = "0x2::coin::Coin<0x2::sui::SUI>"
	suiTransferGasBudget   = uint64(100_000_000)
	suiSweepSafetyReserve  = uint64(200_000_000)
	suiSweepMinSendBalance = uint64(1)
)

// FundSuiWallets sends `amount` MIST of native SUI from mainSigner to each wallet.
// Transfers are executed sequentially to avoid Sui shared-object contention.
func FundSuiWallets(
	ctx context.Context,
	ptbClient *client.PTBClient,
	mainSigner bindutils.SuiSigner,
	mainAddress string,
	wallets []*Wallet,
	amountPerWallet uint64,
) error {
	if len(wallets) == 0 {
		return nil
	}

	// Sui funding from a single signer should be sequential to avoid
	// mutable-object contention on the same gas coin set.
	for _, w := range wallets {
		if err := sendSuiTransfer(ctx, ptbClient, mainSigner, mainAddress, w.Address, amountPerWallet); err != nil {
			return fmt.Errorf("failed to fund Sui wallet %s: %w", w.Address, err)
		}
		slog.Info("Funded Sui wallet", "address", w.Address, "amount", amountPerWallet)
	}

	slog.Info("Funded Sui wallets", "count", len(wallets), "amountPerWallet", amountPerWallet)
	return nil
}

// sendSuiTransfer sends native SUI from the main signer to a recipient address.
func sendSuiTransfer(
	ctx context.Context,
	ptbClient *client.PTBClient,
	mainSigner bindutils.SuiSigner,
	mainAddress string,
	recipientAddress string,
	amount uint64,
) error {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, mainAddress, suiCoinObjectType)
	if err != nil {
		return fmt.Errorf("query SUI coins for sender: %w", err)
	}
	if len(coins) == 0 {
		return fmt.Errorf("no SUI coins found for sender %s", mainAddress)
	}

	mergedCoinID, err := sui.MergeSuiCoins(ctx, ptbClient, mainSigner, mainAddress, coins)
	if err != nil {
		return fmt.Errorf("merge SUI coins before transfer: %w", err)
	}

	// Re-query merged coin balance for the transfer preflight.
	refreshed, err := ptbClient.QueryCoinsByAddress(ctx, mainAddress, suiCoinObjectType)
	if err != nil {
		return fmt.Errorf("re-query SUI coins for sender: %w", err)
	}
	var mergedBalance uint64
	for _, c := range refreshed {
		if c.GetObjectId() == mergedCoinID {
			mergedBalance = c.GetBalance()
			break
		}
	}
	if mergedBalance == 0 {
		return fmt.Errorf("merged coin %s not found for sender %s", mergedCoinID, mainAddress)
	}
	if mergedBalance <= amount+suiTransferGasBudget {
		return fmt.Errorf("insufficient sender balance: coin=%s balance=%d need>%d", mergedCoinID, mergedBalance, amount+suiTransferGasBudget)
	}

	ptb := transaction.NewTransaction()
	gasArg := ptb.Gas()
	transferCoin := ptb.SplitCoins(gasArg, []transaction.Argument{ptb.Pure(amount)})
	ptb.TransferObjects([]transaction.Argument{transferCoin}, ptb.Pure(recipientAddress))

	gasBudget := suiTransferGasBudget
	callOpts := &suibind.CallOpts{
		Signer:           mainSigner,
		GasObject:        mergedCoinID,
		GasBudget:        &gasBudget,
		WaitForExecution: true,
	}

	resp, err := suibind.ExecutePTB(ctx, callOpts, ptbClient, ptb)
	if err != nil {
		return fmt.Errorf("execute Sui funding transfer PTB: %w", err)
	}

	slog.Info("Sui funding transfer sent",
		"from", mainAddress,
		"to", recipientAddress,
		"amount", amount,
		"txDigest", resp.Digest,
	)

	return nil
}

// FundEVMWallets sends `amount` wei of native ETH from mainAuth to each wallet.
// Transfers are executed sequentially to avoid nonce races.
func FundEVMWallets(
	ctx context.Context,
	ethClient *ethclient.Client,
	mainAuth *ethbind.TransactOpts,
	wallets []*Wallet,
	amountPerWallet *big.Int,
) error {
	if len(wallets) == 0 {
		return nil
	}

	// Funding from a single EVM key should be sequential unless explicit nonce
	// management is added. This avoids nonce races under concurrent sends.
	for _, w := range wallets {
		tx, err := sendEthTransfer(ctx, ethClient, mainAuth, common.HexToAddress(w.Address), amountPerWallet)
		if err != nil {
			return fmt.Errorf("failed to fund EVM wallet %s: %w", w.Address, err)
		}
		slog.Info("Funded EVM wallet", "address", w.Address, "tx", tx.Hash().Hex())
	}

	slog.Info("Funded EVM wallets", "count", len(wallets), "amountPerWallet", amountPerWallet)
	return nil
}

// sendEthTransfer sends native ETH from mainAuth to recipient.
func sendEthTransfer(
	ctx context.Context,
	ethClient *ethclient.Client,
	mainAuth *ethbind.TransactOpts,
	recipient common.Address,
	amount *big.Int,
) (*types.Transaction, error) {
	gasLimit := uint64(21000)
	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	nonce, err := ethClient.PendingNonceAt(ctx, mainAuth.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending nonce: %w", err)
	}

	// Build and sign the transaction directly; do not mutate mainAuth.
	tx := types.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, nil)
	signedTx, err := mainAuth.Signer(mainAuth.From, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transfer: %w", err)
	}

	if err := ethClient.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("failed to send transfer: %w", err)
	}

	// Wait for confirmation
	_, err = ethbind.WaitMined(ctx, ethClient, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed waiting for transfer confirmation: %w", err)
	}

	return signedTx, nil
}

// SweepSuiWallets sends remaining SUI balance from each wallet back to mainAddress.
// Best-effort: wallets with insufficient balance for gas are skipped.
func SweepSuiWallets(
	ctx context.Context,
	ptbClient *client.PTBClient,
	wallets []*Wallet,
	mainAddress string,
) error {
	for _, w := range wallets {
		if w.SuiSigner == nil {
			slog.Warn("Skipping Sui sweep for wallet without signer", "address", w.Address)
			continue
		}

		coins, err := ptbClient.QueryCoinsByAddress(ctx, w.Address, suiCoinObjectType)
		if err != nil {
			slog.Warn("Failed to query SUI coins for sweep", "address", w.Address, "error", err)
			continue
		}
		if len(coins) == 0 {
			continue
		}

		mergedCoinID, err := sui.MergeSuiCoins(ctx, ptbClient, w.SuiSigner, w.Address, coins)
		if err != nil {
			slog.Warn("Failed to merge SUI coins for sweep", "address", w.Address, "error", err)
			continue
		}

		refreshed, err := ptbClient.QueryCoinsByAddress(ctx, w.Address, suiCoinObjectType)
		if err != nil {
			slog.Warn("Failed to re-query SUI coins for sweep", "address", w.Address, "error", err)
			continue
		}

		var balance uint64
		for _, c := range refreshed {
			if c.GetObjectId() == mergedCoinID {
				balance = c.GetBalance()
				break
			}
		}
		if balance <= suiSweepSafetyReserve+suiSweepMinSendBalance {
			continue
		}

		amount := balance - suiSweepSafetyReserve
		ptb := transaction.NewTransaction()
		gasArg := ptb.Gas()
		sweepCoin := ptb.SplitCoins(gasArg, []transaction.Argument{ptb.Pure(amount)})
		ptb.TransferObjects([]transaction.Argument{sweepCoin}, ptb.Pure(mainAddress))

		gasBudget := suiTransferGasBudget
		callOpts := &suibind.CallOpts{
			Signer:           w.SuiSigner,
			GasObject:        mergedCoinID,
			GasBudget:        &gasBudget,
			WaitForExecution: true,
		}

		resp, err := suibind.ExecutePTB(ctx, callOpts, ptbClient, ptb)
		if err != nil {
			slog.Warn("Failed to sweep SUI wallet", "address", w.Address, "error", err)
			continue
		}

		slog.Info("Swept SUI wallet",
			"address", w.Address,
			"amount", amount,
			"recipient", mainAddress,
			"txDigest", resp.Digest,
		)
	}

	return nil
}

// SweepEVMWallets sends remaining ETH balance from each wallet back to mainAddress.
// Best-effort: wallets with insufficient balance for gas are skipped.
func SweepEVMWallets(
	ctx context.Context,
	ethClient *ethclient.Client,
	wallets []*Wallet,
	mainAddress common.Address,
) error {
	g := new(sync.WaitGroup)

	for _, w := range wallets {
		g.Add(1)
		go func(wallet *Wallet) {
			defer g.Done()

			balance, err := ethClient.BalanceAt(ctx, common.HexToAddress(wallet.Address), nil)
			if err != nil {
				slog.Warn("Failed to get balance for sweep", "address", wallet.Address, "error", err)
				return
			}

			gasPrice, err := ethClient.SuggestGasPrice(ctx)
			if err != nil {
				slog.Warn("Failed to get gas price for sweep", "address", wallet.Address, "error", err)
				return
			}

			gasCost := new(big.Int).Mul(gasPrice, big.NewInt(21000))
			if balance.Cmp(gasCost) <= 0 {
				slog.Warn("Skipping sweep: balance too low for gas", "address", wallet.Address, "balance", balance)
				return
			}

			amount := new(big.Int).Sub(balance, gasCost)
			tx, err := sendEthTransfer(ctx, ethClient, wallet.EVMTransactOpts, mainAddress, amount)
			if err != nil {
				slog.Warn("Failed to sweep EVM wallet", "address", wallet.Address, "error", err)
				return
			}
			slog.Info("Swept EVM wallet", "address", wallet.Address, "amount", amount, "tx", tx.Hash().Hex())
		}(w)
	}

	g.Wait()
	return nil
}
