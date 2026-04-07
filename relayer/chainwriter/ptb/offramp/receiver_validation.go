package offramp

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// receiverStandardParamCount is the number of standard parameters that every
// ccip_receive callback must begin with:
//
//	[0] expected_message_id: vector<u8>
//	[1] ref: &CCIPObjectRef
//	[2] message: Any2SuiMessage
const receiverStandardParamCount = 3

// ValidateReceiverCallbackSignature validates that a receiver's ccip_receive
// callback does not declare extra parameters whose types belong to known CCIP
// protocol packages. This prevents a malicious receiver from tricking the
// relayer into injecting protocol-owned objects (e.g. OnRampState) as mutable
// PTB inputs in the transmitter-signed transaction.
func ValidateReceiverCallbackSignature(
	lggr logger.Logger,
	functionSig map[string]any,
	decodedParamTypes []string,
	ccipPackageId string,
	offRampPackageId string,
) error {
	if len(decodedParamTypes) < receiverStandardParamCount {
		return fmt.Errorf(
			"receiver callback has %d parameters, expected at least %d standard parameters "+
				"(expected_message_id, &CCIPObjectRef, Any2SuiMessage)",
			len(decodedParamTypes), receiverStandardParamCount,
		)
	}

	parametersRaw, ok := functionSig["parameters"]
	if !ok {
		return fmt.Errorf("missing 'parameters' field in receiver function signature")
	}
	parameters, ok := parametersRaw.([]any)
	if !ok {
		return fmt.Errorf("'parameters' field is not an array in receiver function signature")
	}

	// Walk raw parameters, skipping TxContext (mirroring DecodeParameters),
	// and inspect every extra parameter beyond the standard 3.
	decodedIdx := 0
	for i, rawParam := range parameters {
		meta := decodeParam(lggr, rawParam, "Reference")
		if meta.Name == "TxContext" {
			continue
		}

		if decodedIdx >= receiverStandardParamCount {
			if meta.Reference == "MutableReference" {
				if isDeniedProtocolPackage(meta.Address, ccipPackageId, offRampPackageId) {
					return fmt.Errorf(
						"receiver callback parameter %d declares mutable reference to CCIP protocol type %s::%s::%s; "+
							"receiver callbacks must not accept mutable references to CCIP protocol objects",
						i, meta.Address, meta.Module, meta.Name,
					)
				}
				if isDeniedProtocolModule(meta.Module, meta.Name) {
					return fmt.Errorf(
						"receiver callback parameter %d references denied protocol type %s::%s; "+
							"receiver callbacks must not accept references to CCIP internal objects",
						i, meta.Module, meta.Name,
					)
				}
			}
		}

		decodedIdx++
	}

	return nil
}

// ValidateReceiverObjectIdCount ensures the number of receiverObjectIds matches
// the number of extra parameters declared by the callback beyond the standard 3.
// A mismatch indicates the callback ABI and the message's extra args are
// inconsistent, which is a precondition for the object injection attack.
func ValidateReceiverObjectIdCount(decodedParamTypes []string, receiverObjectIdCount int) error {
	expectedExtraParams := len(decodedParamTypes) - receiverStandardParamCount
	if expectedExtraParams < 0 {
		expectedExtraParams = 0
	}
	if receiverObjectIdCount != expectedExtraParams {
		return fmt.Errorf(
			"receiver callback declares %d extra object parameters but receiverObjectIds contains %d entries; counts must match",
			expectedExtraParams, receiverObjectIdCount,
		)
	}
	return nil
}

// ValidateReceiverObjectIds checks that none of the supplied receiver object
// IDs reference known CCIP protocol objects. Accepting protocol objects as
// receiver callback arguments would let a malicious receiver modify protocol
// state via the transmitter-signed PTB.
func ValidateReceiverObjectIds(objectIds []string, addressMappings *OffRampAddressMappings) error {
	denied := map[string]string{
		addressMappings.CcipObjectRef: "CCIPObjectRef",
		addressMappings.OffRampState:  "OffRampState",
	}
	if addressMappings.CcipOwnerCap != "" {
		denied[addressMappings.CcipOwnerCap] = "CcipOwnerCap"
	}

	for i, objectId := range objectIds {
		if name, found := denied[objectId]; found {
			return fmt.Errorf(
				"receiverObjectIds[%d] (%s) references protocol object %s; "+
					"receiver callbacks must not be passed CCIP protocol objects",
				i, objectId, name,
			)
		}
	}
	return nil
}

func isDeniedProtocolPackage(addr, ccipPackageId, offRampPackageId string) bool {
	return addr != "" && (addr == ccipPackageId || addr == offRampPackageId)
}

// isDeniedProtocolModule provides a defense-in-depth check against known CCIP
// protocol module+type combinations. This catches cases where the attacker's
// package references protocol types whose package ID isn't in addressMappings
// (e.g. the onramp package).
func isDeniedProtocolModule(module, name string) bool {
	denied := map[string]map[string]bool{
		"onramp":               {"OnRampState": true},
		"offramp":              {"OffRampState": true},
		"fee_quoter":           {"FeeQuoterState": true},
		"token_admin_registry": {"TokenAdminRegistryState": true},
		"receiver_registry":    {"ReceiverRegistry": true},
		"nonce_manager":        {"NonceManagerState": true},
		"state_object":         {"CCIPObjectRef": true},
		"offramp_state_helper": {"ReceiverParams": true},
	}
	if names, ok := denied[module]; ok {
		return names[name]
	}
	return false
}
