package load

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
)

// estimateSuiFunding calculates how much SUI (in MIST) each wallet needs.
// Formula: (requiredCoinsPerWallet × splitAmount) + gasBudgetBuffer.
// requiredCoinsPerWallet follows PrepareSuiCoinPool coin math to avoid underfunding.
func estimateSuiFunding(cfg *config.LoadTestConfig, numWallets int, splitAmount uint64) uint64 {
	if numWallets <= 0 {
		numWallets = 1
	}
	msgPerWallet := cfg.MessageCount / numWallets
	if msgPerWallet < 1 {
		msgPerWallet = 1
	}
	requiredCoins := sui.CalculateRequiredSuiCoins(msgPerWallet)
	base := uint64(requiredCoins) * splitAmount //nolint:gosec // bounded by run config input

	// Add a fixed buffer for PTB gas and merge/split operations.
	gasBuffer := uint64(150_000_000) // 0.15 SUI

	return base + gasBuffer
}

// estimateEvmFunding calculates how much ETH (in wei) each wallet needs.
// Formula: msgPerWallet × (estimatedFee + gasCost) × 1.5 buffer
func estimateEvmFunding(cfg *config.LoadTestConfig, numWallets int) *big.Int {
	if numWallets <= 0 {
		numWallets = 1
	}
	msgPerWallet := cfg.MessageCount / numWallets
	if msgPerWallet < 1 {
		msgPerWallet = 1
	}

	// Conservative per-message cost: 0.01 ETH for CCIP fee + 0.001 ETH for gas.
	perMsgCost, ok := new(big.Int).SetString("11000000000000000", 10) // 0.011 ETH
	if !ok {
		// Fallback if the literal above somehow fails to parse.
		perMsgCost = big.NewInt(11_000_000_000_000_000)
	}

	base := new(big.Int).Mul(perMsgCost, big.NewInt(int64(msgPerWallet)))
	// 1.5x safety margin
	withMargin := new(big.Int).Mul(base, big.NewInt(15))
	withMargin.Div(withMargin, big.NewInt(10))

	return withMargin
}
