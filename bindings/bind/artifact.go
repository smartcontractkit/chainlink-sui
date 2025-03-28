package bind

import (
	"encoding/json"
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
