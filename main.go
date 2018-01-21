package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"

	"github.com/uzuna/learn-go-binary-parse/adapter"
)

func main() {

	data, err := ioutil.ReadFile(`./dummy/test.zip`)
	if err != nil {
		panic(err)
	}
	log.Println(data)
	// adapter.binaty_parse.ConvertType.LEUint16
	// log.Printf("%v", adapter.LEUint16)
	readSigneture(data)
}

type ZipCentralDirectoryEndRecord struct {
	Signeture               uint32
	DiskNumber              uint16
	DiskStartNumber         uint16
	EntriesNumberOfThisDisk uint16
	EntriesNumber           uint16
	Size                    uint32
	Offset                  uint32
	CommentLength           uint16
	// Comment         uint16
}

/*
 Signetureを検知して構造を読み出す
*/
func readSigneture(data []byte) (string, error) {
	// 検知Pattern
	sig := []byte{0x50, 0x4B, 0x05, 0x06}

	sigIndex := bytes.Index(data, sig)
	if sigIndex < 1 {
		return "", errors.New("Not Found end of Zip CDR signeture")
	}

	// 入力先の型を定義
	zRecord := ZipCentralDirectoryEndRecord{}

	// パース情報
	Zips := []adapter.BynaryReadDefine{
		{
			Name:   "Signeture",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "DiskNumber",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "DiskStartNumber",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "EntriesNumberOfThisDisk",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "EntriesNumber",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "Size",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "Offset",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "CommentLength",
			Length: 2,
			Endian: adapter.LEUint16,
		},
	}

	// 指定アドレスからデータを読む
	adapter.ReadBinayOffset(data[sigIndex:], Zips, &zRecord)
	log.Printf("%+v\n", zRecord)
	return "", nil
}
