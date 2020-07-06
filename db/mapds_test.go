package db

import (
	"testing"
)

var memds *HLLDatastorage

func TestMapDSSetandGetKey(t *testing.T) {
	var err error
	var val []byte
	memds, _ = NewHLLDatastorage()

	key := []byte("testcounter")

	if err = memds.Add(key, []byte("12345")); err != nil {
		t.Error(err)
		return
	}

	if val, err = memds.Get(key); err != nil {
		t.Error(err)
		return
	}

	if string(val) != "12345" {
		t.Error("Wrong value")
		return
	}

}
