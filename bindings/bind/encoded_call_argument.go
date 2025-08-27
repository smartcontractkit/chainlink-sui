package bind

import (
	"fmt"

	"github.com/block-vision/sui-go-sdk/transaction"
)

// EncodedCallArgument represents an argument in an encoded Move function call.
type EncodedCallArgument struct {
	// Only one of these should be set:
	CallArg  *transaction.CallArg  // For regular values (Pure, Object)
	Argument *transaction.Argument // For PTB results (Result, NestedResult, GasCoin)

	// TypeName stores the original Move type string (e.g., "&mut 0x2::clock::Clock", "&OwnerCap")
	// This is used to preserve type information (including mutability) through the encoding process
	TypeName string
}

func NewEncodedCallArgFromCallArg(callArg *transaction.CallArg) *EncodedCallArgument {
	return &EncodedCallArgument{CallArg: callArg}
}

func NewEncodedCallArgFromCallArgWithType(callArg *transaction.CallArg, typeName string) *EncodedCallArgument {
	return &EncodedCallArgument{CallArg: callArg, TypeName: typeName}
}

func NewEncodedCallArgFromArgument(arg *transaction.Argument) *EncodedCallArgument {
	return &EncodedCallArgument{Argument: arg}
}

func NewEncodedCallArgFromArgumentWithType(arg *transaction.Argument, typeName string) *EncodedCallArgument {
	return &EncodedCallArgument{Argument: arg, TypeName: typeName}
}

// Validate ensures that exactly one field is set
func (e *EncodedCallArgument) Validate() error {
	if e == nil {
		return fmt.Errorf("nil EncodedCallArgument")
	}

	if e.CallArg != nil && e.Argument != nil {
		return fmt.Errorf("EncodedCallArgument has both CallArg and Argument set")
	}

	if e.CallArg == nil && e.Argument == nil {
		return fmt.Errorf("EncodedCallArgument has neither CallArg nor Argument set")
	}

	return nil
}

// IsArgument returns true if this is a transaction.Argument (PTB result)
func (e *EncodedCallArgument) IsArgument() bool {
	return e != nil && e.Argument != nil
}

// IsCallArg returns true if this is a regular CallArg
func (e *EncodedCallArgument) IsCallArg() bool {
	return e != nil && e.CallArg != nil
}
