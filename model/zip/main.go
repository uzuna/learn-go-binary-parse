package zip

import "log"

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
		cdhList[curEntriNo] = cdh
		curEntriNo += uint16(1)
		curOffset += cdh.HeaderLength
		if curEntriNo >= eocd.EntriesNumber {
			break
		}
	}
	zv.EndOfCentralDirectory = eocd
	zv.CentralDirectoryHeaders = cdhList
	return zv, nil
}
