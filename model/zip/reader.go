package zip

import (
	"os"
)

/*
 File Pointerを使いDiskAccessを減らす
 ここでは
*/
type SeekReader interface {
	Info() Summary
	Cur() Page
	Prev() Page
	Next() Page
}

/*
 Recordの有無と対照コンテンツ
*/
type Page struct {
	Content CentralDirectoryHeader
	Next    int
	Prev    int
}

/*
 File情報
*/
type Summary struct {
	Content  EndOfCentralDirectory
	FileInfo os.FileInfo
}

func readEOCD()
