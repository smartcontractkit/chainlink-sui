package client

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_receiver_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/receiver_registry"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	module_counter "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/counter"
	suiSigner "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// ClientReader maps abstract read requests to specific DevInspect calls using bindings.
// It uses a dynamic registration system to automatically route to DevInspect methods.
type ClientReader interface {
	// TryReadFunction attempts to fulfill the read using registered DevInspect bindings.
	// Returns (handled=true) if a mapping exists, otherwise (handled=false) to fall back to generic flow.
	TryReadFunction(
		ctx context.Context,
		signerAddress string,
		packageId string,
		module string,
		function string,
		args []any,
		argTypes []string,
	) (handled bool, results []any, err error)
}

// DevInspectHandler defines a generic interface for calling DevInspect methods on bindings
type DevInspectHandler interface {
	// Call executes the DevInspect method with the provided parameters
	Call(ctx context.Context, signerAddress string, packageId string, client sui.ISuiAPI, args []any, argTypes []string) ([]any, error)
}

// SimpleDevInspectHandler handles DevInspect calls that don't require arguments
type SimpleDevInspectHandler struct {
	BindingFactory func(packageId string, client sui.ISuiAPI) (any, error)
	MethodName     string
	ResultMapper   func(result any) (any, error) // Optional: converts Go struct fields to contract field names
}

func (h *SimpleDevInspectHandler) Call(ctx context.Context, signerAddress string, packageId string, client sui.ISuiAPI, args []any, argTypes []string) ([]any, error) {
	// Create the binding
	binding, err := h.BindingFactory(packageId, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create binding: %w", err)
	}

	// Get DevInspect interface via reflection
	bindingValue := reflect.ValueOf(binding)
	devInspectMethod := bindingValue.MethodByName("DevInspect")
	if !devInspectMethod.IsValid() {
		return nil, fmt.Errorf("binding does not have DevInspect method")
	}

	devInspect := devInspectMethod.Call(nil)[0]

	// Create call options
	devInspectSigner := suiSigner.NewDevInspectSigner(signerAddress)
	callOpts := &bind.CallOpts{
		Signer:           devInspectSigner,
		WaitForExecution: true,
	}

	// Call the target method
	targetMethod := devInspect.MethodByName(h.MethodName)
	if !targetMethod.IsValid() {
		return nil, fmt.Errorf("DevInspect interface does not have method: %s", h.MethodName)
	}

	callArgs := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(callOpts),
	}

	results := targetMethod.Call(callArgs)
	if len(results) != 2 { // expecting (result, error)
		return nil, fmt.Errorf("unexpected number of return values from DevInspect method")
	}

	// Check for error
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}

	result := results[0].Interface()

	// Apply result mapping if provided
	if h.ResultMapper != nil {
		mappedResult, err := h.ResultMapper(result)
		if err != nil {
			return nil, fmt.Errorf("failed to map result: %w", err)
		}
		result = mappedResult
	}

	// Check if the result mapper returned multiple results (tuple expansion)
	if resultSlice, ok := result.([]any); ok && len(resultSlice) > 1 {
		// Return the expanded tuple elements as separate results
		return resultSlice, nil
	}

	return []any{result}, nil
}

// ParameterizedDevInspectHandler handles DevInspect calls that require arguments
type ParameterizedDevInspectHandler struct {
	BindingFactory func(packageId string, client sui.ISuiAPI) (any, error)
	MethodName     string
	ArgumentParser func(args []any, argTypes []string) ([]any, error)
	ResultMapper   func(result any) (any, error) // Optional: converts Go struct fields to contract field names
}

func (h *ParameterizedDevInspectHandler) Call(ctx context.Context, signerAddress string, packageId string, client sui.ISuiAPI, args []any, argTypes []string) ([]any, error) {
	// Parse arguments
	parsedArgs, err := h.ArgumentParser(args, argTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Create the binding
	binding, err := h.BindingFactory(packageId, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create binding: %w", err)
	}

	// Get DevInspect interface
	bindingValue := reflect.ValueOf(binding)
	devInspectMethod := bindingValue.MethodByName("DevInspect")
	if !devInspectMethod.IsValid() {
		return nil, fmt.Errorf("binding does not have DevInspect method")
	}

	devInspect := devInspectMethod.Call(nil)[0]

	// Create call options
	devInspectSigner := suiSigner.NewDevInspectSigner(signerAddress)
	callOpts := &bind.CallOpts{
		Signer:           devInspectSigner,
		WaitForExecution: true,
	}

	// Call the target method with parsed arguments
	targetMethod := devInspect.MethodByName(h.MethodName)
	if !targetMethod.IsValid() {
		return nil, fmt.Errorf("DevInspect interface does not have method: %s", h.MethodName)
	}

	callArgs := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(callOpts),
	}

	// Add parsed arguments
	for _, arg := range parsedArgs {
		callArgs = append(callArgs, reflect.ValueOf(arg))
	}

	results := targetMethod.Call(callArgs)
	if len(results) != 2 { // expecting (result, error)
		return nil, fmt.Errorf("unexpected number of return values from DevInspect method")
	}

	// Check for error
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}

	result := results[0].Interface()

	// Apply result mapping if provided
	if h.ResultMapper != nil {
		mappedResult, err := h.ResultMapper(result)
		if err != nil {
			return nil, fmt.Errorf("failed to map result: %w", err)
		}
		result = mappedResult
	}

	// Check if the result mapper returned multiple results (tuple expansion)
	if resultSlice, ok := result.([]any); ok && len(resultSlice) > 1 {
		// Return the expanded tuple elements as separate results
		return resultSlice, nil
	}

	return []any{result}, nil
}

type clientReader struct {
	client   SuiPTBClient
	handlers map[string]DevInspectHandler
}

func NewClientReader(client SuiPTBClient) ClientReader {
	cr := &clientReader{
		client:   client,
		handlers: make(map[string]DevInspectHandler),
	}

	// Register handlers dynamically
	cr.registerHandlers()

	return cr
}

func (c *clientReader) registerHandlers() {
	// TokenAdminRegistry module handlers
	c.register("token_admin_registry", "type_and_version", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_token_admin_registry.NewTokenAdminRegistry(packageId, client)
		},
		MethodName: "TypeAndVersion",
	})

	c.register("token_admin_registry", "get_token_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_token_admin_registry.NewTokenAdminRegistry(packageId, client)
		},
		MethodName: "GetTokenConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("get_token_config requires 2 args (object_ref_id, token)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			token := fmt.Sprint(args[1])
			return []any{refObj, token}, nil
		},
		ResultMapper: createStructToMapConverter(),
	})

	// RMN Remote module handlers (RMNProxy contract)
	c.register("rmn_remote", "type_and_version", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_rmn_remote.NewRmnRemote(packageId, client)
		},
		MethodName: "TypeAndVersion",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("type_and_version requires 1 arg (object_ref_id)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			return []any{refObj}, nil
		},
	})

	c.register("rmn_remote", "get_report_digest_header", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_rmn_remote.NewRmnRemote(packageId, client)
		},
		MethodName: "GetReportDigestHeader",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_report_digest_header requires 1 arg (object_ref_id)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			return []any{refObj}, nil
		},
	})

	c.register("rmn_remote", "get_versioned_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_rmn_remote.NewRmnRemote(packageId, client)
		},
		MethodName: "GetVersionedConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_versioned_config requires 1 arg (object_ref_id)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			return []any{refObj}, nil
		},
		// Note: ResultTupleToStruct handling would be done at chainreader level
	})

	c.register("rmn_remote", "get_cursed_subjects", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_rmn_remote.NewRmnRemote(packageId, client)
		},
		MethodName: "GetCursedSubjects",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_cursed_subjects requires 1 arg (object_ref_id)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			return []any{refObj}, nil
		},
	})

	// RMNRemote contract
	c.register("rmn_remote", "get_arm", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_rmn_remote.NewRmnRemote(packageId, client)
		},
		MethodName: "GetArm",
	})

	// FeeQuoter module handlers
	c.register("fee_quoter", "type_and_version", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_fee_quoter.NewFeeQuoter(packageId, client)
		},
		MethodName: "TypeAndVersion",
	})

	c.register("fee_quoter", "get_static_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_fee_quoter.NewFeeQuoter(packageId, client)
		},
		MethodName: "GetStaticConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_static_config requires 1 arg (object_ref_id)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			return []any{refObj}, nil
		},
	})

	c.register("fee_quoter", "get_dest_chain_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_fee_quoter.NewFeeQuoter(packageId, client)
		},
		MethodName: "GetDestChainConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("get_dest_chain_config requires 2 args (object_ref_id, destChainSelector)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			destChainSelector := args[1] // uint64
			return []any{refObj, destChainSelector}, nil
		},
		ResultMapper: createStructToMapConverter(), // For field renames at chainreader level
	})

	// OffRamp module handlers
	c.register("offramp", "type_and_version", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_offramp.NewOfframp(packageId, client)
		},
		MethodName: "TypeAndVersion",
	})

	c.register("offramp", "get_static_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_offramp.NewOfframp(packageId, client)
		},
		MethodName: "GetStaticConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_static_config requires 1 arg (off_ramp_state_id)")
			}
			stateObj := convertToBindingObject(args[0], argTypes[0])
			return []any{stateObj}, nil
		},
		ResultMapper: createStructToMapConverter(), // For field renames
	})

	c.register("offramp", "get_dynamic_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_offramp.NewOfframp(packageId, client)
		},
		MethodName: "GetDynamicConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_dynamic_config requires 1 arg (off_ramp_state_id)")
			}
			stateObj := convertToBindingObject(args[0], argTypes[0])
			return []any{stateObj}, nil
		},
		ResultMapper: createStructToMapConverter(),
	})

	c.register("offramp", "get_source_chain_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_offramp.NewOfframp(packageId, client)
		},
		MethodName: "GetSourceChainConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("get_source_chain_config requires 2 args (off_ramp_state_id, sourceChainSelector)")
			}
			stateObj := convertToBindingObject(args[0], argTypes[0])
			sourceChainSelector := args[1] // uint64
			return []any{stateObj, sourceChainSelector}, nil
		},
		ResultMapper: createStructToMapConverter(),
	})

	c.register("offramp", "get_all_source_chain_configs", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_offramp.NewOfframp(packageId, client)
		},
		MethodName: "GetAllSourceChainConfigs",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("get_all_source_chain_configs requires 2 args (object_ref_id, off_ramp_state_id)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			stateObj := convertToBindingObject(args[1], argTypes[1])
			return []any{refObj, stateObj}, nil
		},
		ResultMapper: createStructToMapConverter(),
	})

	// Router module handlers
	c.register("router", "type_and_version", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_router.NewRouter(packageId, client)
		},
		MethodName: "TypeAndVersion",
	})

	c.register("router", "get_on_ramp", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_router.NewRouter(packageId, client)
		},
		MethodName: "GetOnRamp",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("get_on_ramp requires 2 args (off_ramp_state_id, destChainSelector)")
			}
			stateObj := convertToBindingObject(args[0], argTypes[0])
			destChainSelector := args[1] // uint64
			return []any{stateObj, destChainSelector}, nil
		},
	})

	// OnRamp module handlers
	c.register("onramp", "type_and_version", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_onramp.NewOnramp(packageId, client)
		},
		MethodName: "TypeAndVersion",
	})

	c.register("onramp", "get_dynamic_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_onramp.NewOnramp(packageId, client)
		},
		MethodName: "GetDynamicConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_dynamic_config requires 1 arg (on_ramp_state_id)")
			}
			stateObj := convertToBindingObject(args[0], argTypes[0])
			return []any{stateObj}, nil
		},
		ResultMapper: createStructToMapConverter(),
	})

	c.register("onramp", "get_static_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_onramp.NewOnramp(packageId, client)
		},
		MethodName: "GetStaticConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_static_config requires 1 arg (on_ramp_state_id)")
			}
			stateObj := convertToBindingObject(args[0], argTypes[0])
			return []any{stateObj}, nil
		},
		ResultMapper: createStructToMapConverter(),
	})

	c.register("onramp", "get_dest_chain_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_onramp.NewOnramp(packageId, client)
		},
		MethodName: "GetDestChainConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("get_dest_chain_config requires 2 args (on_ramp_state_id, destChainSelector)")
			}
			stateObj := convertToBindingObject(args[0], argTypes[0])
			destChainSelector := args[1] // uint64
			return []any{stateObj, destChainSelector}, nil
		},
		// Note: ResultTupleToStruct handling would be done at chainreader level
	})

	// ReceiverRegistry module handlers
	c.register("receiver_registry", "is_registered_receiver", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_receiver_registry.NewReceiverRegistry(packageId, client)
		},
		MethodName: "IsRegisteredReceiver",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("is_registered_receiver requires 2 args (ref, receiverPackageId)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			receiverPkg := fmt.Sprint(args[1])
			return []any{refObj, receiverPkg}, nil
		},
	})

	c.register("receiver_registry", "get_receiver_config", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_receiver_registry.NewReceiverRegistry(packageId, client)
		},
		MethodName: "GetReceiverConfig",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("get_receiver_config requires 2 args (ref, receiverPackageId)")
			}
			refObj := convertToBindingObject(args[0], argTypes[0])
			receiverPkg := fmt.Sprint(args[1])
			return []any{refObj, receiverPkg}, nil
		},
	})

	// Counter module handlers (for testing)
	c.register("counter", "get_count", &ParameterizedDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_counter.NewCounter(packageId, client)
		},
		MethodName: "GetCount",
		ArgumentParser: func(args []any, argTypes []string) ([]any, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("get_count requires 1 arg (counter)")
			}
			counterObj := convertToBindingObject(args[0], argTypes[0])
			return []any{counterObj}, nil
		},
	})

	c.register("counter", "get_address_list", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_counter.NewCounter(packageId, client)
		},
		MethodName:   "GetAddressList",
		ResultMapper: createStructToMapConverter(),
	})

	c.register("counter", "get_simple_result", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_counter.NewCounter(packageId, client)
		},
		MethodName:   "GetSimpleResult",
		ResultMapper: createStructToMapConverter(),
	})

	c.register("counter", "get_tuple_struct", &SimpleDevInspectHandler{
		BindingFactory: func(packageId string, client sui.ISuiAPI) (any, error) {
			return module_counter.NewCounter(packageId, client)
		},
		MethodName: "GetTupleStruct",
		ResultMapper: func(result any) (any, error) {
			// The tuple result is a slice, we need to convert each element and return as expanded results
			if tupleSlice, ok := result.([]any); ok {
				// Convert each tuple element
				convertedElements := make([]any, len(tupleSlice))
				for i, element := range tupleSlice {
					converted, err := convertStructToMap(element)
					if err != nil {
						return nil, fmt.Errorf("failed to convert tuple element %d: %w", i, err)
					}
					convertedElements[i] = converted
				}
				return convertedElements, nil
			}
			// If not a tuple slice, just convert normally
			return convertStructToMap(result)
		},
	})
}

func (c *clientReader) register(module, function string, handler DevInspectHandler) {
	key := normalizeKey(module, function)
	c.handlers[key] = handler
}

func (c *clientReader) TryReadFunction(
	ctx context.Context,
	signerAddress string,
	packageId string,
	module string,
	function string,
	args []any,
	argTypes []string,
) (bool, []any, error) {
	key := normalizeKey(module, function)
	handler, exists := c.handlers[key]
	if !exists {
		return false, nil, nil
	}

	results, err := handler.Call(ctx, signerAddress, packageId, c.client.GetClient(), args, argTypes)
	if err != nil {
		return true, nil, fmt.Errorf("DevInspect handler failed for %s: %w", key, err)
	}

	return true, results, nil
}

func normalizeKey(module, function string) string {
	return strings.ToLower(strings.TrimSpace(module)) + "::" + strings.ToLower(strings.TrimSpace(function))
}

func convertToBindingObject(val any, typeName string) bind.Object {
	// Check if already a bind.Object
	if obj, ok := val.(bind.Object); ok {
		return obj
	}

	// Convert string to bind.Object
	return bind.Object{Id: fmt.Sprint(val)}
}

// FieldMappingHelpers contains utilities for converting between Go struct fields and contract field names

// createStructToMapConverter creates a result mapper that converts Go structs to maps with lowercase field names
func createStructToMapConverter() func(result any) (any, error) {
	return func(result any) (any, error) {
		return convertStructToMap(result)
	}
}

// convertStructToMap converts Go structs to map[string]any with lowercase field names for compatibility with field renaming
// Uses reflection to preserve original types instead of JSON marshaling which converts all numbers to float64
func convertStructToMap(value any) (any, error) {
	if value == nil {
		return nil, nil
	}

	return convertStructToMapReflection(reflect.ValueOf(value))
}

// convertStructToMapReflection uses reflection to convert structs to maps while preserving types
func convertStructToMapReflection(v reflect.Value) (any, error) {
	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		result := make(map[string]any)
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			// Skip unexported fields
			if !fieldValue.CanInterface() {
				continue
			}

			// Convert field name from Go convention (Title case) to contract convention (camelCase)
			fieldName := strings.ToLower(field.Name[:1]) + field.Name[1:]

			convertedValue, err := convertStructToMapReflection(fieldValue)
			if err != nil {
				return nil, fmt.Errorf("failed to convert field %s: %w", field.Name, err)
			}

			result[fieldName] = convertedValue
		}
		return result, nil

	case reflect.Slice, reflect.Array:
		result := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			convertedValue, err := convertStructToMapReflection(v.Index(i))
			if err != nil {
				return nil, fmt.Errorf("failed to convert array element %d: %w", i, err)
			}
			result[i] = convertedValue
		}
		return result, nil

	case reflect.Map:
		if v.Type().Key().Kind() == reflect.String {
			result := make(map[string]any)
			for _, key := range v.MapKeys() {
				keyStr := key.String()
				mapValue := v.MapIndex(key)
				convertedValue, err := convertStructToMapReflection(mapValue)
				if err != nil {
					return nil, fmt.Errorf("failed to convert map value for key %s: %w", keyStr, err)
				}
				// Convert key to lowercase if it looks like a Go struct field
				if len(keyStr) > 0 && keyStr[0] >= 'A' && keyStr[0] <= 'Z' {
					keyStr = strings.ToLower(keyStr[:1]) + keyStr[1:]
				}
				result[keyStr] = convertedValue
			}
			return result, nil
		}
		// Non-string key maps, return as interface{}
		return v.Interface(), nil

	default:
		// Return primitive types and interfaces as-is
		return v.Interface(), nil
	}
}
