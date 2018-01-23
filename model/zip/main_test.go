package zip

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestDummy(t *testing.T) {
	data, err := ioutil.ReadFile(`testdata/test.zip`)
	if err != nil {
		t.Fatalf("Fail Open Dummy data %#v", err)
	}
	zv, err := Verbose(data)
	if err != nil {
		t.Fatalf("Fail Open Dummy data %#v", err)
	}

	log.Println(zv)
}
