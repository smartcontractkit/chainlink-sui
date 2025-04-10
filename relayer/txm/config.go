package txm

const (
	// DefaultBroadcastChanSize is the default size of the broadcast channel.
	DefaultBroadcastChanSize = 100
	// DefaultConfirmerPoolPeriodSeconds is the default period for the confirmer pool.
	DefaultConfirmerPoolPeriodSeconds = 1

	// DefaultRequestType is the default request type for transactions.
	DefaultRequestType = "WaitForLocalExecution"
)

type Config struct {
	BroadcastChanSize          uint
	RequestType                string
	ConfirmerPoolPeriodSeconds uint
}

var DefaultConfigSet = Config{
	BroadcastChanSize:          DefaultBroadcastChanSize,
	RequestType:                DefaultRequestType,
	ConfirmerPoolPeriodSeconds: DefaultConfirmerPoolPeriodSeconds,
}
