package zip

import (
	"log"
	"os"
)

/*
 任意のデータをUnzipする
*/
func Unzip(data []byte) error {
	// zipのCD終端レコードを取得する
	eocd, err := ReadEOCD(data)
	if err != nil {
		return err
	}
	// log.Printf("%+v\n", eocd)

	// CD終端レコードからCDの開始点を見つける
	curEntriNo := uint16(0)
	curOffset := eocd.Offset
	cdhList := make([]CentralDirectoryHeader, eocd.EntriesNumber)
	for {
		cddata := data[curOffset:]
		cdh, err := ReadCDH(cddata)
		if err != nil {
			panic(err)
		}
		log.Printf("%+v\n", cdh)
		cdhList[curEntriNo] = cdh
		curEntriNo += uint16(1)
		curOffset += cdh.HeaderLength
		if curEntriNo >= eocd.EntriesNumber {
			break
		}
	}

	// CDをもとにFileDataを読む
	curEntriNo = uint16(0)
	for {
		cdh := cdhList[curEntriNo]
		cdf, err := ReadLocalFile(data, cdh)
		if err != nil {
			panic(err)
		}
		log.Printf("%+v\n", cdf)
		log.Println(string(cdf.FileData))
		curEntriNo += uint16(1)
		if curEntriNo >= eocd.EntriesNumber {
			break
		}
	}

	return nil
}

/*
 Fileの概要だけを取得する
*/
func Verbose(data []byte) (ZipVerbose, error) {
	zv := ZipVerbose{}
	// zipのCD終端レコードを取得する
	eocd, err := ReadEOCD(data)
	if err != nil {
		return zv, err
	}
	// log.Printf("%+v\n", eocd)
	readCD4Verbose(data, &zv, eocd)

	return zv, nil
}

func readCD4Verbose(data []byte, zv *ZipVerbose, eocd EndOfCentralDirectory) error {
	// CD終端レコードからCDの開始点を見つける
	curEntriNo := uint16(0)
	curOffset := eocd.Offset
	cdhList := make([]CentralDirectoryHeader, eocd.EntriesNumber)
	for {
		cddata := data[curOffset:]
		cdh, err := ReadCDH(cddata)
		if err != nil {
			return err
		}
		cdhList[curEntriNo] = cdh
		curEntriNo += uint16(1)
		curOffset += cdh.HeaderLength
		if curEntriNo >= eocd.EntriesNumber {
			break
		}
	}
	zv.EndOfCentralDirectory = eocd
	zv.CentralDirectoryHeaders = cdhList
	return nil
}

/*
 File Pointerを受け取りVerboseを作成する
 巨大なファイルの場合に任意のデータだけを取り出す動作に必要
*/
func VerboseSeek(file *os.File) (ZipVerbose, error) {
	z := ZipVerbose{}

	seekSize := 200
	seekCache := make([]byte, seekSize)
	seekOffset := int64(seekSize - 4)
	seekLength := int64(seekSize)
	// Read
	endOfFileLength, err := file.Seek(-seekLength, os.SEEK_END)
	if err != nil {
		return z, err
	}
	endOfFileLength += seekLength
	for {
		file.Read(seekCache)
		if DetectEOCD(seekCache) > -1 {
			break
		}
		seekLength += seekOffset
		pos, err := file.Seek(-seekLength, os.SEEK_END)
		if err != nil {
			return z, err
		}
		// eocdPosition = pos

		if pos < int64(1) {
			log.Fatalln("Not Found EOCD Signeture")
		}
	}
	// Read後は位置がずれて居tるため位置を戻す
	file.Seek(-seekLength, os.SEEK_END)
	headerCache := make([]byte, seekLength)
	file.Read(headerCache)

	eocd, err := ReadEOCD(headerCache)
	if err != nil {
		return z, err
	}
	log.Println(eocd.Offset, endOfFileLength, seekLength)
	file.Seek(int64(eocd.Offset), os.SEEK_SET)
	cdCache := make([]byte, endOfFileLength-int64(eocd.Offset))
	file.Read(cdCache)
	log.Println(cdCache)
	readCD4VerboseSeek(cdCache, &z, eocd)

	return z, nil
}

func readCD4VerboseSeek(data []byte, zv *ZipVerbose, eocd EndOfCentralDirectory) error {
	// CD終端レコードからCDの開始点を見つける
	curEntriNo := uint16(0)
	curOffset := uint32(0)
	cdhList := make([]CentralDirectoryHeader, eocd.EntriesNumber)
	for {
		cddata := data[curOffset:]
		log.Println("cddata", cddata)
		cdh, err := ReadCDH(cddata)
		if err != nil {
			return err
		}
		cdhList[curEntriNo] = cdh
		curEntriNo += uint16(1)
		curOffset += cdh.HeaderLength
		if curEntriNo >= eocd.EntriesNumber {
			break
		}
	}
	zv.EndOfCentralDirectory = eocd
	zv.CentralDirectoryHeaders = cdhList
	return nil
}
