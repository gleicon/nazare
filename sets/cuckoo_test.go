package sets

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gleicon/nazare/db"
)

var sets *CkSet
var memds db.Datastorage

func randomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func TestMain(t *testing.T) {
	memds, _ = db.NewHLLDatastorage()
	sets, _ = NewCkSets(memds)
}

func TestSetCardinality(t *testing.T) {
	var ssBefore, ssAfter uint
	var err error
	setName := []byte("testSet")

	if ssBefore, err = sets.SCard(setName); err != nil {
		t.Error(err)
		return
	}

	if err := sets.SAdd(setName, []byte(randomString())); err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	if ssAfter, err = sets.SCard(setName); err != nil {
		t.Error(err)
		return
	}

	if ssAfter < 1 {
		t.Error("Empty set")
		return
	}

	if ssAfter-ssBefore != 1 {
		t.Error("Wrong estimate")
		return
	}

}

func TestSetIsMember(t *testing.T) {
	var err error
	var ok bool
	var ssCard uint
	ok = false

	setName := []byte("testIsMemberSet")
	memberName := []byte(randomString())
	fmt.Println("member: " + string(memberName))

	if err := sets.SAdd(setName, memberName); err != nil {
		t.Error(err)
		return
	}

	if ssCard, err = sets.SCard(setName); err != nil {
		t.Error(err)
		return
	}

	if ssCard < 1 {
		t.Error("Empty set")
		return
	}

	if ok, err = sets.SisMember(setName, memberName); err != nil {
		t.Error(err)
		return
	}

	if !ok {
		t.Error("membership failed @ set: " + string(setName) + " member: " + string(memberName))
		return
	}
}

func TestSRem(t *testing.T) {
	var err error
	var ok bool
	var ssCard uint
	ok = false

	setName := []byte("testDeleteMember")
	memberName := []byte(randomString())
	fmt.Println("member: " + string(memberName))

	if err := sets.SAdd(setName, memberName); err != nil {
		t.Error(err)
		return
	}

	if ssCard, err = sets.SCard(setName); err != nil {
		t.Error(err)
		return
	}

	if ssCard < 1 {
		t.Error("Empty set")
		return
	}

	if ok, err = sets.SRem(setName, memberName); err != nil {
		t.Error(err)
		return
	}

	if !ok {
		t.Error("membership failed @ set: " + string(setName) + " member: " + string(memberName))
		return
	}

	ssCard = 0
	if ssCard, err = sets.SCard(setName); err != nil {
		t.Error(err)
		return
	}

	if ssCard > 0 {
		t.Error("Non Empty set")
		return
	}
}
