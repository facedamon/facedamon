package containers

// JSONSerializer provides JSON serialization
type JSONSerializer interface {
	ToJSON() ([]byte, error)
}

type JSONDeserializer interface {
	FromJSON([]byte) error
}
