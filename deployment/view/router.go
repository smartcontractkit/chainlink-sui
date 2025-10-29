package view

type RouterView struct {
	IsTestRouter bool              `json:"isTestRouter"`
	OnRamps      map[uint64]string `json:"onRamps"`  // Map of DestinationChainSelector to OnRampAddress
	OffRamps     map[uint64]string `json:"offRamps"` // Map of DestinationChainSelector to OffRampAddress
}

// GenerateRouterView generates a router view for a given router address
// This is a mocked implementation for now
func GenerateRouterView(routerAddress string, offRampAddresses []string, isTestRouter bool) (RouterView, error) {
	// Mock data for now
	onRamps := map[uint64]string{
		123456789: "0x1234567890abcdef",
		987654321: "0xfedcba0987654321",
	}

	offRamps := map[uint64]string{
		123456789: "0xabcdef1234567890",
		987654321: "0x0987654321fedcba",
	}

	return RouterView{
		IsTestRouter: isTestRouter,
		OnRamps:      onRamps,
		OffRamps:     offRamps,
	}, nil
}
