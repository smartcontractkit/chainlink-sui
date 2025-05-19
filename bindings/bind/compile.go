package bind

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type PackageManifest struct {
	Package      any               `toml:"package"`
	Addresses    map[string]string `toml:"addresses"`
	Dependencies any               `toml:"dependencies"`
	DevAddresses any               `toml:"dev-addresses"`
}

func CompilePackage(packageName contracts.Package, namedAddresses map[string]string) (PackageArtifact, error) {
	packageDir, ok := contracts.Contracts[packageName]
	if !ok {
		return PackageArtifact{}, fmt.Errorf("package %s not found", packageName)
	}

	// Create a random temporary directory path
	dstDir, err := os.MkdirTemp("", "sui-temp-*")
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dstDir)

	srcDir := filepath.Join(".")
	dstRoot := filepath.Join(dstDir, "contracts")
	packageRoot := filepath.Join(dstRoot, packageDir)

	// Copy the (embedded) source directories into the temporary directory root.
	// We need to copy all contracts (not just the specified package) as different packages might depend on each other.
	err = writeEFS(contracts.Embed, srcDir, dstRoot)
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to copy embedded files to %q: %w", dstRoot, err)
	}

	// Update the TOML manifest with the named addresses
	tomlPath := filepath.Join(packageRoot, "Move.toml")
	docBytes, err := os.ReadFile(tomlPath)
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to read TOML file %q: %w", tomlPath, err)
	}

	var manifest PackageManifest
	err = toml.Unmarshal(docBytes, &manifest)
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to parse Move TOML Manifest: %w", err)
	}
	manifest.Addresses = namedAddresses

	// Write the updated TOML file back to the temporary directory
	b, err := toml.Marshal(manifest)
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to marshal TOML file %q: %w", tomlPath, err)
	}
	//nolint:mnd
	err = os.WriteFile(tomlPath, b, 0600)
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to write TOML file %q: %w", tomlPath, err)
	}

	args := []string{
		"move", "build",
		"--dump-bytecode-as-base64",
		"--ignore-chain",
	}

	cmd := exec.Command("sui", args...)
	cmd.Dir = packageRoot // Command is run in the temporary destination directory
	output, err := cmd.Output()
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to run sui move build: %w", err)
	}

	contractArtifact, err := ToArtifact(string(output))
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to parse contract artifact: %w", err)
	}

	return contractArtifact, nil
}

func writeEFS(efs embed.FS, srcDir, dstDir string) error {
	return fs.WalkDir(efs, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, path)

		if d.IsDir() {
			e := os.MkdirAll(dstPath, os.ModePerm)
			if e != nil {
				return fmt.Errorf("failed to create directory %q: %w", dstPath, e)
			}

			return nil
		}

		srcFile, err := efs.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open src file %q: %w", path, err)
		}
		defer func(srcFile fs.File) {
			_ = srcFile.Close()
		}(srcFile)

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("failed to create dst file %q: %w", dstPath, err)
		}
		defer func(dstFile *os.File) {
			_ = dstFile.Close()
		}(dstFile)

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy %q to %q: %w", path, dstPath, err)
		}

		return nil
	})
}
