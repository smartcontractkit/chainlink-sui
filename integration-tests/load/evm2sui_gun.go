package load

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/evm"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/wallet"
)

// EVM2SuiMsgGun sends message-only CCIP messages from EVM to Sui.
// Each instance owns one wallet (TransactOpts) exclusively.
type EVM2SuiMsgGun struct {
	wallet            *wallet.Wallet
	ethClient         *ethclient.Client
	routerAddress     common.Address
	destChainSelector uint64
	receiver          [32]byte
	data              []byte
	extraArgs         []byte
	resultsCh         chan<- config.SentMessage
}

// Call implements the wasp.Gun interface.
// It sends a message-only CCIP request using the wallet's dedicated auth.
func (g *EVM2SuiMsgGun) Call(_ *wasp.Generator) *wasp.Response {
	ctx := context.Background()
	group := "evm->sui"

	// No nonce management needed — single goroutine owns this auth.
	messageID, txHash, err := evm.SendMessage(
		ctx,
		g.ethClient,
		g.wallet.EVMTransactOpts,
		g.routerAddress,
		g.destChainSelector,
		g.receiver[:],
		g.data,
		g.extraArgs,
	)

	if err != nil {
		g.pushResult(false, "", "", err)
		return &wasp.Response{Failed: true, Error: err.Error(), Group: group}
	}

	g.pushResult(true, messageID, txHash, nil)
	return &wasp.Response{
		Failed: false,
		Group:  group,
		Data: map[string]any{
			"messageID": messageID,
			"txHash":    txHash,
		},
	}
}

func (g *EVM2SuiMsgGun) pushResult(success bool, messageID, txHash string, err error) {
	msg := config.SentMessage{
		SourceChainSelector: 0, // set by test harness
		DestChainSelector:   g.destChainSelector,
		Timestamp:           time.Now().Format(time.RFC3339),
		Success:             success,
		MessageID:           messageID,
		TransactionHash:     txHash,
	}
	if err != nil {
		msg.Error = err.Error()
	}
	g.resultsCh <- msg
}
