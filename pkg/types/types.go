// Package types abstracts the table and row details making
// it easy to use dynamodb as an atomic counter
package types

// Schema stores all the necessary information to create or update
// a table
type Schema struct {
	HashKey   *string
	RangeKey  *string
	TableName *string
	AtomicKey *string
}

// Row stores all the necessary information to atomically increment a row
type Row struct {
	Schema     *Schema
	HashValue  *string
	RangeValue *string
	Incr       *string
}
