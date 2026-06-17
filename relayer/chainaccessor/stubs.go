package chainaccessor

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// The methods below are required by ccipocr3.ChainAccessor but are not yet
// implemented natively for Sui. They return ErrNotImplemented. See the package
// plan for the intended follow-up (USDC/CCTP support and price-feed reads).

// MessagesByTokenID (USDCMessageReader) is not implemented for Sui.
func (a *SuiAccessor) MessagesByTokenID(
	ctx context.Context,
	source, dest ccipocr3.ChainSelector,
	tokens map[ccipocr3.MessageTokenID]ccipocr3.RampTokenAmount,
) (map[ccipocr3.MessageTokenID]ccipocr3.Bytes, error) {
	return nil, ErrNotImplemented
}

// GetFeedPricesUSD (PriceReader) is not implemented for Sui.
func (a *SuiAccessor) GetFeedPricesUSD(
	ctx context.Context,
	tokens []ccipocr3.UnknownEncodedAddress,
	tokenInfo map[ccipocr3.UnknownEncodedAddress]ccipocr3.TokenInfo,
) (ccipocr3.TokenPriceMap, error) {
	return nil, ErrNotImplemented
}

// GetFeeQuoterTokenUpdates (PriceReader) is not implemented for Sui.
func (a *SuiAccessor) GetFeeQuoterTokenUpdates(
	ctx context.Context,
	tokensBytes []ccipocr3.UnknownAddress,
) (map[ccipocr3.UnknownEncodedAddress]ccipocr3.TimestampedUnixBig, error) {
	return nil, ErrNotImplemented
}
