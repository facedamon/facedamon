package arraylist

import (
	"facedamon/gods/utils"
	"fmt"
	"strings"
)

// List holds whe elements in a slice
type List struct {
	elements []interface{}
	size     int
}

const (
	// growth by 100%
	growthFactor = float32(2.0)
	//shrink when size is 25% of capacity
	shrinkFactor = float32(0.25)
)

func (list *List) resize(cap int) {
	newElements := make([]interface{}, cap, cap)
	copy(newElements, list.elements)
	list.elements = newElements
}

// string returns a string representation of container
func (list *List) String() string {
	str := "ArrayList\n"
	values := []string{}
	for _, value := range list.elements[:list.size] {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// check that the index is within bounds of the list
func (list *List) withRange(index int) bool {
	return index >= 0 && index < list.size
}

// shrink the array if necessary , when size shrinkfactory percent of current capacity
func (list *List) shrink() {
	if shrinkFactor == 0.0 {
		return
	}
	currentCapacity := cap(list.elements)
	if list.size <= int(float32(currentCapacity)*shrinkFactor) {
		list.resize(list.size)
	}
}

// expand the array if necessary, capacity will be reached if we add n element
func (list *List) growBy(n int) {
	currentCapacity := cap(list.elements)
	if list.size+n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		list.resize(newCapacity)
	}
}

// add appends a value at the end of the list
func (list *List) Add(values ...interface{}) {
	list.growBy(len(values))
	for _, value := range values {
		list.elements[list.size] = value
		list.size++
	}
}

// new instantiates a new list and adds the passed values, if any, to the list
func New(values ...interface{}) *List {
	list := &List{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

// get returns the element at index
// second return parameter is true if index is within bounds of the array and array is not empty
// otherwise false
func (list *List) Get(index int) (interface{}, bool) {
	if !list.withRange(index) {
		return nil, false
	}
	return list.elements[index], true
}

// remove removes the element at given index from the list
func (list *List) Remove(index int) {
	if !list.withRange(index) {
		return
	}
	list.elements[index] = nil
	copy(list.elements[index:], list.elements[index+1:list.size])
	list.size--
	list.shrink()
}

// 这个地方可以使用KMP算法优化
// contains checks if elements are present in the set
// all elements have to be present in the set for method to return true
// performance time complexity of n^2
func (list *List) Contains(values ...interface{}) bool {
	for _, searchValue := range values {
		found := false
		for _, element := range list.elements {
			if element == searchValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// values returns all elements in the list
func (list *List) Values() []interface{} {
	newElements := make([]interface{}, list.size, list.size)
	copy(newElements, list.elements[:])
	return newElements
}

//index of returns index of provided element
func (list *List) IndexOf(value interface{}) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.elements {
		if element == value {
			return index
		}
	}
	return -1
}

// empty returns true if list does not contain any elements
func (list *List) Empty() bool {
	return list.size == 0
}

// size returns number of elements within the list
func (list *List) Size() int {
	return list.size
}

//clear removes all elements from the list
func (list *List) Clear() {
	list.size = 0
	list.elements = []interface{}{}
}

//sort sorts values using
func (list *List) Sort(comparator utils.Comparator) {
	if len(list.elements) < 2 {
		return
	}
	utils.Sort(list.elements[:list.size], comparator)
}

// swap swaps the two values at the specified positions
func (list *List) Swap(i, j int) {
	if list.withRange(i) && list.withRange(j) {
		list.elements[i], list.elements[j] = list.elements[j], list.elements[i]
	}
}

// inert inserts values at specified index position shifting the value at that position
// does not do anying if position is negative or bigger than list`s size
// note : position equal to list`s size is valid append
func (list *List) Insert(index int, values ...interface{}) {
	if !list.withRange(index) {
		if index == list.size {
			list.Add(values...)
		}
		return
	}
	l := len(values)
	list.growBy(1)
	list.size++
	copy(list.elements[index+l:], list.elements[index:list.size-l])
	copy(list.elements[index:], values)
}

// set the value at specified index
// does not do anything if position is negative or bigger than list`s size
// note : position equal to list`s size size is valid,append
func (list *List) Set(index int, value interface{}) {
	if !list.withRange(index) {
		if index == list.size {
			list.Add(value)
		}
		return
	}
	list.elements[index] = value
}
