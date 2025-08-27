package testutils

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
)

// CommitReport represents the commit report structure from offramp.move
// Matches the Move struct: public struct CommitReport has store, drop, copy
type CommitReport struct {
	PriceUpdates         PriceUpdates
	BlessedMerkleRoots   []MerkleRoot
	UnblessedMerkleRoots []MerkleRoot
	RMNSignatures        [][]byte
}

// PriceUpdates represents the price updates structure from offramp.move
// Matches the Move struct: public struct PriceUpdates has store, drop, copy
type PriceUpdates struct {
	TokenPriceUpdates []TokenPriceUpdate
	GasPriceUpdates   []GasPriceUpdate
}

// TokenPriceUpdate represents a token price update from offramp.move
// Matches the Move struct: public struct TokenPriceUpdate has store, drop, copy
type TokenPriceUpdate struct {
	SourceToken []byte   // address in Move becomes []byte for 32-byte Sui address
	UsdPerToken *big.Int // u256 in Move becomes *big.Int in Go
}

// GasPriceUpdate represents a gas price update from offramp.move
// Matches the Move struct: public struct GasPriceUpdate has store, drop, copy
type GasPriceUpdate struct {
	DestChainSelector uint64   // u64 in Move
	UsdPerUnitGas     *big.Int // u256 in Move becomes *big.Int in Go
}

// MerkleRoot represents a merkle root structure from offramp.move
// Matches the Move struct: public struct MerkleRoot has store, drop, copy
type MerkleRoot struct {
	SourceChainSelector uint64 // u64 in Move
	OnRampAddress       []byte // vector<u8> in Move
	MinSeqNr            uint64 // u64 in Move
	MaxSeqNr            uint64 // u64 in Move
	MerkleRoot          []byte // vector<u8> in Move (32 bytes)
}

// Helper function to create a CommitReport with proper byte slice initialization
func NewCommitReport() *CommitReport {
	return &CommitReport{
		PriceUpdates: PriceUpdates{
			TokenPriceUpdates: make([]TokenPriceUpdate, 0),
			GasPriceUpdates:   make([]GasPriceUpdate, 0),
		},
		BlessedMerkleRoots:   make([]MerkleRoot, 0),
		UnblessedMerkleRoots: make([]MerkleRoot, 0),
		RMNSignatures:        make([][]byte, 0),
	}
}

// Helper function to create a TokenPriceUpdate with proper address formatting
func NewTokenPriceUpdate(sourceToken string, usdPerToken *big.Int) TokenPriceUpdate {
	// Convert hex string to 32-byte address
	tokenBytes := make([]byte, DefaultByteSize)
	if len(sourceToken) > 2 && sourceToken[:2] == "0x" {
		// Remove 0x prefix and decode hex
		if decoded, err := hex.DecodeString(sourceToken[2:]); err == nil {
			copy(tokenBytes[32-len(decoded):], decoded) // Right-pad to 32 bytes
		}
	}

	return TokenPriceUpdate{
		SourceToken: tokenBytes,
		UsdPerToken: usdPerToken,
	}
}

// Helper function to create a GasPriceUpdate
func NewGasPriceUpdate(destChainSelector uint64, usdPerUnitGas *big.Int) GasPriceUpdate {
	return GasPriceUpdate{
		DestChainSelector: destChainSelector,
		UsdPerUnitGas:     usdPerUnitGas,
	}
}

// Helper function to create a MerkleRoot with proper byte slice initialization
func NewMerkleRoot(sourceChainSelector uint64, onRampAddress []byte, minSeqNr, maxSeqNr uint64, merkleRoot []byte) MerkleRoot {
	// Ensure merkle root is exactly 32 bytes
	root := make([]byte, DefaultByteSize)
	if len(merkleRoot) > 0 {
		copy(root, merkleRoot)
	}

	return MerkleRoot{
		SourceChainSelector: sourceChainSelector,
		OnRampAddress:       onRampAddress,
		MinSeqNr:            minSeqNr,
		MaxSeqNr:            maxSeqNr,
		MerkleRoot:          root,
	}
}

func GetCommitReport(
	onRampAddress []byte,
	merkleRoot []byte,
	tokenID string,
	price *big.Int,
	chainSelector uint64,
	seqNumStart uint64,
	seqNumEnd uint64,
	gasPrice *big.Int,
	r []byte,
	s []byte,
) CommitReport {
	// Ensure merkle root is 32 bytes
	merkleRootBytes := make([]byte, DefaultByteSize)
	copy(merkleRootBytes, merkleRoot)

	// Ensure signatures are 32 bytes and combine R and S into 64-byte signatures
	rBytes := make([]byte, DefaultByteSize)
	sBytes := make([]byte, DefaultByteSize)
	copy(rBytes, r)
	copy(sBytes, s)

	// Combine R and S into a single 64-byte signature
	signature := make([]byte, DefaultByteSize*SignatureComponents)
	copy(signature[:32], rBytes)
	copy(signature[32:], sBytes)

	return CommitReport{
		BlessedMerkleRoots: []MerkleRoot{},
		UnblessedMerkleRoots: []MerkleRoot{
			NewMerkleRoot(
				chainSelector,
				onRampAddress,
				seqNumStart,
				seqNumEnd,
				merkleRootBytes,
			),
		},
		PriceUpdates: PriceUpdates{
			TokenPriceUpdates: []TokenPriceUpdate{
				NewTokenPriceUpdate(tokenID, price),
			},
			GasPriceUpdates: []GasPriceUpdate{
				NewGasPriceUpdate(chainSelector, gasPrice),
			},
		},
		RMNSignatures: [][]byte{},
	}
}

// SerializeCommitReport serializes a CommitReport using BCS format to match the Move contract's expected deserialization format.
// The Move contract expects the following order:
// 1. TokenPriceUpdates (vector of TokenPriceUpdate)
// 2. GasPriceUpdates (vector of GasPriceUpdate)
// 3. BlessedMerkleRoots (vector of MerkleRoot)
// 4. UnblessedMerkleRoots (vector of MerkleRoot)
// 5. RMNSignatures (vector of fixed 64-byte vectors)
func SerializeCommitReport(report CommitReport) ([]byte, error) {
	s := &bcs.Serializer{}

	// Serialize TokenPriceUpdates
	bcs.SerializeSequenceWithFunction(report.PriceUpdates.TokenPriceUpdates, s, func(s *bcs.Serializer, item TokenPriceUpdate) {
		if len(item.SourceToken) != DefaultByteSize {
			s.SetError(fmt.Errorf("source token address must be exactly 32 bytes, got %d", len(item.SourceToken)))
			return
		}

		// Serialize as vector<u8> instead of address (as in the contract)
		s.FixedBytes(item.SourceToken)

		// Serialize usd_per_token as u256
		if item.UsdPerToken == nil {
			s.U256(*big.NewInt(0))
		} else {
			s.U256(*item.UsdPerToken)
		}
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize TokenPriceUpdates: %w", s.Error())
	}

	// Serialize GasPriceUpdates
	bcs.SerializeSequenceWithFunction(report.PriceUpdates.GasPriceUpdates, s, func(s *bcs.Serializer, item GasPriceUpdate) {
		// Serialize dest_chain_selector as u64
		s.U64(item.DestChainSelector)
		// Serialize usd_per_unit_gas as u256
		if item.UsdPerUnitGas == nil {
			s.U256(*big.NewInt(0))
		} else {
			s.U256(*item.UsdPerUnitGas)
		}
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize GasPriceUpdates: %w", s.Error())
	}

	// Helper function to serialize MerkleRoot vector
	serializeMerkleRoots := func(merkleRoots []MerkleRoot, name string) error {
		bcs.SerializeSequenceWithFunction(merkleRoots, s, func(s *bcs.Serializer, item MerkleRoot) {
			// Serialize source_chain_selector as u64
			s.U64(item.SourceChainSelector)
			// Serialize on_ramp_address as vector<u8>
			s.WriteBytes(item.OnRampAddress)
			// Serialize min_seq_nr as u64
			s.U64(item.MinSeqNr)
			// Serialize max_seq_nr as u64
			s.U64(item.MaxSeqNr)
			// Serialize merkle_root as fixed 32-byte vector
			if len(item.MerkleRoot) != DefaultByteSize {
				s.SetError(fmt.Errorf("merkle root must be exactly 32 bytes, got %d", len(item.MerkleRoot)))
				return
			}
			s.FixedBytes(item.MerkleRoot)
		})
		if s.Error() != nil {
			return fmt.Errorf("failed to serialize %s: %w", name, s.Error())
		}

		return nil
	}

	// Serialize BlessedMerkleRoots
	if err := serializeMerkleRoots(report.BlessedMerkleRoots, "BlessedMerkleRoots"); err != nil {
		return nil, err
	}

	// Serialize UnblessedMerkleRoots
	if err := serializeMerkleRoots(report.UnblessedMerkleRoots, "UnblessedMerkleRoots"); err != nil {
		return nil, err
	}

	// Serialize RMNSignatures
	bcs.SerializeSequenceWithFunction(report.RMNSignatures, s, func(s *bcs.Serializer, item []byte) {
		// Each RMN signature should be exactly 64 bytes (32 bytes R + 32 bytes S)
		if len(item) != DefaultByteSize*SignatureComponents {
			s.SetError(fmt.Errorf("RMN signature must be exactly 64 bytes, got %d", len(item)))
			return
		}
		s.FixedBytes(item)
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize RMNSignatures: %w", s.Error())
	}

	return s.ToBytes(), nil
}
