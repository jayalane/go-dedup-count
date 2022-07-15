// -*- tab-width: 2 -*-

package dedupcount

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var nameFrags = []string{"hello",
	"db",
	"chris",
	"whatever/is/it/called",
	"dlv",
	"s3-bucket-control",
}

func makeName() string {
	values := []string{}
	for i := 0; i < 10; i++ {
		values = append(values, nameFrags[rand.Intn(len(nameFrags))])
	}
	return strings.Join(values, "-")
}

func TestSet(t *testing.T) {
	d := New("test1")
	defer d.Close()
	saveName := makeName()
	fmt.Println(saveName)
	d.Set(saveName, "hi there")
	d.Set(saveName, "hi there 2")
	for i := 0; i < 100; i++ {
		d.Set(makeName(), makeName())
	}
	time.Sleep(1100 * time.Millisecond)
	for i := 0; i < 100; i++ {
		fatal := ""
		if d.KeySetP(makeName()) {
			fatal = "Random new string was in set, unlikely"
		}
		if !d.KeySetP(saveName) {
			fatal = "saveName was not still in set"
		}
		if fatal != "" {
			t.Log(fatal)
			t.Fail()
		}
	}
	dups := d.GetDups()
	// check that hi there is present with 2 values
	v, ok := dups[saveName]
	if !ok {
		t.Log("savedName not in dups list")
		t.Fail()
	}
	if len(v) != 2 {
		t.Log("saveName not in dups list with 2 entries")
		t.Fail()
	}
}
