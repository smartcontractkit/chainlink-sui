module test::generics {

    // Generic struct with single type parameter
    public struct Box<T: store> has key, store {
        id: UID,
        value: T,
    }

    // Generic struct with phantom type parameter (like Coin)
    public struct Token<phantom T> has key, store {
        id: UID,
        balance: u64,
    }

    // Generic struct with multiple type parameters
    public struct Pair<T: store + drop, U: store + drop> has store, drop {
        first: T,
        second: U,
    }

    // Function with generic type parameter
    public fun create_box<T: store>(value: T, ctx: &mut TxContext): Box<T> {
        Box {
            id: object::new(ctx),
            value
        }
    }

    // Function returning generic type
    public fun unbox<T: store>(box: Box<T>): T {
        let Box { id, value } = box;
        object::delete(id);
        value
    }

    public fun deposit<T>(token: &mut Token<T>, coin: Token<T>) {
        let Token { id, balance } = coin;
        object::delete(id);
        token.balance = token.balance + balance;
    }

    // Function returning amount from generic token
    public fun balance<T>(token: &Token<T>): u64 {
        token.balance
    }

    // Function with multiple generic parameters
    public fun create_pair<T: store + drop, U: store + drop>(first: T, second: U): Pair<T, U> {
        Pair { first, second }
    }

    // Function that uses specific instantiation
    public fun create_sui_token(ctx: &mut TxContext): Token<sui::sui::SUI> {
        Token {
            id: object::new(ctx),
            balance: 0,
        }
    }

    // Entry function that creates and transfers a SUI token
    public entry fun create_and_transfer_sui_token(ctx: &mut TxContext) {
        let token = Token<sui::sui::SUI> {
            id: object::new(ctx),
            balance: 0,
        };
        transfer::transfer(token, ctx.sender());
    }

    // Generic function that creates and transfers a token of any type
    public fun create_and_transfer_token<T>(ctx: &mut TxContext) {
        let token = Token<T> {
            id: object::new(ctx),
            balance: 0,
        };
        transfer::transfer(token, ctx.sender());
    }
    
    // Test function with non-phantom generic that creates and transfers a box
    public fun create_and_transfer_box<T: store>(value: T, ctx: &mut TxContext) {
        let box = Box {
            id: object::new(ctx),
            value
        };
        transfer::transfer(box, ctx.sender());
    }
}
