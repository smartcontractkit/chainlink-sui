package util

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const (
	packageAddressCachePrefix = "pkg_addr:"
	packageIdsCachePrefix     = "pkg_ids:"
	latestPackageIdPrefix     = "latest_pkg:"
	defaultCacheTTL           = 10 * time.Minute
	identifierParts           = 3
)

// PackageResolver handles module name to package address resolution and package ID management
type PackageResolver struct {
	log              logger.Logger
	client           client.SuiPTBClient
	packageAddresses map[string]string
	cache            *cache.Cache
	mutex            sync.RWMutex
}

// ResolvedIdentifier represents a parsed Sui identifier
type ResolvedIdentifier struct {
	PackageId    string
	ModuleName   string
	FunctionName string
	Identifier   string
}

// NewPackageResolver creates a new package resolver instance
func NewPackageResolver(log logger.Logger, client client.SuiPTBClient) *PackageResolver {
	return &PackageResolver{
		log:              log,
		client:           client,
		packageAddresses: make(map[string]string),
		cache:            client.GetCache(),
	}
}

// BindPackage binds a module name to its package address
func (pr *PackageResolver) BindPackage(moduleName string, packageAddress string) error {
	moduleName = pr.normalizeName(moduleName)

	if !isValidSuiAddress(packageAddress) {
		return fmt.Errorf("invalid Sui package address format: %s", packageAddress)
	}

	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	pr.packageAddresses[moduleName] = packageAddress

	cacheKey := packageAddressCachePrefix + moduleName
	pr.cache.Set(cacheKey, packageAddress, defaultCacheTTL)

	pr.log.Debugw("Package bound", "module", moduleName, "address", packageAddress)
	return nil
}

// UnbindPackage removes a module binding
func (pr *PackageResolver) UnbindPackage(moduleName string) error {
	moduleName = pr.normalizeName(moduleName)

	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	if _, exists := pr.packageAddresses[moduleName]; !exists {
		return fmt.Errorf("no binding exists for module: %s", moduleName)
	}

	delete(pr.packageAddresses, moduleName)

	cacheKey := packageAddressCachePrefix + moduleName
	pr.cache.Delete(cacheKey)

	pr.log.Debugw("Package unbound", "module", moduleName)
	return nil
}

// ResolvePackageAddress resolves a module name to its package address
func (pr *PackageResolver) ResolvePackageAddress(moduleName string) (string, error) {
	moduleName = pr.normalizeName(moduleName)

	pr.mutex.RLock()
	address, exists := pr.packageAddresses[moduleName]
	pr.mutex.RUnlock()

	if exists {
		return address, nil
	}

	cacheKey := packageAddressCachePrefix + moduleName
	if cachedAddr, found := pr.cache.Get(cacheKey); found {
		if addr, ok := cachedAddr.(string); ok {
			pr.mutex.Lock()
			pr.packageAddresses[moduleName] = addr
			pr.mutex.Unlock()
			return addr, nil
		}
	}

	return "", fmt.Errorf("no package address found for module: %s", moduleName)
}

// ResolvePackageIds gets all package IDs for a module (including upgrades)
func (pr *PackageResolver) ResolvePackageIds(ctx context.Context, moduleName string, signerAddress string) ([]string, error) {
	moduleName = pr.normalizeName(moduleName)

	packageAddress, err := pr.ResolvePackageAddress(moduleName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve package address for module %s: %w", moduleName, err)
	}

	cacheKey := packageIdsCachePrefix + packageAddress + ":" + moduleName
	if cachedIds, found := pr.cache.Get(cacheKey); found {
		if ids, ok := cachedIds.([]string); ok {
			return ids, nil
		}
	}

	packageIds, err := pr.client.LoadModulePackageIds(ctx, packageAddress, moduleName, signerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to load package IDs for module %s: %w", moduleName, err)
	}

	pr.cache.Set(cacheKey, packageIds, defaultCacheTTL)

	pr.log.Debugw("Package IDs resolved",
		"module", moduleName,
		"address", packageAddress,
		"packageIds", packageIds)

	return packageIds, nil
}

// ResolveLatestPackageId gets the latest (most recent) package ID for a module
func (pr *PackageResolver) ResolveLatestPackageId(ctx context.Context, moduleName string, signerAddress string) (string, error) {
	moduleName = pr.normalizeName(moduleName)

	packageAddress, err := pr.ResolvePackageAddress(moduleName)
	if err != nil {
		return "", fmt.Errorf("failed to resolve package address for module %s: %w", moduleName, err)
	}

	cacheKey := latestPackageIdPrefix + packageAddress + ":" + moduleName
	if cachedId, found := pr.cache.Get(cacheKey); found {
		if id, ok := cachedId.(string); ok {
			return id, nil
		}
	}

	latestId, err := pr.client.GetLatestPackageId(ctx, packageAddress, moduleName, signerAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get latest package ID for module %s: %w", moduleName, err)
	}

	pr.cache.Set(cacheKey, latestId, defaultCacheTTL)

	pr.log.Debugw("Latest package ID resolved",
		"module", moduleName,
		"address", packageAddress,
		"latestId", latestId)

	return latestId, nil
}

// ParseIdentifier parses a Sui identifier in the format "packageId::moduleName::functionName"
func (pr *PackageResolver) ParseIdentifier(identifier string) (*ResolvedIdentifier, error) {
	parts := strings.Split(identifier, "::")
	if len(parts) != identifierParts {
		return nil, fmt.Errorf("invalid identifier format, expected 'packageId::moduleName::functionName', got: %s", identifier)
	}

	packageId := strings.TrimSpace(parts[0])
	moduleName := strings.TrimSpace(parts[1])
	functionName := strings.TrimSpace(parts[2])

	if packageId == "" || moduleName == "" || functionName == "" {
		return nil, fmt.Errorf("identifier parts cannot be empty: %s", identifier)
	}

	if !isValidSuiAddress(packageId) {
		return nil, fmt.Errorf("invalid package ID format: %s", packageId)
	}

	return &ResolvedIdentifier{
		PackageId:    packageId,
		ModuleName:   moduleName,
		FunctionName: functionName,
		Identifier:   identifier,
	}, nil
}

// ResolveIdentifier resolves a module name to latest package ID and creates full identifier
func (pr *PackageResolver) ResolveIdentifier(ctx context.Context, moduleName string, functionName string, signerAddress string) (*ResolvedIdentifier, error) {
	latestPackageId, err := pr.ResolveLatestPackageId(ctx, moduleName, signerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve latest package ID: %w", err)
	}

	identifier := fmt.Sprintf("%s::%s::%s", latestPackageId, moduleName, functionName)

	return &ResolvedIdentifier{
		PackageId:    latestPackageId,
		ModuleName:   moduleName,
		FunctionName: functionName,
		Identifier:   identifier,
	}, nil
}

// GetBoundModules returns all currently bound module names
func (pr *PackageResolver) GetBoundModules() []string {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	modules := make([]string, 0, len(pr.packageAddresses))
	for module := range pr.packageAddresses {
		modules = append(modules, module)
	}

	return modules
}

// InvalidateCache removes cached entries for a specific module
func (pr *PackageResolver) InvalidateCache(moduleName string) {
	moduleName = pr.normalizeName(moduleName)

	pr.mutex.RLock()
	packageAddress := pr.packageAddresses[moduleName]
	pr.mutex.RUnlock()

	keys := []string{
		packageAddressCachePrefix + moduleName,
		packageIdsCachePrefix + packageAddress + ":" + moduleName,
		latestPackageIdPrefix + packageAddress + ":" + moduleName,
	}

	for _, key := range keys {
		pr.cache.Delete(key)
	}

	pr.log.Debugw("Cache invalidated for module", "module", moduleName)
}

func (pr *PackageResolver) ValidateBinding(moduleName string, packageAddress string) error {
	moduleName = pr.normalizeName(moduleName)

	pr.mutex.RLock()
	boundPackageAddress := pr.packageAddresses[moduleName]
	pr.mutex.RUnlock()

	// If the key exists but the address does not match
	if boundPackageAddress == packageAddress {
		return nil
	}

	// If the key does not exist, check the cache
	if boundPackageAddress == "" {
		cacheKey := packageAddressCachePrefix + moduleName
		if cachedAddr, found := pr.cache.Get(cacheKey); found {
			if cachedAddr.(string) == packageAddress {
				return nil
			}
		}
	}

	return fmt.Errorf("invalid binding for module: %s and address: %s", moduleName, packageAddress)
}

// ClearCache clears all cached entries
func (pr *PackageResolver) ClearCache() {
	items := pr.cache.Items()
	for key := range items {
		if strings.HasPrefix(key, packageAddressCachePrefix) ||
			strings.HasPrefix(key, packageIdsCachePrefix) ||
			strings.HasPrefix(key, latestPackageIdPrefix) {
			pr.cache.Delete(key)
		}
	}

	pr.log.Debug("All package resolver cache entries cleared")
}

// isValidSuiAddress checks if an address has the proper Sui format (starts with 0x)
func isValidSuiAddress(address string) bool {
	return strings.HasPrefix(address, "0x") && len(address) > 2
}

func (pr *PackageResolver) normalizeName(moduleName string) string {
	return strings.ToLower(moduleName)
}

// String returns string representation of ResolvedIdentifier
func (ri *ResolvedIdentifier) String() string {
	return ri.Identifier
}

// ToIdentifier converts ResolvedIdentifier back to identifier string
func (ri *ResolvedIdentifier) ToIdentifier() string {
	return fmt.Sprintf("%s::%s::%s", ri.PackageId, ri.ModuleName, ri.FunctionName)
}
