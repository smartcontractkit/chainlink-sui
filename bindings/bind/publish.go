package bind

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/fardream/go-bcs/bcs"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	sui_pattokan "github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
)

type ObjectID = string

// TODO: Should make use of opts
func PublishPackage(
	ctx context.Context,
	opts TxOpts,
	// TODO: Replace by a Signer common interface
	signer signer.Signer,
	client sui.ISuiAPI,
	req models.PublishRequest,
) (ObjectID, *models.SuiTransactionBlockResponse, error) {
	var modules [][]byte
	for _, encodedModule := range req.CompiledModules {
		decodedModule := decodeBase64(encodedModule)
		modules = append(modules, decodedModule)
	}

	var deps []*sui_pattokan.Address
	for _, dep := range req.Dependencies {
		suiAddressDep, err := ToSuiAddress(dep)
		if err != nil {
			// TODO: Give better error desc
			return "", nil, err
		}
		deps = append(deps, suiAddressDep)
	}

	signerAddress, err := ToSuiAddress(signer.Address)
	if err != nil {
		// TODO: Give better error desc
		return "", nil, err
	}

	// Construct the Transaction using PTB
	ptb := suiptb.NewTransactionDataTransactionBuilder()
	arg := ptb.PublishUpgradeable(modules, deps)
	// The program object is transferred to the signer address once deployed
	recArg, err := ptb.Pure(signerAddress)
	ptb.Command(suiptb.Command{
		TransferObjects: &suiptb.ProgrammableTransferObjects{
			Objects: []suiptb.Argument{arg},
			Address: recArg,
		}})

	// Finish transaction details and encode it
	txBytes, err := FinishTransactionFromBuilder(ctx, ptb, signer.Address, client)
	if err != nil {
		return "", nil, err
	}

	// Sign and send Transaction
	tx, err := SignAndSendTx(ctx, signer, client, txBytes)
	if err != nil {
		msg := fmt.Errorf("failed to execute tx when publishing: %v", err)
		return "", nil, msg
	}

	if tx.Effects.Status.Status == "failure" {
		return "", nil, fmt.Errorf("transaction failed: %v", tx.Effects.Status.Error)
	}

	// Find the object ID from the transaction
	objectId, err := findObjectIdFromPublishTx(*tx)
	if err != nil {
		return "", nil, err
	}

	return objectId, tx, err
}

func findObjectIdFromPublishTx(tx models.SuiTransactionBlockResponse) (string, error) {
	for _, change := range tx.ObjectChanges {
		if change.Type == "published" {
			return change.ObjectId, nil
		}
	}

	return "", errors.New("object ID not found in transaction")
}

func decodeBase64(encoded string) []byte {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatalf("Error decoding base64: %v", err)
	}
	return data
}

func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Fetches every coin owned by the address. Assumes the first one in the list is the SUI coin
// TODO: Should check this is the coin we want to use, maybe through args
func getGasCoinData(ctx context.Context, address string, client sui.ISuiAPI) (*sui_pattokan.ObjectRef, error) {
	coin, err := client.SuiXGetAllCoins(ctx, models.SuiXGetAllCoinsRequest{
		Owner: address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get coins: %v", err)
	}

	if len(coin.Data) == 0 {
		return nil, fmt.Errorf("no coin data found for signer: %s", address)
	}

	coinAddress, err := ToSuiAddress(coin.Data[0].CoinObjectId)
	if err != nil {
		return nil, fmt.Errorf("failed to get coin address: %v", err)
	}

	version, err := strconv.ParseUint(coin.Data[0].Version, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid version: %v", err)
	}

	digest, err := sui_pattokan.NewDigest(coin.Data[0].Digest)
	if err != nil {
		return nil, fmt.Errorf("invalid coin digest: %v", err)
	}

	return &sui_pattokan.ObjectRef{
		ObjectId: coinAddress,
		Version:  version,
		Digest:   digest,
	}, nil
}

func FinishTransactionFromBuilder(ctx context.Context, ptb *suiptb.ProgrammableTransactionBuilder, signer string, client sui.ISuiAPI) ([]byte, error) {
	pt := ptb.Finish()

	address, err := ToSuiAddress(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signer address")
	}

	coinData, err := getGasCoinData(ctx, signer, client)
	if err != nil {
		return nil, err
	}

	// TODO: Should be configurable with txopts
	gasBudget := uint64(20000000)
	txData := suiptb.NewTransactionData(
		address,
		pt,
		[]*sui_pattokan.ObjectRef{coinData},
		gasBudget,
		suiclient.DefaultGasPrice,
	)

	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}
