package sui

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GenericExtraArgsV2Tag is the 4-byte tag prefix for GenericExtraArgsV2.
const GenericExtraArgsV2Tag = "0x181dcf10"

// MakeBCSEVMExtraArgsV2 constructs GenericExtraArgsV2 for Sui→EVM messages.
//
// Format: tag (4 bytes) + gasLimit (u256, 32-byte little-endian) + allowOOO (bool, 1 byte)
//
// This is the BCS encoding used by the Sui CCIP onramp for EVM destination chains.
func MakeBCSEVMExtraArgsV2(gasLimit *big.Int, allowOOO bool) []byte {
	// Encode gasLimit as 32-byte little-endian
	glBytes := make([]byte, 32)
	gasLimit.FillBytes(glBytes) // big-endian fill
	// Reverse to little-endian
	for i, j := 0, 31; i < j; i, j = i+1, j-1 {
		glBytes[i], glBytes[j] = glBytes[j], glBytes[i]
	}

	oooByte := byte(0)
	if allowOOO {
		oooByte = 1
	}

	payload := append(glBytes, oooByte)
	return append(hexutil.MustDecode(GenericExtraArgsV2Tag), payload...)
}
