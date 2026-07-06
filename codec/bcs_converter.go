package codec

import (
	"fmt"
	"strconv"

	aptosBCS "github.com/aptos-labs/aptos-go-sdk/bcs"
)

// BCSPrimitiveHandler defines a function that reads a primitive type from a BCS deserializer
type BCSPrimitiveHandler func(*aptosBCS.Deserializer) (any, error)

// BCSVectorHandler defines a function that reads a vector type from a BCS deserializer
type BCSVectorHandler func(*aptosBCS.Deserializer, uint64) (any, error)

// BCSTypeConverter provides a registry-based approach to converting BCS types to JSON-compatible values
type BCSTypeConverter struct {
	primitiveHandlers map[string]BCSPrimitiveHandler
	vectorHandlers    map[string]BCSVectorHandler
}

// NewBCSTypeConverter creates a new BCS type converter with all standard Sui types registered
func NewBCSTypeConverter() *BCSTypeConverter {
	c := &BCSTypeConverter{
		primitiveHandlers: make(map[string]BCSPrimitiveHandler),
		vectorHandlers:    make(map[string]BCSVectorHandler),
	}

	// Register primitive type handlers
	c.registerPrimitiveHandlers()
	c.registerVectorHandlers()

	return c
}

// registerPrimitiveHandlers registers all standard Sui primitive type handlers
func (c *BCSTypeConverter) registerPrimitiveHandlers() {
	// U8 - return as-is
	c.RegisterPrimitive("U8", func(d *aptosBCS.Deserializer) (any, error) {
		return d.U8(), nil
	})
	c.RegisterPrimitive("u8", func(d *aptosBCS.Deserializer) (any, error) {
		return d.U8(), nil
	})

	// U16 - return as-is
	c.RegisterPrimitive("U16", func(d *aptosBCS.Deserializer) (any, error) {
		return d.U16(), nil
	})
	c.RegisterPrimitive("u16", func(d *aptosBCS.Deserializer) (any, error) {
		return d.U16(), nil
	})

	// U32 - return as-is
	c.RegisterPrimitive("U32", func(d *aptosBCS.Deserializer) (any, error) {
		return d.U32(), nil
	})
	c.RegisterPrimitive("u32", func(d *aptosBCS.Deserializer) (any, error) {
		return d.U32(), nil
	})

	// U64 - return as string for JSON compatibility
	c.RegisterPrimitive("U64", func(d *aptosBCS.Deserializer) (any, error) {
		value := d.U64()
		return strconv.FormatUint(value, base10), nil
	})
	c.RegisterPrimitive("u64", func(d *aptosBCS.Deserializer) (any, error) {
		value := d.U64()
		return strconv.FormatUint(value, base10), nil
	})

	// U128 - return as string for JSON compatibility
	c.RegisterPrimitive("U128", func(d *aptosBCS.Deserializer) (any, error) {
		value := d.U128()
		return value.String(), nil
	})
	c.RegisterPrimitive("u128", func(d *aptosBCS.Deserializer) (any, error) {
		value := d.U128()
		return value.String(), nil
	})

	// U256 - return as string for JSON compatibility
	c.RegisterPrimitive("U256", func(d *aptosBCS.Deserializer) (any, error) {
		value := d.U256()
		return value.String(), nil
	})
	c.RegisterPrimitive("u256", func(d *aptosBCS.Deserializer) (any, error) {
		value := d.U256()
		return value.String(), nil
	})

	// Bool - return as-is
	c.RegisterPrimitive("Bool", func(d *aptosBCS.Deserializer) (any, error) {
		return d.Bool(), nil
	})
	c.RegisterPrimitive("bool", func(d *aptosBCS.Deserializer) (any, error) {
		return d.Bool(), nil
	})

	// Address - return as byte array
	c.RegisterPrimitive("Address", func(d *aptosBCS.Deserializer) (any, error) {
		addressBytesLen := 32
		return d.ReadFixedBytes(addressBytesLen), nil
	})
	c.RegisterPrimitive("address", func(d *aptosBCS.Deserializer) (any, error) {
		addressBytesLen := 32
		return d.ReadFixedBytes(addressBytesLen), nil
	})
}

// registerVectorHandlers registers standard vector type handlers
func (c *BCSTypeConverter) registerVectorHandlers() {
	// Vector<U8> - read as byte array
	c.RegisterVector("U8", func(d *aptosBCS.Deserializer, length uint64) (any, error) {
		bytes := make([]byte, length)
		for i := range length {
			bytes[i] = d.U8()
		}
		return bytes, nil
	})
	c.RegisterVector("u8", func(d *aptosBCS.Deserializer, length uint64) (any, error) {
		bytes := make([]byte, length)
		for i := range length {
			bytes[i] = d.U8()
		}
		return bytes, nil
	})

	// Vector<U64> - read as uint64 array
	c.RegisterVector("U64", func(d *aptosBCS.Deserializer, length uint64) (any, error) {
		uint64s := make([]uint64, length)
		for i := range length {
			uint64s[i] = d.U64()
		}
		return uint64s, nil
	})
	c.RegisterVector("u64", func(d *aptosBCS.Deserializer, length uint64) (any, error) {
		uint64s := make([]uint64, length)
		for i := range length {
			uint64s[i] = d.U64()
		}
		return uint64s, nil
	})

	// Vector<Address> - read as address array
	c.RegisterVector("Address", func(d *aptosBCS.Deserializer, length uint64) (any, error) {
		addresses := make([]any, length)
		for i := range length {
			addressBytesLen := 32
			addresses[i] = d.ReadFixedBytes(addressBytesLen)
		}
		return addresses, nil
	})
	c.RegisterVector("address", func(d *aptosBCS.Deserializer, length uint64) (any, error) {
		addresses := make([]any, length)
		for i := range length {
			addressBytesLen := 32
			addresses[i] = d.ReadFixedBytes(addressBytesLen)
		}
		return addresses, nil
	})
}

// RegisterPrimitive registers a handler for a primitive type
func (c *BCSTypeConverter) RegisterPrimitive(typeName string, handler BCSPrimitiveHandler) {
	c.primitiveHandlers[typeName] = handler
}

// RegisterVector registers a handler for a vector type
func (c *BCSTypeConverter) RegisterVector(elementType string, handler BCSVectorHandler) {
	c.vectorHandlers[elementType] = handler
}

// DecodePrimitive decodes a primitive type using the registered handler
func (c *BCSTypeConverter) DecodePrimitive(deserializer *aptosBCS.Deserializer, primitiveType string) (any, error) {
	handler, ok := c.primitiveHandlers[primitiveType]
	if !ok {
		return nil, fmt.Errorf("unsupported primitive type: %s", primitiveType)
	}
	return handler(deserializer)
}

// DecodeVector decodes a vector type using the registered handler
func (c *BCSTypeConverter) DecodeVector(deserializer *aptosBCS.Deserializer, elementType string) (any, error) {
	length := deserializer.Uleb128()

	handler, ok := c.vectorHandlers[elementType]
	if !ok {
		// Fall back to generic primitive vector handling
		primitiveHandler, primitiveOk := c.primitiveHandlers[elementType]
		if !primitiveOk {
			return nil, fmt.Errorf("unsupported vector element type: %s", elementType)
		}

		// Generic vector of primitives
		result := make([]any, length)
		for i := range length {
			value, err := primitiveHandler(deserializer)
			if err != nil {
				return nil, fmt.Errorf("failed to decode vector element at index %d: %w", i, err)
			}
			result[i] = value
		}
		return result, nil
	}

	return handler(deserializer, uint64(length))
}

// HasPrimitiveHandler checks if a primitive type handler is registered
func (c *BCSTypeConverter) HasPrimitiveHandler(typeName string) bool {
	_, ok := c.primitiveHandlers[typeName]
	return ok
}

// HasVectorHandler checks if a vector type handler is registered
func (c *BCSTypeConverter) HasVectorHandler(elementType string) bool {
	_, ok := c.vectorHandlers[elementType]
	return ok
}

// Global BCS converter instance for package-wide use (lazy initialization)
var defaultBCSConverter *BCSTypeConverter

// getDefaultBCSConverter returns the global BCS converter, initializing it if necessary
func getDefaultBCSConverter() *BCSTypeConverter {
	if defaultBCSConverter == nil {
		defaultBCSConverter = NewBCSTypeConverter()
	}
	return defaultBCSConverter
}

// DecodeBCSPrimitive decodes a BCS primitive type using the default converter
func DecodeBCSPrimitive(deserializer *aptosBCS.Deserializer, primitiveType string) (any, error) {
	return getDefaultBCSConverter().DecodePrimitive(deserializer, primitiveType)
}

// DecodeBCSVector decodes a BCS vector type using the default converter
func DecodeBCSVector(deserializer *aptosBCS.Deserializer, elementType string) (any, error) {
	return getDefaultBCSConverter().DecodeVector(deserializer, elementType)
}
