package bind

import (
	"context"
	"errors"
	"fmt"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
)

// Fetches every coin owned by the address. Looks for a SUI object, returns the first it finds
func FetchDefaultGasCoinRef(ctx context.Context, client suiclient.ClientImpl, address string) (*sui.ObjectRef, error) {
	suiCoins, err := fetchOwnedSuiCoins(ctx, client, address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch owned SUI coins: %w", err)
	}
	if len(suiCoins) == 0 {
		return nil, fmt.Errorf("no SUI coins found for address: %s", address)
	}

	return suiCoins[0], nil
}

func ToSuiObjectRef(ctx context.Context, client suiclient.ClientImpl, objectId string, address string) (*sui.ObjectRef, error) {
	// Convert the object ID to a Sui address
	suiAddress, err := ToSuiAddress(objectId)
	if err != nil {
		return nil, fmt.Errorf("failed to convert object ID to address: %w", err)
	}

	refs, err := fetchOwnedSuiCoins(ctx, client, address)
	if err != nil {
		return nil, fmt.Errorf("failed to convert object ID to address: %w", err)
	}

	for _, ref := range refs {
		if ref.ObjectId == suiAddress {
			return ref, nil
		}
	}

	return nil, errors.New("gas object ID not found in SUI owned coins")
}

func fetchOwnedSuiCoins(ctx context.Context, client suiclient.ClientImpl, address string) ([]*sui.ObjectRef, error) {
	suiAddress, err := ToSuiAddress(address)
	if err != nil {
		return nil, err
	}
	coin, err := client.GetAllCoins(ctx, &suiclient.GetAllCoinsRequest{
		Owner: suiAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get default Gas coins: %w", err)
	}

	if len(coin.Data) == 0 {
		return nil, fmt.Errorf("no coin data found for signer: %s", address)
	}

	var coinRefs []*sui.ObjectRef
	for _, coinData := range coin.Data {
		if isSuiCoin(*coinData) {
			coinRefs = append(coinRefs, coinData.Ref())
		}
	}

	return coinRefs, nil
}

func isSuiCoin(coin suiclient.Coin) bool {
	return coin.CoinType == "0x2::sui::SUI"
}
