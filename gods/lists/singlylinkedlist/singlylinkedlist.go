// Package singlylinkedlist
// 单链表
// 查找比array效率低
// 插入比array效率高e
package singlylinkedlist

import (
	"fmt"
	"github.com/facedamon/gods/utils"
	"strings"
)

type element struct {
	value interface{}
	next  *element
}

// List holds the elements, where each element points to the next element
type List struct {
	first *element
	last  *element
	size  int
}

// Add appends a value at the end of the list
func (list *List) Add(values ...interface{}) {
	for _, value := range values {
		newElement := &element{value: value}
		if list.size == 0 {
			list.first = newElement
			list.last = newElement
		} else {
			list.last.next = newElement
			list.last = newElement
		}
		list.size++
	}
}

// New instantiated a new list and adds the passed values, if any
// to the list
func New(values ...interface{}) *List {
	list := &List{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

// Prepend prepends a values
func (list *List) Prepend(values ...interface{}) {
	for v := len(values) - 1; v >= 0; v-- {
		newElement := &element{value: values[v], next: list.first}
		list.first = newElement
		if list.size == 0 {
			list.last = newElement
		}
		list.size++
	}
}

func (list *List) withinRange(index int) bool {
	return index >= 0 && index < list.size
}

// Get returns the element at index
// second return parameter is true if index within bounds of the array
// and array is not empty
func (list *List) Get(index int) (interface{}, bool) {
	if !list.withinRange(index) {
		return nil, false
	}
	element := list.first
	for e := 0; e != index; e, element = e+1, element.next {
	}
	return element.value, true
}

// Clear removes all elements from the list
func (list *List) Clear() {
	list.size = 0
	list.first = nil
	list.last = nil
}

// Remove removes the element at the given index from the list
func (list *List) Remove(index int) {
	if !list.withinRange(index) {
		return
	}
	if list.size == 1 {
		list.Clear()
		return
	}
	// 将要删除的节点的前一位指针
	var beforeElement *element
	element := list.first
	for e := 0; e != index; e, element = e+1, element.next {
		beforeElement = element
	}
	if element == list.first {
		list.first = element.next
	}
	if element == list.last {
		list.last = beforeElement
	}
	// 指针绕过要删除的节点
	if beforeElement != nil {
		beforeElement.next = element.next
	}
	element = nil
	list.size--
}

// Contains check if values are present int the set
// all values have to be present in the set for method to return true
// performance time complexity of n^2
// returns true if no arguments are passed at all
func (list *List) Contains(values ...interface{}) bool {
	if len(values) == 0 {
		return true
	}
	if list.size == 0 {
		return false
	}
	for _, value := range values {
		found := false
		for element := list.first; element != nil; element = element.next {
			if element.value == value {
				found = true
				break
			}
		}
		//只要有一个不存在就返回false
		if !found {
			return false
		}
	}
	return true
}

// Values returns all elements int the list
func (list *List) Values() []interface{} {
	values := make([]interface{}, list.size, list.size)
	for e, element := 0, list.first; element != nil; e, element = e+1, element.next {
		values[e] = element.value
	}
	return values
}

// IndexOf returns index of provided element
func (list *List) IndexOf(value interface{}) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.Values() {
		if element == value {
			return index
		}
	}
	return -1
}

// Empty returns true is list does not contain any elements
func (list *List) Empty() bool {
	return list.size == 0
}

// Size returns number of elements within the list
func (list *List) Size() int {
	return list.size
}

// Sort values using
func (list *List) Sort(comparator utils.Comparator) {
	if list.size < 2 {
		return
	}
	values := list.Values()
	utils.Sort(values, comparator)
	list.Clear()
	list.Add(values...)
}

// Swap values of two elements at the given indices
func (list *List) Swap(i, j int) {
	if list.withinRange(i) && list.withinRange(j) && i != j {
		var element1, element2 *element
		// 同时检索两个节点
		for e, currentElement := 0, list.first; element1 == nil || element2 == nil; e, currentElement = e+1, currentElement.next {
			switch e {
			case i:
				element1 = currentElement
			case j:
				element2 = currentElement
			}
		}
		element1.value, element2.value = element2.value, element1.value
	}
}

// Insert inserts values at specified index position
func (list *List) Insert(index int, values ...interface{}) {
	if !list.withinRange(index) {
		if index == list.size {
			list.Add(values...)
		}
		return
	}
	list.size += len(values)
	var beforeElement *element
	foundElement := list.first
	for e := 0; e != index; e, foundElement = e+1, foundElement.next {
		beforeElement = foundElement
	}
	// 如果是第一个节点，那么首次beforeElement = nil
	if foundElement == list.first {
		oldNextElement := list.first
		for i, value := range values {
			newElement := &element{value: value}
			if i == 0 {
				list.first = newElement
			} else {
				//从第二次开始链接后面的节点
				beforeElement.next = newElement
			}
			beforeElement = newElement
		}
		beforeElement.next = oldNextElement
	} else {
		oldNextElement := beforeElement.next
		for _, value := range values {
			newElement := &element{value: value}
			beforeElement.next = newElement
			beforeElement = newElement
		}
		beforeElement.next = oldNextElement
	}
}

// Set value at specified index
// Does not do anything if position is negativr or bigger than list`s size
func (list *List) Set(index int, value interface{}) {
	if !list.withinRange(index) {
		if index == list.size {
			list.Add(value)
		}
		return
	}
	foundElement := list.first
	for e := 0; e != index; {
		e, foundElement = e+1, foundElement.next
	}
	foundElement.value = value
}

// String returns a string representation of container
func (list *List) String() string {
	str := "SinglyLinkedList\n"
	values := []string{}
	for element := list.first; element != nil; element = element.next {
		values = append(values, fmt.Sprintf("%v", element.value))
	}
	str += strings.Join(values, ", ")
	return str
}
