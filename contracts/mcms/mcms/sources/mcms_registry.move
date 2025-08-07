module mcms::mcms_registry {

    use mcms::params;
    use std::string::String;
    use std::type_name::{Self, TypeName};
    use sui::bag::{Self, Bag};
    use sui::event;
    use sui::address;
    
    public struct Registry has key {
        id: UID,
        /// Maps account address -> package cap
        /// Only one cap per account address/package
        package_caps: Bag,
    }

    /// `ExecutingCallbackParams` is created when an operation is ready to be executed from MCMS
    public struct ExecutingCallbackParams {
        target: address,
        module_name: String,
        function_name: String,
        data: vector<u8>,
    }

    public struct EntrypointRegistered has copy, drop {
        registry_id: ID,
        account_address: address,
        module_name: String,
    }

    const EPackageCapAlreadyRegistered: u64 = 1;
    const EPackageCapNotRegistered: u64 = 2;
    const EPackageIdMismatch: u64 = 3;
    const EModuleNameMismatch: u64 = 4;

    public struct MCMS_REGISTRY has drop {}

    fun init(_witness: MCMS_REGISTRY, ctx: &mut TxContext) {
        let registry = Registry {
            id: object::new(ctx),
            package_caps: bag::new(ctx),
        };

        transfer::share_object(registry);
    }

    public fun register_entrypoint<T: drop, C: key + store>(
        registry: &mut Registry,
        _proof: T,
        package_cap: C,
        _ctx: &TxContext,
    ) {
        let proof_type = type_name::get<T>();
        let (proof_account_address, proof_module_name) = params::get_account_address_and_module_name(
            proof_type,
        );

        assert!(!registry.package_caps.contains(proof_account_address), EPackageCapAlreadyRegistered);

        // Register package cap for package address
        registry.package_caps.add(proof_account_address, package_cap);

        event::emit(EntrypointRegistered {
            registry_id: object::id(registry),
            account_address: proof_account_address,
            module_name: proof_module_name,
        });
    }

    public fun get_callback_params<T: drop, C: key + store>(
        registry: &mut Registry,
        _proof: T,
        params: ExecutingCallbackParams,
    ): (&C, String, vector<u8>) {
        let ExecutingCallbackParams { target, module_name, function_name, data } = params;

        let proof_type = type_name::get<T>();
        let (proof_account_address, proof_module_name) = params::get_account_address_and_module_name(
            proof_type,
        );

        assert!(target == proof_account_address, EPackageIdMismatch);
        assert!(module_name == proof_module_name, EModuleNameMismatch);

        // Validate the proof comes from same package ID
        assert!(registry.package_caps.contains(proof_account_address), EPackageCapNotRegistered);

        let package_cap = registry.package_caps.borrow(proof_account_address);
        (package_cap, function_name, data)
    }

    public fun release_cap<T: drop, C: key + store>(registry: &mut Registry, _witness: T): C {
        let proof_type = type_name::get<T>();
        let (proof_account_address, _) = params::get_account_address_and_module_name(
            proof_type,
        );

        assert!(registry.package_caps.contains(proof_account_address), EPackageCapNotRegistered);

        registry.package_caps.remove(proof_account_address)
    }

    public(package) fun borrow_owner_cap<C: key + store>(registry: &Registry): &C {
        registry.package_caps.borrow(get_multisig_address())
    }

    public fun get_callback_params_for_mcms<T: drop>(
        params: ExecutingCallbackParams,
        _proof: T,
    ): (address, String, String, vector<u8>) {
        let ExecutingCallbackParams { target, module_name, function_name, data } = params;

        let proof_type = type_name::get<T>();
        let (proof_account_address, proof_module_name) = params::get_account_address_and_module_name(
            proof_type,
        );

        assert!(target == proof_account_address, EPackageIdMismatch);
        assert!(module_name == proof_module_name, EModuleNameMismatch);

        (target, module_name, function_name, data)
    }

    public(package) fun get_callback_params_from_mcms(
        params: ExecutingCallbackParams,
    ): (address, String, String, vector<u8>) {
        let ExecutingCallbackParams { target, module_name, function_name, data } = params;
        (target, module_name, function_name, data)
    }

    public(package) fun create_executing_callback_params(
        target: address,
        module_name: String,
        function_name: String,
        data: vector<u8>,
    ): ExecutingCallbackParams {
        ExecutingCallbackParams {
            target,
            module_name,
            function_name,
            data,
        }
    }

    public fun is_package_registered(registry: &Registry, package_address: address): bool {
        registry.package_caps.contains(package_address)
    }

    public fun target(params: &ExecutingCallbackParams): address {
        params.target
    }

    public fun module_name(params: &ExecutingCallbackParams): String {
        params.module_name
    }

    public fun function_name(params: &ExecutingCallbackParams): String {
        params.function_name
    }

    public fun data(params: &ExecutingCallbackParams): vector<u8> {
        params.data
    }

    public fun get_multisig_address(): address {
        address::from_ascii_bytes(&type_name::get<McmsProof>().get_address().into_bytes())
    }

    public struct McmsProof has drop {}

    public(package) fun create_mcms_proof(): McmsProof {
        McmsProof {}
    }

    // ===================== TESTS =====================

    #[test_only]
    /// Initialize the registry for testing
    public fun test_init(ctx: &mut TxContext) {
        init(MCMS_REGISTRY {}, ctx)
    }

    #[test_only]
    /// Create executing callback params for testing
    public fun test_create_executing_callback_params(
        target: address,
        module_name: String,
        function_name: String,
        data: vector<u8>,
    ): ExecutingCallbackParams {
        create_executing_callback_params(target, module_name, function_name, data)
    }

}