package main

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/mockcontracts/deploy"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Hello, World!")

	// Build and publish contracts
	gasBudget := int(200000000)

	lggr, err := logger.NewWith(func(cfg *zap.Config) {
		cfg.Level.SetLevel(zap.InfoLevel)
	})
	if err != nil {
		lggr.Errorw("Failed to create logger", "error", err)
		return
	}

	contractPath, err := deploy.BuildSetup(lggr, "contracts/test/sources/")
	if err != nil {
		lggr.Errorw("Failed to build contract", "error", err)
		return
	}

	lggr.Infow("Published contracts", "contractPath", contractPath)

	onrampRouterPackageId, onrampRouterPublishOutput, err := deploy.PublishContract(lggr, "test", contractPath, &gasBudget)
	if err != nil {
		lggr.Errorw("Failed to publish contract", "error", err)
		return
	}

	lggr.Infow("Published contracts", "onrampRouterPackageId", onrampRouterPackageId, "onrampRouterPublishOutput", onrampRouterPublishOutput)
}
