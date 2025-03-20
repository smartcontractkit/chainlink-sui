package bindings

type Counter struct {
	address string // find sui address type
	client  any    // find sui client type
	// find relevant types
}

func NewCounter(address string, client any) *Counter {
	return &Counter{
		address: address,
		client:  client,
	}
}

func (c *Counter) Increment() (string, error) {
	// calls increment function
	return "", nil
}

func (c *Counter) GetState() (string, error) {
	// finds the current count
	return "", nil
}

func (c *Counter) IncrementMult() (string, error) {
	// calls increment_mult function
	return "", nil
}

func (c *Counter) Initialize() (string, error) {
	// calls initialize function
	return "", nil
}
