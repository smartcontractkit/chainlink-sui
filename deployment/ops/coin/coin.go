package coin

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/block-vision/sui-go-sdk/models"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// GetCoinMetadataOutput is a clone of models.CoinMetadataResponse from sui-go-sdk
type GetCoinMetadataOutput struct {
	Id          string
	Decimals    int
	Name        string
	Symbol      string
	Description string
	IconUrl     string
}

var GetCoinSymbolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "get-coin-symbol"),
	semver.MustParse("0.1.0"),
	"Retrieves the symbol for a SUI coin",
	getCoinSymbolHandler,
)

var getCoinSymbolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, coinObjectTypeArg string) (GetCoinMetadataOutput, error) {
	ctx := context.Background()

	rsp, err := deps.Client.SuiXGetCoinMetadata(ctx, models.SuiXGetCoinMetadataRequest{
		CoinType: coinObjectTypeArg,
	})
	if err != nil {
		return GetCoinMetadataOutput{}, err
	}

	return GetCoinMetadataOutput{
		Id:          rsp.Id,
		Decimals:    rsp.Decimals,
		Name:        rsp.Name,
		Symbol:      rsp.Symbol,
		Description: rsp.Description,
		IconUrl:     rsp.IconUrl,
	}, nil
}
