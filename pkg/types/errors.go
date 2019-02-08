package types

import "errors"

var (
	// ErrAtomicAttribute means that the AtomicKey was not among the attributes
	// returned by DynamoDB
	ErrAtomicAttribute = errors.New("could not get attribute value for the atomic field")
)
