package containers

// IteratorWithIndex is stateful iterator for ordered containers whose values can be fetched by an index
type IteratorWithIndex interface {
	// next moves the iterator to the next element and returns true if there was a next element in the container.
	// if next() returns true, then next element`s index and value can be retrieved by index and value
	// if next() was called for the first time, then it will point the iterator to the first element if it exists.
	// modifies the state of the iterator
	Next() bool

	// value returns the current element`s value
	// does not modify the state of the iterator
	Value() interface{}

	// index returns the current element`s index
	// does not modify the state of the iterator
	Index() int

	// begin resets the iterator to its initial state(one-before-first)
	// call next() to fetch the first element if any
	Begin()

	// first moves the iterator to the first element and returns true if there was a first element in the container
	// if first() returns true, then first element`s index and value be retrieved by index() and value()
	// modifies the state of the iterator
	First() bool
}

type ReverseIteratorWithIndex interface {
	Prev() bool

	End()

	Last() bool

	IteratorWithIndex
}
