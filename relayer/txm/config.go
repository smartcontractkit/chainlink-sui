package txm

const (
	// DefaultBroadcastChanSize is the default size of the broadcast channel.
	DefaultBroadcastChanSize = 100
)

type Config struct {
	BroadcastChanSize uint
	IsExecutionLocal  bool
}

var DefaultConfigSet = Config{
	BroadcastChanSize: DefaultBroadcastChanSize,
	IsExecutionLocal:  true,
}
