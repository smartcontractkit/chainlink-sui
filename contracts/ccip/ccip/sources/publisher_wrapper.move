module ccip::publisher_wrapper;

use sui::address;
use sui::package::Publisher;

const EProofNotAtPublisherAddressAndModule: u64 = 1;

public struct PublisherWrapper<phantom T: drop> {
    package_address: address,
}

/// Having reference to Publisher means you have access to `Publisher` object.
/// This is only sent to the package deployer, therefore we know only the owner can call this.
public fun create<T: drop>(publisher: &Publisher, _proof: T): PublisherWrapper<T> {
    assert!(publisher.from_module<T>(), EProofNotAtPublisherAddressAndModule);
    let package_bytes = publisher.package().as_bytes();
    PublisherWrapper<T> { package_address: address::from_ascii_bytes(package_bytes) }
}

public(package) fun get_package_address<T: drop>(publisher_wrapper: PublisherWrapper<T>): address {
    let PublisherWrapper { package_address } = publisher_wrapper;
    package_address
}

#[test_only]
public fun destroy<T: drop>(publisher_wrapper: PublisherWrapper<T>): address {
    let PublisherWrapper { package_address } = publisher_wrapper;
    package_address
}
