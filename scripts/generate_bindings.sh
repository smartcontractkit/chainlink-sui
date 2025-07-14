
echo "Generating bindings for Move Sui contracts..."

# Build the bindings (add the path to contracts you want to generate bindings for)

# Test Package
go run bindgen/main.go --moveConfig ./contracts/test/ --input ./contracts/test/sources/counter.move --output ./bindings/generated/test/counter
go run bindgen/main.go --moveConfig ./contracts/test/ --input ./contracts/test/sources/complex.move --output ./bindings/generated/test/complex
go run bindgen/main.go --moveConfig ./contracts/test/ --input ./contracts/test/sources/generics.move --output ./bindings/generated/test/generics

# CCIP
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip --input ./contracts/ccip/ccip/sources/fee_quoter.move --output ./bindings/generated/ccip/ccip/fee_quoter
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip --input ./contracts/ccip/ccip/sources/nonce_manager.move --output ./bindings/generated/ccip/ccip/nonce_manager
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip --input ./contracts/ccip/ccip/sources/receiver_registry.move --output ./bindings/generated/ccip/ccip/receiver_registry
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip --input ./contracts/ccip/ccip/sources/rmn_remote.move --output ./bindings/generated/ccip/ccip/rmn_remote
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip --input ./contracts/ccip/ccip/sources/token_admin_registry.move --output ./bindings/generated/ccip/ccip/token_admin_registry

# CCIP - Onramp
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip_onramp --input ./contracts/ccip/ccip_onramp/sources/onramp.move --output ./bindings/generated/ccip/ccip_onramp/onramp

# CCIP - Offramp
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip_offramp --input ./contracts/ccip/ccip_offramp/sources/offramp.move --output ./bindings/generated/ccip/ccip_offramp/offramp

# CCIP - LINK
go run bindgen/main.go --moveConfig ./contracts/ccip/link_token --input ./contracts/ccip/link_token/sources/link_token.move --output ./bindings/generated/ccip/link_token/link_token

# CCIP - Managed Token
go run bindgen/main.go --moveConfig ./contracts/ccip/managed_token --input ./contracts/ccip/managed_token/sources/managed_token.move --output ./bindings/generated/ccip/managed_token/managed_token

# CCIP - Token Pool
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip_token_pools/token_pool --input ./contracts/ccip/ccip_token_pools/token_pool/sources/token_pool.move --output ./bindings/generated/ccip/ccip_token_pools/token_pool

# CCIP - Lock Release Token Pool
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip_token_pools/lock_release_token_pool --input ./contracts/ccip/ccip_token_pools/lock_release_token_pool/sources/lock_release_token_pool.move --output ./bindings/generated/ccip/ccip_token_pools/lock_release_token_pool

# CCIP - Burn Mint Token Pool
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip_token_pools/burn_mint_token_pool --input ./contracts/ccip/ccip_token_pools/burn_mint_token_pool/sources/burn_mint_token_pool.move --output ./bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool

# CCIP - Managed Token Pool
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip_token_pools/managed_token_pool --input ./contracts/ccip/ccip_token_pools/managed_token_pool/sources/managed_token_pool.move --output ./bindings/generated/ccip/ccip_token_pools/managed_token_pool

# CCIP - USDCTokenPool
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip_token_pools/usdc_token_pool --input ./contracts/ccip/ccip_token_pools/usdc_token_pool/sources/usdc_token_pool.move --output ./bindings/generated/ccip/ccip_token_pools/usdc_token_pool

# CCIP Router
go run bindgen/main.go --moveConfig ./contracts/ccip/ccip_router --input ./contracts/ccip/ccip_router/sources/router.move --output ./bindings/generated/ccip/ccip_router/

# MCMS
go run bindgen/main.go --moveConfig ./contracts/mcms/mcms --input ./contracts/mcms/mcms/sources/mcms.move --output ./bindings/generated/mcms/mcms
go run bindgen/main.go --moveConfig ./contracts/mcms/mcms --input ./contracts/mcms/mcms/sources/mcms_account.move --output ./bindings/generated/mcms/mcms_account
go run bindgen/main.go --moveConfig ./contracts/mcms/mcms --input ./contracts/mcms/mcms/sources/mcms_deployer.move --output ./bindings/generated/mcms/mcms_deployer
go run bindgen/main.go --moveConfig ./contracts/mcms/mcms --input ./contracts/mcms/mcms/sources/mcms_registry.move --output ./bindings/generated/mcms/mcms_registry
