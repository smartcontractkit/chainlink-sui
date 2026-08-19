package wallet

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/hkdf"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// Wallet holds a funded account ready for load testing.
// For Sui wallets, EVMTransactOpts and EVMPrivKey are nil.
// For EVM wallets, SuiSigner is nil.
type Wallet struct {
	Address string

	// Sui fields
	SuiSigner bindutils.SuiSigner

	// EVM fields
	EVMTransactOpts *bind.TransactOpts
	EVMPrivKey      *ecdsa.PrivateKey
}

// GenerateSuiWallets creates N Ed25519 keypairs.
// If seed is non-nil, keys are deterministically derived via HKDF so the same
// seed always yields the same wallet set; otherwise keys are random.
func GenerateSuiWallets(n int, seed []byte) ([]*Wallet, error) {
	wallets := make([]*Wallet, 0, n)
	for i := 0; i < n; i++ {
		privateKey, err := deriveEd25519Key(seed, i)
		if err != nil {
			return nil, fmt.Errorf("failed to derive Ed25519 key %d: %w", i, err)
		}

		s := signer.NewPrivateKeySigner(privateKey)
		address, err := s.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get Sui address for wallet %d: %w", i, err)
		}

		wallets = append(wallets, &Wallet{
			Address:   address,
			SuiSigner: s,
		})
	}
	return wallets, nil
}

// GenerateEVMWallets creates N secp256k1 keypairs with TransactOpts.
// If seed is non-nil, keys are deterministically derived via HKDF so the same
// seed always yields the same wallet set; otherwise keys are random.
func GenerateEVMWallets(n int, chainID *big.Int, seed []byte) ([]*Wallet, error) {
	wallets := make([]*Wallet, 0, n)
	for i := 0; i < n; i++ {
		privateKey, err := deriveECDSAKey(seed, i)
		if err != nil {
			return nil, fmt.Errorf("failed to derive secp256k1 key %d: %w", i, err)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			return nil, fmt.Errorf("failed to create transactor for wallet %d: %w", i, err)
		}

		wallets = append(wallets, &Wallet{
			Address:         auth.From.Hex(),
			EVMTransactOpts: auth,
			EVMPrivKey:      privateKey,
		})
	}
	return wallets, nil
}

// deriveEd25519Key returns a deterministic Ed25519 key from (seed, index) or a
// random key when seed is nil.
func deriveEd25519Key(seed []byte, index int) (ed25519.PrivateKey, error) {
	if len(seed) == 0 {
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		return priv, err
	}
	info := []byte(fmt.Sprintf("sui-load-wallet-%d", index))
	r := hkdf.New(sha256.New, seed, nil, info)
	derived := make([]byte, ed25519.SeedSize)
	if _, err := r.Read(derived); err != nil {
		return nil, fmt.Errorf("hkdf read: %w", err)
	}
	return ed25519.NewKeyFromSeed(derived), nil
}

// deriveECDSAKey returns a deterministic secp256k1 key from (seed, index) or a
// random key when seed is nil.
func deriveECDSAKey(seed []byte, index int) (*ecdsa.PrivateKey, error) {
	if len(seed) == 0 {
		return crypto.GenerateKey()
	}
	info := []byte(fmt.Sprintf("evm-load-wallet-%d", index))
	r := hkdf.New(sha256.New, seed, nil, info)
	derived := make([]byte, 32)
	if _, err := r.Read(derived); err != nil {
		return nil, fmt.Errorf("hkdf read: %w", err)
	}
	return crypto.ToECDSAUnsafe(derived), nil
}

// ParseSeed decodes a hex-encoded 32-byte wallet seed.
func ParseSeed(s string) ([]byte, error) {
	seed, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid wallet seed hex: %w", err)
	}
	if len(seed) != 32 {
		return nil, fmt.Errorf("wallet seed must be 32 bytes, got %d", len(seed))
	}
	return seed, nil
}
