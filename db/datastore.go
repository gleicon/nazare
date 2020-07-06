package db

/*
Datastorage describes the generic interface to store counters.
*/
type Datastorage interface {
	Add([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Delete([]byte) error
	Close()
	Flush()
}
