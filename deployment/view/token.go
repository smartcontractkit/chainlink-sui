package view

type TokenView struct {
	ContractMetaData

	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Decimals   uint8  `json:"decimals"`
	IconURI    string `json:"iconURI,omitempty"`
	ProjectURI string `json:"projectURI,omitempty"`
	Supply     uint64 `json:"supply"`

	Burners                   []string `json:"burners"`
	Minters                   []string `json:"minters"`
	ManagedTokenObjectAddress string   `json:"managedTokenObjectAddress,omitempty"`
}

// GenerateTokenView generates a token view for a given managed token
// This is a mocked implementation for now
func GenerateTokenView(managedTokenObjectAddress string) (TokenView, error) {
	//  TODO: implement token view logic
	return TokenView{}, nil
}
