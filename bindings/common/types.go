package common

import "github.com/smartcontractkit/chainlink-sui/bindings/bind"

// Common types used across CCIP packages. These are objects in Move, in Bindings only references
type CCIPObjectRef = bind.Object
type OwnerCap = bind.Object
type NonceManagerCap = bind.Object
type SourceTransferCap = bind.Object
type DestTransferCap = bind.Object
type FeeQuoterCap = bind.Object
type TokenParams = bind.Object
