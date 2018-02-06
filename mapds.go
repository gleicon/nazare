package main

type HLLDatastorage struct {
	Datastorage
	bytemap map[string][]byte
}

/*
NewHLLCounters is a high level HLL counter abstraction.
Wraps the hll implementation, keeps an active counter register, knows how to save it to disk.
*/
func NewHLLDatastorage() (*HLLDatastorage, error) {

	hll := HLLDatastorage{}
	hll.bytemap = make(map[string][]byte)

	return &hll, nil
}

func (hds *HLLDatastorage) Add(name string, payload []byte) error {
	hds.bytemap[name] = payload
	return nil
}

func (hds *HLLDatastorage) Get(name string) ([]byte, error) {
	return hds.bytemap[name], nil

}

func (hds *HLLDatastorage) Delete(string) {} // NOOP
func (hds *HLLDatastorage) Close()        {} // NOOP
func (hds *HLLDatastorage) Flush()        {} // NOOP
