package chainaccessor

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// callView resolves the latest package ID for the bound contract and performs a
// devInspect Move view call, returning the JSON-decoded return values.
func (a *SuiAccessor) callView(
	ctx context.Context,
	contractName, module, function string,
	args []any,
	argTypes, typeArgs []string,
) ([]any, error) {
	packageAddr, err := a.bindings.getPackageAddress(contractName)
	if err != nil {
		return nil, err
	}

	// Pin the latest (post-upgrade) package ID for the module being called.
	latestPackageID, err := a.client.GetLatestPackageId(ctx, packageAddr, module)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve latest package id for %s::%s: %w", contractName, module, err)
	}

	values, err := a.client.ReadFunction(ctx, latestPackageID, module, function, args, argTypes, typeArgs)
	if err != nil {
		return nil, fmt.Errorf("view call %s::%s failed: %w", module, function, err)
	}
	return values, nil
}

// --- SourceAccessor (state reads) ---

// GetExpectedNextSequenceNumber returns the expected next sequence number for
// messages being sent to the provided destination.
//
// Method(onramp::get_expected_next_sequence_number) on OnRamp.
func (a *SuiAccessor) GetExpectedNextSequenceNumber(ctx context.Context, dest ccipocr3.ChainSelector) (ccipocr3.SeqNum, error) {
	stateID, err := a.bindings.getStateObjectID("onramp")
	if err != nil {
		return 0, err
	}

	values, err := a.callView(ctx, ContractNameOnRamp, "onramp", "get_expected_next_sequence_number",
		[]any{stateID, uint64(dest)},
		[]string{"object_id", "u64"}, nil)
	if err != nil {
		return 0, err
	}

	v, err := firstReturn(values)
	if err != nil {
		return 0, err
	}
	n, err := asUint64(v)
	if err != nil {
		return 0, fmt.Errorf("decoding expected next sequence number: %w", err)
	}
	return ccipocr3.SeqNum(n), nil
}

// GetTokenPriceUSD looks up a token price in USD for the provided address.
//
// Method(fee_quoter::get_token_price) on FeeQuoter.
func (a *SuiAccessor) GetTokenPriceUSD(ctx context.Context, address ccipocr3.UnknownAddress) (ccipocr3.TimestampedUnixBig, error) {
	refID, err := a.bindings.getCCIPObjectRefID()
	if err != nil {
		return ccipocr3.TimestampedUnixBig{}, err
	}

	token := suiAddressFromBytes(address)
	values, err := a.callView(ctx, ContractNameFeeQuoter, "fee_quoter", "get_token_price",
		[]any{refID, token},
		[]string{"object_id", "address"}, nil)
	if err != nil {
		return ccipocr3.TimestampedUnixBig{}, err
	}

	return decodeTimestampedPrice(values)
}

// GetFeeQuoterDestChainConfig returns the fee quoter destination chain config.
//
// Method(fee_quoter::get_dest_chain_config) on FeeQuoter.
func (a *SuiAccessor) GetFeeQuoterDestChainConfig(ctx context.Context, dest ccipocr3.ChainSelector) (ccipocr3.FeeQuoterDestChainConfig, error) {
	refID, err := a.bindings.getCCIPObjectRefID()
	if err != nil {
		return ccipocr3.FeeQuoterDestChainConfig{}, err
	}

	values, err := a.callView(ctx, ContractNameFeeQuoter, "fee_quoter", "get_dest_chain_config",
		[]any{refID, uint64(dest)},
		[]string{"object_id", "u64"}, nil)
	if err != nil {
		return ccipocr3.FeeQuoterDestChainConfig{}, err
	}

	v, err := firstReturn(values)
	if err != nil {
		return ccipocr3.FeeQuoterDestChainConfig{}, err
	}
	m, err := asMap(v)
	if err != nil {
		return ccipocr3.FeeQuoterDestChainConfig{}, fmt.Errorf("decoding dest chain config: %w", err)
	}
	return decodeFeeQuoterDestChainConfig(m)
}

// --- DestinationAccessor (state reads) ---

// NextSeqNum reads the source chain config from the OffRamp to get the next
// expected sequence number (MinSeqNr) for each given source chain.
//
// Method(offramp::get_source_chain_config) on OffRamp.
func (a *SuiAccessor) NextSeqNum(ctx context.Context, sources []ccipocr3.ChainSelector) (map[ccipocr3.ChainSelector]ccipocr3.SeqNum, error) {
	refID, err := a.bindings.getCCIPObjectRefID()
	if err != nil {
		return nil, err
	}
	stateID, err := a.bindings.getStateObjectID("offramp")
	if err != nil {
		return nil, err
	}

	out := make(map[ccipocr3.ChainSelector]ccipocr3.SeqNum, len(sources))
	for _, src := range sources {
		values, err := a.callView(ctx, ContractNameOffRamp, "offramp", "get_source_chain_config",
			[]any{refID, stateID, uint64(src)},
			[]string{"object_id", "object_id", "u64"}, nil)
		if err != nil {
			return nil, fmt.Errorf("reading source chain config for %d: %w", src, err)
		}

		v, err := firstReturn(values)
		if err != nil {
			return nil, err
		}
		m, err := asMap(v)
		if err != nil {
			return nil, fmt.Errorf("decoding source chain config for %d: %w", src, err)
		}
		minSeqNr, err := uint64Field(m, "min_seq_nr")
		if err != nil {
			return nil, fmt.Errorf("decoding min_seq_nr for %d: %w", src, err)
		}
		out[src] = ccipocr3.SeqNum(minSeqNr)
	}
	return out, nil
}

// GetLatestPriceSeqNr returns the latest price (OCR) sequence number for the
// destination chain.
//
// Method(offramp::get_latest_price_sequence_number) on OffRamp.
func (a *SuiAccessor) GetLatestPriceSeqNr(ctx context.Context) (ccipocr3.SeqNum, error) {
	stateID, err := a.bindings.getStateObjectID("offramp")
	if err != nil {
		return 0, err
	}

	values, err := a.callView(ctx, ContractNameOffRamp, "offramp", "get_latest_price_sequence_number",
		[]any{stateID},
		[]string{"object_id"}, nil)
	if err != nil {
		return 0, err
	}

	v, err := firstReturn(values)
	if err != nil {
		return 0, err
	}
	n, err := asUint64(v)
	if err != nil {
		return 0, fmt.Errorf("decoding latest price sequence number: %w", err)
	}
	return ccipocr3.SeqNum(n), nil
}

// GetChainFeePriceUpdate returns the latest chain (gas) fee price update for the
// provided selectors.
//
// Method(fee_quoter::get_dest_chain_gas_price) on FeeQuoter.
func (a *SuiAccessor) GetChainFeePriceUpdate(ctx context.Context, selectors []ccipocr3.ChainSelector) (map[ccipocr3.ChainSelector]ccipocr3.TimestampedUnixBig, error) {
	refID, err := a.bindings.getCCIPObjectRefID()
	if err != nil {
		return nil, err
	}

	out := make(map[ccipocr3.ChainSelector]ccipocr3.TimestampedUnixBig, len(selectors))
	for _, sel := range selectors {
		values, err := a.callView(ctx, ContractNameFeeQuoter, "fee_quoter", "get_dest_chain_gas_price",
			[]any{refID, uint64(sel)},
			[]string{"object_id", "u64"}, nil)
		if err != nil {
			return nil, fmt.Errorf("reading dest chain gas price for %d: %w", sel, err)
		}
		price, err := decodeTimestampedPrice(values)
		if err != nil {
			return nil, fmt.Errorf("decoding dest chain gas price for %d: %w", sel, err)
		}
		out[sel] = price
	}
	return out, nil
}

// Nonces is not implemented: the Sui nonce_manager exposes get_outbound_nonce
// but no inbound-nonce view, which is what DestinationAccessor.Nonces
// (GetInboundNonce) requires. Implementing this requires a contract-side view
// (see plan risks).
//
// Method(GetInboundNonce) on NonceManager.
func (a *SuiAccessor) Nonces(
	ctx context.Context,
	addresses map[ccipocr3.ChainSelector][]ccipocr3.UnknownEncodedAddress,
) (map[ccipocr3.ChainSelector]map[string]uint64, error) {
	return nil, fmt.Errorf("%w: Sui nonce_manager has no inbound-nonce view", ErrNotImplemented)
}

// decodeTimestampedPrice decodes a Move TimestampedPrice { value: u256, timestamp: u64 }
// return value into a ccipocr3.TimestampedUnixBig.
func decodeTimestampedPrice(values []any) (ccipocr3.TimestampedUnixBig, error) {
	v, err := firstReturn(values)
	if err != nil {
		return ccipocr3.TimestampedUnixBig{}, err
	}
	m, err := asMap(v)
	if err != nil {
		return ccipocr3.TimestampedUnixBig{}, fmt.Errorf("decoding timestamped price: %w", err)
	}

	rawValue, err := field(m, "value")
	if err != nil {
		return ccipocr3.TimestampedUnixBig{}, err
	}
	value, err := asBigInt(rawValue)
	if err != nil {
		return ccipocr3.TimestampedUnixBig{}, fmt.Errorf("decoding price value: %w", err)
	}

	ts, err := uint64Field(m, "timestamp")
	if err != nil {
		return ccipocr3.TimestampedUnixBig{}, err
	}

	return ccipocr3.TimestampedUnixBig{
		Value:     value,
		Timestamp: uint32(ts),
	}, nil
}

// decodeFeeQuoterDestChainConfig maps the Sui fee_quoter DestChainConfig struct
// onto the family-agnostic ccipocr3.FeeQuoterDestChainConfig. The Sui struct
// fields map 1:1 onto the common type (see fee_quoter.move).
func decodeFeeQuoterDestChainConfig(m map[string]any) (ccipocr3.FeeQuoterDestChainConfig, error) {
	u := func(key string) (uint64, error) { return uint64Field(m, key) }

	isEnabled, err := boolField(m, "is_enabled")
	if err != nil {
		return ccipocr3.FeeQuoterDestChainConfig{}, err
	}
	enforceOutOfOrder, err := boolField(m, "enforce_out_of_order")
	if err != nil {
		return ccipocr3.FeeQuoterDestChainConfig{}, err
	}

	cfg := ccipocr3.FeeQuoterDestChainConfig{IsEnabled: isEnabled, EnforceOutOfOrder: enforceOutOfOrder}

	type binding struct {
		key string
		set func(uint64)
	}
	for _, b := range []binding{
		{"max_number_of_tokens_per_msg", func(v uint64) { cfg.MaxNumberOfTokensPerMsg = uint16(v) }},
		{"max_data_bytes", func(v uint64) { cfg.MaxDataBytes = uint32(v) }},
		{"max_per_msg_gas_limit", func(v uint64) { cfg.MaxPerMsgGasLimit = uint32(v) }},
		{"dest_gas_overhead", func(v uint64) { cfg.DestGasOverhead = uint32(v) }},
		{"dest_gas_per_payload_byte_base", func(v uint64) { cfg.DestGasPerPayloadByteBase = uint32(v) }},
		{"dest_gas_per_payload_byte_high", func(v uint64) { cfg.DestGasPerPayloadByteHigh = uint32(v) }},
		{"dest_gas_per_payload_byte_threshold", func(v uint64) { cfg.DestGasPerPayloadByteThreshold = uint32(v) }},
		{"dest_data_availability_overhead_gas", func(v uint64) { cfg.DestDataAvailabilityOverheadGas = uint32(v) }},
		{"dest_gas_per_data_availability_byte", func(v uint64) { cfg.DestGasPerDataAvailabilityByte = uint16(v) }},
		{"dest_data_availability_multiplier_bps", func(v uint64) { cfg.DestDataAvailabilityMultiplierBps = uint16(v) }},
		{"default_token_fee_usd_cents", func(v uint64) { cfg.DefaultTokenFeeUSDCents = uint16(v) }},
		{"default_token_dest_gas_overhead", func(v uint64) { cfg.DefaultTokenDestGasOverhead = uint32(v) }},
		{"default_tx_gas_limit", func(v uint64) { cfg.DefaultTxGasLimit = uint32(v) }},
		{"gas_multiplier_wei_per_eth", func(v uint64) { cfg.GasMultiplierWeiPerEth = v }},
		{"network_fee_usd_cents", func(v uint64) { cfg.NetworkFeeUSDCents = uint32(v) }},
		{"gas_price_staleness_threshold", func(v uint64) { cfg.GasPriceStalenessThreshold = uint32(v) }},
	} {
		val, err := u(b.key)
		if err != nil {
			return ccipocr3.FeeQuoterDestChainConfig{}, err
		}
		b.set(val)
	}

	return cfg, nil
}

func boolField(m map[string]any, key string) (bool, error) {
	v, err := field(m, key)
	if err != nil {
		return false, err
	}
	return asBool(v)
}
