package bind

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/crypto/blake2b"

	"github.com/smartcontractkit/chainlink-sui/contracts"
)

const env = "docker"

type PackageManifest struct {
	Package      any               `toml:"package"`
	Addresses    map[string]string `toml:"addresses"`
	Dependencies any               `toml:"dependencies"`
	DevAddresses any               `toml:"dev-addresses"`
}

type BuildInfo struct {
	CompiledPackageInfo struct {
		PackageName string `yaml:"package_name"`
	} `yaml:"compiled_package_info"`
}

type RPCResponse struct {
	Result struct {
		Data struct {
			ObjectID string            `json:"objectId"`
			Owner    map[string]string `json:"owner"`
		} `json:"data"`
	} `json:"result"`
}

// Root object
type TransactionData struct {
	V1 TransactionDataV1 `json:"V1"`
}

// TransactionDataV1
type TransactionDataV1 struct {
	Kind       TransactionKind `json:"kind"`
	Sender     string          `json:"sender"`
	GasData    GasData         `json:"gas_data"`
	Expiration string          `json:"expiration"`
}

// TransactionKind
type TransactionKind struct {
	ProgrammableTransaction ProgrammableTransaction `json:"ProgrammableTransaction"`
}

// ProgrammableTransaction
type ProgrammableTransaction struct {
	Inputs   []CallArg `json:"inputs"`
	Commands []Command `json:"commands"`
}

// CallArg — mimics SDK CallArg but matches CLI JSON
type CallArg struct {
	Pure []byte `json:"Pure,omitempty"`
	// You could later add Object, UnresolvedPure, etc.
}

// Command — mimics SDK Command but matches CLI JSON
type Command struct {
	MoveCall *ProgrammableMoveCall `json:"MoveCall,omitempty"`
	Publish  [][]interface{}       `json:"Publish,omitempty"`
	Upgrade  []interface{}         `json:"Upgrade,omitempty"`
}

// ArgumentRef — small helper for fields like {"Result": 0}
type ArgumentRef struct {
	Result *uint16 `json:"Result,omitempty"`
	Input  *uint16 `json:"Input,omitempty"`
}

// TransferObjects
type TransferObjects struct {
	Objects []ArgumentRef `json:"TransferObjects,omitempty"`
	Address ArgumentRef   `json:"Address,omitempty"`
}

// ProgrammableMoveCall — kept simple for now
type ProgrammableMoveCall struct {
	Package       string        `json:"package"`
	Module        string        `json:"module"`
	Function      string        `json:"function"`
	TypeArguments []interface{} `json:"type_arguments"`
	Arguments     []ArgumentRef `json:"arguments"`
}

// GasData
type GasData struct {
	Payment [][]interface{} `json:"payment"`
	Owner   string          `json:"owner"`
	Price   uint64          `json:"price"`
	Budget  uint64          `json:"budget"`
}

func CompilePackage(packageName contracts.Package, namedAddresses map[string]string, isUpgrade bool, suiRPC string) (PackageArtifact, error) {
	var rpcURL string
	// 1️. Detect dynamic RPC from Docker
	if suiRPC == "" {
		var err error
		rpcURL, err = getDynamicSuiRPC()
		if err != nil {
			return PackageArtifact{}, fmt.Errorf("failed to detect sui rpc url: %w", err)
		}
	} else {
		rpcURL = suiRPC
	}

	fmt.Printf("Using Sui RPC: %s\n", rpcURL)

	// before you set the temp dir
	prevConfigDir := os.Getenv("SUI_CONFIG_DIR")

	// Create isolated config
	tempConfigDir, err := os.MkdirTemp("", "sui-config-*")
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to create temp sui config dir: %w", err)
	}
	defer os.RemoveAll(tempConfigDir)

	os.Setenv("SUI_CONFIG_DIR", tempConfigDir)

	// ➜ IMPORTANT: restore when we leave CompilePackage
	defer func() {
		if prevConfigDir == "" {
			_ = os.Unsetenv("SUI_CONFIG_DIR")
		} else {
			_ = os.Setenv("SUI_CONFIG_DIR", prevConfigDir)
		}
	}()

	// Initialize config non-interactively
	initCmd := exec.Command("sui", "client", "--yes", "--json")
	initCmd.Env = append(os.Environ(), fmt.Sprintf("SUI_CONFIG_DIR=%s", tempConfigDir))
	if out, err := initCmd.CombinedOutput(); err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to init sui client: %w\n%s", err, out)
	}

	// 2. Create or update a sui env alias (in current config)
	if err := setupSuiEnv(env, rpcURL); err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to create sui env alias: %w", err)
	}

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

	if packageName == contracts.Test {
		testSecondaryAddr := namedAddresses["test_secondary"]
		if !isZeroAddress(testSecondaryAddr) {
			testSecondaryDir := filepath.Join(dstRoot, "test_secondary")
			if err := managePackage(testSecondaryDir, 1, rpcURL, env, testSecondaryAddr, testSecondaryAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage Test Secondary dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for Test Secondary (no published address found)")
		}
	}

	if packageName == contracts.ManagedToken {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}
	}

	if packageName == contracts.ManagedTokenFaucet {
		managedTokenAddr := namedAddresses["managed_token"]
		if !isZeroAddress(managedTokenAddr) {
			managedTokenDir := filepath.Join(dstRoot, "ccip", "managed_token")
			if err := managePackage(managedTokenDir, 1, rpcURL, env, managedTokenAddr, managedTokenAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage ManagedToken dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for ManagedToken (no published address found)")
		}

		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}
	}

	if packageName == contracts.MCMSUser {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}
	}

	if packageName == contracts.MCMSUserV2 {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}
	}

	if packageName == contracts.LockReleaseTokenPool || packageName == contracts.BurnMintTokenPool || packageName == contracts.ManagedTokenPool || packageName == contracts.USDCTokenPool {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}

		ccipAddr := namedAddresses["ccip"]
		if !isZeroAddress(ccipAddr) {
			ccipDir := filepath.Join(dstRoot, "ccip", "ccip")
			if err := managePackage(ccipDir, 1, rpcURL, env, ccipAddr, ccipAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage CCIP dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for CCIP (no published address found)")
		}
	}

	if packageName == contracts.ManagedTokenPool {
		managedTokenAddr := namedAddresses["managed_token"]
		if !isZeroAddress(managedTokenAddr) {
			managedTokenDir := filepath.Join(dstRoot, "ccip", "managed_token")
			if err := managePackage(managedTokenDir, 1, rpcURL, env, managedTokenAddr, managedTokenAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage managed token dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for managed token (no published address found)")
		}
	}

	if packageName == contracts.CCIPDummyReceiver {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}

		ccipAddr := namedAddresses["ccip"]
		if !isZeroAddress(ccipAddr) {
			ccipDir := filepath.Join(dstRoot, "ccip", "ccip")
			if err := managePackage(ccipDir, 1, rpcURL, env, ccipAddr, ccipAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage CCIP dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for CCIP (no published address found)")
		}

	}

	if packageName == contracts.CCIPRouter {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}
	}

	if packageName == contracts.CCIP {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}

		// if upgrade it needs to move.lock in it's own pkg
		if isUpgrade {
			// Replace fee_quoter.move inside the temp sui-temp-* workspace with upgraded mock version
			upgradeSrc := filepath.Join(dstRoot, "ccip", "mock_ccip_v2", "fee_quoter.move")

			// Path inside the temp workspace (automatically created)
			upgradeDst := filepath.Join(packageRoot, "sources", "fee_quoter.move")

			input, err := os.ReadFile(upgradeSrc)
			if err != nil {
				return PackageArtifact{}, fmt.Errorf("reading feequoter upgrade mock %q: %w", upgradeSrc, err)
			}

			// Overwrite the onramp.move in the sui-temp workspace
			if err := os.WriteFile(upgradeDst, input, 0o644); err != nil {
				return PackageArtifact{}, fmt.Errorf("replacing feequoter.move inside sui-temp workspace: %w", err)
			}

			ccipAddr := namedAddresses["original_ccip_pkg"]
			if !isZeroAddress(mcmsAddr) {
				ccipDir := filepath.Join(dstRoot, "ccip", "ccip")
				if err := managePackage(ccipDir, 1, rpcURL, env, ccipAddr, ccipAddr); err != nil {
					return PackageArtifact{}, fmt.Errorf("failed to manage CCIP dependency: %w", err)
				}
			} else {
				fmt.Println("Skipping manage-package for CCIP (no published address found)")
			}
		}
	}

	if packageName == contracts.CCIPOnramp {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}

		ccipAddr := namedAddresses["ccip"]
		if !isZeroAddress(ccipAddr) {
			ccipDir := filepath.Join(dstRoot, "ccip", "ccip")
			if err := managePackage(ccipDir, 1, rpcURL, env, ccipAddr, ccipAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage CCIP dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for CCIP (no published address found)")
		}

		// TODO: make this only for mock test upgrade
		if isUpgrade {
			// Replace onramp.move inside the temp sui-temp-* workspace with upgraded mock version
			upgradeSrc := filepath.Join(dstRoot, "ccip", "mock_onramp_v2", "onramp.move")
			upgradeDst := filepath.Join(packageRoot, "sources", "onramp.move")

			// Read the mock upgrade file from repo
			input, err := os.ReadFile(upgradeSrc)
			if err != nil {
				return PackageArtifact{}, fmt.Errorf("reading onramp upgrade mock %q: %w", upgradeSrc, err)
			}

			// Overwrite the onramp.move in the sui-temp workspace
			if err := os.WriteFile(upgradeDst, input, 0o644); err != nil {
				return PackageArtifact{}, fmt.Errorf("replacing onramp.move inside sui-temp workspace: %w", err)
			}

			ccipOnRampAddr := namedAddresses["original_onramp_pkg"]
			if !isZeroAddress(ccipOnRampAddr) {
				ccipOnRampDir := filepath.Join(dstRoot, "ccip", "ccip_onramp")
				if err := managePackage(ccipOnRampDir, 1, rpcURL, env, ccipOnRampAddr, ccipOnRampAddr); err != nil {
					return PackageArtifact{}, fmt.Errorf("failed to manage CCIP OnRamp dependency: %w", err)
				}
			} else {
				fmt.Println("Skipping manage-package for CCIP OnRamp (no published address found)")
			}

			// also upgrade ccip move.Lock with updated values
			ccipLatestAddr := namedAddresses["latest_ccip_pkg"]
			ccipOriginalAddr := namedAddresses["ccip"]
			if !isZeroAddress(ccipLatestAddr) && !isZeroAddress(ccipOriginalAddr) {
				ccipDir := filepath.Join(dstRoot, "ccip", "ccip")
				if err := managePackage(ccipDir, 2, rpcURL, env, ccipOriginalAddr, ccipLatestAddr); err != nil {
					return PackageArtifact{}, fmt.Errorf("failed to manage CCIP dependency for onRamp: %w", err)
				}
			} else {
				fmt.Println("Skipping manage-package for CCIP Dependency for OnRamp (no published address found)")
			}

		}
	}

	if packageName == contracts.CCIPOfframp {
		mcmsAddr := namedAddresses["mcms"]
		if !isZeroAddress(mcmsAddr) {
			mcmsDir := filepath.Join(dstRoot, "mcms", "mcms")
			if err := managePackage(mcmsDir, 1, rpcURL, env, mcmsAddr, mcmsAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage MCMS dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for MCMS (no published address found)")
		}

		ccipAddr := namedAddresses["ccip"]
		if !isZeroAddress(mcmsAddr) {
			ccipDir := filepath.Join(dstRoot, "ccip", "ccip")
			if err := managePackage(ccipDir, 1, rpcURL, env, ccipAddr, ccipAddr); err != nil {
				return PackageArtifact{}, fmt.Errorf("failed to manage CCIP dependency: %w", err)
			}
		} else {
			fmt.Println("Skipping manage-package for CCIP (no published address found)")
		}

		// TODO: make this only for mock test upgrade
		if isUpgrade {
			// Replace offramp.move inside the temp sui-temp-* workspace with upgraded mock version
			upgradeSrc := filepath.Join(dstRoot, "ccip", "mock_offramp_v2", "offramp.move")
			upgradeDst := filepath.Join(packageRoot, "sources", "offramp.move")

			// Read the mock upgrade file from repo
			input, err := os.ReadFile(upgradeSrc)
			if err != nil {
				return PackageArtifact{}, fmt.Errorf("reading offramp upgrade mock %q: %w", upgradeSrc, err)
			}

			// Overwrite the offramp.move in the sui-temp workspace
			if err := os.WriteFile(upgradeDst, input, 0o644); err != nil {
				return PackageArtifact{}, fmt.Errorf("replacing offramp.move inside sui-temp workspace: %w", err)
			}

			fmt.Printf(" Using upgraded offramp.move inside sui-temp workspace:\n  SRC: %s\n  DST: %s\n", upgradeSrc, upgradeDst)

			ccipOffRampAddr := namedAddresses["original_offramp_pkg"]
			if !isZeroAddress(ccipOffRampAddr) {
				ccipOffRampDir := filepath.Join(dstRoot, "ccip", "ccip_offramp")
				if err := managePackage(ccipOffRampDir, 1, rpcURL, env, ccipOffRampAddr, ccipOffRampAddr); err != nil {
					return PackageArtifact{}, fmt.Errorf("failed to manage CCIP OffRamp dependency: %w", err)
				}
			} else {
				fmt.Println("Skipping manage-package for CCIP OffRamp (no published address found)")
			}

			// also upgrade ccip move.Lock with updated values
			ccipLatestAddr := namedAddresses["latest_ccip_pkg"]
			ccipOriginalAddr := namedAddresses["ccip"]
			if !isZeroAddress(ccipLatestAddr) && !isZeroAddress(ccipOriginalAddr) {
				ccipDir := filepath.Join(dstRoot, "ccip", "ccip")
				if err := managePackage(ccipDir, 2, rpcURL, env, ccipOriginalAddr, ccipLatestAddr); err != nil {
					return PackageArtifact{}, fmt.Errorf("failed to manage CCIP dependency for offramp: %w", err)
				}
			} else {
				fmt.Println("Skipping manage-package for CCIP Dependency for OffRamp (no published address found)")
			}

		}

	}

	var cmd *exec.Cmd
	var digest []byte
	var deps []string
	var modules []string
	if isUpgrade {
		cmd = exec.Command("sui", "client", "upgrade", "--upgrade-capability", namedAddresses["upgrade_cap"], "--serialize-unsigned-transaction", "--sender", namedAddresses["signer"], "--json")
		cmd.Dir = packageRoot
		output, err := cmd.Output()
		if err != nil {
			return PackageArtifact{}, fmt.Errorf("sui client upgrade --upgrade-capability **addr** --dry-run (%s): %w\nOutput:\n%s", cmd.Dir, err, output)
		}

		var resp TransactionData
		if err := json.Unmarshal(output, &resp); err != nil {
			return PackageArtifact{}, err
		}

		//  dependencies
		depsInput := resp.V1.Kind.ProgrammableTransaction.Commands[1].Upgrade
		depsAny, ok := depsInput[1].([]interface{})
		if !ok {
			return PackageArtifact{}, fmt.Errorf("unexpected dependencies format")
		}
		for i, v := range depsAny {
			addrStr, ok := v.(string)
			if !ok {
				fmt.Printf("dep[%d] not a string, got %T\n", i, v)
				continue
			}
			deps = append(deps, addrStr)
		}

		// modules
		modulesInput := resp.V1.Kind.ProgrammableTransaction.Commands[1].Upgrade
		modulesAny, ok := modulesInput[0].([]interface{})
		if !ok {
			fmt.Printf("Upgrade[0] is not []interface{}, got %T\n", modulesInput[0])
			return PackageArtifact{}, fmt.Errorf("unexpected dependencies format")
		}
		var base64Modules []string
		for i, modAny := range modulesAny {
			byteArr, ok := modAny.([]interface{})
			if !ok {
				fmt.Printf("module[%d] is not []interface{}, got %T\n", i, modAny)
				continue
			}

			// Convert []interface{} → []byte
			moduleBytes := make([]byte, len(byteArr))
			for j, b := range byteArr {
				moduleBytes[j] = byte(b.(float64)) // JSON numbers come in as float64
			}

			// Encode module bytes to Base64
			b64 := base64.StdEncoding.EncodeToString(moduleBytes)
			base64Modules = append(base64Modules, b64)
		}

		modules = base64Modules

		// digest
		digestInput := resp.V1.Kind.ProgrammableTransaction.Inputs[2].Pure
		// first element is the length of an array in bcs byte so removing it
		digest = digestInput[1:]

	} else {
		cmd = exec.Command("sui", "client", "publish", "--serialize-unsigned-transaction", "--sender", namedAddresses["signer"], "--json")
		cmd.Dir = packageRoot
		output, err := cmd.Output()
		if err != nil {
			return PackageArtifact{}, fmt.Errorf("sui client publish --serialize-unsigned-transaction (%s): %w\nOutput:\n%s", cmd.Dir, err, output)
		}

		var resp TransactionData
		if err := json.Unmarshal(output, &resp); err != nil {
			return PackageArtifact{}, err
		}

		//  dependencies
		depsInput := resp.V1.Kind.ProgrammableTransaction.Commands[0].Publish[1]
		for i, v := range depsInput {
			addrStr, ok := v.(string)
			if !ok {
				fmt.Printf("dep[%d] not a string, got %T\n", i, v)
				continue
			}
			deps = append(deps, addrStr)
		}

		// modules
		modulesInput := resp.V1.Kind.ProgrammableTransaction.Commands[0].Publish[0]
		// Prepare a slice to store Base64 strings
		var base64Modules []string

		for i, modAny := range modulesInput {
			byteArr, ok := modAny.([]interface{})
			if !ok {
				fmt.Printf("module[%d] is not []interface{}, got %T\n", i, modAny)
				continue
			}

			// Convert []interface{} → []byte
			moduleBytes := make([]byte, len(byteArr))
			for j, b := range byteArr {
				moduleBytes[j] = byte(b.(float64)) // JSON numbers are float64
			}

			// Encode module bytes to Base64
			b64 := base64.StdEncoding.EncodeToString(moduleBytes)
			base64Modules = append(base64Modules, b64)
		}

		modules = base64Modules
	}

	artifact := PackageArtifact{
		Modules:      modules,
		Dependencies: deps,
		Digest:       digest,
	}

	return artifact, nil
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

func isZeroAddress(addr string) bool {
	// Remove 0x prefix if present
	addr = strings.TrimPrefix(addr, "0x")
	for _, c := range addr {
		if c != '0' {
			return false
		}
	}
	return true
}

func managePackage(packageRoot string, version int, rpcURL, env, originalPkgId, latestPkgId string) error {
	//  Fetch chain identifier directly from the node
	chainID, err := getChainIdentifier(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to query chain identifier from %s: %w", rpcURL, err)
	}

	// Run manage-package
	cmd := exec.Command(
		"sui", "move", "manage-package",
		"--environment", env,
		"--network-id", chainID,
		"--original-id", originalPkgId,
		"--latest-id", latestPkgId,
		"--version-number", fmt.Sprintf("%d", version),
	)
	cmd.Dir = packageRoot
	cmd.Env = os.Environ() // includes SUI_CONFIG + PATH

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sui move manage-package failed: %w\nOutput:\n%s", err, string(out))
	}

	return nil
}

func getChainIdentifier(rpcURL string) (string, error) {
	req := `{"jsonrpc":"2.0","id":1,"method":"sui_getChainIdentifier"}`
	cmd := exec.Command("curl", "-s", "-X", "POST", "-H", "Content-Type: application/json", "-d", req, rpcURL)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to query chain identifier: %w", err)
	}
	var resp struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(out, &resp); err != nil {
		return "", fmt.Errorf("failed to parse chain identifier: %w\nResponse:\n%s", err, string(out))
	}
	return resp.Result, nil
}

func getDynamicSuiRPC() (string, error) {
	// only used for internal tests
	if envRPC := os.Getenv("SUI_RPC_URL"); envRPC != "" {
		return envRPC, nil
	}

	cmd := exec.Command("docker", "ps", "--filter", "ancestor=mysten/sui-tools:devnet-v1.61.0", "--format", "{{.Ports}}")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("docker ps failed: %w", err)
	}

	// Example: "0.0.0.0:14097->9000/tcp, 0.0.0.0:14098->9123/tcp"
	for _, part := range strings.Split(strings.TrimSpace(string(out)), ",") {
		p := strings.TrimSpace(part)
		if strings.Contains(p, "->9000") {
			hostPort := strings.Split(strings.Split(p, ":")[1], "->")[0]
			return fmt.Sprintf("http://127.0.0.1:%s", hostPort), nil
		}
	}
	return "", fmt.Errorf("could not find sui rpc port mapping for port 9000")
}

// setupSuiEnv ensures a Sui CLI environment alias exists for the given RPC.
// If the alias already exists, it removes it directly from client.yaml before recreating.
type suiEnv struct {
	Alias string `json:"alias"`
	RPC   string `json:"rpc"`
}

func setupSuiEnv(alias, rpcURL string) error {
	// Step 1 — Fetch all current envs via CLI
	cmd := exec.Command("sui", "client", "envs", "--json")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list Sui environments: %w", err)
	}

	var parsed []any
	if err := json.Unmarshal(out, &parsed); err != nil {
		return fmt.Errorf("failed to parse envs JSON: %w\nOutput:\n%s", err, string(out))
	}

	var envList []suiEnv
	if arr, ok := parsed[0].([]any); ok {
		for _, e := range arr {
			data, _ := json.Marshal(e)
			var env suiEnv
			if err := json.Unmarshal(data, &env); err == nil {
				envList = append(envList, env)
			}
		}
	}

	// Step 2 — Check for existing alias and remove it
	for _, e := range envList {
		if e.Alias == alias {
			if err := removeAliasFromClientYAML(alias); err != nil {
				return fmt.Errorf("failed to remove existing alias: %w", err)
			}
			break
		}
	}

	// Step  — Create new alias
	newCmd := exec.Command("sui", "client", "new-env",
		"--rpc", rpcURL,
		"--alias", alias,
	)
	newCmd.Env = os.Environ()
	newOut, err := newCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed to create sui env '%s': %v\nOutput:\n%s", alias, err, string(newOut))
	}

	// Step 4️ — Switch to new env
	switchCmd := exec.Command("sui", "client", "switch", "--env", alias)
	switchCmd.Env = os.Environ()
	switchOut, err := switchCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to env '%s': %w\nOutput:\n%s", alias, err, string(switchOut))
	}

	// Step 5️ — Verify
	activeCmd := exec.Command("sui", "client", "active-env")
	activeCmd.Env = os.Environ()
	return nil
}

func removeAliasFromClientYAML(alias string) error {
	configDir := os.Getenv("SUI_CONFIG_DIR")
	if configDir == "" {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".sui", "sui_config")
	}
	configPath := filepath.Join(configDir, "client.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", configPath, err)
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	skip := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- alias:") && strings.Contains(trimmed, alias) {
			skip = true // start skipping this env block
			continue
		}
		if skip && strings.HasPrefix(trimmed, "- alias:") {
			skip = false // stop skipping at next alias
		}
		if !skip {
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func ComputeDigestForModulesAndDeps(modules [][]byte, objectIDs [][]byte) [32]byte {
	var components [][]byte

	// Hash each module individually
	moduleDigests := make([][32]byte, 0, len(modules))
	for _, module := range modules {
		digest := blake2b.Sum256(module)
		moduleDigests = append(moduleDigests, digest)
	}

	// Convert digests to byte slices for sorting
	for i := range moduleDigests {
		components = append(components, moduleDigests[i][:])
	}

	// Add object IDs to components
	components = append(components, objectIDs...)

	// Sort components so order doesn't matter
	sort.Slice(components, func(i, j int) bool {
		return string(components[i]) < string(components[j])
	})

	// Hash all components together
	hasher, _ := blake2b.New256(nil)
	for _, c := range components {
		hasher.Write(c)
	}

	var result [32]byte
	copy(result[:], hasher.Sum(nil))
	return result
}
