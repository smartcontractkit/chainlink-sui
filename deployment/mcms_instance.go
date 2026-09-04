package deployment

import (
	"fmt"
	"strings"

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

	// CLLCCIPQualifier and RMNMCMSQualifier are the datastore qualifiers for the Sui MCMS
	// instances: slow -> CLLCCIP, fastcurse -> RMNMCMS.
	CLLCCIPQualifier = "CLLCCIP"
	RMNMCMSQualifier = "RMNMCMS"
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

// DatastoreQualifier returns the datastore qualifier for this instance:
// slow -> CLLCCIP, fastcurse -> RMNMCMS. The fastcurse
// address-book label is still added separately (see AddressBookLabel) so the Sui
// loader, which special-cases it, keeps working unchanged.
func (i MCMSInstance) DatastoreQualifier() string {
	if i == MCMSInstanceFastCurse {
		return RMNMCMSQualifier
	}
	return CLLCCIPQualifier
}

// MCMSInstanceFromQualifier is the inverse of DatastoreQualifier: it recovers the instance
// from a datastore ref's qualifier, reporting false for a qualifier that names no instance.
//
// It lives next to DatastoreQualifier so the two directions cannot drift; the round trip is
// pinned by test. MCMSFastCurseLabel is accepted because rows written before the purpose
// qualifiers existed carry "fastcurse" in the qualifier field.
func MCMSInstanceFromQualifier(qualifier string) (MCMSInstance, bool) {
	switch strings.TrimSpace(qualifier) {
	case RMNMCMSQualifier, MCMSFastCurseLabel:
		return MCMSInstanceFastCurse, true
	case CLLCCIPQualifier:
		return MCMSInstanceSlow, true
	default:
		return "", false
	}
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

// mcmsObjectTypes maps one MCMS deployment's report to the seven (contract type, object ID)
// pairs recorded for it. StoreMCMSInAddressBook writes exactly these and PlannedMCMSRefs
// plans exactly these, so the write and the pre-deploy conflict check cannot drift.
func mcmsObjectTypes(mcmsReport mcmsops.DeployMCMSSeqOutput) []struct {
	name    string
	typ     cldf.ContractType
	address string
} {
	return []struct {
		name    string
		typ     cldf.ContractType
		address string
	}{
		{"package", SuiMcmsPackageIDType, mcmsReport.PackageId},
		{"state object", SuiMcmsObjectIDType, mcmsReport.Objects.McmsMultisigStateObjectId},
		{"registry", SuiMcmsRegistryObjectIDType, mcmsReport.Objects.McmsRegistryObjectId},
		{"account state", SuiMcmsAccountStateObjectIDType, mcmsReport.Objects.McmsAccountStateObjectId},
		{"account owner cap", SuiMcmsAccountOwnerCapObjectIDType, mcmsReport.Objects.McmsAccountOwnerCapObjectId},
		{"timelock", SuiMcmsTimelockObjectIDType, mcmsReport.Objects.TimelockObjectId},
		{"deployer state", SuiMcmsDeployerObjectIDType, mcmsReport.Objects.McmsDeployerStateObjectId},
	}
}

// PlannedMCMSRefs returns the datastore keys a deployment of this MCMS instance will occupy,
// so a changeset's VerifyPreconditions can check them before anything is deployed.
func PlannedMCMSRefs(instance MCMSInstance) []PlannedRef {
	qualifier := instance.DatastoreQualifier()
	objects := mcmsObjectTypes(mcmsops.DeployMCMSSeqOutput{})
	planned := make([]PlannedRef, 0, len(objects))
	for _, o := range objects {
		planned = append(planned, PlannedRef{Type: o.typ, Qualifier: qualifier})
	}
	return planned
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

	for _, o := range mcmsObjectTypes(mcmsReport) {
		if err := SaveSuiAddress(ab, ds, chainSelector, o.address, addLabel(cldf.NewTypeAndVersion(o.typ, Version1_0_0)), instance.DatastoreQualifier()); err != nil {
			return fmt.Errorf("save MCMS %s for %s instance on chain %d: %w", o.name, instance, chainSelector, err)
		}
	}
	return nil
}
