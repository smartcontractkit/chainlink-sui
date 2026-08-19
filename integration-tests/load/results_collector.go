package load

import (
	"time"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
)

// collectResults drains the results channel and builds a RunResults struct.
// It should be called after the WASP profile has finished.
// We stop when no new messages arrive for idleTimeout.
func collectResults(cfg *config.LoadTestConfig, ch <-chan config.SentMessage) *config.RunResults {
	results := &config.RunResults{
		RunName:             cfg.RunName,
		EnvName:             cfg.EnvName,
		SourceChainSelector: cfg.SourceChainSelector,
		DestChainSelector:   cfg.DestChainSelector,
		TotalMessages:       cfg.MessageCount,
		RunStarted:          time.Now().Format(time.RFC3339),
		Messages:            make([]config.SentMessage, 0, cfg.MessageCount),
	}

	idleTimeout := 2 * time.Second
	timer := time.NewTimer(idleTimeout)
	defer timer.Stop()

	for {
		select {
		case msg := <-ch:
			msg.SourceChainSelector = cfg.SourceChainSelector
			results.Messages = append(results.Messages, msg)
			if msg.Success {
				results.SuccessfulMessages++
			} else {
				results.FailedMessages++
			}

			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(idleTimeout)
		case <-timer.C:
			results.RunEnded = time.Now().Format(time.RFC3339)
			return results
		}
	}

}
