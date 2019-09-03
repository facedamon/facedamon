package handler

// Handler interface is a dosomething handler
type Handler interface {
	Do(k, v interface{})
}

// HandlerFunc func is a decorator for Do interface
type HandlerFunc func(k, v interface{})

// Do implment the handler interface, decorator the method of do
func (f HandlerFunc) Do(k, v interface{}) {
	f(k, v)
}

// EachFunc func decorator the Each func
func EachFunc(m map[interface{}]interface{}, f func(k, v interface{})) {
	Each(m, HandlerFunc(f))
}

// Each func each the map with handler
func Each(m map[interface{}]interface{}, h Handler) {
	if m != nil && len(m) > 0 {
		for k, v := range m {
			h.Do(k, v)
		}
	}
}
