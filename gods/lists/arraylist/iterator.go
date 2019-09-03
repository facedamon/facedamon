package arraylist

// Iterator holding the iterator`s state
type Iterator struct {
	list  *List
	index int
}

// iterator return a stateful iterator where values can be fetched by an index
func (list *List) Iterator() Iterator {
	return Iterator{list: list, index: -1}
}

// next moves the iterator to the next element and returns true
// if there was a next element in the container.
// if next() returns true,then next element`s index and value can be retrieved
// if next() was called for the first time,then it will point the iterator to the first element if it exists
// Modifies the state of the iterator
func (iterator *Iterator) Next() bool {
	if iterator.index < iterator.list.size {
		iterator.index++
	}
	return iterator.list.withRange(iterator.index)
}

// Prev moves the iterator to previous element and returns true if there was a previous element in the container
// if prev() returns true,then previous element`s index and value can be retrieved
// Modifies the state of the iterator
func (iterator *Iterator) Prev() bool {
	if iterator.index >= 0 {
		iterator.index--
	}
	return iterator.list.withRange(iterator.index)
}

// value returns the current elements`s value
// Does not modify the state of the iterator
func (iterator *Iterator) Value() interface{} {
	return iterator.list.elements[iterator.index]
}

// Index returns the current elements`s index
// Does not modify the state of the iterator
func (iterator *Iterator) Index() int {
	return iterator.index
}

// Begin resets the iterator to its initial state
// call next() to fetch the first element if any
func (iterator *Iterator) Begin() {
	iterator.index = -1
}

// End moves the iterator past the last element
// call prev() to fetch the last element if any
func (iterator *Iterator) End() {
	iterator.index = iterator.list.size
}

// First moves the iterator to the first element and returns true if there was a first element in the container
// if first() returns true, then first element`s index and value can be retrieved by index() and value)(
// modifies the state of the iterator
func (iterator *Iterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}

//Last moves the iterator to the last element and returns true if there was a last element in the container
// if last() returns true. then first element`s index and value can be retrieved by index() and value()
// modifies the state of the iterator
func (iterator *Iterator) Last() bool {
	iterator.End()
	return iterator.Prev()
}
