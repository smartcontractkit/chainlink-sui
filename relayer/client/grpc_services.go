package client

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/common/grpcconn"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
)

// Close closes the underlying gRPC connection. JSON-RPC clients are stateless and need no cleanup.
func (c *PTBClient) Close() error {
	if c.grpcClient == nil {
		return nil
	}

	c.grpcServicesMu.Lock()
	defer c.grpcServicesMu.Unlock()

	c.ledgerService = nil
	c.stateService = nil
	c.txExecService = nil
	c.movePkgService = nil
	c.subscriptionService = nil

	err := c.grpcClient.Close()
	c.grpcClient = nil
	return err
}

// HealthCheckGrpc verifies the gRPC connection by calling LedgerService.GetServiceInfo.
func (c *PTBClient) HealthCheckGrpc(ctx context.Context) (chainID string, err error) {
	service, err := c.getLedgerService(ctx)
	if err != nil {
		return "", err
	}

	resp, err := service.GetServiceInfo(ctx, &suirpcv2.GetServiceInfoRequest{})
	if err != nil {
		return "", fmt.Errorf("GetServiceInfo failed: %w", err)
	}

	if resp.ChainId != nil {
		chainID = *resp.ChainId
	}

	return chainID, nil
}

// VerifyGrpcServices initializes all gRPC service stubs to verify connectivity.
func (c *PTBClient) VerifyGrpcServices(ctx context.Context) error {
	if _, err := c.getLedgerService(ctx); err != nil {
		return fmt.Errorf("LedgerService: %w", err)
	}
	if _, err := c.getStateService(ctx); err != nil {
		return fmt.Errorf("StateService: %w", err)
	}
	if _, err := c.getTransactionExecutionService(ctx); err != nil {
		return fmt.Errorf("TransactionExecutionService: %w", err)
	}
	if _, err := c.getMovePackageService(ctx); err != nil {
		return fmt.Errorf("MovePackageService: %w", err)
	}

	return nil
}

func (c *PTBClient) requireGrpcClient() (*grpcconn.SuiGrpcClient, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client is not configured")
	}
	return c.grpcClient, nil
}

// getLedgerService returns the LedgerService client, lazily initializing on first use.
func (c *PTBClient) getLedgerService(ctx context.Context) (suirpcv2.LedgerServiceClient, error) {
	c.grpcServicesMu.Lock()
	defer c.grpcServicesMu.Unlock()

	if c.ledgerService != nil {
		return c.ledgerService, nil
	}

	grpcClient, err := c.requireGrpcClient()
	if err != nil {
		return nil, err
	}

	service, err := grpcClient.LedgerService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger service: %w", err)
	}

	c.ledgerService = service
	return service, nil
}

// getStateService returns the StateService client, lazily initializing on first use.
func (c *PTBClient) getStateService(ctx context.Context) (suirpcv2.StateServiceClient, error) {
	c.grpcServicesMu.Lock()
	defer c.grpcServicesMu.Unlock()

	if c.stateService != nil {
		return c.stateService, nil
	}

	grpcClient, err := c.requireGrpcClient()
	if err != nil {
		return nil, err
	}

	service, err := grpcClient.StateService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get state service: %w", err)
	}

	c.stateService = service
	return service, nil
}

// getTransactionExecutionService returns the TransactionExecutionService client.
func (c *PTBClient) getTransactionExecutionService(ctx context.Context) (suirpcv2.TransactionExecutionServiceClient, error) {
	c.grpcServicesMu.Lock()
	defer c.grpcServicesMu.Unlock()

	if c.txExecService != nil {
		return c.txExecService, nil
	}

	grpcClient, err := c.requireGrpcClient()
	if err != nil {
		return nil, err
	}

	service, err := grpcClient.TransactionExecutionService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction execution service: %w", err)
	}

	c.txExecService = service
	return service, nil
}

// getMovePackageService returns the MovePackageService client.
func (c *PTBClient) getMovePackageService(ctx context.Context) (suirpcv2.MovePackageServiceClient, error) {
	c.grpcServicesMu.Lock()
	defer c.grpcServicesMu.Unlock()

	if c.movePkgService != nil {
		return c.movePkgService, nil
	}

	grpcClient, err := c.requireGrpcClient()
	if err != nil {
		return nil, err
	}

	service, err := grpcClient.MovePackageService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get move package service: %w", err)
	}

	c.movePkgService = service
	return service, nil
}

func (c *PTBClient) getSubscriptionService(ctx context.Context) (suirpcv2.SubscriptionServiceClient, error) {
	c.grpcServicesMu.Lock()
	defer c.grpcServicesMu.Unlock()

	if c.subscriptionService != nil {
		return c.subscriptionService, nil
	}

	grpcClient, err := c.requireGrpcClient()
	if err != nil {
		return nil, err
	}

	service, err := grpcClient.SubscriptionService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription service: %w", err)
	}

	c.subscriptionService = service
	return service, nil
}

// resetGrpcServices clears cached service stubs so the next call re-initializes them.
// Used after connection errors when the underlying client reconnects.
func (c *PTBClient) resetGrpcServices() {
	c.grpcServicesMu.Lock()
	defer c.grpcServicesMu.Unlock()

	c.ledgerService = nil
	c.stateService = nil
	c.txExecService = nil
	c.movePkgService = nil
	c.subscriptionService = nil
}
