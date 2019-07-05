package singlylinkedlist

// Iterator golding the iterator`s state
type Iterator struct {
	list    *List
	index   int
	element *element
}

// Iterator returns a stateful iterator whose values can be fetched by an index
func (list *List) Iterator() Iterator {
	return Iterator{list: list, index: -1, element: nil}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container
// if Next() returns true, then next element`s index and value can be retrieved by index() and Value().
// if Next()was called for the first time, then it will point the iterator to the first element if it exists.
func (iterator *Iterator) Next() bool {
	if iterator.index < iterator.list.size {
		iterator.index++
	}
	// out of the position
	if !iterator.list.withinRange(iterator.index) {
		iterator.element = nil
		return false
	}
	if iterator.index == 0 {
		iterator.element = iterator.list.first
	} else {
		iterator.element = iterator.element.next
	}
	return true
}
