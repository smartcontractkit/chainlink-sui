package bind

import (
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml"
	"golang.org/x/crypto/blake2b"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-sui/contracts"
)

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

func CompilePackage(packageName contracts.Package, namedAddresses map[string]string, isUpgrade bool) (PackageArtifact, error) {

	// 1️Detect dynamic RPC from Docker
	rpcURL, err := getDynamicSuiRPC()
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to detect sui rpc url: %w", err)
	}
	fmt.Printf("Using dynamic Sui RPC: %s\n", rpcURL)

	// 2 Create or update a sui env alias (in current config)
	env := "docker"
	if err := setupSuiEnv(env, rpcURL); err != nil {
		return PackageArtifact{}, fmt.Errorf("failed to create sui env alias: %w", err)
	}

	packageDir, ok := contracts.Contracts[packageName]
	if !ok {
		return PackageArtifact{}, fmt.Errorf("unknown package: %s", packageName)
	}

	// Create temp dir for isolated compilation
	dstDir, err := os.MkdirTemp("/Users/stackman/Desktop/chainlink-sui/", "sui-temp-*")
	if err != nil {
		return PackageArtifact{}, fmt.Errorf("creating temp dir: %w", err)
	}
	// defer os.RemoveAll(dstDir)

	dstRoot := filepath.Join(dstDir, "contracts")
	packageRoot := filepath.Join(dstRoot, packageDir)

	// Copy embedded contract files to temp workspace
	if err = writeEFS(contracts.Embed, ".", dstRoot); err != nil {
		return PackageArtifact{}, fmt.Errorf("copying embedded files to %q: %w", dstRoot, err)
	}

	// Load and patch Move.toml
	// tomlPath := filepath.Join(packageRoot, "Move.toml")
	// manifest, err := loadManifest(tomlPath)
	// if err != nil {
	// 	return PackageArtifact{}, fmt.Errorf("loading manifest: %w", err)
	// }

	// //  Filter namedAddresses: exclude env metadata keys
	// manifest.Addresses = filterMoveAddresses(namedAddresses)

	// if err = writeManifest(tomlPath, manifest); err != nil {
	// 	return PackageArtifact{}, fmt.Errorf("writing manifest: %w", err)
	// }
	// fmt.Println("Move.toml updated with valid Move addresses only")

	// if packageName == contracts.CCIPRouter {
	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
	// 	}
	// }

	// Special-case: update published-at of CCIP if this is the onramp package
	// if packageName == contracts.CCIPOnramp {
	// 	if isUpgrade {
	// 		// Replace onramp.move inside the temp sui-temp-* workspace with upgraded mock version
	// 		upgradeSrc := filepath.Join(dstRoot, "ccip", "mock_onramp_v2", "onramp.move")
	// 		upgradeDst := filepath.Join(packageRoot, "sources", "onramp.move")

	// 		input, err := os.ReadFile(upgradeSrc)
	// 		if err != nil {
	// 			return PackageArtifact{}, fmt.Errorf("reading onramp upgrade mock %q: %w", upgradeSrc, err)
	// 		}
	// 		if err := os.WriteFile(upgradeDst, input, 0o644); err != nil {
	// 			return PackageArtifact{}, fmt.Errorf("replacing onramp.move inside sui-temp workspace: %w", err)
	// 		}
	// 		fmt.Printf(" Using upgraded onramp.move inside sui-temp workspace:\n  SRC: %s\n  DST: %s\n", upgradeSrc, upgradeDst)

	// 		// --- (1) Write env.local for OnRamp ---
	// 		if envBlock, ok := buildEnvBlock("onramp", namedAddresses); ok {
	// 			onrampLock := filepath.Join(packageRoot, "Move.lock")
	// 			if err := appendEnvBlock(onrampLock, envBlock); err != nil {
	// 				return PackageArtifact{}, fmt.Errorf("writing env block to OnRamp Move.lock: %w", err)
	// 			}
	// 			fmt.Printf(" Added [env.local] metadata to OnRamp Move.lock:\n%s\n", envBlock)
	// 		} else {
	// 			fmt.Println(" No env metadata provided for OnRamp — skipping [env.local] block")
	// 		}

	// 		// --- (2) Write env.local for CCIP ---
	// 		if envBlock, ok := buildEnvBlock("ccip", namedAddresses); ok {
	// 			ccipLock := filepath.Join(dstRoot, "ccip", "ccip", "Move.lock")
	// 			if _, err := os.Stat(ccipLock); err == nil {
	// 				if err := appendEnvBlock(ccipLock, envBlock); err != nil {
	// 					return PackageArtifact{}, fmt.Errorf("writing env block to CCIP Move.lock: %w", err)
	// 				}
	// 				fmt.Printf("Added [env.local] metadata to CCIP Move.lock:\n%s\n", envBlock)
	// 			} else {
	// 				fmt.Printf("CCIP Move.lock not found at %s — skipping\n", ccipLock)
	// 			}
	// 		} else {
	// 			fmt.Println(" No env metadata provided for CCIP — skipping [env.local] block")
	// 		}
	// 	}

	// 	if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
	// 	}

	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
	// 	}

	// }

	// // Special-case: update published-at of CCIP & MCMs if it's a offRamp package
	// if packageName == contracts.CCIPOfframp {
	// 	if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
	// 	}
	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
	// 	}
	// }

	// if packageName == contracts.CCIP {
	// 	if isUpgrade {
	// 		// Replace fee_quoter.move inside the temp sui-temp-* workspace with upgraded mock version
	// 		upgradeSrc := filepath.Join(dstRoot, "ccip", "mock_ccip_v2", "fee_quoter.move")

	// 		// Path inside the temp workspace (automatically created)
	// 		upgradeDst := filepath.Join(packageRoot, "sources", "fee_quoter.move")

	// 		input, err := os.ReadFile(upgradeSrc)
	// 		if err != nil {
	// 			return PackageArtifact{}, fmt.Errorf("reading feequoter upgrade mock %q: %w", upgradeSrc, err)
	// 		}

	// 		// Overwrite the onramp.move in the sui-temp workspace
	// 		if err := os.WriteFile(upgradeDst, input, 0o644); err != nil {
	// 			return PackageArtifact{}, fmt.Errorf("replacing feequoter.move inside sui-temp workspace: %w", err)
	// 		}

	// 		fmt.Printf(" Using upgraded feequoter.move inside sui-temp workspace:\n  SRC: %s\n  DST: %s\n", upgradeSrc, upgradeDst)

	// 		// --- (2) Write env.local for CCIP ---
	// 		if envBlock, ok := buildEnvBlock("ccip", namedAddresses); ok {
	// 			ccipLock := filepath.Join(dstRoot, "ccip", "ccip", "Move.lock")
	// 			if _, err := os.Stat(ccipLock); err == nil {
	// 				if err := appendEnvBlock(ccipLock, envBlock); err != nil {
	// 					return PackageArtifact{}, fmt.Errorf("writing env block to CCIP Move.lock: %w", err)
	// 				}
	// 				fmt.Printf(" Added [env.local] metadata to CCIP Move.lock:\n%s\n", envBlock)
	// 			} else {
	// 				fmt.Printf("⚠️ CCIP Move.lock not found at %s — skipping\n", ccipLock)
	// 			}
	// 		} else {
	// 			fmt.Println("ℹNo env metadata provided for CCIP — skipping [env.local] block")
	// 		}

	// 	}

	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
	// 	}
	// }

	// if packageName == contracts.ManagedToken {
	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
	// 	}
	// }

	// if packageName == contracts.CCIPRouter {
	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
	// 	}
	// }

	// if packageName == contracts.MCMSUser {
	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
	// 	}
	// }
	// if packageName == contracts.MCMSUserV2 {
	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMsV2 published-at: %w", err)
	// 	}
	// }

	// // Special-case: update published-at of CCIP, CCIP Token Pool, & MCMs if it's a token pool package
	// if packageName == contracts.LockReleaseTokenPool || packageName == contracts.BurnMintTokenPool || packageName == contracts.ManagedTokenPool || packageName == contracts.USDCTokenPool {
	// 	if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
	// 	}

	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMs published-at: %w", err)
	// 	}
	// }

	// // Special-case: update published-at of Managed Token if it's a managed token pool package
	// if packageName == contracts.ManagedTokenPool {
	// 	if err = updatePublishedAt(dstRoot, contracts.ManagedToken, namedAddresses["managed_token"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating Managed Token published-at: %w", err)
	// 	}
	// }

	// // Special-case: update published-at of CCIP & MCMS if it's a dummy receiver package
	// if packageName == contracts.CCIPDummyReceiver {
	// 	if err = updatePublishedAt(dstRoot, contracts.CCIP, namedAddresses["ccip"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating CCIP published-at: %w", err)
	// 	}
	// 	if err = updatePublishedAt(dstRoot, contracts.MCMS, namedAddresses["mcms"]); err != nil {
	// 		return PackageArtifact{}, fmt.Errorf("updating MCMS published-at: %w", err)
	// 	}
	// }

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

			fmt.Printf(" Using upgraded feequoter.move inside sui-temp workspace:\n  SRC: %s\n  DST: %s\n", upgradeSrc, upgradeDst)

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

			fmt.Printf(" Using upgraded onramp.move inside sui-temp workspace:\n  SRC: %s\n  DST: %s\n", upgradeSrc, upgradeDst)

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
	}

	// Compile the Move package
	//nolint:noctx // we do not need a context here

	var cmd *exec.Cmd
	var digest []byte
	var deps []string
	var modules []string
	if isUpgrade {
		fmt.Println("UGPRADE HAPPENING IN: ", packageRoot)
		fmt.Println("NAMED ADDRESSES: ", namedAddresses)

		cmd = exec.Command("sui", "client", "upgrade", "--upgrade-capability", namedAddresses["upgrade_cap"], "--serialize-unsigned-transaction", "--sender", namedAddresses["signer"], "--json")
		cmd.Dir = packageRoot
		output, err := cmd.Output()
		if err != nil {
			return PackageArtifact{}, fmt.Errorf("sui client upgrade --upgrade-capability **addr** --dry-run (%s): %w\nOutput:\n%s", cmd.Dir, err, output)
		}

		err = os.WriteFile("upgrade_output.json", output, 0644)
		if err != nil {
			log.Fatalf("failed to write file: %v", err)
		}
		fmt.Println("Wrote UPGRADE OUTPUT to upgrade_output.json")

		var resp TransactionData
		if err := json.Unmarshal(output, &resp); err != nil {
			fmt.Println("ERROR DURING UNMARSHAL")
			return PackageArtifact{}, err
		}

		//  dependencies
		depsInput := resp.V1.Kind.ProgrammableTransaction.Commands[1].Upgrade
		depsAny, ok := depsInput[1].([]interface{})
		if !ok {
			fmt.Printf("Upgrade[1] is not []interface{}, got %T\n", depsInput[1])
			return PackageArtifact{}, fmt.Errorf("unexpected dependencies format")
		}
		for i, v := range depsAny {
			addrStr, ok := v.(string)
			if !ok {
				fmt.Printf("dep[%d] not a string, got %T\n", i, v)
				continue
			}
			fmt.Println("Dependency for Upgrade:", addrStr)
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
		fmt.Println("SIGNER ADDR: ", namedAddresses["signer"])
		cmd = exec.Command("sui", "client", "publish", "--serialize-unsigned-transaction", "--sender", namedAddresses["signer"], "--json")
		cmd.Dir = packageRoot
		output, err := cmd.Output()
		if err != nil {
			return PackageArtifact{}, fmt.Errorf("sui client publish --serialize-unsigned-transaction (%s): %w\nOutput:\n%s", cmd.Dir, err, output)
		}

		var resp TransactionData
		if err := json.Unmarshal(output, &resp); err != nil {
			fmt.Println("ERROR DURING UNMARSHAL")
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
			fmt.Println("Dependency:", addrStr)
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

// buildEnvBlock builds the TOML env.local block for a given package prefix
func buildEnvBlock(prefix string, meta map[string]string) (string, bool) {
	chainID := meta["chain_id_"+prefix]
	originalID := meta["original_published_id_"+prefix]
	latestID := meta["latest_published_id_"+prefix]
	version := meta["published_version_"+prefix]

	if originalID == "" && latestID == "" && version == "" {
		return "", false
	}

	envBlock := "[env]\n\n[env.local]\n"
	if chainID != "" {
		envBlock += fmt.Sprintf(`chain-id = "%s"`+"\n", chainID)
	}
	if originalID != "" {
		envBlock += fmt.Sprintf(`original-published-id = "%s"`+"\n", originalID)
	}
	if latestID != "" {
		envBlock += fmt.Sprintf(`latest-published-id = "%s"`+"\n", latestID)
	}
	if version != "" {
		envBlock += fmt.Sprintf(`published-version = "%s"`+"\n", version)
	}
	return envBlock, true
}

// appendEnvBlock appends a TOML env.local section to the specified Move.lock file
func appendEnvBlock(lockPath, envBlock string) error {
	f, err := os.OpenFile(lockPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("opening %s: %w", lockPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + envBlock + "\n"); err != nil {
		return fmt.Errorf("writing env block to %s: %w", lockPath, err)
	}
	return nil
}

// filterMoveAddresses returns only the valid Move.toml addresses (excludes metadata)
func filterMoveAddresses(all map[string]string) map[string]string {
	moveAddresses := make(map[string]string)
	for k, v := range all {
		// Skip metadata keys that begin with known prefixes
		if strings.HasPrefix(k, "chain_id_") ||
			strings.HasPrefix(k, "original_published_id_") ||
			strings.HasPrefix(k, "latest_published_id_") ||
			strings.HasPrefix(k, "published_version_") {
			continue
		}
		moveAddresses[k] = v
	}
	return moveAddresses
}

func parseBuildInfo(path string) (*BuildInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var info BuildInfo
	if err := yaml.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// encodeModules reads and base64-encodes each Move .mv file in the given order (if provided).
// If order is empty, it falls back to alphabetical ordering.
// Logs each module file name and its base64 prefix for debugging.
func encodeModules(dir string, order []string) ([]string, error) {
	var modules []string

	fmt.Printf("🔍 Reading bytecode modules from: %s\n", dir)

	// --- Case 1: Use provided topological order from dry-run ---
	if len(order) > 0 {
		fmt.Printf("📦 Using provided module order (%d modules): %v\n", len(order), order)

		for i, modName := range order {
			path := filepath.Join(dir, modName+".mv")

			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read module %s: %w", path, err)
			}

			encoded := base64.StdEncoding.EncodeToString(data)
			modules = append(modules, encoded)

			// Log for debugging: show filename and first 64 chars of Base64 (truncated)
			prefix := encoded
			if len(prefix) > 64 {
				prefix = prefix[:64] + "..."
			}
			fmt.Printf("[%2d] 🧩 %s.mv → Base64 prefix: %s\n", i, modName, prefix)
		}

		fmt.Printf("✅ Encoded %d modules in topological order.\n", len(modules))
		return modules, nil
	}

	// --- Case 2: No order provided, fall back to alphabetical ---
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading bytecode dir: %w", err)
	}

	fmt.Printf("⚠️ No module order provided — using alphabetical order.\n")

	for i, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".mv" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}

		encoded := base64.StdEncoding.EncodeToString(data)
		modules = append(modules, encoded)

		prefix := encoded
		if len(prefix) > 64 {
			prefix = prefix[:64] + "..."
		}
		fmt.Printf("[%2d] 🧩 %s → Base64 prefix: %s\n", i, entry.Name(), prefix)
	}

	if len(modules) == 0 {
		return nil, fmt.Errorf("no .mv modules found in %s", dir)
	}

	fmt.Printf("✅ Encoded %d modules alphabetically.\n", len(modules))
	return modules, nil
}

func getPackageNameFromMoveToml(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, "Move.toml"))
	if err != nil {
		return "", err
	}
	type moveToml struct {
		Package struct {
			Name string `toml:"name"`
		} `toml:"package"`
	}
	var cfg moveToml
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	if cfg.Package.Name == "" {
		return "", fmt.Errorf("no [package].name in Move.toml")
	}
	return cfg.Package.Name, nil
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
	fmt.Printf(" Chain ID: %s\n", chainID)

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

	fmt.Printf(" Managed package %s (env=%s, rpc=%s, chain-id=%s)\n%s\n",
		originalPkgId, latestPkgId, env, rpcURL, chainID, string(out))

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
	cmd := exec.Command("docker", "ps", "--filter", "ancestor=mysten/sui-tools:devnet", "--format", "{{.Ports}}")
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
	fmt.Printf("Setting up Sui env alias '%s' for RPC: %s\n", alias, rpcURL)

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
			fmt.Printf("🧹 Removing existing alias '%s' (RPC: %s)\n", alias, e.RPC)
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
		return fmt.Errorf("failed to create sui env '%s': %w\nOutput:\n%s", alias, err, string(newOut))
	}
	fmt.Printf("Registered new Sui alias '%s' with RPC %s\n", alias, rpcURL)

	// Step 4️⃣ — Switch to new env
	switchCmd := exec.Command("sui", "client", "switch", "--env", alias)
	switchCmd.Env = os.Environ()
	switchOut, err := switchCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to env '%s': %w\nOutput:\n%s", alias, err, string(switchOut))
	}
	fmt.Printf("🔁 Switched active Sui env to '%s'\n", alias)

	// Step 5️⃣ — Verify
	activeCmd := exec.Command("sui", "client", "active-env")
	activeCmd.Env = os.Environ()
	activeOut, _ := activeCmd.Output()
	fmt.Printf("Active Sui environment: %s\n", strings.TrimSpace(string(activeOut)))

	return nil
}

func removeAliasFromClientYAML(alias string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot resolve home directory: %w", err)
	}
	configPath := fmt.Sprintf("%s/.sui/sui_config/client.yaml", homeDir)

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

// createOrUpdateEnv runs `sui client new-env --rpc <url> --alias <alias>`
func createOrUpdateEnv(alias, rpcURL string) error {
	cmd := exec.Command("sui", "client", "new-env",
		"--rpc", rpcURL,
		"--alias", alias,
	)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create/update Sui env '%s': %w\nOutput:\n%s", alias, err, string(out))
	}
	fmt.Printf("Registered/updated Sui alias '%s' with RPC %s\n", alias, rpcURL)
	return nil
}

// switchEnv switches to an existing alias
func switchEnv(alias string) error {
	switchCmd := exec.Command("sui", "client", "switch", "--env", alias)
	switchCmd.Env = os.Environ()
	out, err := switchCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to env %s: %w\nOutput:\n%s", alias, err, string(out))
	}
	fmt.Printf("🔁 Switched active Sui env to '%s'\n", alias)

	// Verify switch worked
	verifyCmd := exec.Command("sui", "client", "active-env")
	verifyCmd.Env = os.Environ()
	activeOut, _ := verifyCmd.Output()
	fmt.Printf(" Active Sui environment: %s\n", strings.TrimSpace(string(activeOut)))

	return nil
}

// detectSuiContainer finds the container ID/name for mysten/sui-tools:devnet
func detectSuiContainer() (string, error) {
	cmd := exec.Command("docker", "ps",
		"--filter", "ancestor=mysten/sui-tools:devnet",
		"--format", "{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("docker ps failed: %w", err)
	}
	name := strings.TrimSpace(string(out))
	if name == "" {
		return "", fmt.Errorf("no container found for image mysten/sui-tools:devnet")
	}
	return name, nil
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

// convert hex addresses like 0xabc123... to [][]byte
func hexAddressesToBytes(addrs []string) ([][]byte, error) {
	var out [][]byte
	for _, addr := range addrs {
		a := strings.TrimPrefix(addr, "0x")
		// ensure even length for hex.DecodeString
		if len(a)%2 != 0 {
			a = "0" + a
		}
		b, err := hex.DecodeString(a)
		if err != nil {
			return nil, fmt.Errorf("failed to decode address %s: %w", addr, err)
		}
		out = append(out, b)
	}
	return out, nil
}
