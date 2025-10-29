package view

type MCMSWithTimelockView struct {
	ContractMetaData

	Bypasser  MCMSConfig `json:"bypasser"`
	Proposer  MCMSConfig `json:"proposer"`
	Canceller MCMSConfig `json:"canceller"`

	TimelockMinDelay         uint64                    `json:"timelockMinDelay"`
	TimelockBlockedFunctions []TimelockBlockedFunction `json:"timelockBlockedFunctions"`
}

type MCMSConfig struct {
	Signers     []MCMSSigner `json:"signers"`
	GroupQuorum uint8        `json:"group_quorum"`
	GroupParent uint8        `json:"group_parent"`
}

type MCMSSigner struct {
	Signer     string `json:"signer"`
	EvmSigner  string `json:"evm_signer"`
	Index      uint64 `json:"index"`
	GroupIndex uint8  `json:"group_index"`
}

type TimelockBlockedFunction struct {
	Target       string `json:"target"`
	ModuleName   string `json:"moduleName"`
	FunctionName string `json:"functionName"`
}

// GenerateMCMSWithTimelockView generates an MCMS with timelock view for a given MCMS address
// This is a mocked implementation for now
func GenerateMCMSWithTimelockView(mcmsAddress string) (MCMSWithTimelockView, error) {
	// Mock data for now
	bypasserConfig := MCMSConfig{
		Signers: []MCMSSigner{
			{
				Signer:     "0xbypasser_signer1",
				EvmSigner:  "0xevmbypasser1",
				Index:      0,
				GroupIndex: 0,
			},
		},
		GroupQuorum: 1,
		GroupParent: 0,
	}

	proposerConfig := MCMSConfig{
		Signers: []MCMSSigner{
			{
				Signer:     "0xproposer_signer1",
				EvmSigner:  "0xevmproposer1",
				Index:      0,
				GroupIndex: 0,
			},
			{
				Signer:     "0xproposer_signer2",
				EvmSigner:  "0xevmproposer2",
				Index:      1,
				GroupIndex: 0,
			},
		},
		GroupQuorum: 2,
		GroupParent: 0,
	}

	cancellerConfig := MCMSConfig{
		Signers: []MCMSSigner{
			{
				Signer:     "0xcanceller_signer1",
				EvmSigner:  "0xevmcanceller1",
				Index:      0,
				GroupIndex: 0,
			},
		},
		GroupQuorum: 1,
		GroupParent: 0,
	}

	blockedFunctions := []TimelockBlockedFunction{
		{
			Target:       "0xtarget1",
			ModuleName:   "module1",
			FunctionName: "blocked_function1",
		},
	}

	return MCMSWithTimelockView{
		ContractMetaData: ContractMetaData{
			Address:        mcmsAddress,
			Owner:          "0xmock_mcms_owner",
			TypeAndVersion: "MCMSWithTimelock 1.0.0",
		},
		Bypasser:                 bypasserConfig,
		Proposer:                 proposerConfig,
		Canceller:                cancellerConfig,
		TimelockMinDelay:         3600,
		TimelockBlockedFunctions: blockedFunctions,
	}, nil
}
