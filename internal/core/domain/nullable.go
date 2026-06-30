package domain

/*
Nullable needs to specify:
- field not provided
- field provided: value
- field provided: nil
*/
type Nullable[T any] struct {
	Value *T
	Set   bool
}
