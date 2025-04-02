package bind

import (
	"fmt"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suisigner"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

// Utilities around the PTB and its types
func ToSuiAddress(address string) (*sui.Address, error) {
	return sui.AddressFromHex(address)
}

func ToSuiSignatures(signatures []string) ([]*suisigner.Signature, error) {
	suiSignatures := make([]*suisigner.Signature, 0)
	for _, sig := range signatures {
		bytes, err := toSuiSignature(sig)
		if err != nil {
			return nil, err
		}
		suiSignatures = append(suiSignatures, bytes)
	}

	return suiSignatures, nil
}

func toSuiSignature(sig string) (*suisigner.Signature, error) {
	decoded, err := codec.DecodeBase64(sig)
	if err != nil {
		return nil, err
	}
	const lengthSig = 97

	if len(decoded) != lengthSig {
		return nil, fmt.Errorf("invalid signature length: expected 97, got %d", len(decoded))
	}
	var sigArray [lengthSig]byte
	copy(sigArray[:], decoded)

	return &suisigner.Signature{
		Ed25519SuiSignature: &suisigner.Ed25519SuiSignature{
			Signature: sigArray,
		},
	}, nil
}
