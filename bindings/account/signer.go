package bindings

import (
	sui "github.com/block-vision/sui-go-sdk"
	"github.com/block-vision/sui-go-sdk/signer"
)

func NewSigner(privatekey string) *signer.Signer {
	c, err := sui.NewRpcClient("http://localhost:8080")
	r, err := c.SignedTransaction
	s, err := signer.NewSignertWithMnemonic(privatekey)
	if err != nil {
		panic(err)
	}
	return s
}
