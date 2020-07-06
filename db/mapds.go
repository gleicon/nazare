package db

/*
HLLDatastorage is a memory based data storage implementation
*/
type HLLDatastorage struct {
	Datastorage
	bytemap map[string][]byte
}

/*
NewHLLDatastorage is a high level HLL counter abstraction.
Wraps the hll implementation, keeps an active counter register, knows how to save it to disk.
*/
func NewHLLDatastorage() (*HLLDatastorage, error) {

	hll := HLLDatastorage{}
	hll.bytemap = make(map[string][]byte)

	return &hll, nil
}

/*
Add a new key
*/
func (hds *HLLDatastorage) Add(key []byte, payload []byte) error {
	hds.bytemap[string(key)] = payload
	return nil
}

/*
Get data
*/
func (hds *HLLDatastorage) Get(key []byte) ([]byte, error) {
	return hds.bytemap[string(key)], nil

}

/*
Delete a key
*/
func (hds *HLLDatastorage) Delete(key []byte) error {
	delete(hds.bytemap, string(key))
	return nil
}

/*
Close the db - NOOP
*/
func (hds *HLLDatastorage) Close() {} // NOOP

/*
Flush the db - NOOP
*/
func (hds *HLLDatastorage) Flush() {} // NOOP
