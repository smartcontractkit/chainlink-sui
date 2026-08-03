package evm

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
)

// SuiExtraArgsV1Tag is the 4-byte tag prefix for SuiExtraArgsV1.
const SuiExtraArgsV1Tag = "0x21ea4ca9"

// suiExtraArgsABI is the parsed ABI for encoding SuiExtraArgsV1.
var suiExtraArgsABI = func() abi.ABI {
	a, err := abi.JSON(strings.NewReader(message_hasher.MessageHasherABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse MessageHasher ABI: %v", err))
	}
	return a
}()

// SerializeClientSUIExtraArgsV1 encodes SuiExtraArgsV1 for EVM→Sui messages.
//
// Format: tag (4 bytes) + ABI-encoded ClientSuiExtraArgsV1
//
// Uses the message_hasher.MessageHasherABI for ABI encoding.
func SerializeClientSUIExtraArgsV1(data message_hasher.ClientSuiExtraArgsV1) ([]byte, error) {
	tag := hexutil.MustDecode(SuiExtraArgsV1Tag)
	v, err := suiExtraArgsABI.Methods["encodeSUIExtraArgsV1"].Inputs.Pack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode SuiExtraArgsV1: %w", err)
	}
	return append(tag, v...), nil
}
