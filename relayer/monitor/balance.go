package monitor

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	aptosBalanceMonitor "github.com/smartcontractkit/chainlink-aptos/relayer/monitor"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

const SuiDecimals = 9
const SuiDecimalsDenominator = 1_000_000_000
const BalanceCheckTimeout = 30 * time.Second

// BalanceMonitorOpts contains the options for creating a new Sui account balance monitor.
type BalanceMonitorOpts struct {
	ChainInfo aptosBalanceMonitor.ChainInfo

	Config    aptosBalanceMonitor.GenericBalanceConfig
	Logger    logger.Logger
	Keystore  core.Keystore
	NewClient func() (client.SuiPTBClient, error)
}

// NewBalanceMonitor returns a balance monitoring services.Service which reports balance of all Keystore accounts.
func NewBalanceMonitor(opts BalanceMonitorOpts) (services.Service, error) {
	return aptosBalanceMonitor.NewGenericBalanceMonitor(aptosBalanceMonitor.GenericBalanceMonitorOpts{
		ChainInfo:           opts.ChainInfo,
		ChainNativeCurrency: "SUI",
		Config:              opts.Config,
		Logger:              opts.Logger,
		Keystore:            opts.Keystore,
		NewGenericBalanceClient: func() (aptosBalanceMonitor.GenericBalanceClient, error) {
			ptbClient, err := opts.NewClient()
			if err != nil {
				return nil, fmt.Errorf("failed to get new client: %w", err)
			}

			return balanceClient{client: ptbClient}, nil
		},
		KeyToAccountMapper: func(ctx context.Context, pubKey string) (string, error) {
			// We need to convert the Sui public key to an account address
			hexBytes, err := hex.DecodeString(pubKey)
			if err != nil {
				return "", fmt.Errorf("failed to decode public key: %w", err)
			}

			return client.GetAddressFromPublicKey(hexBytes)
		},
	})
}

// Sui balance reader client implementation
type balanceClient struct {
	client client.SuiPTBClient
}

// GetAccountBalance returns the account balance in SUI
func (c balanceClient) GetAccountBalance(address string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), BalanceCheckTimeout)
	defer cancel()

	// Get the account balance
	balance, err := c.client.GetSUIBalance(ctx, address)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return mistToSui(balance.Uint64()), nil
}

// Convert MIST to SUI as 1/10^9 SUI
// Source: https://docs.sui.io/references/framework/sui/sui
func mistToSui(mist uint64) float64 {
	return float64(mist) / SuiDecimalsDenominator
}
