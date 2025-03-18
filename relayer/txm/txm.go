package txm

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

const expectedFunctionTokens = 3

type TxManager interface {
	Enqueue(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerAddress, function string, typeArgs []string, paramTypes []string, paramValues []any, simulateTx bool) error
}

type SuiTxm struct {
	lggr               logger.Logger
	suiGateway         client.SuiClient
	keyStoreRepository keystore.Keystore
	IsExecutionLocal   bool
	signer             signer.SuiSigner
}

func NewSuiTxm(lggr logger.Logger, gateway client.SuiClient, k keystore.Keystore, isExecutionLocal bool, sig signer.SuiSigner) (*SuiTxm, error) {
	return &SuiTxm{
		lggr:               logger.Named(lggr, "SuiTxm"),
		suiGateway:         gateway,
		keyStoreRepository: k,
		IsExecutionLocal:   isExecutionLocal,
		signer:             sig,
	}, nil
}

func (s *SuiTxm) Enqueue(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerAddress, function string, typeArgs []string, paramTypes []string, paramValues []any, simulateTx bool) error {
	functionTokens := strings.Split(function, "::")
	if len(functionTokens) != expectedFunctionTokens {
		msg := fmt.Sprintf("unexpected function name, expected 3 tokens, got %d", len(functionTokens))
		s.lggr.Error(msg)

		return errors.New(msg)
	}

	packageObjectId := functionTokens[0]
	moduleName := functionTokens[1]
	functionName := functionTokens[2]

	if len(paramTypes) != len(paramValues) {
		msg := fmt.Sprintf("unexpected number of parameters, expected %d, got %d", len(paramTypes), len(paramValues))
		s.lggr.Error(msg)

		return errors.New(msg)
	}

	functionValues := make([]any, len(paramValues))
	for i, v := range paramValues {
		value, err := codec.EncodeToSuiValue(paramTypes[i], v)
		if err != nil {
			s.lggr.Errorf("failed to encode value: %v", err)
			return err
		}

		functionValues[i] = value
	}

	rsp, err := s.suiGateway.MoveCall(ctx, models.MoveCallRequest{
		Signer:          signerAddress,
		PackageObjectId: packageObjectId,
		Module:          moduleName,
		Function:        functionName,
		// We will only need to pass the type arguments if the function is generic
		// TODO: must implement logic to check if the function is generic
		TypeArguments: []any{},
		Arguments:     functionValues,
		GasBudget:     txMetadata.GasLimit.String(),
	})

	if err != nil {
		msg := fmt.Sprintf("failed to move call: %v", err)
		s.lggr.Error(msg)

		return errors.New(msg)
	}

	txBytes, err := base64.StdEncoding.DecodeString(rsp.TxBytes)
	if err != nil {
		msg := fmt.Sprintf("failed to decode tx bytes: %v", err)
		s.lggr.Error(msg)

		return errors.New(msg)
	}

	signatures, err := s.signer.Sign(txBytes)
	if err != nil {
		log.Fatalf("Error signing transaction: %v", err)
	}

	var requestType string

	if s.IsExecutionLocal {
		requestType = "WaitForLocalExecution"
	} else {
		requestType = "WaitForCommit"
	}

	payload := client.TransactionBlockRequest{
		TxBytes:    rsp.TxBytes,
		Signatures: signatures,
		Options: client.TransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: requestType,
	}

	rsp2, err := s.suiGateway.SendTransaction(ctx, payload)

	if err != nil {
		msg := fmt.Sprintf("failed to send transaction: %v", err)
		s.lggr.Error(msg)

		return errors.New(msg)
	}

	s.lggr.Debugw("Transaction sent: %+v", rsp2)

	return nil
}
