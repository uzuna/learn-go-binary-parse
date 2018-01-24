package zip

import (
	"bytes"
	"errors"
	"log"

	"github.com/uzuna/learn-go-binary-parse/adapter"
)

/*
 EOCDを探す
*/
func DetectEOCD(data []byte) int {
	sig := []byte{0x50, 0x4B, 0x05, 0x06}
	return bytes.Index(data, sig)
}

/*
 Signetureを検知して構造を読み出す
*/
func ReadEOCD(data []byte) (EndOfCentralDirectory, error) {
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

/*
 Signetureを検知して構造を読み出す
*/
func ReadCDH(data []byte) (CentralDirectoryHeader, error) {
	// 検知Pattern
	sig := []byte{0x50, 0x4B, 0x01, 0x02}

	// 入力先の型を定義
	zRecord := CentralDirectoryHeader{}
	sigIndex := bytes.Index(data, sig)
	if sigIndex < 0 {
		return zRecord, errors.New("Not Found end of CDH")
	}
	// log.Println("SIGH", sigIndex)

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
			Name:   "VersionNeededToExtract",
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
			Name:   "InternalFileAttributes",
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

/*
 read Local File data
*/
func ReadLocalFile(binarydata []byte, cdh CentralDirectoryHeader) (LocalFile, error) {
	// 検知Pattern
	sig := []byte{0x50, 0x4B, 0x03, 0x04}

	data := binarydata[cdh.RelativeOffsetOfLocalHeader:]

	// 入力先の型を定義
	zRecord := LocalFile{}
	sigIndex := bytes.Index(data, sig)
	if sigIndex < 0 {
		return zRecord, errors.New("Not Found end of LFD")
	}
	log.Println("SIGH", sigIndex)

	// parse data
	Zips := []adapter.BynaryReadDefine{
		{
			Name:   "Signeture",
			Length: 4,
			Endian: adapter.LEUint32,
		},
		{
			Name:   "VersionNeededToExtract",
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
	}

	// 指定アドレスからデータを読む
	offset, err := adapter.ReadBinayOffset(data[sigIndex:], Zips, &zRecord)
	if err != nil {
		return zRecord, err
	}

	zRecord.CompressedSize = cdh.CompressedSize
	zRecord.UnCompressedSize = cdh.UnCompressedSize

	zRecord.FileName = string(data[offset : offset+uint32(zRecord.FileNameLength)])
	offset = offset + uint32(zRecord.FileNameLength)
	zRecord.ExtraField = data[offset : offset+uint32(zRecord.ExtraFieldLength)]
	offset = offset + uint32(zRecord.ExtraFieldLength)
	zRecord.FileData = data[offset : offset+uint32(zRecord.CompressedSize)]
	log.Println("ReadData", offset, zRecord.CompressedSize)
	return zRecord, nil
}
