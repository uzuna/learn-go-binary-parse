package model

import (
	"bytes"
	"errors"
	"log"

	"github.com/uzuna/learn-go-binary-parse/adapter"
)

/*
 任意のデータをUnzipする
*/
func Unzip(data []byte) error {
	// zipのCD終端レコードを取得する
	eocd, err := readEOCD(data)
	if err != nil {
		return err
	}
	// log.Printf("%+v\n", eocd)

	// CD終端レコードからCDの開始点を見つける
	curEntriNo := uint16(0)
	curOffset := eocd.Offset
	for {
		cddata := data[curOffset:]
		cdh, err := readCDH(cddata)
		if err != nil {
			panic(err)
		}
		log.Printf("%+v\n", cdh)
		curEntriNo += uint16(1)
		curOffset += cdh.HeaderLength
		if curEntriNo >= eocd.EntriesNumber {
			break
		}
	}

	return nil
}

// central directory
type EndOfCentralDirectory struct {
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
func readEOCD(data []byte) (EndOfCentralDirectory, error) {
	// 検知Pattern
	sig := []byte{0x50, 0x4B, 0x05, 0x06}

	// 入力先の型を定義
	zRecord := EndOfCentralDirectory{}
	sigIndex := bytes.Index(data, sig)
	if sigIndex < 0 {
		return zRecord, errors.New("Not Found end of Zip CDR signeture")
	}

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
	_, err := adapter.ReadBinayOffset(data[sigIndex:], Zips, &zRecord)
	if err != nil {
		return zRecord, err
	}
	return zRecord, nil
}

type CentralDirectoryHeader struct {
	Signeture                   uint32
	VersionMadeBy               []uint8
	VersionNeededExtract        uint16
	GeneralPurposeBitFlag       uint16
	CompressionMethod           uint16
	LastModFileTime             uint16
	LastModFileDate             uint16
	CRC32                       uint32
	CompressedSize              uint32
	UnCompressedSize            uint32
	FileNameLength              uint16
	ExtraFieldLength            uint16
	FileCommentLength           uint16
	DiskNumberStart             uint16
	InternalFIleAttributes      uint16
	ExternalFileAttributes      uint32
	RelativeOffsetOfLocalHeader uint32

	FileName    string
	ExtraField  []uint8
	FileComment []uint8

	// 終端位置を取得
	HeaderLength uint32
}

/*
 Signetureを検知して構造を読み出す
*/
func readCDH(data []byte) (CentralDirectoryHeader, error) {
	// 検知Pattern
	sig := []byte{0x50, 0x4B, 0x01, 0x02}

	// 入力先の型を定義
	zRecord := CentralDirectoryHeader{}
	sigIndex := bytes.Index(data, sig)
	if sigIndex < 0 {
		return zRecord, errors.New("Not Found end of CDH")
	}
	log.Println("SIGH", sigIndex)

	// パース情報
	Zips := []adapter.BynaryReadDefine{
		{
			Name:   "Signeture",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "VersionMadeBy",
			Length: 2,
			Endian: adapter.LEUint8,
		},
		{
			Name:   "VersionNeededExtract",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "GeneralPurposeBitFlag",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "CompressionMethod",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "LastModFileTime",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "LastModFileDate",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "CRC32",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "CompressedSize",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "UnCompressedSize",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "FileNameLength",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "ExtraFieldLength",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "FileCommentLength",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "DiskNumberStart",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "InternalFIleAttributes",
			Length: 2,
			Endian: adapter.LEUint16,
		},
		{
			Name:   "ExternalFileAttributes",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "RelativeOffsetOfLocalHeader",
			Length: 4,
			Endian: adapter.LEUint32,
		},
	}

	// 指定アドレスからデータを読む
	offset, err := adapter.ReadBinayOffset(data[sigIndex:], Zips, &zRecord)
	if err != nil {
		return zRecord, err
	}

	zRecord.FileName = string(data[offset : offset+uint32(zRecord.FileNameLength)])
	zRecord.HeaderLength = offset +
		uint32(zRecord.FileNameLength) +
		uint32(zRecord.ExtraFieldLength) +
		uint32(zRecord.FileCommentLength)
	return zRecord, nil
}
