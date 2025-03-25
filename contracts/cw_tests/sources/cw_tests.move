/// Module: cw_tests
module cw_tests::cw_tests {

    use std::string::String;

    /// List of todos. Can be managed by the owner and shared with others.
    public struct TodoList has key, store {
        id: UID,
        items: vector<String>
    }

    fun init(ctx: &mut TxContext) {
        // Initialize the contract with an empty todo list.
        let lst = TodoList{
            id: object::new(ctx),
            items: vector[]
        };

        transfer::share_object(lst);
    }

    /// Add a new todo item to the list.
    public entry fun add(list: &mut TodoList, item: String) {
        list.items.push_back(item);
    }

    public entry fun replace_items(list: &mut TodoList, items: vector<String>) {
        list.items = items;
    }

    /// Remove a todo item from the list by index.
    public fun remove(list: &mut TodoList, index: u64): String {
        list.items.remove(index)
    }

    /// Delete the list and the capability to manage it.
    public fun delete(list: TodoList) {
        let TodoList { id, items: _ } = list;
        id.delete();
    }

    /// Get the number of items in the list.
    public entry fun length(list: &TodoList): u64 {
        list.items.length()
    }
}