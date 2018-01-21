package adapter

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// Enable Convert Type
type ConvertType int

// Enable Convert TypeList
const (
	LEUint32 ConvertType = iota
	LEUint16
	LEUint8
	// ZipFileDateTime	// 専用のパーサを追加する
)

// 読み取2りを行う時の定義
type BynaryReadDefine struct {
	Name   string
	Length uint
	Endian ConvertType
}

/*
 汎用binary read
 定義に従いデータを読んで渡された構造体に値を入れる
*/
func ReadBinayOffset(data []byte,
	defines []BynaryReadDefine,
	record interface{},
) (uint32, error) {

	offset := uint32(0)

	// 構造体のPointerから値の書き込み先要素を得る
	structValue := reflect.ValueOf(record).Elem()
	for _, value := range defines {
		// Read
		name := value.Name
		nextOffset := offset + uint32(value.Length)
		vd := data[offset:nextOffset]
		// log.Println(key, name, vd)
		offset = nextOffset

		// byteデータを任意の型にCast
		dataEncoded := readBytes(vd, value.Endian)
		// log.Println(name, vd, dataEncoded)

		// Reflectを使ってSet https://github.com/oleiade/reflections より。
		structFieldValue := structValue.FieldByName(name)
		if !structFieldValue.IsValid() {
			return offset, fmt.Errorf("No such field: %s in obj", name)
		}

		// If obj field value is not settable an error is thrown
		if !structFieldValue.CanSet() {
			return offset, fmt.Errorf("Cannot set %s field value", name)
		}

		val := reflect.ValueOf(dataEncoded)

		// Cast済みのためそのままSet
		structFieldValue.Set(val)
	}
	return offset, nil
}

/*
 endianの指定を読んで変換
*/
func readBytes(data []byte, endian ConvertType) interface{} {
	switch endian {
	case LEUint32:
		return binary.LittleEndian.Uint32(data)
	case LEUint16:
		return binary.LittleEndian.Uint16(data)
	case LEUint8:
		return data
	default:
		panic(fmt.Errorf("Undefined endian string: [%v]", endian))
	}
}
