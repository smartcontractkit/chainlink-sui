package offramp

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

func AnyPointer[T any](v T) *T {
	return &v
}

type OffRampAddressMappings struct {
	CcipPackageId    string `json:"ccipPackageId"`
	CcipObjectRef    string `json:"ccipObjectRef"`
	CcipOwnerCap     string `json:"ccipOwnerCap"`
	ClockObject      string `json:"clockObject"`
	OffRampPackageId string `json:"offRampPackageId"`
	OffRampState     string `json:"offRampState"`
}

// GetOfframpAddressMappings initializes and populates all required address mappings for PTB expansion operations.
//
// This function performs discovery and resolution of critical CCIP infrastructure addresses by:
// 1. Using the provided OffRamp package ID to query and discover the CCIP package ID
// 2. Reading owned objects to locate the OffRamp state pointer and extract the state address
// 3. Reading CCIP package objects to find the CCIP object reference and owner capability addresses
// 4. Assembling a complete address mapping required for subsequent PTB operations
//
// Parameters:
//   - ctx: Context for the operation, used for request lifecycle management
//   - lggr: Logger instance for debugging and operational visibility
//   - ptbClient: Sui PTB client for reading blockchain state and objects
//   - offRampPackageId: The OffRamp package identifier to start discovery from
//   - publicKey: Public key bytes for generating signer address for read operations
//
// Returns:
//   - OffRampAddressMappings: A struct containing all resolved addresses
//   - error: Error if any discovery step fails, objects are missing, or network issues occur
func GetOfframpAddressMappings(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	offRampPackageId string,
	publicKey []byte,
) (OffRampAddressMappings, error) {
	// address mappings for the expander
	addressMappings := OffRampAddressMappings{
		CcipPackageId:    "",
		CcipObjectRef:    "",
		CcipOwnerCap:     "",
		ClockObject:      "0x6",
		OffRampPackageId: offRampPackageId,
		OffRampState:     "",
	}

	// Use the `toAddress` (offramp package ID) from the config overrides to get the offramp pointer object
	signerAddress, err := client.GetAddressFromPublicKey(publicKey)
	if err != nil {
		lggr.Errorw("Error getting signer address", "error", err)
		return OffRampAddressMappings{}, err
	}
	getCCIPPackageIdResponse, err := ptbClient.ReadFunction(ctx, signerAddress, addressMappings.OffRampPackageId, "offramp", "get_ccip_package_id", []any{}, []string{})
	if err != nil {
		lggr.Errorw("Error reading ccip package id", "error", err)
		return OffRampAddressMappings{}, err
	}
	lggr.Debugw("getCCIPPackageIdResponse", "getCCIPPackageIdResponse", getCCIPPackageIdResponse)
	// Parse the response to get the returned address as a hex string
	var addressBytes []byte

	// Handle both byte slice and base64 string responses
	switch v := getCCIPPackageIdResponse[0].(type) {
	case []byte:
		// Response is already raw bytes ([]byte and []uint8 are the same type)
		addressBytes = v
	case string:
		// Response is base64-encoded string, decode it
		var decodeErr error
		addressBytes, decodeErr = base64.StdEncoding.DecodeString(v)
		if decodeErr != nil {
			lggr.Errorw("Error decoding base64 ccip package id", "error", decodeErr)
			return OffRampAddressMappings{}, decodeErr
		}
	default:
		lggr.Errorw("Unexpected type for ccip package id response", "type", fmt.Sprintf("%T", getCCIPPackageIdResponse[0]))
		return OffRampAddressMappings{}, fmt.Errorf("unexpected type for ccip package id response, got %T", getCCIPPackageIdResponse[0])
	}
	// Convert bytes to hex string with "0x" prefix
	ccipPackageId := "0x" + hex.EncodeToString(addressBytes)
	addressMappings.CcipPackageId = ccipPackageId

	lggr.Debugw("ccipPackageId", "ccipPackageId", addressMappings.CcipPackageId)
	lggr.Debugw("offRampPackageId", "offrampPackageId", addressMappings.OffRampPackageId)

	// get the offramp state object
	offrampOwnedObjects, err := ptbClient.ReadOwnedObjects(ctx, addressMappings.OffRampPackageId, nil)
	if err != nil {
		lggr.Errorw("Error reading offramp state object", "error", err)
		return OffRampAddressMappings{}, err
	}
	for _, ccipOwnedObject := range offrampOwnedObjects {
		if ccipOwnedObject.Data.Type != "" && strings.Contains(ccipOwnedObject.Data.Type, "offramp::OffRampStatePointer") {
			lggr.Debugw("Found offramp state object pointer", "fields", ccipOwnedObject.Data.Content.Fields)
			// parse the object into a map
			parsedObject := ccipOwnedObject.Data.Content.Fields
			lggr.Debugw("offRampStatePointer Parsed", "offRampStatePointer", parsedObject)
			addressMappings.OffRampState = parsedObject["off_ramp_state_id"].(string)

			break
		}
	}
	if addressMappings.OffRampState == "" {
		lggr.Errorw("Address mappings are not populated", "addressMappings", addressMappings)
		return OffRampAddressMappings{}, fmt.Errorf("address mappings are missing required fields for expander (offRampState)")
	}

	// Get the object pointer present in the CCIP package ID
	ccipOwnedObjects, err := ptbClient.ReadOwnedObjects(ctx, addressMappings.CcipPackageId, nil)
	if err != nil {
		lggr.Errorw("Error reading ccip object ref", "error", err)
		return OffRampAddressMappings{}, err
	}
	for _, ccipOwnedObject := range ccipOwnedObjects {
		if ccipOwnedObject.Data.Type != "" && strings.Contains(ccipOwnedObject.Data.Type, "state_object::CCIPObjectRefPointer") {
			// parse the object into a map
			parsedObject := ccipOwnedObject.Data.Content.Fields
			if err != nil {
				lggr.Errorw("Error parsing ccip object ref", "error", err)
				return OffRampAddressMappings{}, err
			}
			lggr.Debugw("ccipObjectRefPointer", "ccipObjectRefPointer", parsedObject)
			addressMappings.CcipObjectRef = parsedObject["object_ref_id"].(string)
			addressMappings.CcipOwnerCap = parsedObject["owner_cap_id"].(string)

			break
		}
	}
	// check that address mappings are populated
	if addressMappings.CcipObjectRef == "" || addressMappings.CcipOwnerCap == "" {
		lggr.Errorw("Address mappings are not populated", "addressMappings", addressMappings)
		return OffRampAddressMappings{}, fmt.Errorf("address mappings are missing required fields for expander (ccipObjectRef, ccipOwnerCap)")
	}

	lggr.Debugw("Address mappings for expander", "addressMappings", addressMappings)

	return addressMappings, nil
}
