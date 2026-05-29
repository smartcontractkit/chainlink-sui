package bind

import (
	"context"
	"errors"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const suiCoinType = "0x2::coin::Coin<0x2::sui::SUI>"

func FetchDefaultGasCoinRef(ctx context.Context, chainClient client.BindingsClient, address string) (*models.SuiObjectRef, error) {
	suiCoins, err := fetchOwnedSuiCoins(ctx, chainClient, address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch owned SUI coins: %w", err)
	}
	if len(suiCoins) == 0 {
		return nil, fmt.Errorf("no SUI coins found for address: %s", address)
	}

	return suiCoins[0], nil
}

func ToSuiObjectRef(ctx context.Context, chainClient client.BindingsClient, objectId string, address string) (*models.SuiObjectRef, error) {
	normalizedObjId, err := bindutils.ConvertAddressToString(objectId)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID %v: %w", objectId, err)
	}

	refs, err := fetchOwnedSuiCoins(ctx, chainClient, address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch owned SUI coins: %w", err)
	}

	for _, ref := range refs {
		if ref.ObjectId == normalizedObjId {
			return ref, nil
		}
	}

	return nil, errors.New("gas object ID not found in SUI owned coins")
}

func fetchOwnedSuiCoins(ctx context.Context, chainClient client.BindingsClient, address string) ([]*models.SuiObjectRef, error) {
	normalizedAddr, err := bindutils.ConvertAddressToString(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address %v: %w", address, err)
	}

	coins, err := chainClient.QueryCoinsByAddress(ctx, normalizedAddr, suiCoinType)
	if err != nil {
		return nil, fmt.Errorf("failed to get default Gas coins for address: %s: %w", address, err)
	}

	if len(coins) == 0 {
		return nil, fmt.Errorf("no coin data found for signer: %s", address)
	}

	coinRefs := make([]*models.SuiObjectRef, 0, len(coins))
	for _, coin := range coins {
		coinRefs = append(coinRefs, mapGrpcCoinToObjectRef(coin))
	}

	return coinRefs, nil
}
