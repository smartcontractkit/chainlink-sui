package codec

// SuiFunctionParam defines a parameter for a Sui function call
type SuiFunctionParam struct {
	// Name of the parameter
	Name string
	// Type of the parameter (e.g., "u64", "String", "vector<u8>")
	Type string
	// Whether the parameter is required
	Required bool
	// Default value to use if not provided
	DefaultValue interface{}
}
