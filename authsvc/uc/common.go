package uc

// Hashable is an interface for objects that can create a hash of themselves.
type Hashable interface {
	Hash() string
}
