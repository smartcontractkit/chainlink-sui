module managed_token::mint_allowance;

const EOverflow: u64 = 0;
const EInsufficientAllowance: u64 = 1;

/// A MintAllowance for a coin of type T.
/// Used for minting and burning.
public struct MintAllowance<phantom T> has copy, drop, store {
    value: u64,
    is_unlimited: bool,
}

/// [Package private] Gets the current allowance of the MintAllowance object.
public(package) fun value<T>(self: &MintAllowance<T>): u64 {
    self.value
}

public(package) fun is_unlimited<T>(self: &MintAllowance<T>): bool {
    self.is_unlimited
}

public(package) fun allowance_info<T>(self: &MintAllowance<T>): (u64, bool) {
    (self.value, self.is_unlimited)
}

/// [Package private] Create a new MintAllowance for type T.
public(package) fun new<T>(): MintAllowance<T> {
    MintAllowance { value: 0, is_unlimited: false }
}

/// [Package private] Set allowance to `value`
public(package) fun set<T>(self: &mut MintAllowance<T>, value: u64, is_unlimited: bool) {
    self.value = value;
    self.is_unlimited = is_unlimited;
}

/// [Package private] Increase the allowance by `value`
public(package) fun increase<T>(self: &mut MintAllowance<T>, value: u64) {
    assert!(value < (18446744073709551615u64 - self.value), EOverflow);
    self.value = self.value + value;
}

/// [Package private] Decrease the allowance by `value`
public(package) fun decrease<T>(self: &mut MintAllowance<T>, value: u64) {
    assert!(self.value >= value, EInsufficientAllowance);
    self.value = self.value - value;
}

/// [Package private] Destroy object
public(package) fun destroy<T>(self: MintAllowance<T>) {
    let MintAllowance { value: _, is_unlimited: _ } = self;
}
