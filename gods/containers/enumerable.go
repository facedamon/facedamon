package containers

// enumerable provides functions for ordered containers whose values
// can be fetch by an index
type EnumerableWithIndex interface {
	//each calls the given function once for each element,
	//passing that element`s index and value
	Each(func(index int, value interface{}))

	// map invokes the given function once for each element and
	// returns a container containing the values returned by the given function
	// don`t want to type assert when chaining
	Map(func(index int, value interface{}) interface{}) Container

	// select returns a new container containing all elements for which the given
	// function returns a true value
	Select(func(index int, value interface{}) bool) Container

	// any passes each element of the container to the given function
	// and returns true if the function ever returns true for any element
	Any(func(index int, value interface{}) bool) bool

	// all passes each element of the container to the given function
	// and returns true if the function returns true for all element
	All(func(index int, value interface{}) bool) bool

	// find passes each element of the container to the given function
	// and returns the first(index,value) for which the function is true or -1,nil otherwise
	// if no element matches the criteria
	Find(func(index int, value interface{}) bool) (int, interface{})
}
