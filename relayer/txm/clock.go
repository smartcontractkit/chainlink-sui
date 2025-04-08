package txm

import "time"

func GetCurrentUnixTimestamp() uint64 {
	unixTime := time.Now().UTC().Unix()
	if unixTime < 0 {
		// Handle unexpected negative time or return a default
		return 0
	}

	return uint64(unixTime)
}
