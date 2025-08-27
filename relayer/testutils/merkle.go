package testutils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk/bcs"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/test-go/testify/require"
)

type MerkleTree [][32]byte

type RootMetadata struct {
	Role                 uint8
	ChainID              *big.Int
	MultiSig             []byte
	PreOpCount           uint64
	PostOpCount          uint64
	OverridePreviousRoot bool
}

type Op struct {
	Role         uint8
	ChainID      *big.Int
	MultiSig     []byte
	Nonce        uint64
	To           []byte
	ModuleName   string
	FunctionName string
	Data         []byte
}

type TimelockOperation struct {
	Target       []byte
	ModuleName   string
	FunctionName string
	Data         []byte
}

func HashPair(left, right [32]byte) [32]byte {
	if bytes.Compare(left[:], right[:]) < 0 {
		return crypto.Keccak256Hash(left[:], right[:])
	}

	return crypto.Keccak256Hash(right[:], left[:])
}

// MerkleTree methods
func (mt MerkleTree) GetRoot() [32]byte {
	return mt[len(mt)-1]
}

func (mt MerkleTree) GetProof(index int) [][32]byte {
	proof := [][32]byte{}

	//nolint:all
	for index < len(mt)-1 {
		siblingIndex := index ^ 1
		proof = append(proof, mt[siblingIndex])
		index = (len(mt) + 1 + index) / 2
	}

	return proof
}

func (mt MerkleTree) VerifyProof(proof [][32]byte, leaf [32]byte) bool {
	computedHash := leaf
	for _, p := range proof {
		computedHash = HashPair(computedHash, p)
	}

	return bytes.Equal(computedHash[:], mt[len(mt)-1][:])
}

func NewMerkleTree(leaves [][32]byte) (MerkleTree, error) {
	if len(leaves) == 0 {
		return nil, errors.New("empty leaf set")
	}

	// Calculate the next power of 2
	leafCount := len(leaves)
	treeSize := 1
	for treeSize < leafCount {
		treeSize *= 2
	}

	// Create a new slice with the correct size
	//nolint:all
	paddedLeaves := make([][32]byte, treeSize)
	copy(paddedLeaves, leaves)

	// Fill the rest with zero leaves
	zeroLeaf := [32]byte{}
	for i := leafCount; i < treeSize; i++ {
		paddedLeaves[i] = zeroLeaf
	}

	tree := make(MerkleTree, treeSize)
	copy(tree, paddedLeaves)

	index := 0
	for levelSize := treeSize; levelSize > 1; levelSize /= 2 {
		for i := index; i < index+levelSize; i += 2 {
			//nolint:all
			tree = append(tree, HashPair(tree[i], tree[i+1]))
		}
		index += levelSize
	}

	return tree, nil
}

func HashRootMetadata(metadata RootMetadata) common.Hash {
	// TODO: change value to SUI when contract is changed
	MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA := crypto.Keccak256([]byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA_APTOS"))
	ser := bcs.Serializer{}
	ser.FixedBytes(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA)
	ser.U8(metadata.Role)
	ser.U256(*metadata.ChainID)
	ser.WriteBytes(metadata.MultiSig)
	ser.U64(metadata.PreOpCount)
	ser.U64(metadata.PostOpCount)
	ser.Bool(metadata.OverridePreviousRoot)

	return crypto.Keccak256Hash(ser.ToBytes())
}

func HashOp(op *Op) common.Hash {
	// TODO: change value to SUI when contract is changed
	MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP := crypto.Keccak256([]byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP_APTOS"))
	ser := bcs.Serializer{}
	ser.FixedBytes(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP)
	ser.U8(op.Role)
	ser.U256(*op.ChainID)
	ser.WriteBytes(op.MultiSig)
	ser.U64(op.Nonce)
	ser.WriteBytes(op.To)
	ser.WriteString(op.ModuleName)
	ser.WriteString(op.FunctionName)
	ser.WriteBytes(op.Data)

	return crypto.Keccak256Hash(ser.ToBytes())
}

func CalculateSignedHash(rootHash [32]byte, validUntil uint64) [32]byte {
	// Equivalent to Solidity's abi.encode(bytes32, uint64)
	//nolint:all
	data := make([]byte, 64)
	copy(data[:32], rootHash[:])
	binary.BigEndian.PutUint64(data[56:], validUntil)

	// Keccak256 hash of the ABI encoded parameters
	hashedEncodedParams := crypto.Keccak256(data)

	// Prepare the Ethereum signed message
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	ethMsg := append(prefix, hashedEncodedParams...)

	// Final Keccak256 hash
	return crypto.Keccak256Hash(ethMsg)
}

func GenerateSignatures(t *testing.T, signers []ecdsa.PrivateKey, signedHash [32]byte) [][]byte {
	t.Helper()

	signatures := make([][]byte, len(signers))
	for i, signer := range signers {
		signature, err := crypto.Sign(signedHash[:], &signer)
		require.NoError(t, err)

		// Adjust the v value, we need to read 27.
		// ref: https://github.com/ethereum/go-ethereum/blob/b590cae89232299d54aac8aada88c66d00c5b34c/crypto/signature_nocgo.go#L90
		v := signature[crypto.RecoveryIDOffset]
		//nolint:all
		require.True(t, v >= 0 && v <= 3, "v should be between 0 and 3")
		signature[crypto.RecoveryIDOffset] += 27

		signatures[i] = signature
	}

	return signatures
}

func GenerateMerkleTree(ops []Op, rootMetadata RootMetadata) (MerkleTree, error) {
	leaves := make([][32]byte, len(ops)+1)
	leaves[0] = HashRootMetadata(rootMetadata)
	for i, op := range ops {
		leaves[i+1] = HashOp(&op)
	}

	return NewMerkleTree(leaves)
}

func SerializeScheduleBatchParams(ops []TimelockOperation, predecessor []byte, salt []byte, delay uint64) ([]byte, error) {
	return bcs.SerializeSingle(func(ser *bcs.Serializer) {
		// Serialize targets vector
		//nolint:gosec
		ser.Uleb128(uint32(len(ops)))
		for _, op := range ops {
			ser.FixedBytes(op.Target)
		}

		// Write module names
		//nolint:gosec
		ser.Uleb128(uint32(len(ops)))
		for _, op := range ops {
			ser.WriteString(op.ModuleName)
		}

		// Write function names
		//nolint:gosec
		ser.Uleb128(uint32(len(ops)))
		for _, op := range ops {
			ser.WriteString(op.FunctionName)
		}

		// Write data
		//nolint:gosec
		ser.Uleb128(uint32(len(ops)))
		for _, op := range ops {
			ser.WriteBytes(op.Data)
		}

		ser.WriteBytes(predecessor)
		ser.WriteBytes(salt)
		ser.U64(delay)
	})
}
