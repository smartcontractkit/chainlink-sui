package bind

import (
	"bytes"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
)

type CallArgManager struct {
	inputs []transaction.CallArg
}

func NewCallArgManager() *CallArgManager {
	return &CallArgManager{
		inputs: []transaction.CallArg{},
	}
}

func NewCallArgManagerWithExisting(existingInputs []*transaction.CallArg) *CallArgManager {
	inputs := make([]transaction.CallArg, 0, len(existingInputs))
	for _, input := range existingInputs {
		if input != nil {
			inputs = append(inputs, *input)
		}
	}

	return &CallArgManager{
		inputs: inputs,
	}
}

func (m *CallArgManager) AddCallArg(arg *transaction.CallArg) transaction.Argument {
	if arg.Object != nil {
		if existingIndex := m.findObjectIndex(arg.Object); existingIndex != nil {
			if arg.Object.SharedObject != nil && arg.Object.SharedObject.Mutable {
				if m.inputs[*existingIndex].Object != nil &&
					m.inputs[*existingIndex].Object.SharedObject != nil {
					m.inputs[*existingIndex].Object.SharedObject.Mutable = true
				}
			}

			return transaction.Argument{
				Input: existingIndex,
			}
		}
	}

	// for unresolved objects, check if it already exists
	if arg.UnresolvedObject != nil {
		if existingIndex := m.findUnresolvedObjectIndex(&arg.UnresolvedObject.ObjectId); existingIndex != nil {
			return transaction.Argument{
				Input: existingIndex,
			}
		}
	}

	index := len(m.inputs)
	indexUint16 := uint16(index) // #nosec G115
	m.inputs = append(m.inputs, *arg)

	return transaction.Argument{
		Input: &indexUint16,
	}
}

func (m *CallArgManager) GetInputs() []*transaction.CallArg {
	result := make([]*transaction.CallArg, len(m.inputs))
	for i := range m.inputs {
		result[i] = &m.inputs[i]
	}

	return result
}

func (m *CallArgManager) findObjectIndex(obj *transaction.ObjectArg) *uint16 {
	for i, input := range m.inputs {
		if input.Object == nil {
			continue
		}

		if obj.SharedObject != nil && input.Object.SharedObject != nil {
			if bytes.Equal(obj.SharedObject.ObjectId[:], input.Object.SharedObject.ObjectId[:]) {
				index := uint16(i) // #nosec G115
				return &index
			}
		}

		if obj.ImmOrOwnedObject != nil && input.Object.ImmOrOwnedObject != nil {
			if bytes.Equal(obj.ImmOrOwnedObject.ObjectId[:], input.Object.ImmOrOwnedObject.ObjectId[:]) {
				index := uint16(i) // #nosec G115
				return &index
			}
		}

		if obj.Receiving != nil && input.Object.Receiving != nil {
			if bytes.Equal(obj.Receiving.ObjectId[:], input.Object.Receiving.ObjectId[:]) {
				index := uint16(i) // #nosec G115
				return &index
			}
		}
	}

	return nil
}

func (m *CallArgManager) findUnresolvedObjectIndex(objectId *models.SuiAddressBytes) *uint16 {
	for i, input := range m.inputs {
		if input.UnresolvedObject != nil {
			if bytes.Equal(objectId[:], input.UnresolvedObject.ObjectId[:]) {
				index := uint16(i) // #nosec G115
				return &index
			}
		}
	}

	return nil
}

func (m *CallArgManager) ConvertCallArgsToArguments(callArgs []*transaction.CallArg) ([]*transaction.Argument, error) {
	arguments := make([]*transaction.Argument, len(callArgs))

	for i, callArg := range callArgs {
		if callArg == nil {
			return nil, fmt.Errorf("nil CallArg at index %d", i)
		}
		arg := m.AddCallArg(callArg)
		arguments[i] = &arg
	}

	return arguments, nil
}

func (m *CallArgManager) ConvertEncodedCallArgsToArguments(encodedArgs []*EncodedCallArgument) ([]*transaction.Argument, error) {
	arguments := make([]*transaction.Argument, len(encodedArgs))

	for i, encArg := range encodedArgs {
		if encArg == nil {
			return nil, fmt.Errorf("nil EncodedCallArgument at index %d", i)
		}

		if err := encArg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid EncodedCallArgument at index %d: %w", i, err)
		}

		if encArg.Argument != nil {
			// No deduplication needed for transaction.Argument (expected to be Result, NestedResult, or GasCoin)
			// TODO: validate that it's not Input?
			arguments[i] = encArg.Argument
		} else if encArg.CallArg != nil {
			// Normal CallArg processing with deduplication
			arg := m.AddCallArg(encArg.CallArg)
			arguments[i] = &arg
		}
	}

	return arguments, nil
}
