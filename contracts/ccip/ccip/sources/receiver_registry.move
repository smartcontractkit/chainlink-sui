module ccip::receiver_registry {
    use std::type_name::{Self, TypeName};
    use std::string::{Self, String};

    use sui::event;
    use sui::vec_map::{Self, VecMap};

    use ccip::state_object::{Self, CCIPObjectRef, OwnerCap};

    public struct FunctionInfo has store, copy, drop {
        module_address: address,
        module_name: String,
        function_name: String
    }

    public struct ReceiverConfig has store, copy, drop {
        ccip_receive_function: FunctionInfo,
        proof_typename: TypeName,
    }

    public struct ReceiverRegistry has key, store {
        id: UID,
        receiver_configs: VecMap<address, ReceiverConfig>
    }

    public struct ReceiverRegistered has copy, drop {
        receiver_address: address,
        receiver_module_name: vector<u8>
    }

    const RECEIVER_REGISTRY: vector<u8> = b"ReceiverRegistry";
    const E_ALREADY_REGISTERED: u64 = 1;
    const E_ALREADY_INITIALIZED: u64 = 2;

    public fun type_and_version(): String {
        string::utf8(b"ReceiverRegistry 1.6.0")
    }

    public fun intialize(
        ref: &mut CCIPObjectRef,
        _: &OwnerCap,
        ctx: &mut TxContext
    ) {
        assert!(
            !state_object::contains(ref, RECEIVER_REGISTRY),
            E_ALREADY_INITIALIZED
        );
        let state = ReceiverRegistry {
            id: object::new(ctx),
            receiver_configs: vec_map::empty()
        };

        state_object::add(ref, RECEIVER_REGISTRY, state, ctx);
    }

    public fun register_receiver<ProofType: drop>(
        ref: &mut CCIPObjectRef,
        receiver_module_name: vector<u8>,
        _proof: ProofType,
        ctx: &mut TxContext
    ) {
        let receiver_address = ctx.sender();
        let registry = state_object::borrow_mut_with_ctx<ReceiverRegistry>(ref, RECEIVER_REGISTRY, ctx);
        assert!(!registry.receiver_configs.contains(&receiver_address), E_ALREADY_REGISTERED);

        let ccip_receive_function =
            FunctionInfo {
                module_address: receiver_address,
                module_name: string::utf8(receiver_module_name),
                function_name: string::utf8(b"ccip_receive")
            };
        let proof_typename = type_name::get<ProofType>();

        let receiver_config = ReceiverConfig {
            ccip_receive_function,
            proof_typename,
        };
        registry.receiver_configs.insert(receiver_address, receiver_config);

        event::emit(ReceiverRegistered { receiver_address, receiver_module_name });
    }

    public fun is_registered_receiver(ref: &CCIPObjectRef, receiver_address: address): bool {
        let registry = state_object::borrow<ReceiverRegistry>(ref, RECEIVER_REGISTRY);
        return registry.receiver_configs.contains(&receiver_address)
    }
}