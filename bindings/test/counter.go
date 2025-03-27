package counter

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// This should be auto-generated when compiling, same as the bindings
// Currently built with `sui move build --dump-bytecode-as-base64 --with-unpublished-dependencies`
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

func PublishCounter(ctx context.Context, opts bind.TxOpts, signer signer.Signer, client sui.ISuiAPI) (*Counter, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.ToArtifact(CounterJSON)
	if err != nil {
		return nil, nil, err
	}

	req := bind.BuildPublishRequest(artifact, opts, signer.Address)
	objId, tx, err := bind.PublishPackage(ctx, signer, client, req)
	if err != nil {
		return nil, nil, err
	}

	return NewCounter(objId, client), tx, nil
}

type ICounter interface {
	Increment(objectId string) (bind.IMethod, error)
	IncrementMult(objectId string, a, b uint64) (bind.IMethod, error)
}

type Counter struct {
	objectID bind.ObjectID
	client   sui.ISuiAPI
	// TODO: find relevant information to store
	// Might need package Id instead of objectId? Or both?
}

var _ ICounter = (*Counter)(nil)

func NewCounter(objId string, client sui.ISuiAPI) *Counter {
	return &Counter{
		objectID: objId,
		// TODO: Remove after txs are built locally
		client: client,
	}
}

func (c *Counter) EncodeIncrement(counterObjectId string) (encodedArgs []any, err error) {
	return bind.Encode(
		[]string{
			"address",
		},
		[]any{counterObjectId},
	)
}

func (c *Counter) EncodeIncrementMult(counterObjectId string, a, b uint64) (encodedArgs []any, err error) {
	return bind.Encode(
		[]string{
			"address",
			"u64",
			"u64",
		},
		[]any{counterObjectId, a, b},
	)
}

func (c *Counter) Increment(counterObjectId string) (bind.IMethod, error) {
	build := func(opts bind.TxOpts, signer string) (models.TxnMetaData, error) {
		payload, err := c.EncodeIncrement(counterObjectId)
		if err != nil {
			return models.TxnMetaData{}, err
		}
		req := bind.BuildCallRequest(opts, signer, c.objectID, "Counter", "increment", payload)
		// TODO: This should not use client (should be a local computation), and shouldn't need context
		unsignedTx, err := c.client.MoveCall(context.Background(), req)
		return unsignedTx, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build)), nil
}

func (c *Counter) IncrementMult(counterObjectId string, a, b uint64) (bind.IMethod, error) {
	build := func(opts bind.TxOpts, signer string) (models.TxnMetaData, error) {
		payload, err := c.EncodeIncrementMult(counterObjectId, a, b)
		if err != nil {
			return models.TxnMetaData{}, err
		}
		req := bind.BuildCallRequest(opts, signer, c.objectID, "Counter", "increment_mult", payload)
		// TODO: This should not use client (should be a local computation), and shouldn't need context
		unsignedTx, err := c.client.MoveCall(context.Background(), req)
		return unsignedTx, nil
	}
	return bind.NewMethod(build, bind.MakeExecute(build)), nil
}
