package load

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/wallet"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// Sui2EVMMsgGun sends message-only CCIP messages from Sui to EVM.
// Each instance owns one wallet and its gas coin pool exclusively.
type Sui2EVMMsgGun struct {
	wallet            *wallet.Wallet
	gasPool           *sui.SuiCoinPool
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
// It pops two SUI coins (gas + fee) from the wallet's pool and sends a message-only CCIP request.
func (g *Sui2EVMMsgGun) Call(_ *wasp.Generator) *wasp.Response {
	ctx := context.Background()
	group := "sui->evm"

	gasCoin, err := g.gasPool.Pop(ctx)
	if err != nil {
		g.pushResult(false, "", "", 0, err)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("gas coin pop: %v", err), Group: group}
	}
	feeCoin, err := g.gasPool.Pop(ctx)
	if err != nil {
		g.pushResult(false, "", "", 0, err)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("fee coin pop: %v", err), Group: group}
	}
	if gasCoin == feeCoin {
		g.pushResult(false, "", "", 0, fmt.Errorf("gas and fee coin IDs are equal: %s", gasCoin))
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
		gasCoin,
		g.feeTokenType,
		g.suiCoinMetaID,
		feeCoin,
		g.destChainSelector,
		g.receiver,
		g.data,
		g.gasBudget,
		g.evmCallbackGas,
	)

	if err != nil {
		g.pushResult(false, "", "", 0, err)
		return &wasp.Response{Failed: true, Error: err.Error(), Group: group}
	}

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
	g.resultsCh <- msg
}
