package txm

import (
	"crypto/rand"
	"math/big"
	"time"
)

const (
	// RandomPrecision is the number used for generating random values with sufficient precision
	// when calculating jitter. Higher values give better precision in the resulting fractional values.
	RandomPrecision = 1000000
)

func GetCurrentUnixTimestamp() uint64 {
	unixTime := time.Now().UTC().Unix()
	if unixTime < 0 {
		// Handle unexpected negative time or return a default
		return 0
	}

	return uint64(unixTime)
}

// GetTicker creates a ticker that fires at regular intervals based on the provided period.
// The function applies jitter to the base period to prevent synchronized timing across multiple instances.
//
// Parameters:
//   - basePeriod: The base polling frequency in seconds
//
// Returns:
//   - *time.Ticker: A ticker that triggers at the jittered interval
//   - time.Duration: The actual jittered duration used for the ticker
//
// The returned ticker should be stopped when no longer needed using ticker.Stop()
// to prevent resource leaks.
func GetTicker(basePeriod uint) (*time.Ticker, time.Duration) {
	// Convert uint seconds to time.Duration safely
	var baseDuration time.Duration

	// Maximum value that can be safely converted to time.Duration and multiplied by time.Second
	maxSafeSeconds := uint(int64(1<<63-1) / int64(time.Second))

	if basePeriod <= maxSafeSeconds {
		// Safe conversion when within limits
		baseDuration = time.Duration(basePeriod) * time.Second
	} else {
		// Cap at maximum safe value if input is too large
		baseDuration = time.Duration(maxSafeSeconds) * time.Second
	}

	// Add jitter to the base duration
	jitteredDuration := AddJitter(baseDuration)

	// Create and return a ticker with the jittered duration
	return time.NewTicker(jitteredDuration), jitteredDuration
}

func AddJitter(d time.Duration) time.Duration {
	// Apply up to ±25% jitter
	jitterFactor := 0.25

	// Generate a secure random float64 between -jitterFactor and +jitterFactor
	// using crypto/rand instead of math/rand for better security
	maxRand := big.NewInt(RandomPrecision) // Use a large number for precision
	randInt, err := rand.Int(rand.Reader, maxRand)
	if err != nil {
		// If we can't generate a secure random number, return the original duration
		return d
	}

	// Convert to float64 between 0 and 1
	randFloat := float64(randInt.Int64()) / float64(maxRand.Int64())

	// Convert to range -jitterFactor to +jitterFactor
	jitter := (randFloat*2*jitterFactor - jitterFactor)

	// Apply jitter to base duration
	return time.Duration(float64(d) * (1 + jitter))
}
