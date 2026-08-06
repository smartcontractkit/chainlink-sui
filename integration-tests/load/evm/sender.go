package evm

import (
	"context"
	"fmt"
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
	return SendTokenMessage(
		ctx,
		client,
		auth,
		routerAddress,
		destChainSelector,
		receiver,
		data,
		common.Address{},
		big.NewInt(0),
		extraArgs,
	)
}

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
