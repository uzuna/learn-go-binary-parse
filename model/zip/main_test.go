package zip

import (
	"io/ioutil"
	"log"
	"os"
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

func TestSeek(t *testing.T) {
	file, err := os.Open(`testdata/test.zip`)
	defer file.Close()
	if err != nil {
		t.Fatalf("Fail Open Dummy data %#v", err)
	}
	z, err := VerboseSeek(file)
	if err != nil {
		t.Fatalf("%#v", err)
	}
	log.Printf("%+v", z)
}
