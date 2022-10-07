package adapter

type Adapter interface {
	Read() (interface{}, error)
	Close()
}
