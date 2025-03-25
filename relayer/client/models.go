package client

type TransactionBlockOptions struct {
	ShowInput          bool `json:"showInput,omitempty"`
	ShowRawInput       bool `json:"showRawInput,omitempty"`
	ShowEffects        bool `json:"showEffects,omitempty"`
	ShowEvents         bool `json:"showEvents,omitempty"`
	ShowObjectChanges  bool `json:"showObjectChanges,omitempty"`
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
}

// TransactionBlockRequest represents the request the SuiExecuteTransactionBlock endpoint.
// https://docs.sui.io/sui-api-ref#sui_executetransactionblock
type TransactionBlockRequest struct {
	// BCS serialized transaction data bytes without its type tag, as base-64 encoded string.
	TxBytes string `json:"txBytes"`
	// A list of signatures (`flag || signature || pubkey` bytes, as base-64 encoded string).
	// Signature is committed to the intent message of the transaction data, as base-64 encoded string.
	Signatures []string `json:"signature"`
	// Options for specifying the content to be returned
	Options TransactionBlockOptions `json:"options"`
	// The request type, derived from `SuiTransactionBlockResponseOptions` if None.
	// The optional enumeration values are: `WaitForEffectsCert`, or `WaitForLocalExecution`
	RequestType string `json:"requestType"`
}
