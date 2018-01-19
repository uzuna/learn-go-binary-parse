package main

import (
	"bytes"
	"log"
)

func main() {
	buf := bytes.NewBuffer([]byte{1, 2, 3})

	buf.Write([]byte{4, 5, 6})

	b := make([]byte, 3)
	buf.Read(b)

	log.Println(b, buf.Bytes())

}
