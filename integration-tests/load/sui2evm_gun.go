package load

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/wallet"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// Sui2EVMMsgGun sends message-only CCIP messages from Sui to EVM.
// Each instance owns one wallet and two dedicated coins:
// - gasCoinID for PTB gas
// - feeCoinID for ccip_send fee splits
type Sui2EVMMsgGun struct {
	wallet            *wallet.Wallet
	gasCoinID         string
	feeCoinID         string
	feeAmount         uint64
	callSeq           uint64
	mu                sync.Mutex
	ptbClient         *client.PTBClient
	ccipPkgID         string
	onRampPkgID       string
	ccipObjectRefID   string
	onRampStateID     string
	feeTokenType      string
	suiCoinMetaID     string
	destChainSelector uint64
	receiver          []byte
	data              []byte
	gasBudget         uint64
	evmCallbackGas    uint64
	resultsCh         chan<- config.SentMessage
}

// Call implements the wasp.Gun interface.
// It sends a message-only CCIP request using one persistent gas coin and one persistent fee coin.
func (g *Sui2EVMMsgGun) Call(_ *wasp.Generator) *wasp.Response {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.callSeq++
	seq := g.callSeq

	ctx := context.Background()
	group := "sui->evm"
	start := time.Now()

	slog.Info("Sui message send started",
		"wallet", g.wallet.Address,
		"callSeq", seq,
		"destChainSelector", g.destChainSelector,
	)

	if g.gasCoinID == "" || g.feeCoinID == "" {
		err := fmt.Errorf("missing gas/fee coin IDs: gas=%q fee=%q", g.gasCoinID, g.feeCoinID)
		g.pushResult(false, "", "", 0, err)
		return &wasp.Response{Failed: true, Error: err.Error(), Group: group}
	}
	if g.gasCoinID == g.feeCoinID {
		g.pushResult(false, "", "", 0, fmt.Errorf("gas and fee coin IDs are equal: %s", g.gasCoinID))
		return &wasp.Response{Failed: true, Error: "gas and fee coin IDs are equal", Group: group}
	}

	messageID, txDigest, seqNum, err := sui.SendMessage(
		ctx,
		g.ptbClient,
		g.wallet.SuiSigner,
		g.ccipPkgID,
		g.onRampPkgID,
		g.ccipObjectRefID,
		g.onRampStateID,
		g.gasCoinID,
		g.feeTokenType,
		g.suiCoinMetaID,
		g.feeCoinID,
		g.feeAmount,
		g.destChainSelector,
		g.receiver,
		g.data,
		g.gasBudget,
		g.evmCallbackGas,
	)

	if err != nil {
		slog.Error("Sui message send failed",
			"wallet", g.wallet.Address,
			"callSeq", seq,
			"elapsedMs", time.Since(start).Milliseconds(),
			"error", err,
		)
		g.pushResult(false, "", "", 0, err)
		return &wasp.Response{Failed: true, Error: err.Error(), Group: group}
	}

	slog.Info("Sui message send completed",
		"wallet", g.wallet.Address,
		"callSeq", seq,
		"elapsedMs", time.Since(start).Milliseconds(),
		"txDigest", txDigest,
		"messageID", messageID,
	)

	g.pushResult(true, messageID, txDigest, seqNum, nil)
	return &wasp.Response{
		Failed: false,
		Group:  group,
		Data: map[string]any{
			"messageID": messageID,
			"txDigest":  txDigest,
			"seqNum":    seqNum,
		},
	}
}

func (g *Sui2EVMMsgGun) pushResult(success bool, messageID, txDigest string, seqNum uint64, err error) {
	msg := config.SentMessage{
		SourceChainSelector: 0, // set by test harness
		DestChainSelector:   g.destChainSelector,
		Timestamp:           time.Now().Format(time.RFC3339),
		Success:             success,
		MessageID:           messageID,
		TransactionHash:     txDigest,
		SequenceNumber:      fmt.Sprintf("%d", seqNum),
	}
	if err != nil {
		msg.Error = err.Error()
	}
	select {
	case g.resultsCh <- msg:
	default:
		slog.Warn("Dropping result due to full results channel", "destChainSelector", g.destChainSelector)
	}
}
