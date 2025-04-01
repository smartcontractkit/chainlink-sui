package bind

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	sui_pattokan "github.com/pattonkan/sui-go/sui"
)

// Utilities around the PTB and its types
func ToSuiAddress(address string) (*sui_pattokan.Address, error) {
	return sui_pattokan.AddressFromHex(address)
}

// Fetches every coin owned by the address. Looks for a SUI object, returns the first it finds
func FetchDefaultGasCoinRef(ctx context.Context, client sui.ISuiAPI, address string) (*sui_pattokan.ObjectRef, error) {
	suiCoins, err := fetchOwnedSuiCoins(ctx, client, address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch owned SUI coins: %w", err)
	}
	if len(suiCoins) == 0 {
		return nil, fmt.Errorf("no SUI coins found for address: %s", address)
	}

	return suiCoins[0], nil
}

func ToSuiObjectRef(ctx context.Context, client sui.ISuiAPI, objectId string, address string) (*sui_pattokan.ObjectRef, error) {
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

func fetchOwnedSuiCoins(ctx context.Context, client sui.ISuiAPI, address string) ([]*sui_pattokan.ObjectRef, error) {
	coin, err := client.SuiXGetAllCoins(ctx, models.SuiXGetAllCoinsRequest{
		Owner: address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get default Gas coins: %w", err)
	}

	if len(coin.Data) == 0 {
		return nil, fmt.Errorf("no coin data found for signer: %s", address)
	}

	var coinRefs []*sui_pattokan.ObjectRef
	for _, data := range coin.Data {
		if isSuiCoin(data) {
			coinAddress, err := ToSuiAddress(data.CoinObjectId)
			if err != nil {
				return nil, fmt.Errorf("failed to get coin address: %w", err)
			}

			version, err := strconv.ParseUint(data.Version, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid version: %w", err)
			}

			digest, err := sui_pattokan.NewDigest(data.Digest)
			if err != nil {
				return nil, fmt.Errorf("invalid coin digest: %w", err)
			}
			coinRefs = append(coinRefs, &sui_pattokan.ObjectRef{
				ObjectId: coinAddress,
				Version:  version,
				Digest:   digest,
			})
		}
	}

	return coinRefs, nil
}

func isSuiCoin(coin models.CoinData) bool {
	return coin.CoinType == "0x2::sui::SUI"
}
