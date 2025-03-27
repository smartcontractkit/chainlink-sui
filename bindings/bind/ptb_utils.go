package bind

import sui_pattokan "github.com/pattonkan/sui-go/sui"

// Utilities around the PTB and its types
func ToSuiAddress(address string) (*sui_pattokan.Address, error) {
	return sui_pattokan.AddressFromHex(address)
}
