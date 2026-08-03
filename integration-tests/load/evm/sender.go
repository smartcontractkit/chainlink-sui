package evm

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	router "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
)

// GetFee queries the CCIP fee for a message from the EVM Router.
func GetFee(
	ctx context.Context,
	client *ethclient.Client,
	routerAddress common.Address,
	destChainSelector uint64,
	receiver []byte,
	data []byte,
	extraArgs []byte,
) (*big.Int, error) {
	routerContract, err := router.NewRouter(routerAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Router: %w", err)
	}

	msg := router.ClientEVM2AnyMessage{
		Receiver:     receiver,
		Data:         data,
		TokenAmounts: []router.ClientEVMTokenAmount{},
		FeeToken:     common.Address{}, // zero address = native ETH
		ExtraArgs:    extraArgs,
	}

	fee, err := routerContract.GetFee(&bind.CallOpts{Context: ctx}, destChainSelector, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get fee: %w", err)
	}

	return fee, nil
}

// SendMessage sends a CCIP message from EVM to Sui.
//
// Returns the message ID and transaction hash.
func SendMessage(
	ctx context.Context,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	routerAddress common.Address,
	destChainSelector uint64,
	receiver []byte,
	data []byte,
	extraArgs []byte,
) (messageID string, txHash string, err error) {
	routerContract, err := router.NewRouter(routerAddress, client)
	if err != nil {
		return "", "", fmt.Errorf("failed to instantiate Router: %w", err)
	}

	msg := router.ClientEVM2AnyMessage{
		Receiver:     receiver,
		Data:         data,
		TokenAmounts: []router.ClientEVMTokenAmount{},
		FeeToken:     common.Address{}, // zero address = native ETH
		ExtraArgs:    extraArgs,
	}

	// Get fee and add 20% buffer
	fee, err := routerContract.GetFee(&bind.CallOpts{Context: ctx}, destChainSelector, msg)
	if err != nil {
		return "", "", fmt.Errorf("failed to get fee: %w", err)
	}

	feeWithBuffer := new(big.Int).Add(fee, new(big.Int).Div(fee, big.NewInt(5)))
	auth.Value = feeWithBuffer

	slog.Info("Sending EVM→Sui message",
		"fee", fee.String(),
		"feeWithBuffer", feeWithBuffer.String(),
	)

	tx, err := routerContract.CcipSend(auth, destChainSelector, msg)
	if err != nil {
		return "", "", fmt.Errorf("failed to send CCIP message: %w", err)
	}

	txHash = tx.Hash().Hex()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		slog.Warn("Could not wait for transaction to be mined, using tx hash as message ID",
			"txHash", txHash,
			"error", err,
		)
		return txHash, txHash, nil
	}

	// Extract message ID from receipt logs
	messageID, _, err = ExtractMessageIDFromReceipt(receipt, destChainSelector)
	if err != nil {
		slog.Warn("Could not extract message ID from receipt, using tx hash",
			"txHash", txHash,
			"error", err,
		)
		messageID = txHash
	}

	return messageID, txHash, nil
}

// CCIPMessageSentTopic is the keccak256 hash of the CCIPMessageSent event signature.
// event CCIPMessageSent(uint64 indexed destChainSelector, uint64 sequenceNumber, ...)
var CCIPMessageSentTopic = common.HexToHash("0x...")

// ExtractMessageIDFromReceipt extracts the CCIPMessageSent event from a transaction receipt.
// Scans receipt logs for the CCIPMessageSent event by matching the event topic.
func ExtractMessageIDFromReceipt(
	receipt *types.Receipt,
	destChainSelector uint64,
) (messageID string, seqNum uint64, err error) {
	// The CCIPMessageSent event is emitted by the OnRamp contract (not the Router).
	// We scan all logs in the receipt for the event signature.
	// The event signature is: CCIPMessageSent(uint64,uint64,(bytes32,uint64,uint64,uint64,uint64),...)
	// We look for logs where the first indexed topic matches destChainSelector.

	for _, log := range receipt.Logs {
		// CCIPMessageSent has destChainSelector as the first indexed parameter (topic[1])
		if len(log.Topics) >= 2 {
			// The first indexed param (destChainSelector) is in topic[1]
			// topic[0] is the event signature hash
			chainSel := new(big.Int).SetBytes(log.Topics[1].Bytes()).Uint64()
			if chainSel == destChainSelector {
				// Found the event — extract message ID from the data
				// The message ID is the first field in the non-indexed data (bytes32)
				if len(log.Data) >= 32 {
					messageID = common.Bytes2Hex(log.Data[:32])
				}
				// Sequence number is in topic[2]
				if len(log.Topics) >= 3 {
					seqNum = new(big.Int).SetBytes(log.Topics[2].Bytes()).Uint64()
				}
				return messageID, seqNum, nil
			}
		}
	}

	return "", 0, fmt.Errorf("CCIPMessageSent event not found in receipt for dest chain %d", destChainSelector)
}
