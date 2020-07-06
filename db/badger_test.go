package db

import (
	"os"
	"testing"
)

var badgerds *BadgerDatastorage

func TestBadgerDSSetandGetKey(t *testing.T) {
	var err error
	var val []byte

	err = os.RemoveAll("/tmp/test_badger.db")
	badgerds, _ = NewBadgerDatastorage("/tmp/test_badger.db")

	key := []byte("testcounter")

	if err = badgerds.Add(key, []byte("12345")); err != nil {
		t.Error(err)
		return
	}

	if val, err = badgerds.Get(key); err != nil {
		t.Error(err)
		return
	}

	if string(val) != "12345" {
		t.Error("Wrong value")
		return
	}
	err = os.RemoveAll("/tmp/test_badger.db")
	if err != nil {
		t.Fatal(err)
	}

}
