package operations

import (
	"fmt"

	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type OpTxInput[I any] struct {
	GasOverrides bind.TxOpts
	Input        I
}

type OpTxResult[O any] struct {
	Digest    string
	PackageId string
	Objects   O
}

type OpTxDeps struct {
	Client suiclient.ClientImpl
	Signer rel.SuiSigner
	// We could have some logic to modify the gas based on input
	GetTxOpts func() bind.TxOpts
}

func NewSuiOperationName(pkg string, module string, action string) string {
	return fmt.Sprintf("sui-%s-%s-%s", pkg, module, action)
}
