package client

import (
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
)

type TransactionBlockOptions struct {
	ShowInput          bool `json:"showInput,omitempty"`
	ShowRawInput       bool `json:"showRawInput,omitempty"`
	ShowEffects        bool `json:"showEffects,omitempty"`
	ShowEvents         bool `json:"showEvents,omitempty"`
	ShowObjectChanges  bool `json:"showObjectChanges,omitempty"`
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
}

// TransactionRequestType defines the possible request types for transaction execution
type TransactionRequestType string

const (
	WaitForEffectsCert    TransactionRequestType = "WaitForEffectsCert"
	WaitForLocalExecution TransactionRequestType = "WaitForLocalExecution"
)

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

type MoveCallRequest struct {
	// the transaction signer's Sui address
	Signer string `json:"signer"`
	// the package containing the module and function
	PackageObjectId string `json:"packageObjectId"`
	// the specific module in the package containing the function
	Module string `json:"module"`
	// the function to be called
	Function string `json:"function"`
	// the type arguments to the function
	TypeArguments []any `json:"typeArguments"`
	// the arguments to the function
	Arguments []any `json:"arguments"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas *string `json:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget"`
}

type TxnMetaData struct {
	TxBytes string `json:"txBytes"`
}

type SuiTransactionBlockResponse struct {
	TxDigest      string                                                `json:"txDigest"`
	Status        SuiExecutionStatus                                    `json:"status"`
	Effects       suiclient.SuiTransactionBlockEffectsV1                `json:"effects"`
	Events        []*suiclient.Event                                    `json:"events,omitempty"`
	Timestamp     uint64                                                `json:"timestamp"`
	Height        uint64                                                `json:"height"`
	ObjectChanges []suiclient.WrapperTaggedJson[suiclient.ObjectChange] `json:"objectChanges,omitempty"`
}

type EventFilterByMoveEventModule struct {
	Package string `json:"package"`
	Module  string `json:"module"`
	Event   string `json:"event"`
}

type EventData struct {
	Id struct {
		TxDigest string `json:"txDigest"`
		EventSeq string `json:"eventSeq"`
	} `json:"id"`
	PackageId         string `json:"packageId"`
	TransactionModule string `json:"transactionModule"`
	Sender            string `json:"sender"`
	Type              struct {
		Address string `json:"address"`
		Module  string `json:"module"`
		Name    string `json:"name"`
	} `json:"type"`
	ParsedJson  any    `json:"parsedJson"`
	Bcs         string `json:"bcs"`
	TimestampMs string `json:"timestampMs"`
}

type PaginatedEventsResponse struct {
	Data        []EventData `json:"data"`
	NextCursor  string      `json:"nextCursor"`
	HasNextPage bool        `json:"hasNextPage"`
}

type EventId struct {
	TxDigest string      `json:"txDigest"`
	EventSeq *sui.BigInt `json:"eventSeq"`
}

type CoinData struct {
	CoinType     string `json:"coinType"`
	CoinObjectId string `json:"coinObjectId"`
	Version      string `json:"version"`
	Digest       string `json:"digest"`
	Balance      string `json:"balance"`
	PreviousTx   string `json:"previousTx"`
}

type SuiExecutionStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type QuerySortOptions struct {
	Descending bool `json:"descending"`
}

type TransactionResult struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
