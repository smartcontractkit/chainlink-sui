module ccip_dummy_receiver::dummy_receiver {
    use std::string::{Self, String};
    use sui::event;

    use ccip::client;
    use ccip::dynamic_dispatcher as dd;
    use ccip::receiver_registry;
    use ccip::state_object::CCIPObjectRef;

    public struct OwnerCap has key, store {
        id: UID,
        receiver_address: address,
    }

    public struct ReceivedMessage has copy, drop {
        message_id: vector<u8>,
        source_chain_selector: u64,
        sender: vector<u8>,
        data: vector<u8>,
        dest_token_transfer_length: u64,
    }

    public struct CCIPReceiverState has key {
        id: UID,
        counter: u64,
        message_id: vector<u8>,
        source_chain_selector: u64,
        sender: vector<u8>,
        data: vector<u8>,
        dest_token_transfer_length: u64,
    }

    public struct DummyReceiverProof has drop {}

    public fun type_and_version(): String {
        string::utf8(b"DummyReceiver 1.6.0")
    }

    fun init(ctx: &mut TxContext) {
       let state = CCIPReceiverState {
           id: object::new(ctx),
           counter: 0,
           message_id: vector[],
           source_chain_selector: 0,
           sender: vector[],
           data: vector[],
           dest_token_transfer_length: 0,
        };

        let owner_cap = OwnerCap {
            id: object::new(ctx),
            receiver_address: object::id_to_address(object::borrow_id(&state)),
        };

        transfer::share_object(state);
        transfer::transfer(owner_cap, ctx.sender());
    }

    // this can be gated by owner cap
    public fun register_receiver(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        receiver_registry::register_receiver(
            ref, b"dummy_receiver", DummyReceiverProof {}, ctx
        );
    }

    public fun get_counter(state: &CCIPReceiverState): u64 {
        state.counter
    }

    // this requires the caller (DON or user) to extract message from potato
    public fun ccip_receive(state: &mut CCIPReceiverState, message: client::Any2SuiMessage) {
        state.counter = state.counter + 1;
        state.message_id = client::get_message_id(&message);
        state.source_chain_selector = client::get_source_chain_selector(&message);
        state.sender = client::get_sender(&message);
        state.data = client::get_data(&message);
        state.dest_token_transfer_length = client::get_dest_token_amounts(&message).length();

        event::emit(
            ReceivedMessage {
                message_id: state.message_id,
                source_chain_selector: state.source_chain_selector,
                sender: state.sender,
                data: state.data,
                dest_token_transfer_length: state.dest_token_transfer_length,
            }
        );
    }

    // TODO:
    // this requires the receiver to extract the msg from the potato. is this any safer?
    // however, this means the receiver will always get called no matter if the message is presented
    // and this function must return the hot potato.
    public fun ccip_receive_potato(state: &mut CCIPReceiverState, receiver_params: dd::ReceiverParams): dd::ReceiverParams {
        let (message_op, receiver_params) = dd::extract_any2sui_message(receiver_params);
        if (message_op.is_none()) {
            return receiver_params
        };
        let message = message_op.borrow();
        state.counter = state.counter + 1;
        state.message_id = client::get_message_id(message);
        state.source_chain_selector = client::get_source_chain_selector(message);
        state.sender = client::get_sender(message);
        state.data = client::get_data(message);
        state.dest_token_transfer_length = client::get_dest_token_amounts(message).length();

        event::emit(
            ReceivedMessage {
                message_id: state.message_id,
                source_chain_selector: state.source_chain_selector,
                sender: state.sender,
                data: state.data,
                dest_token_transfer_length: state.dest_token_transfer_length,
            }
        );
        receiver_params
    }
}
