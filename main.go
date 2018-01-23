package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/uzuna/learn-go-binary-parse/model/zip"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	zv, err := zip.Verbose(data)
	if err != nil {
		panic(err)
	}
	for _, v := range zv.CentralDirectoryHeaders {
		log.Println(v.FileName,
			v.RelativeOffsetOfLocalHeader,
			v.CompressedSize,
		)
	}
}
