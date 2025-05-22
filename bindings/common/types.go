package common

// Common types used across CCIP packages
type CCIPObjectRef struct {
	Id string `move:"sui::object::UID"`
}

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}
