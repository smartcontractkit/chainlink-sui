package load

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
)

// estimateSuiFunding calculates how much SUI (in MIST) each wallet needs.
// Formula:
//   msgPerWallet × splitAmount   (CCIP fee reserve)
// + msgPerWallet × 5_000_000     (network gas for send PTBs)
// + 400_000_000                  (merge/split setup + safety buffer)
func estimateSuiFunding(cfg *config.LoadTestConfig, numWallets int, splitAmount uint64) uint64 {
	if numWallets <= 0 {
		numWallets = 1
	}
	msgPerWallet := cfg.MessageCount / numWallets
	if msgPerWallet < 1 {
		msgPerWallet = 1
	}
	ccipFees := uint64(msgPerWallet) * splitAmount //nolint:gosec
	networkGas := uint64(msgPerWallet) * 5_000_000 //nolint:gosec
	setupBuffer := uint64(400_000_000)

	return ccipFees + networkGas + setupBuffer
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
