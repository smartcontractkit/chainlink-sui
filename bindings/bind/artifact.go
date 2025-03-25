package bind

import (
	"encoding/json"

	"github.com/block-vision/sui-go-sdk/models"
)

type PackageArtifact struct {
	Modules      []string `json:"modules"`
	Dependencies []string `json:"dependencies"`
	Digest       []int    `json:"digest"`
}

func ToArtifact(artifactJSON string) (PackageArtifact, error) {
	var artifact PackageArtifact
	if err := json.Unmarshal([]byte(artifactJSON), &artifact); err != nil {
		return PackageArtifact{}, err
	}
	return artifact, nil
}

func BuildPublishRequest(artifact PackageArtifact, opts TxOpts, signer string) models.PublishRequest {
	return models.PublishRequest{
		Sender:          signer,
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
		Gas:             &opts.GasObject,
		GasBudget:       opts.GasBudget,
	}
}
