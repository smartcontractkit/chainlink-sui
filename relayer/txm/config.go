package txm

const (
	// DefaultBroadcastChanSize is the default size of the broadcast channel.
	DefaultBroadcastChanSize = 100
	// DefaultConfirmPollSecs is the default period for the confirmer pool.
	DefaultConfirmPollSecs = 10

	// DefaultRequestType is the default request type for transactions.
	DefaultRequestType = "WaitForLocalExecution"

	DefaultMaxGasAmount          = 200000
	DefaultMaxTxRetryAttempts    = 5
	DefaultTransactionTimeout    = "30s"
	DefaultMaxConcurrentRequests = 100
)

type Config struct {
	BroadcastChanSize     uint
	RequestType           string
	ConfirmPollSecs       uint
	DefaultMaxGasAmount   uint64
	MaxTxRetryAttempts    uint64
	TransactionTimeout    string
	MaxConcurrentRequests uint64
}

var DefaultConfigSet = Config{
	BroadcastChanSize: DefaultBroadcastChanSize,
	RequestType:       DefaultRequestType,
	ConfirmPollSecs:   DefaultConfirmPollSecs,

	DefaultMaxGasAmount: DefaultMaxGasAmount,
	MaxTxRetryAttempts:  DefaultMaxTxRetryAttempts,

	TransactionTimeout:    DefaultTransactionTimeout,
	MaxConcurrentRequests: DefaultMaxConcurrentRequests,
}
