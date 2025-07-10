package utils

// SuiSigner defines an interface for signing messages in the Sui blockchain format.
// This is a copy of the interface from relayer/signer to make bindings self-contained.
type SuiSigner interface {
	// Sign signs the given message and returns the serialized signature.
	Sign(message []byte) ([]string, error)

	// GetAddress returns the Sui address derived from the signer's public key
	GetAddress() (string, error)
}
