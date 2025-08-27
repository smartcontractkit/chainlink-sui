package bind

import (
	"context"
	"errors"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

// Fetches every coin owned by the address. Looks for a SUI object, returns the first it finds
func FetchDefaultGasCoinRef(ctx context.Context, client sui.ISuiAPI, address string) (*models.SuiObjectRef, error) {
	suiCoins, err := fetchOwnedSuiCoins(ctx, client, address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch owned SUI coins: %w", err)
	}
	if len(suiCoins) == 0 {
		return nil, fmt.Errorf("no SUI coins found for address: %s", address)
	}

	return suiCoins[0], nil
}

func ToSuiObjectRef(ctx context.Context, client sui.ISuiAPI, objectId string, address string) (*models.SuiObjectRef, error) {
	// Normalize the object ID
	normalizedObjId, err := bindutils.ConvertAddressToString(objectId)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID %v: %w", objectId, err)
	}

	refs, err := fetchOwnedSuiCoins(ctx, client, address)
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

func fetchOwnedSuiCoins(ctx context.Context, client sui.ISuiAPI, address string) ([]*models.SuiObjectRef, error) {
	// Normalize address
	normalizedAddr, err := bindutils.ConvertAddressToString(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address %v: %w", address, err)
	}

	// Create request for getting all coins
	var limit uint64 = 50
	req := models.SuiXGetAllCoinsRequest{
		Owner: normalizedAddr,
		Limit: limit,
	}

	coins, err := client.SuiXGetAllCoins(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get default Gas coins for address: %s: %w", address, err)
	}

	if len(coins.Data) == 0 {
		return nil, fmt.Errorf("no coin data found for signer: %s", address)
	}

	var coinRefs []*models.SuiObjectRef
	for _, coinData := range coins.Data {
		if isSuiCoin(coinData) {
			version, err := parseVersionString(coinData.Version)
			if err != nil {
				return nil, fmt.Errorf("failed to parse coin version: %w", err)
			}

			coinRef := &models.SuiObjectRef{
				ObjectId: coinData.CoinObjectId,
				Version:  version,
				Digest:   coinData.Digest,
			}
			coinRefs = append(coinRefs, coinRef)
		}
	}

	return coinRefs, nil
}

func isSuiCoin(coin models.CoinData) bool {
	return coin.CoinType == "0x2::sui::SUI"
}
