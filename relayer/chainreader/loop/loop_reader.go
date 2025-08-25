package loop

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/loop"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

const (
	READ_COMPONENTS_COUNT = 3
)

func NewLoopChainReader(log logger.Logger, reader types.ContractReader) types.ContractReader {
	return &loopChainReader{logger: log, reader: reader, moduleAddresses: map[string]string{}}
}

type loopChainReader struct {
	services.Service
	types.UnimplementedContractReader
	logger          logger.Logger
	reader          types.ContractReader
	moduleAddresses map[string]string
}

func (s *loopChainReader) Name() string {
	return s.reader.Name()
}

func (s *loopChainReader) Ready() error {
	return s.reader.Ready()
}

func (s *loopChainReader) HealthReport() map[string]error {
	return s.reader.HealthReport()
}

func (s *loopChainReader) Start(ctx context.Context) error {
	return s.reader.Start(ctx)
}

func (s *loopChainReader) Close() error {
	return s.reader.Close()
}

func (s *loopChainReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	readComponents := strings.Split(readIdentifier, "-")
	if len(readComponents) != READ_COMPONENTS_COUNT {
		return fmt.Errorf("invalid read identifier: %s", readIdentifier)
	}

	_, contractName, _ := readComponents[0], readComponents[1], readComponents[2]

	_, ok := s.moduleAddresses[contractName]
	if !ok {
		return fmt.Errorf("no such contract: %s", contractName)
	}

	convertedResult := []byte{}

	jsonParamBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %+w", err)
	}

	// we always bind before calling query functions, because the LOOP plugin may have restarted.
	err = s.reader.Bind(ctx, s.getBindings())
	if err != nil {
		return fmt.Errorf("failed to re-bind before GetLatestValue: %w", err)
	}

	err = s.reader.GetLatestValue(ctx, readIdentifier, confidenceLevel, &jsonParamBytes, &convertedResult)
	if err != nil {
		return fmt.Errorf("failed to call GetLatestValue over LOOP: %w", err)
	}

	err = s.decodeGLVReturnValue(readIdentifier, convertedResult, returnVal)
	if err != nil {
		return fmt.Errorf("failed to decode GetLatestValue return value: %w", err)
	}

	return nil
}

func (s *loopChainReader) BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	convertedRequest := types.BatchGetLatestValuesRequest{}
	for contract, requestBatch := range request {
		convertedBatch := []types.BatchRead{}
		for _, read := range requestBatch {
			jsonParamBytes, err := json.Marshal(read.Params)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal params: %+w", err)
			}
			convertedBatch = append(convertedBatch, types.BatchRead{
				ReadName:  read.ReadName,
				Params:    jsonParamBytes,
				ReturnVal: &[]byte{},
			})
		}
		convertedRequest[contract] = convertedBatch
	}

	// we always bind before calling query functions, because the LOOP plugin may have restarted.
	err := s.reader.Bind(ctx, s.getBindings())
	if err != nil {
		return nil, fmt.Errorf("failed to re-bind before BatchGetLatestValues: %w", err)
	}

	result, err := s.reader.BatchGetLatestValues(ctx, convertedRequest)
	if err != nil {
		return nil, err
	}

	convertedResult := types.BatchGetLatestValuesResult{}
	for contract, resultBatch := range result {
		requestBatch := request[contract]
		convertedBatch := []types.BatchReadResult{}
		for i, result := range resultBatch {
			read := requestBatch[i]
			resultValue, resultError := result.GetResult()
			convertedResult := types.BatchReadResult{ReadName: result.ReadName}
			if resultError == nil {
				resultPointer := resultValue.(*[]byte)
				err := s.decodeGLVReturnValue(result.ReadName, *resultPointer, read.ReturnVal)
				if err != nil {
					resultError = fmt.Errorf("failed to decode BatchGetLatestValue return value: %w", err)
				}
			}
			convertedResult.SetResult(read.ReturnVal, resultError)
			convertedBatch = append(convertedBatch, convertedResult)
		}
		convertedResult[contract] = convertedBatch
	}

	s.logger.Debugw("BatchGetLatestValues result", "result", convertedResult)

	return convertedResult, nil
}

func (s *loopChainReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	err := s.reader.Bind(ctx, s.getBindings())
	if err != nil {
		return nil, fmt.Errorf("failed to re-bind before BatchGetLatestValues: %w", err)
	}

	convertedExpressions, err := loop.SerializeExpressions(filter.Expressions)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize QueryKey expressions: %w", err)
	}

	convertedFilter := query.KeyFilter{
		Key:         filter.Key,
		Expressions: convertedExpressions,
	}

	sequences, err := s.reader.QueryKey(ctx, contract, convertedFilter, limitAndSort, &[]byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to call QueryKey over LOOP: %w", err)
	}

	for i, sequence := range sequences {
		jsonBytes := sequence.Data.(*[]byte)
		jsonData := map[string]any{}
		err := json.Unmarshal(*jsonBytes, &jsonData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal LOOP sourced JSON event data (`%s`): %w", string(*jsonBytes), err)
		}

		eventData := reflect.New(reflect.TypeOf(sequenceDataType).Elem()).Interface()
		err = codec.DecodeSuiJsonValue(jsonData, eventData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode LOOP sourced event data (`%s`) into a Sui value: %+w", string(*jsonBytes), err)
		}

		sequences[i].Data = eventData
	}

	return sequences, nil
}

func (s *loopChainReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	for _, binding := range bindings {
		s.moduleAddresses[binding.Name] = binding.Address
	}

	return s.reader.Bind(ctx, bindings)
}

func (s *loopChainReader) Unbind(ctx context.Context, bindings []types.BoundContract) error {
	for _, binding := range bindings {
		key := binding.Name
		if _, ok := s.moduleAddresses[key]; !ok {
			return fmt.Errorf("no such binding: %s", key)
		}

		delete(s.moduleAddresses, key)
	}

	// we ignore unbind errors, because if the LOOP plugin restarted, the binding would not exist.
	_ = s.reader.Unbind(ctx, bindings)

	return nil
}

func (s *loopChainReader) getBindings() []types.BoundContract {
	bindings := []types.BoundContract{}

	for name, address := range s.moduleAddresses {
		bindings = append(bindings, types.BoundContract{
			Address: address,
			Name:    name,
		})
	}

	return bindings
}

func (s *loopChainReader) decodeGLVReturnValue(label string, jsonBytes []byte, returnVal any) error {
	var result any
	err := json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s GetLatestValue result (`%s`): %w", label, string(jsonBytes), err)
	}

	err = codec.DecodeSuiJsonValue(result, returnVal)
	if err != nil {
		return fmt.Errorf("failed to decode %s GetLatestValue JSON value (`%s`) to %T: %w", label, string(jsonBytes), returnVal, err)
	}

	return nil
}
