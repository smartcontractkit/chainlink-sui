// / DO NOT EDIT - this will be removed
package expander

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// PTBExpander defines a generic interface for expanding PTB (Programmable Transaction Block) commands
// for various operations in the Sui blockchain.
//
// This interface uses Go generics to provide type-safe PTB expansion operations
// for CCIP (Cross-Chain Interoperability Protocol) on Sui.
// The expansion process involves translating high-level requests
// into low-level Sui Move function calls within a PTB structure.
type PTBExpander[T any, R any] interface {
	// Expand performs PTB expansion with generic types.
	Expand(
		ctx context.Context,
		lggr logger.Logger,
		args T,
		signerPublicKey []byte,
	) (R, error)
}
