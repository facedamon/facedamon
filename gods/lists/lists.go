package lists

import (
	"github.com/facedamon/gods/containers"
	"github.com/facedamon/gods/utils"
)

// List interface that all lists implement
type List interface {
	Get(index int) (interface{}, bool)
	Remove(index int)
	Add(values ...interface{})
	Contains(values ...interface{}) bool
	Sort(comparator utils.Comparator)
	Swap(index1, index2 int)
	Insert(index int, values ...interface{})
	Set(index int, value interface{})

	containers.Container
}
