package operations

import (
	"fmt"

	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type OpTxInput[I any] struct {
	Input I
}

type OpTxResult[O any] struct {
	Digest    string
	PackageId string
	Objects   O
}

type OpTxDeps struct {
	Client sui.ISuiAPI
	Signer rel.SuiSigner
	// We could have some logic to modify the gas based on input
	GetCallOpts func() *bind.CallOpts
}

func NewSuiOperationName(pkg string, module string, action string) string {
	return fmt.Sprintf("sui-%s-%s-%s", pkg, module, action)
}
