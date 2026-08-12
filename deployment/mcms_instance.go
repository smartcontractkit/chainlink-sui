package deployment

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

// MCMSInstance identifies which MCMS deployment on a chain an operation targets.
// Slow and fastcurse instances share the same address-book type constants and are
// distinguished by the fastcurse label on address-book entries.
type MCMSInstance string

const (
	MCMSInstanceSlow      MCMSInstance = "slow"
	MCMSInstanceFastCurse MCMSInstance = MCMSFastCurseLabel
)

// MCMSInstanceFromFastCurseFlag maps the legacy isFastCurse config flag to an instance.
func MCMSInstanceFromFastCurseFlag(isFastCurse bool) MCMSInstance {
	if isFastCurse {
		return MCMSInstanceFastCurse
	}
	return MCMSInstanceSlow
}

// AddressBookLabel returns the address-book label for this instance.
// The slow instance uses no label.
func (i MCMSInstance) AddressBookLabel() string {
	if i == MCMSInstanceFastCurse {
		return MCMSFastCurseLabel
	}
	return ""
}

// IsFastCurse reports whether this instance is the fastcurse MCMS.
func (i MCMSInstance) IsFastCurse() bool {
	return i == MCMSInstanceFastCurse
}

func (i MCMSInstance) String() string {
	return string(i)
}

// HasMCMSInstance reports whether the chain state has a deployed MCMS for the instance.
func (s CCIPChainState) HasMCMSInstance(instance MCMSInstance) bool {
	fields := s.MCMSStateByInstance(instance)
	return fields.StateObjectID != "" && fields.RegistryObjectID != ""
}

// MCMSStateByInstance returns the MCMS object IDs for the requested instance.
func (s CCIPChainState) MCMSStateByInstance(instance MCMSInstance) MCMSStateFields {
	switch instance {
	case MCMSInstanceFastCurse:
		return MCMSStateFields{
			PackageID:               s.FastCurseMCMSPackageID,
			StateObjectID:           s.FastCurseMCMSStateObjectID,
			RegistryObjectID:        s.FastCurseMCMSRegistryObjectID,
			DeployerStateObjectID:   s.FastCurseMCMSDeployerStateObjectID,
			AccountStateObjectID:    s.FastCurseMCMSAccountStateObjectID,
			AccountOwnerCapObjectID: s.FastCurseMCMSAccountOwnerCapObjectID,
			TimelockObjectID:        s.FastCurseMCMSTimelockObjectID,
		}
	case MCMSInstanceSlow:
		return MCMSStateFields{
			PackageID:               s.MCMSPackageID,
			StateObjectID:           s.MCMSStateObjectID,
			RegistryObjectID:        s.MCMSRegistryObjectID,
			DeployerStateObjectID:   s.MCMSDeployerStateObjectID,
			AccountStateObjectID:    s.MCMSAccountStateObjectID,
			AccountOwnerCapObjectID: s.MCMSAccountOwnerCapObjectID,
			TimelockObjectID:        s.MCMSTimelockObjectID,
		}
	default:
		return MCMSStateFields{}
	}
}

// StoreMCMSInAddressBook saves one MCMS deployment's seven object IDs to both the
// address book and the datastore under the correct instance label, so slow and
// fastcurse entries can coexist on one chain. Each entry is dual-written via
// SaveSuiAddress.
func StoreMCMSInAddressBook(ab *cldf.AddressBookMap, ds fdatastore.MutableAddressRefStore, chainSelector uint64, mcmsReport mcmsops.DeployMCMSSeqOutput, instance MCMSInstance) error {
	addLabel := func(tv cldf.TypeAndVersion) cldf.TypeAndVersion {
		if label := instance.AddressBookLabel(); label != "" {
			tv.Labels.Add(label)
		}
		return tv
	}

	save := func(addr string, typ cldf.ContractType) error {
		return SaveSuiAddress(ab, ds, chainSelector, addr, addLabel(cldf.NewTypeAndVersion(typ, Version1_0_0)))
	}

	if err := save(mcmsReport.PackageId, SuiMcmsPackageIDType); err != nil {
		return fmt.Errorf("save MCMS package for %s instance on chain %d: %w", instance, chainSelector, err)
	}
	if err := save(mcmsReport.Objects.McmsMultisigStateObjectId, SuiMcmsObjectIDType); err != nil {
		return fmt.Errorf("save MCMS state object for %s instance on chain %d: %w", instance, chainSelector, err)
	}
	if err := save(mcmsReport.Objects.McmsRegistryObjectId, SuiMcmsRegistryObjectIDType); err != nil {
		return fmt.Errorf("save MCMS registry for %s instance on chain %d: %w", instance, chainSelector, err)
	}
	if err := save(mcmsReport.Objects.McmsAccountStateObjectId, SuiMcmsAccountStateObjectIDType); err != nil {
		return fmt.Errorf("save MCMS account state for %s instance on chain %d: %w", instance, chainSelector, err)
	}
	if err := save(mcmsReport.Objects.McmsAccountOwnerCapObjectId, SuiMcmsAccountOwnerCapObjectIDType); err != nil {
		return fmt.Errorf("save MCMS account owner cap for %s instance on chain %d: %w", instance, chainSelector, err)
	}
	if err := save(mcmsReport.Objects.TimelockObjectId, SuiMcmsTimelockObjectIDType); err != nil {
		return fmt.Errorf("save MCMS timelock for %s instance on chain %d: %w", instance, chainSelector, err)
	}
	if err := save(mcmsReport.Objects.McmsDeployerStateObjectId, SuiMcmsDeployerObjectIDType); err != nil {
		return fmt.Errorf("save MCMS deployer state for %s instance on chain %d: %w", instance, chainSelector, err)
	}
	return nil
}
