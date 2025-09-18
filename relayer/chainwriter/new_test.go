package chainwriter_test

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	suicodec "github.com/smartcontractkit/chainlink-sui/relayer/codec"
	keystore "github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

func PrivKeyFromHex(hexKey string) ed25519.PrivateKey {
	bytes, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Fatalf("failed to decode hex private key: %v", err)
	}

	switch len(bytes) {
	case ed25519.SeedSize: // 32 bytes
		return ed25519.NewKeyFromSeed(bytes)
	case ed25519.PrivateKeySize: // 64 bytes
		return ed25519.PrivateKey(bytes)
	default:
		log.Fatalf("invalid private key length: got %d, want 32 (seed) or 64 (private key)", len(bytes))
		return nil
	}
}

func TestCWlocal(t *testing.T) {

	logger, err := logger.New()
	if err != nil {
		return
	}

	// Setup new PTB client
	keystoreInstance := keystore.NewTestKeystore(t)

	privKey := PrivKeyFromHex("")
	keystoreInstance.AddKey(privKey)

	pubKey := privKey.Public().(ed25519.PublicKey)
	fmt.Println("PUBKEY ", pubKey)
	pubKeyBytes := []byte(pubKey)

	fmt.Printf("pubKey len=%d bytes: %x\n", len(pubKeyBytes), pubKeyBytes)

	relayerClient, err := client.NewPTBClient(logger, "", nil, 30*time.Second, keystoreInstance, 5, "WaitForEffectsCert")
	if err != nil {
		return
	}

	store := txm.NewTxmStoreImpl(logger)
	conf := txm.DefaultConfigSet

	retryManager := txm.NewDefaultRetryManager(5)
	gasLimit := big.NewInt(30000000)
	gasManager := txm.NewSuiGasManager(logger, relayerClient, *gasLimit, 0)

	txManager, err := txm.NewSuiTxm(logger, relayerClient, keystoreInstance, conf, store, retryManager, gasManager)
	if err != nil {
		return
	}

	nonMutable := false

	chainWriterConfig := config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			"OffRamp": {
				Name:     "offramp",
				ModuleID: "0x123",
				Functions: map[string]*config.ChainWriterFunction{
					// "commit": {
					// 	Name:      "commit",
					// 	PublicKey: pubKeyBytes,
					// 	Params:    []suicodec.SuiFunctionParam{},
					// 	PTBCommands: []config.ChainWriterPTBCommand{{
					// 		Type:      suicodec.SuiPTBCommandMoveCall,
					// 		PackageId: strPtr("0xcea836b828b80a8dc5eee57a7d61feb8e1ce96a6a88d3cf08e393e193b9aabb3"),
					// 		ModuleId:  strPtr("offramp"),
					// 		Function:  strPtr("commit"),
					// 		Params: []suicodec.SuiFunctionParam{
					// 			{
					// 				Name:     "ref",
					// 				Type:     "object_id",
					// 				Required: true,
					// 			},
					// 			{
					// 				Name:     "state",
					// 				Type:     "object_id",
					// 				Required: true,
					// 			},
					// 			{
					// 				Name:      "clock",
					// 				Type:      "object_id",
					// 				Required:  true,
					// 				IsMutable: &nonMutable,
					// 			},
					// 			{
					// 				Name:     "ReportContext",
					// 				Type:     "vector<vector<u8>>",
					// 				Required: true,
					// 			},
					// 			{
					// 				Name:     "Report",
					// 				Type:     "vector<u8>",
					// 				Required: true,
					// 			},
					// 			{
					// 				Name:     "Signatures",
					// 				Type:     "vector<vector<u8>>",
					// 				Required: true,
					// 			},
					// 		},
					// 	}},
					// },
					"Execute": {
						Name:      "execute",
						PublicKey: pubKeyBytes,
						Params:    []suicodec.SuiFunctionParam{},
						PTBCommands: []config.ChainWriterPTBCommand{
							{
								Type:     suicodec.SuiPTBCommandMoveCall,
								ModuleId: strPtr("offramp"),
								Function: strPtr("init_execute"),
								Params: []suicodec.SuiFunctionParam{
									{
										Name:      "ref",
										Type:      "object_id",
										Required:  true,
										IsMutable: &nonMutable,
									},
									{
										Name:     "state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:      "clock",
										Type:      "object_id",
										Required:  true,
										IsMutable: &nonMutable,
									},
									{
										Name:     "ReportContext",
										Type:     "vector<vector<u8>>",
										Required: true,
									},
									{
										Name:     "Report",
										Type:     "vector<u8>",
										Required: true,
									},
									{
										Name:     "token_receiver",
										Type:     "vector<u8>",
										Required: true,
									},
								},
							},
							{
								Type:     suicodec.SuiPTBCommandMoveCall,
								ModuleId: strPtr("offramp"),
								Function: strPtr("finish_execute"),
								Params: []suicodec.SuiFunctionParam{
									{
										Name:     "state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "receiver_params",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: uint16(0),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	chainWriter, err := chainwriter.NewSuiChainWriter(logger, txManager, chainWriterConfig, false)
	if err != nil {
		return
	}

	// This is for commit
	_ = cwConfig.Arguments{
		Args: map[string]any{
			"ReportContext": [][]byte{
				{0, 10, 169, 132, 194, 171, 244, 18, 74, 216, 170, 183, 250, 241, 171, 44, 151, 254, 55, 175, 165, 208, 249, 215, 58, 9, 173, 188, 85, 211, 63, 215},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 61, 49},
			},
			"Report": []byte{
				0, 1, 217, 26, 217, 201, 79, 186, 65, 222, 160, 55, 185, 82, 3, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			"Signatures": [][]byte{
				{
					116, 160, 9, 60, 137, 237, 163, 23, 146, 240, 160, 163, 6, 27, 75, 138,
					17, 211, 228, 195, 55, 27, 56, 55, 137, 252, 47, 17, 170, 228, 5, 108,
					23, 101, 185, 218, 32, 30, 134, 187, 238, 192, 133, 125, 42, 250, 71, 147,
					80, 91, 202, 33, 19, 112, 197, 124, 235, 175, 112, 0, 246, 136, 3, 87,
					56, 236, 213, 119, 220, 246, 128, 108, 47, 153, 193, 154, 85, 44, 168, 123,
					56, 144, 253, 140, 93, 1, 149, 113, 60, 70, 119, 149, 227, 97, 183, 13,
				},
				{
					27, 247, 166, 192, 23, 21, 138, 125, 110, 35, 60, 217, 246, 210, 34, 152,
					7, 190, 169, 32, 177, 99, 193, 165, 72, 153, 101, 20, 103, 220, 90, 192,
					40, 179, 27, 0, 44, 180, 16, 251, 215, 63, 84, 101, 115, 225, 136, 210,
					221, 198, 171, 95, 200, 51, 245, 235, 6, 241, 184, 174, 197, 77, 146, 4,
					114, 32, 22, 97, 215, 110, 73, 43, 225, 206, 180, 190, 21, 68, 193, 244,
					228, 3, 214, 174, 110, 18, 172, 104, 252, 47, 78, 197, 175, 142, 24, 9,
				},
			},
		},
	}

	// This is for exec
	ptbArgsExec := map[string]any{
		"ReportContext": [][]byte{
			{
				0, 10, 147, 52, 37, 122, 128, 196, 95, 63, 5, 140, 211, 39, 185, 153,
				4, 163, 93, 212, 130, 251, 160, 123, 129, 221, 146, 248, 82, 252, 174, 243,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 180,
			},
		},
		"Report": []byte{
			217, 26, 217, 201, 79, 186, 65, 222, 61, 30, 149, 232, 195, 168, 248, 164,
			104, 226, 106, 3, 127, 187, 225, 62, 88, 108, 3, 225, 252, 112, 153, 78,
			229, 42, 130, 102, 110, 166, 98, 139,
			217, 26, 217, 201, 79, 186, 65, 222, 236, 17, 130, 250, 167, 194, 123, 135,
			2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			20, 213, 204, 218, 191, 17, 227, 222, 141, 47, 100, 2, 46, 35, 42, 193,
			128, 1, 184, 172, 172,
			18, 104, 101, 108, 108, 111, 32, 115, 117, 105, 32, 102, 114, 111, 109, 32,
			101, 118, 109,
			240, 213, 109, 42, 135, 243, 90, 176, 203, 242, 31, 21, 233, 255, 2, 32,
			85, 54, 51, 223, 165, 133, 100, 116, 247, 7, 186, 242, 177, 72, 197, 161,
			64, 66, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1,
			153, 221, 87, 219, 246, 46, 233, 79, 160, 146, 109, 23, 192, 39, 176, 153,
			232, 16, 166, 243, 80, 94, 96, 51, 47, 162, 158, 96, 59, 240, 151, 134,
		},
		"Info": map[string]any{
			"AbstractReports": []map[string]any{
				{
					"sourceChainSelector": uint64(16015286601757825753),
					"messages": []map[string]any{
						{
							"data": []byte("hello sui from evm"),
							"extraArgs": []byte{
								33, 234, 76, 169, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 15, 66, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 6, 255, 165, 93, 243, 140, 118, 46, 60, 74, 198, 97, 175,
								68, 29, 25, 218, 91, 210, 161, 191, 190, 29, 99, 41, 194, 76, 193, 11, 75, 177, 25, 190,
							},
							"feeToken": []byte{
								9, 125, 144, 201, 211, 224, 181, 12, 166, 14,
								26, 228, 95, 106, 129, 1, 15, 159, 181, 52,
							},
							"feeTokenAmount": map[string]any{},
							"feeValueJuels":  map[string]any{},
							"header": map[string]any{
								"destChainSelector": uint64(9762610643973837292),
								"messageId": []byte{
									61, 30, 149, 232, 195, 168, 248, 164, 104, 226, 106, 3,
									127, 187, 225, 62, 88, 108, 3, 225, 252, 112, 153, 78,
									229, 42, 130, 102, 110, 166, 98, 139,
								},
								"msgHash": []byte{},
								"nonce":   uint64(0),
								"onRamp": []byte{
									224, 32, 205, 209, 35, 238, 188, 96, 173, 201,
									183, 216, 68, 230, 202, 4, 15, 144, 200, 119,
								},
								"seqNum":              uint64(2),
								"sourceChainSelector": uint64(16015286601757825753),
								"txHash":              "",
							},
							"receiver": []byte{
								240, 213, 109, 42, 135, 243, 90, 176, 203, 242,
								31, 21, 233, 255, 2, 32, 85, 54, 51, 223,
								165, 133, 100, 116, 247, 7, 186, 242, 177, 72,
								197, 161,
							},
							"sender": []byte{
								213, 204, 218, 191, 17, 227, 222, 141, 47, 100,
								2, 46, 35, 42, 193, 128, 1, 184, 172, 172,
							},
							"tokenAmounts": []any{},
						},
					},
					"offchainTokenData": []any{nil},
					"proofs": [][]byte{
						{
							153, 221, 87, 219, 246, 46, 233, 79, 160, 146,
							109, 23, 192, 39, 176, 153, 232, 16, 166, 243,
							80, 94, 96, 51, 47, 162, 158, 96, 59, 240, 151, 134,
						},
					},
					"proofFlagBits": map[string]any{},
				},
			},
			"MerkleRoots": []map[string]any{
				{
					"chain": uint64(16015286601757825753),
					"onRampAddress": []byte{
						224, 32, 205, 209, 35, 238, 188, 96, 173, 201,
						183, 216, 68, 230, 202, 4, 15, 144, 200, 119,
					},
					"seqNumsRange": [2]uint64{1, 2},
					"merkleRoot": []byte{
						99, 248, 65, 155, 77, 64, 18, 119, 196, 3,
						166, 56, 168, 12, 184, 198, 93, 38, 236, 105,
						135, 79, 97, 77, 174, 241, 69, 201, 251, 19, 249, 29,
					},
				},
			},
		},
		"ExtraData": map[string]any{
			"ExtraArgsDecoded": map[string]any{
				"allowOutOfOrderExecution": true,
				"gasLimit":                 uint64(1000000), // ensure *big.Int
				"receiverObjectIds": [][]byte{
					{
						0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 6,
					},
					{
						255, 165, 93, 243, 140, 118, 46, 60,
						74, 198, 97, 175, 68, 29, 25, 218,
						91, 210, 161, 191, 190, 29, 99, 41,
						194, 76, 193, 11, 75, 177, 25, 190,
					},
				},
				"tokenReceiver": []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			"DestExecDataDecoded": []map[string]any{},
		},
	}

	ctx := context.Background()
	txId := "offramp_test"
	err = chainWriter.SubmitTransaction(ctx,
		"OffRamp",
		"Execute",
		&ptbArgsExec,
		txId,
		"0xcea836b828b80a8dc5eee57a7d61feb8e1ce96a6a88d3cf08e393e193b9aabb3",
		&commonTypes.TxMeta{},
		nil,
	)
	if err != nil {
		return
	}

	tx, err := txManager.TransactionRepository.GetTransaction(txId)
	if err != nil {
		fmt.Println("FAILED GETTING TX: ", tx)
		return
	}

	payload := client.TransactionBlockRequest{
		TxBytes:    tx.Payload,
		Signatures: tx.Signatures,
		Options: client.TransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
		},
		RequestType: tx.RequestType,
	}

	resp, err := txManager.SuiGateway.SendTransaction(ctx, payload)
	if err != nil {
		fmt.Println("FAILED TO SEND TX: ", err)
		return
	}

	fmt.Println("RES: ", resp)

	// require.Eventually(t, func() bool {
	// 	status, statusErr := chainWriter.GetTransactionStatus(ctx, txId)
	// 	if statusErr != nil {
	// 		return false
	// 	}
	// 	return status == commonTypes.Finalized
	// }, 5*time.Second, 1*time.Second, "Transaction final state not reached")
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}

func TestPTBClientNew(t *testing.T) {
	ctx := t.Context()
	// var result []models.SuiObjectResponse

	client := sui.NewSuiClient("")

	ownedObjectsReq := models.SuiXGetOwnedObjectsRequest{
		Address: "0x479161ba654faab2eeb6a08c9df5c4ebceb1dda7e1e7d17af71634da3f65574d",
		Query: models.SuiObjectResponseQuery{
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
				ShowType:    true,
				ShowOwner:   true,
			},
		},
		Limit: uint64(50),
	}

	res, err := client.SuiXGetOwnedObjects(ctx, ownedObjectsReq)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	ownedObjects := res.Data

	pointers := []string{"state_object::CCIPObjectRefPointer"}
	pointersValuesMap := make(map[string]map[string]any)

	// check each returned object
	for _, ownedObject := range ownedObjects {
		// fmt.Println("NEW: ", ownedObject.Data.Content.Fields)
		// check if it matches any of the pointers
		fmt.Println(ownedObject.Data.Type)
		for _, pointer := range pointers {
			if ownedObject.Data.Type != "" && strings.Contains(ownedObject.Data.Type, pointer) {
				pointersValuesMap[pointer] = ownedObject.Data.Content.Fields
			}
		}
	}

	fmt.Println("pointerValueMap: ", pointersValuesMap)
	// fmt.Println("RESULT: ", res[0].Data)
}

// {"level":"info","ts":"2025-08-06T15:27:49.375Z","logger":"Sui.2.RelayerService.PluginRelayerClient.PluginSui.SuiRelayer.SuiTxm","caller":"loop/logger.go:155","msg":"Broadcasting transaction","version":"2.26.0@fd25ca0","chainID":"2","chain":"sui","caller":"txm/broadcaster.go:76","txID":"OffRamp-0x5281e9647db50ef820b88c5cbcb0d3f16c78ff82090f6fe6afa36942a51b5aaf-f4a1a065-ed8a-41c2-98b5-cccfc07eb333","payload":{"Attempt":0,"Digest":"","Functions":[],"LastUpdatedAt":1754494069,"Metadata":{"GasLimit":null,"WorkflowExecutionID":null},"Payload":"AAAAAGv9jalQVaDETyAe9rW3B0enlvSzf8PF0mTUvl+w5HYzAXwlqdTiQBBN079EVaYx8Qbzm2yyadYDb68WmrN1ugHb6X1yHAAAAAAgP4gOupOpdn9Qp/82yDsyB3p0OZNogVS3rplGFyr/DsNr/Y2pUFWgxE8gHva1twdHp5b0s3/DxdJk1L5fsOR2M+gDAAAAAAAAAMLrCwAAAAAA","RequestType":"WaitForEffectsCert","Sender":"0x6bfd8da95055a0c44f201ef6b5b70747a796f4b37fc3c5d264d4be5fb0e47633","Signatures":["APk612A5jG/vedhhVdlQ6+UHwYkkXXsQpXsrhiViWu+T+1Xqe3C+9Qn/2fNPFxhoVd9NCTjnVF+5aPgDQUmElQYNXJC99NVmN137ACHFnSFW6WyEYDASs1q3ZhD52Nto3A=="],"State":0,"Timestamp":1754494069,"TransactionID":"OffRamp-0x5281e9647db50ef820b88c5cbcb0d3f16c78ff82090f6fe6afa36942a51b5aaf-f4a1a065-ed8a-41c2-98b5-cccfc07eb333","TxError":null},"timestamp":"2025-08-06T15:27:49.375Z"}

// new
// {"level":"info","ts":"2025-08-06T17:30:03.063Z","logger":"Sui.2.RelayerService.PluginRelayerClient.PluginSui.SuiRelayer.SuiTxm","caller":"loop/logger.go:155","msg":"Broadcasting transaction","version":"2.26.0@049e4dd","caller":"txm/broadcaster.go:76","txID":"OffRamp-0x5281e9647db50ef820b88c5cbcb0d3f16c78ff82090f6fe6afa36942a51b5aaf-25ba3607-a9ee-469a-bf34-387e3065acb4","payload":{"Attempt":0,"Digest":"","Functions":[],"LastUpdatedAt":1754501403,"Metadata":{"GasLimit":null,"WorkflowExecutionID":null},"Payload":"AAAAANypa19lw2Wqrx5tnxXvmra482k4OlQZWzeOByZk0z1MAfOi0PnCqWolSFlCbWAHCPcXtq62kvc989W38K5YldFt531yHAAAAAAgk8WXrX+q15hYuZ+0cvxWYcdQKyGWoCRyi+LRSccJtXLcqWtfZcNlqq8ebZ8V75q2uPNpODpUGVs3jgcmZNM9TADC6wsAAAAAAA==","RequestType":"WaitForEffectsCert","Sender":"0xdca96b5f65c365aaaf1e6d9f15ef9ab6b8f369383a54195b378e072664d33d4c","Signatures":["AIAHwuasda7U7S29v4489hLs/hY+4AiUrnMdSw3Jz6UIHfzZVF+LdRK1qKR6I3nZX6B4k1P+TtCstDtGlpeFRAnsnaMm21g0v3AZR5xeciDdLfoa5oU1zzGJJqAaMhosJw=="],"State":0,"Timestamp":1754501403,"TransactionID":"OffRamp-0x5281e9647db50ef820b88c5cbcb0d3f16c78ff82090f6fe6afa36942a51b5aaf-25ba3607-a9ee-469a-bf34-387e3065acb4","TxError":null},"chainID":"2","chain":"sui","timestamp":"2025-08-06T17:30:03.062Z"}
func TestVerifyMessage(t *testing.T) {

	// sigBase64 := "APk612A5jG/vedhhVdlQ6+UHwYkkXXsQpXsrhiViWu+T+1Xqe3C+9Qn/2fNPFxhoVd9NCTjnVF+5aPgDQUmElQYNXJC99NVmN137ACHFnSFW6WyEYDASs1q3ZhD52Nto3A=="
	messageWithIntentBase64 := "AAAAANypa19lw2Wqrx5tnxXvmra482k4OlQZWzeOByZk0z1MAfOi0PnCqWolSFlCbWAHCPcXtq62kvc989W38K5YldFt531yHAAAAAAgk8WXrX+q15hYuZ+0cvxWYcdQKyGWoCRyi+LRSccJtXLcqWtfZcNlqq8ebZ8V75q2uPNpODpUGVs3jgcmZNM9TADC6wsAAAAAAA=="

	rawBcs, err := base64.StdEncoding.DecodeString(messageWithIntentBase64)
	if err != nil {
		return
	}

	fmt.Println(rawBcs)
	// const (
	// 	sigSize    = ed25519.SignatureSize // 64
	// 	pubKeySize = ed25519.PublicKeySize // 32
	// )

	// decodedMsg, err := base64.StdEncoding.DecodeString(messageWithIntentBase64)
	// require.NoError(t, err)

	// // decode **unpadded** signature
	// decodedSig, err := base64.RawStdEncoding.DecodeString(sigBase64)
	// require.NoError(t, err)

	// // sanity
	// want := 1 + sigSize + pubKeySize
	// require.Equal(t, want, len(decodedSig), "serialized signature length")

	// scheme := decodedSig[0]
	// sig := decodedSig[1 : 1+sigSize]
	// pubKey := decodedSig[1+sigSize : 1+sigSize+pubKeySize]

	// // already has intent
	// digest := blake2b.Sum256(decodedMsg)
	// valid := ed25519.Verify(pubKey, digest[:], sig)

	// t.Logf("scheme=0x%02x, valid=%v", scheme, valid)
	// require.True(t, valid, "signature must verify")

	// signer, pass, err := models.VerifyTransaction(message, sig)
	// require.NoError(t, err)

	// fmt.Println(signer, pass)
}
