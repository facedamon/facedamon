package containers

import "facedamon/gods/utils"

// container is base interface that all data structures implements
type Container interface {
	Empty() bool
	Size() int
	Clear()
	Values() []interface{}
}

// GetSortedValues returns sorted container`s elements with respect to the passed comparator
// does not effect the ordering of elements within the container
func GetSortedValues(container Container, comparator utils.Comparator) []interface{} {
	values := container.Values()
	if len(values) < 2 {
		return values
	}
	utils.Sort(values, comparator)
	return values
}
