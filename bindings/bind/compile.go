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
		return PackageArtifact{}, fmt.Errorf("unknown package: %s", packageName)
	}

	// Create temp dir for isolated compilation
	dstDir, err := os.MkdirTemp("", "sui-temp-*")
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(dstDir)

	dstRoot := filepath.Join(dstDir, "contracts")
	packageRoot := filepath.Join(dstRoot, packageDir)

	// Copy embedded contract files to temp workspace
	if err = writeEFS(contracts.Embed, ".", dstRoot); err != nil {
		return PackageArtifact{}, fmt.Errorf("copying embedded files to %q: %w", dstRoot, err)
	}

	// Load and patch Move.toml
	tomlPath := filepath.Join(packageRoot, "Move.toml")
	manifest, err := loadManifest(tomlPath)
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("loading manifest: %w", err)
	}
	manifest.Addresses = namedAddresses
	if err = writeManifest(tomlPath, manifest); err != nil {
		return PackageArtifact{}, fmt.Errorf("writing manifest: %w", err)
	}

	// Special-case: update published-at of CCIP if this is the onramp package
	if packageName == contracts.CCIPOnramp {
		if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
		}

		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
		}
	}

	// Special-case: update published-at of CCIP & MCMS if this is the TokenPool package
	if packageName == contracts.CCIPTokenPool {
		if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
		}

		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
		}
	}

	// Special-case: update published-at of CCIP & MCMs if it's a offRamp package
	if packageName == contracts.CCIPOfframp {
		if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
		}
		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
		}
	}

	if packageName == contracts.CCIP {
		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
		}
	}

	if packageName == contracts.ManagedToken {
		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
		}
	}

	if packageName == contracts.CCIPRouter {
		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
		}
	}

	if packageName == contracts.MCMSUser {
		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
		}
	}

	// Special-case: update published-at of CCIP, CCIP Token Pool, & MCMs if it's a token pool package
	if packageName == contracts.LockReleaseTokenPool || packageName == contracts.BurnMintTokenPool || packageName == contracts.ManagedTokenPool || packageName == contracts.USDCTokenPool {
		if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
		}

		if err = updatePublishedAt(dstRoot, contracts.CCIPTokenPool, namedAddresses["ccip_token_pool"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating CCIP Token Pool published-at: %w", err)
		}

		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
		}
	}

	// Special-case: update published-at of Managed Token if it's a managed token pool package
	if packageName == contracts.ManagedTokenPool {
		if err = updatePublishedAt(dstRoot, contracts.ManagedToken, namedAddresses["managed_token"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating Managed Token published-at: %w", err)
		}
	}

	// Special-case: update published-at of CCIP & MCMS if it's a dummy receiver package
	if packageName == contracts.CCIPDummyReceiver {
		if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
		}
		if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
			return PackageArtifact{}, fmt.Errorf("updating MCMS published-at: %w", err)
		}
	}

	// Compile the Move package
	cmd := exec.Command("sui", "move", "build", "--dump-bytecode-as-base64", "--ignore-chain")
	cmd.Dir = packageRoot

	output, err := cmd.Output()
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("sui move build failed (%s): %w\nOutput:\n%s", cmd.Dir, err, output)
	}

	return ToArtifact(string(output))
}

func loadManifest(path string) (PackageManifest, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return PackageManifest{}, fmt.Errorf("reading %s: %w", path, err)
	}

	var manifest PackageManifest
	if err := toml.Unmarshal(bytes, &manifest); err != nil {
		return PackageManifest{}, fmt.Errorf("unmarshaling %s: %w", path, err)
	}

	return manifest, nil
}

func writeManifest(path string, manifest PackageManifest) error {
	bytes, err := toml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("marshaling TOML: %w", err)
	}

	//nolint:mnd
	return os.WriteFile(path, bytes, 0600)
}

func updatePublishedAt(root string, pkg contracts.Package, addr string) error {
	dir, ok := contracts.Contracts[pkg]
	if !ok {
		return fmt.Errorf("unknown package: %s", pkg)
	}
	path := filepath.Join(root, dir, "Move.toml")

	manifest, err := loadManifest(path)
	if err != nil {
		return err
	}

	var pkgTable map[string]any
	if pkgTable, ok = manifest.Package.(map[string]any); !ok {
		return fmt.Errorf("[package] table is not a map")
	}
	pkgTable["published-at"] = addr

	return writeManifest(path, manifest)
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
