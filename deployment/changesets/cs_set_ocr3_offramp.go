package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var _ cldf.ChangeSetV2[SetOCR3OffRampConfig] = SetOCR3Offramp{}

type SetOCR3Offramp struct{}

// Ed25519Scheme Ed25519 signature scheme flag
// https://docs.sui.io/concepts/cryptography/transaction-auth/keys-addresses#address-format
const Ed25519Scheme byte = 0x00

// TODO: We need common OCR3 utilities from `chainlink/deployments`
func (s SetOCR3Offramp) Apply(e cldf.Environment, config SetOCR3OffRampConfig) (cldf.ChangesetOutput, error) {
	return cldf.ChangesetOutput{}, fmt.Errorf("not implemented")
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (s SetOCR3Offramp) VerifyPreconditions(e cldf.Environment, config SetOCR3OffRampConfig) error {
	return nil
}
