package main

import (
	"io/ioutil"

	"github.com/uzuna/learn-go-binary-parse/model/zip"
)

func main() {

	data, err := ioutil.ReadFile(`./dummy/test.zip`)
	if err != nil {
		panic(err)
	}
	// log.Println(data)
	errUnzip := zip.Unzip(data)
	if errUnzip != nil {
		panic(errUnzip)
	}
	// adapter.binaty_parse.ConvertType.LEUint16
	// log.Printf("%v", adapter.LEUint16)

}
