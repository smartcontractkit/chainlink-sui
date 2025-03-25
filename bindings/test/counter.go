package counter

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// This should be auto-generated when compiling, same as the bindings
const CounterJSON = `{
	"modules": [
		"oRzrCwYAAAAKAQAIAggMAxQfBDMCBTUkB1mBAQjaAUAKmgIIDKICVg34AgIAAwEKAQwBDQAADAABAgQAAwECAAAIAAEAAAYCAQAABwMBAAAEBAUAAQkABgACCwgBAQgFBwEHCAIAAgcIAAcIAgQHCAADAwcIAgEGCAABAwEIAQEIAAEJAAdDb3VudGVyCVR4Q29udGV4dANVSUQHY291bnRlcglnZXRfY291bnQCaWQJaW5jcmVtZW50DmluY3JlbWVudF9tdWx0CmluaXRpYWxpemUDbmV3Bm9iamVjdAxzaGFyZV9vYmplY3QIdHJhbnNmZXIKdHhfY29udGV4dAV2YWx1ZQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAgIFCAEOAwABBAABBgsAEQQGAAAAAAAAAAASADgAAgEBBAABCQoAEAAUBgEAAAAAAAAAFgsADwAVAgIBBAABCwoAEAAUCwELAhgWCwAPABUCAwEEAAEECwAQABQCAAEA",
		"oRzrCwYAAAAJAQAMAgwcAyg5BGEIBWk0B50BowIIwANgCqAEGgy6BHQABwEXAhACFQIZAhoAAgMAAAADAAAEAwAAAQgAAQMHAAMGBAAFBQIAABIAAQAADgIBAAANAwMAAAsEBAAADAUFAAAKBgYAAAgHBwAACQgIAAIPCwEBAwMTAAkABBYLAQEICgoIDAgNCA4BBwgGAAUGCAMDCAQKAgcIBgEDAQ8CDgMBCAQBCgIBCgoCAQgFAQgDAQkAAQgAAQgBAQgCEERvdWJsZVZhbHVlRXZlbnQKRXZlbnRTdG9yZRBTaW5nbGVWYWx1ZUV2ZW50BlN0cmluZxBUcmlwbGVWYWx1ZUV2ZW50CVR4Q29udGV4dANVSUQEZWNobxBlY2hvX2J5dGVfdmVjdG9yF2VjaG9fYnl0ZV92ZWN0b3JfdmVjdG9yC2VjaG9fc3RyaW5nCWVjaG9fdTI1NhJlY2hvX3UzMl91NjRfdHVwbGUIZWNob191NjQQZWNob193aXRoX2V2ZW50cwRlbWl0BWV2ZW50AmlkBGluaXQDbmV3Bm51bWJlcgZvYmplY3QMc2hhcmVfb2JqZWN0BnN0cmluZwR0ZXh0CHRyYW5zZmVyCnR4X2NvbnRleHQFdmFsdWUGdmFsdWVzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAIBGwMBAgIUAxgIBAICARwKCgIDAgERCAUAAAAAAQULABEJEgM4AAIBAQQACBAKARIAOAELAQsCEgE4AkAHAAAAAAAAAAAMBQ0FCwNEBwsFEgI4AwICAQAAAQILAAIDAQAAAQILAAIEAQAAAQMLAAsBAgUBAAABAgsAAgYBAAABAgsAAgcBAAABAgsAAgA="
	],
	"dependencies": [
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000000000000000000000000000002"
	],
	"digest": [
		205, 51, 124, 253, 124, 136, 134, 192, 240, 141, 61, 232, 201, 214, 35, 229,
		116, 23, 208, 55, 0, 255, 80, 13, 43, 253, 114, 188, 176, 168, 115, 210
	]
}`

type Counter struct {
	objectID bind.ObjectID
	client   sui.ISuiAPI
	signer   signer.Signer
	// TODO: find relevant information to store
}

func PublishCounter(ctx context.Context, opts bind.TxOpts, signer signer.Signer, client sui.ISuiAPI) (*Counter, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.ToArtifact(CounterJSON)
	if err != nil {
		return nil, nil, err
	}

	req := bind.BuildPublishRequest(artifact, opts, string(signer.PubKey))
	objId, tx, err := bind.PublishPackage(signer, client, req)
	if err != nil {
		return nil, nil, err
	}

	return NewCounter(objId, signer, client), tx, nil
}

func NewCounter(objId string, signer signer.Signer, client sui.ISuiAPI) *Counter {
	return &Counter{
		objectID: objId,
		client:   client,
		signer:   signer,
	}
}

func (c *Counter) Increment(ctx context.Context, opts bind.TxOpts) (*models.SuiTransactionBlockResponse, error) {
	// calls increment function
	tx, err := bind.CallMethod(ctx, c.signer, c.client, models.MoveCallRequest{
		PackageObjectId: c.objectID,
		Signer:          string(c.signer.PubKey),
		Module:          "Counter",
		Function:        "Increment",
		Arguments:       []any{},
		Gas:             &opts.GasObject,
		GasBudget:       opts.GasBudget,
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *Counter) Initialize(ctx context.Context, opts bind.TxOpts) (*models.SuiTransactionBlockResponse, error) {
	// calls increment function
	tx, err := bind.CallMethod(ctx, c.signer, c.client, models.MoveCallRequest{
		PackageObjectId: c.objectID,
		Signer:          string(c.signer.PubKey),
		Module:          "Counter",
		Function:        "Initialize",
		Arguments:       []any{},
		Gas:             &opts.GasObject,
		GasBudget:       opts.GasBudget,
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// TODO: Build the rest of funtions
func (c *Counter) IncrementMult() (string, error) {
	// calls increment_mult function
	return "", nil
}
