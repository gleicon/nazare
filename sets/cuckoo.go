package sets

import (
	"errors"
	"sync"

	"github.com/gleicon/nazare/db"
	cuckoo "github.com/seiflotfy/cuckoofilter"
)

/*
CkSetStats general stats about a set of counters and its db
*/
type CkSetStats struct {
	ActiveItems uint64
}

/*
CkSet wraps a cuckoo filter impl and provide locks, stats
*/
type CkSet struct {
	ckrwlocks   map[string]*sync.RWMutex
	datastorage db.Datastorage
	stats       CkSetStats
}

/*
NewCkSets is a high level wrapper and set impl that uses cuckoo filter. It receives a datastorage impl
*/
func NewCkSets(ds db.Datastorage) (*CkSet, error) {

	ckh := CkSet{}
	ckh.datastorage = ds
	ckh.ckrwlocks = make(map[string]*sync.RWMutex)

	return &ckh, nil
}

func (ckset *CkSet) lockKey(key []byte) *sync.RWMutex {
	if ckset.ckrwlocks[string(key)] == nil {
		ckset.ckrwlocks[string(key)] = new(sync.RWMutex)
	}

	localMutex := ckset.ckrwlocks[string(key)]
	return localMutex
}

/*
SAdd a member to a set
*/
func (ckset *CkSet) SAdd(key, member []byte) error {
	var sts *cuckoo.Filter

	localMutex := ckset.lockKey(key)
	localMutex.Lock()
	defer localMutex.Unlock()

	value, err := ckset.datastorage.Get(key)
	if err != nil {
		return errors.New("Error fetching set: " + string(key))
	}
	if value == nil {
		// tunable cuckoo size
		sts = cuckoo.NewFilter(1024 * 1024)
	} else {
		sts, err = cuckoo.Decode(value)
		if err != nil {
			return errors.New("Error decoding filter set: " + string(key))
		}
	}
	sts.InsertUnique(key)
	ckset.datastorage.Add(key, sts.Encode())

	return nil
}

/*
SisMember tells if an item belongs to a set
*/
func (ckset *CkSet) SisMember(key, member []byte) (bool, error) {
	value, err := ckset.datastorage.Get(key)
	if err != nil {
		return false, errors.New("Error fetching set: " + string(key))
	}
	if key != nil {
		sts, err := cuckoo.Decode(value)
		if err != nil {
			return false, errors.New("Error decoding filter set: " + string(key))
		}
		if sts.Lookup(member) {
			return true, nil
		}
	}
	return false, nil
}

/*
SRem removes an item from the set
*/
func (ckset *CkSet) SRem(key, member []byte) (bool, error) {
	localMutex := ckset.lockKey(key)
	localMutex.Lock()
	defer localMutex.Unlock()

	value, err := ckset.datastorage.Get(key)
	if err != nil {
		return false, errors.New("Error fetching set: " + string(key))
	}
	if key != nil {
		sts, err := cuckoo.Decode(value)
		if err != nil {
			return false, errors.New("Error decoding filter set: " + string(key))
		}
		if sts.Delete(member) {
			return true, nil
		}
		if err = ckset.datastorage.Add(key, sts.Encode()); err != nil {
			return false, errors.New("Error encoding filter set: " + string(key))
		}
	}
	return false, nil
}

/*
SCard is cardinality of a set (SCount items on a set)
*/
func (ckset *CkSet) SCard(key []byte) (uint, error) {
	value, err := ckset.datastorage.Get(key)
	if err != nil {
		return 0, errors.New("Error fetching set: " + string(key))
	}
	if key != nil {
		sts, err := cuckoo.Decode(value)
		if err != nil {
			return 0, errors.New("Error decoding filter set for counting: " + string(key))
		}
		return sts.Count(), nil
	}
	return 0, nil
}
