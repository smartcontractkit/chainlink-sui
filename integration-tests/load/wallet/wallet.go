package wallet

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"math/big"

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

// GenerateSuiWallets creates N independent Ed25519 keypairs.
// Returns a slice of Wallets with SuiSigner and Address populated.
func GenerateSuiWallets(n int) ([]*Wallet, error) {
	wallets := make([]*Wallet, 0, n)
	for i := 0; i < n; i++ {
		_, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate Ed25519 key %d: %w", i, err)
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

// GenerateEVMWallets creates N independent secp256k1 keypairs with TransactOpts.
// Returns a slice of Wallets with EVMTransactOpts, EVMPrivKey, and Address populated.
func GenerateEVMWallets(n int, chainID *big.Int) ([]*Wallet, error) {
	wallets := make([]*Wallet, 0, n)
	for i := 0; i < n; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate secp256k1 key %d: %w", i, err)
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
