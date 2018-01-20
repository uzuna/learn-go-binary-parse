package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

func main() {
	// file read
	data, err := ioutil.ReadFile(`./dummy/test.zip`)
	if err != nil {
		panic(err)
	}
	log.Println(data)

	// sig := []byte{0x50, 0x4B, 0x05, 0x06}
	// log.Println(bytes.Index(data, sig))

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

// 読み取りを行う時の定義
type BynaryReadDefine struct {
	name   string
	length int
	endian string
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
	Zips := []BynaryReadDefine{
		{
			name:   "Signeture",
			length: 4,
			endian: "LittleEndian.Uint32",
		},
		{
			name:   "DiskNumber",
			length: 2,
			endian: "LittleEndian.Uint16",
		},
		{
			name:   "DiskStartNumber",
			length: 2,
			endian: "LittleEndian.Uint16",
		},
		{
			name:   "EntriesNumberOfThisDisk",
			length: 2,
			endian: "LittleEndian.Uint16",
		},
		{
			name:   "EntriesNumber",
			length: 2,
			endian: "LittleEndian.Uint16",
		},
		{
			name:   "Size",
			length: 4,
			endian: "LittleEndian.Uint32",
		},
		{
			name:   "Offset",
			length: 4,
			endian: "LittleEndian.Uint32",
		},
		{
			name:   "CommentLength",
			length: 2,
			endian: "LittleEndian.Uint16",
		},
	}

	// 指定アドレスからデータを読む
	readBinayOffset(data[sigIndex:], Zips, &zRecord)
	log.Printf("%+v\n", zRecord)
	return "", nil
}

/*
 汎用binary read
 定義に従いデータを読んで渡された構造体に値を入れる
*/
func readBinayOffset(data []byte,
	defines []BynaryReadDefine,
	record interface{},
) error {

	var offset = 0

	// 構造体のPointerから値の書き込み先要素を得る
	structValue := reflect.ValueOf(record).Elem()
	for _, value := range defines {
		// Read
		name := value.name
		vd := data[offset : offset+value.length]
		// log.Println(key, name, vd)
		offset += value.length

		// byteデータを任意の型にCast
		dataEncoded := readBytes(vd, value.endian)
		log.Println(value.name, vd, dataEncoded)

		// Reflectを使ってSet https://github.com/oleiade/reflections より。
		structFieldValue := structValue.FieldByName(value.name)
		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in obj", name)
		}

		// If obj field value is not settable an error is thrown
		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", name)
		}

		val := reflect.ValueOf(dataEncoded)

		// Cast済みのためそのままSet
		structFieldValue.Set(val)
	}
	return nil
}

/*
 endianの指定を読んで変換
*/
func readBytes(data []byte, endian string) interface{} {
	switch endian {
	case "LittleEndian.Uint32":
		return binary.LittleEndian.Uint32(data)
	case "LittleEndian.Uint16":
		return binary.LittleEndian.Uint16(data)
	default:
		panic(fmt.Errorf("Undefined endian string: [%v]", endian))
	}
}
