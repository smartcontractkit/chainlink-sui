# smartcontractkit Go modules
## Main module
```mermaid
flowchart LR

	chain-selectors
	click chain-selectors href "https://github.com/smartcontractkit/chain-selectors"
	chainlink-aptos --> chainlink-aptos/codec
	chainlink-aptos --> chainlink-framework/metrics
	click chainlink-aptos href "https://github.com/smartcontractkit/chainlink-aptos"
	chainlink-aptos/codec --> chainlink-common
	click chainlink-aptos/codec href "https://github.com/smartcontractkit/chainlink-aptos"
	chainlink-ccip --> chainlink-common
	chainlink-ccip --> chainlink-protos/rmn/v1.6/go
	click chainlink-ccip href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-common --> chainlink-common/pkg/chipingress
	chainlink-common --> chainlink-protos/billing/go
	chainlink-common --> chainlink-protos/cre/go
	chainlink-common --> chainlink-protos/linking-service/go
	chainlink-common --> chainlink-protos/metering/go
	chainlink-common --> chainlink-protos/node-platform
	chainlink-common --> chainlink-protos/storage-service
	chainlink-common --> chainlink-protos/workflows/go
	chainlink-common --> freeport
	chainlink-common --> grpc-proxy
	chainlink-common --> libocr
	click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-common/pkg/chipingress
	click chainlink-common/pkg/chipingress href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-framework/metrics --> chainlink-common
	click chainlink-framework/metrics href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-protos/billing/go
	click chainlink-protos/billing/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/cre/go --> chain-selectors
	click chainlink-protos/cre/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/linking-service/go
	click chainlink-protos/linking-service/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/metering/go
	click chainlink-protos/metering/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/node-platform
	click chainlink-protos/node-platform href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/rmn/v1.6/go
	click chainlink-protos/rmn/v1.6/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/storage-service
	click chainlink-protos/storage-service href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/workflows/go
	click chainlink-protos/workflows/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-sui --> chainlink-aptos
	chainlink-sui --> chainlink-ccip
	chainlink-sui --> chainlink-sui/codec
	click chainlink-sui href "https://github.com/smartcontractkit/chainlink-sui"
	chainlink-sui/codec --> chainlink-aptos/codec
	click chainlink-sui/codec href "https://github.com/smartcontractkit/chainlink-sui"
	freeport
	click freeport href "https://github.com/smartcontractkit/freeport"
	go-sumtype2
	click go-sumtype2 href "https://github.com/smartcontractkit/go-sumtype2"
	grpc-proxy
	click grpc-proxy href "https://github.com/smartcontractkit/grpc-proxy"
	libocr --> go-sumtype2
	click libocr href "https://github.com/smartcontractkit/libocr"

	subgraph chainlink-aptos-repo[chainlink-aptos]
		 chainlink-aptos
		 chainlink-aptos/codec
	end
	click chainlink-aptos-repo href "https://github.com/smartcontractkit/chainlink-aptos"

	subgraph chainlink-common-repo[chainlink-common]
		 chainlink-common
		 chainlink-common/pkg/chipingress
	end
	click chainlink-common-repo href "https://github.com/smartcontractkit/chainlink-common"

	subgraph chainlink-protos-repo[chainlink-protos]
		 chainlink-protos/billing/go
		 chainlink-protos/cre/go
		 chainlink-protos/linking-service/go
		 chainlink-protos/metering/go
		 chainlink-protos/node-platform
		 chainlink-protos/rmn/v1.6/go
		 chainlink-protos/storage-service
		 chainlink-protos/workflows/go
	end
	click chainlink-protos-repo href "https://github.com/smartcontractkit/chainlink-protos"

	subgraph chainlink-sui-repo[chainlink-sui]
		 chainlink-sui
		 chainlink-sui/codec
	end
	click chainlink-sui-repo href "https://github.com/smartcontractkit/chainlink-sui"

	classDef outline stroke-dasharray:6,fill:none;
	class chainlink-aptos-repo,chainlink-common-repo,chainlink-protos-repo,chainlink-sui-repo outline
```
## All modules
```mermaid
flowchart LR

	ccip-owner-contracts
	click ccip-owner-contracts href "https://github.com/smartcontractkit/ccip-owner-contracts"
	chain-selectors
	click chain-selectors href "https://github.com/smartcontractkit/chain-selectors"
	chainlink-aptos --> chainlink-aptos/codec
	chainlink-aptos --> chainlink-framework/metrics
	click chainlink-aptos href "https://github.com/smartcontractkit/chainlink-aptos"
	chainlink-aptos/codec --> chainlink-common
	click chainlink-aptos/codec href "https://github.com/smartcontractkit/chainlink-aptos"
	chainlink-canton
	click chainlink-canton href "https://github.com/smartcontractkit/chainlink-canton"
	chainlink-canton/authentication
	click chainlink-canton/authentication href "https://github.com/smartcontractkit/chainlink-canton"
	chainlink-canton/contracts/v2 --> go-daml
	click chainlink-canton/contracts/v2 href "https://github.com/smartcontractkit/chainlink-canton"
	chainlink-ccip --> chainlink-common
	chainlink-ccip --> chainlink-protos/rmn/v1.6/go
	click chainlink-ccip href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-ccip/chains/evm
	click chainlink-ccip/chains/evm href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-ccip/chains/solana --> chainlink-ccip/chains/solana/gobindings
	chainlink-ccip/chains/solana --> chainlink-common
	click chainlink-ccip/chains/solana href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-ccip/chains/solana/gobindings
	click chainlink-ccip/chains/solana/gobindings href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-ccip/deployment
	click chainlink-ccip/deployment href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-common --> chainlink-common/pkg/chipingress
	chainlink-common --> chainlink-protos/billing/go
	chainlink-common --> chainlink-protos/cre/go
	chainlink-common --> chainlink-protos/linking-service/go
	chainlink-common --> chainlink-protos/metering/go
	chainlink-common --> chainlink-protos/node-platform
	chainlink-common --> chainlink-protos/storage-service
	chainlink-common --> chainlink-protos/workflows/go
	chainlink-common --> freeport
	chainlink-common --> grpc-proxy
	chainlink-common --> libocr
	click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-common/pkg/chipingress
	click chainlink-common/pkg/chipingress href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-common/pkg/monitoring
	click chainlink-common/pkg/monitoring href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-deployments-framework --> ccip-owner-contracts
	chainlink-deployments-framework --> chainlink-ccip/chains/evm
	chainlink-deployments-framework --> chainlink-evm/gethwrappers
	chainlink-deployments-framework --> chainlink-protos/op-catalog
	chainlink-deployments-framework --> chainlink-testing-framework/seth
	chainlink-deployments-framework --> mcms
	click chainlink-deployments-framework href "https://github.com/smartcontractkit/chainlink-deployments-framework"
	chainlink-evm/gethwrappers
	click chainlink-evm/gethwrappers href "https://github.com/smartcontractkit/chainlink-evm"
	chainlink-framework/metrics --> chainlink-common
	click chainlink-framework/metrics href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-protos/billing/go
	click chainlink-protos/billing/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/cre/go --> chain-selectors
	click chainlink-protos/cre/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/job-distributor
	click chainlink-protos/job-distributor href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/linking-service/go
	click chainlink-protos/linking-service/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/metering/go
	click chainlink-protos/metering/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/node-platform
	click chainlink-protos/node-platform href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/op-catalog
	click chainlink-protos/op-catalog href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/rmn/v1.6/go
	click chainlink-protos/rmn/v1.6/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/storage-service
	click chainlink-protos/storage-service href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/workflows/go
	click chainlink-protos/workflows/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-stellar
	click chainlink-stellar href "https://github.com/smartcontractkit/chainlink-stellar"
	chainlink-stellar/bindings
	click chainlink-stellar/bindings href "https://github.com/smartcontractkit/chainlink-stellar"
	chainlink-sui --> chainlink-aptos
	chainlink-sui --> chainlink-ccip
	chainlink-sui --> chainlink-sui/codec
	click chainlink-sui href "https://github.com/smartcontractkit/chainlink-sui"
	chainlink-sui/codec --> chainlink-aptos/codec
	click chainlink-sui/codec href "https://github.com/smartcontractkit/chainlink-sui"
	chainlink-sui/deployment
	click chainlink-sui/deployment href "https://github.com/smartcontractkit/chainlink-sui"
	chainlink-sui/integration-tests --> chainlink-sui/deployment
	click chainlink-sui/integration-tests href "https://github.com/smartcontractkit/chainlink-sui"
	chainlink-testing-framework/framework
	click chainlink-testing-framework/framework href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/seth
	click chainlink-testing-framework/seth href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-ton --> chainlink-ccip
	chainlink-ton --> chainlink-common/pkg/monitoring
	chainlink-ton --> chainlink-framework/metrics
	chainlink-ton --> chainlink-ton/cciplib
	click chainlink-ton href "https://github.com/smartcontractkit/chainlink-ton"
	chainlink-ton/cciplib --> chainlink-common
	click chainlink-ton/cciplib href "https://github.com/smartcontractkit/chainlink-ton"
	chainlink-tron/relayer --> chainlink-common
	click chainlink-tron/relayer href "https://github.com/smartcontractkit/chainlink-tron"
	freeport
	click freeport href "https://github.com/smartcontractkit/freeport"
	go-daml --> freeport
	click go-daml href "https://github.com/smartcontractkit/go-daml"
	go-sumtype2
	click go-sumtype2 href "https://github.com/smartcontractkit/go-sumtype2"
	grpc-proxy
	click grpc-proxy href "https://github.com/smartcontractkit/grpc-proxy"
	libocr --> go-sumtype2
	click libocr href "https://github.com/smartcontractkit/libocr"
	mcms --> chainlink-canton
	mcms --> chainlink-canton/authentication
	mcms --> chainlink-canton/contracts/v2
	mcms --> chainlink-ccip/chains/solana
	mcms --> chainlink-ccip/deployment
	mcms --> chainlink-deployments-framework
	mcms --> chainlink-protos/job-distributor
	mcms --> chainlink-stellar
	mcms --> chainlink-stellar/bindings
	mcms --> chainlink-sui
	mcms --> chainlink-testing-framework/framework
	mcms --> chainlink-ton
	mcms --> chainlink-tron/relayer
	click mcms href "https://github.com/smartcontractkit/mcms"

	subgraph chainlink-aptos-repo[chainlink-aptos]
		 chainlink-aptos
		 chainlink-aptos/codec
	end
	click chainlink-aptos-repo href "https://github.com/smartcontractkit/chainlink-aptos"

	subgraph chainlink-canton-repo[chainlink-canton]
		 chainlink-canton
		 chainlink-canton/authentication
		 chainlink-canton/contracts/v2
	end
	click chainlink-canton-repo href "https://github.com/smartcontractkit/chainlink-canton"

	subgraph chainlink-ccip-repo[chainlink-ccip]
		 chainlink-ccip
		 chainlink-ccip/chains/evm
		 chainlink-ccip/chains/solana
		 chainlink-ccip/chains/solana/gobindings
		 chainlink-ccip/deployment
	end
	click chainlink-ccip-repo href "https://github.com/smartcontractkit/chainlink-ccip"

	subgraph chainlink-common-repo[chainlink-common]
		 chainlink-common
		 chainlink-common/pkg/chipingress
		 chainlink-common/pkg/monitoring
	end
	click chainlink-common-repo href "https://github.com/smartcontractkit/chainlink-common"

	subgraph chainlink-protos-repo[chainlink-protos]
		 chainlink-protos/billing/go
		 chainlink-protos/cre/go
		 chainlink-protos/job-distributor
		 chainlink-protos/linking-service/go
		 chainlink-protos/metering/go
		 chainlink-protos/node-platform
		 chainlink-protos/op-catalog
		 chainlink-protos/rmn/v1.6/go
		 chainlink-protos/storage-service
		 chainlink-protos/workflows/go
	end
	click chainlink-protos-repo href "https://github.com/smartcontractkit/chainlink-protos"

	subgraph chainlink-stellar-repo[chainlink-stellar]
		 chainlink-stellar
		 chainlink-stellar/bindings
	end
	click chainlink-stellar-repo href "https://github.com/smartcontractkit/chainlink-stellar"

	subgraph chainlink-sui-repo[chainlink-sui]
		 chainlink-sui
		 chainlink-sui/codec
		 chainlink-sui/deployment
		 chainlink-sui/integration-tests
	end
	click chainlink-sui-repo href "https://github.com/smartcontractkit/chainlink-sui"

	subgraph chainlink-testing-framework-repo[chainlink-testing-framework]
		 chainlink-testing-framework/framework
		 chainlink-testing-framework/seth
	end
	click chainlink-testing-framework-repo href "https://github.com/smartcontractkit/chainlink-testing-framework"

	subgraph chainlink-ton-repo[chainlink-ton]
		 chainlink-ton
		 chainlink-ton/cciplib
	end
	click chainlink-ton-repo href "https://github.com/smartcontractkit/chainlink-ton"

	classDef outline stroke-dasharray:6,fill:none;
	class chainlink-aptos-repo,chainlink-canton-repo,chainlink-ccip-repo,chainlink-common-repo,chainlink-protos-repo,chainlink-stellar-repo,chainlink-sui-repo,chainlink-testing-framework-repo,chainlink-ton-repo outline
```
