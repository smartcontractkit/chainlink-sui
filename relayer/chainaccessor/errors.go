package chainaccessor

import "errors"

// ErrNotImplemented is returned by ChainAccessor methods that are not yet
// implemented natively for Sui. The methods are present so that *SuiAccessor
// satisfies the ccipocr3.ChainAccessor interface, but callers should treat this
// error as "this capability is not available on Sui yet".
var ErrNotImplemented = errors.New("sui chain accessor: not implemented")

// ErrNotBound is returned when a read is attempted against a contract that has
// not been bound via Sync.
var ErrNotBound = errors.New("sui chain accessor: contract not bound")
