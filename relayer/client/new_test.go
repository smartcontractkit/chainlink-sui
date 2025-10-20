package client_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/test-go/testify/require"
)

// TestRateLimit_OffRampRead exercises PTBClient.ReadFunction with the same
// parameters used by the chainreader for offramp.latest_config_details.
// It stresses concurrency + rate limiter under real testnet RPC conditions.
func TestRateLimit_OffRampRead_MultiNode(t *testing.T) {
	log := logger.Test(t)

	rpcURL := ""

	// Simulation parameters
	numNodes := 4               // number of independent nodes
	maxConcurrent := int64(100) // per-node semaphore size (currently used)
	numRequestsPerNode := 50    // total requests per node

	keystore := testutils.NewTestKeystore(t)
	signer, _ := testutils.GetAccountAndKeyFromSui(keystore)

	t.Logf("Simulating %d nodes × %d requests (each node concurrency=%d) → shared RPC %s",
		numNodes, numRequestsPerNode, maxConcurrent, rpcURL)

	// Function call definitions (same as in real logs)
	offrampPkg := "0x0bd389e3b489efccc1e12b4c611a7f3a771602b73633d083b8c36c028b6eb618"
	feeQuoterPkg := "0x8b2c8b7f4527ce8643e209fcbfeb9ab8bdb4dc6d42b81acebde18183b2116215"
	argTokenPool := "0x98048b2884b75b0454ae853ef727047e94005ae527aba35a8831f529025130fe"
	argOfframp := "0x4ed2a35187613d3b0770d4a3c3ba59ea54c084570be99be4ee4d7198e11b92e3"

	type readCall struct {
		packageID string
		module    string
		function  string
		args      []any
		argTypes  []string
	}

	readFns := []readCall{
		{offrampPkg, "offramp", "latest_config_details", []any{argOfframp, uint8(0)}, []string{"object_id", "u8"}},
		{offrampPkg, "offramp", "get_static_config", []any{argTokenPool, argOfframp}, []string{"object_id", "object_id"}},
		{offrampPkg, "offramp", "get_dynamic_config", []any{argTokenPool, argOfframp}, []string{"object_id", "object_id"}},
		{offrampPkg, "offramp", "get_source_chain_config", []any{argTokenPool, argOfframp, uint64(16015286601757825753)}, []string{"object_id", "object_id", "u64"}},
		{feeQuoterPkg, "fee_quoter", "get_static_config", []any{argTokenPool}, []string{"object_id"}},
	}

	// Simulate multiple independent nodes
	var globalWG sync.WaitGroup
	for nodeID := 0; nodeID < numNodes; nodeID++ {
		globalWG.Add(1)
		go func(nid int) {
			defer globalWG.Done()

			ptb, err := client.NewPTBClient(
				log,
				rpcURL,
				nil,
				120*time.Second,
				keystore,
				maxConcurrent,
				"WaitForLocalExecution",
			)
			require.NoError(t, err, "node %d: failed to init PTBClient", nid)

			var wg sync.WaitGroup
			for i := 0; i < numRequestsPerNode; i++ {
				wg.Add(1)
				go func(reqID int) {
					defer wg.Done()

					call := readFns[reqID%len(readFns)]
					ctx := context.Background()
					start := time.Now()

					_, err := ptb.ReadFunction(ctx, signer, call.packageID, call.module, call.function, call.args, call.argTypes)
					elapsed := time.Since(start)

					if err != nil {
						t.Logf("[Node %02d | Req %03d] %s::%s::%s failed after %v: %v",
							nid, reqID, call.packageID[:10], call.module, call.function, elapsed, err)
					} else {
						t.Logf("[Node %02d | Req %03d]  %s::%s::%s succeeded in %v",
							nid, reqID, call.packageID[:10], call.module, call.function, elapsed)
					}
				}(i)
			}

			wg.Wait()
			t.Logf("Node %02d finished all requests", nid)
		}(nodeID)
	}

	globalWG.Wait()
	t.Log("All simulated nodes completed requests")
}
