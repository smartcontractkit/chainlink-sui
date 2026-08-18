//go:build integration

package offramp

// ValidateReceiverObjectOwner exports the unexported transmitter-ownership guard for the external
// integration test in execute_integration_test.go (package offramp_test). The external package is
// required because an internal test file cannot import testutils without forming an import cycle
// (testutils -> txm -> offramp). Test-only; not part of the public API.
var ValidateReceiverObjectOwner = validateReceiverObjectOwner
