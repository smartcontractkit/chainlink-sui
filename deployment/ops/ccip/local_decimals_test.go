package ccipops

import (
	"context"
	"errors"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
)

type stubCoinMetadataClient struct {
	decimals int
	err      error
}

func (s stubCoinMetadataClient) GetCoinMetadata(context.Context, string) (models.CoinMetadataResponse, error) {
	if s.err != nil {
		return models.CoinMetadataResponse{}, s.err
	}
	return models.CoinMetadataResponse{Decimals: s.decimals}, nil
}

func TestResolveLocalDecimalsFetchesFromChain(t *testing.T) {
	t.Parallel()

	decimals, err := ResolveLocalDecimals(
		context.Background(),
		stubCoinMetadataClient{decimals: 9},
		"abc::link::LINK",
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, byte(9), decimals)
}

func TestResolveLocalDecimalsAddsHexPrefix(t *testing.T) {
	t.Parallel()

	var capturedType string
	client := stubCoinMetadataClient{decimals: 6}
	_, err := ResolveLocalDecimals(
		context.Background(),
		coinMetadataClientFunc(func(_ context.Context, coinType string) (models.CoinMetadataResponse, error) {
			capturedType = coinType
			return client.GetCoinMetadata(context.Background(), coinType)
		}),
		"abc::usdc::USDC",
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, "0xabc::usdc::USDC", capturedType)
}

func TestResolveLocalDecimalsRejectsMismatch(t *testing.T) {
	t.Parallel()

	supplied := byte(8)
	_, err := ResolveLocalDecimals(
		context.Background(),
		stubCoinMetadataClient{decimals: 9},
		"0xabc::link::LINK",
		&supplied,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "local decimals mismatch")
}

func TestResolveLocalDecimalsRequiresTokenType(t *testing.T) {
	t.Parallel()

	_, err := ResolveLocalDecimals(context.Background(), stubCoinMetadataClient{decimals: 9}, "", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "token type is required")
}

func TestValidateRegisterPoolLocalDecimals(t *testing.T) {
	t.Parallel()

	err := ValidateRegisterPoolLocalDecimals(
		context.Background(),
		stubCoinMetadataClient{decimals: 9},
		"0xabc::link::LINK",
		9,
	)
	require.NoError(t, err)

	err = ValidateRegisterPoolLocalDecimals(
		context.Background(),
		stubCoinMetadataClient{decimals: 9},
		"0xabc::link::LINK",
		8,
	)
	require.Error(t, err)
}

type coinMetadataClientFunc func(context.Context, string) (models.CoinMetadataResponse, error)

func (f coinMetadataClientFunc) GetCoinMetadata(ctx context.Context, coinType string) (models.CoinMetadataResponse, error) {
	return f(ctx, coinType)
}

func TestResolveLocalDecimalsPropagatesClientError(t *testing.T) {
	t.Parallel()

	_, err := ResolveLocalDecimals(
		context.Background(),
		stubCoinMetadataClient{err: errors.New("rpc unavailable")},
		"0xabc::link::LINK",
		nil,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "rpc unavailable")
}
