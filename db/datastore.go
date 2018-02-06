package db

/*
Datastorage describes the generic interface to store counters.
*/
type Datastorage interface {
	Add(string, []byte) error
	Get(string) ([]byte, error)
	Delete(string)
	Close()
	Flush()
}
