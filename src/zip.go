package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func main() {
	// 引数が不足する場合は戻る
	if len(os.Args) < 3 {
		fmt.Println("Usage: %s [zipname] [src] [src2] ...")
		return
	}

	// 出力Fileを作成
	dest, err := os.Create(os.Args[1])
	if err != nil {
		panic(err)
	}

	zipWriter := zip.NewWriter(dest)
	defer zipWriter.Close()

	for _, s := range os.Args[2:] {
		if err := addToZip(s, zipWriter); err != nil {
			panic(err)
		}
	}

}

func addToZip(filename string, zipWriter *zip.Writer) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	hdr, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	f, err := zipWriter.CreateHeader(hdr)

	// writer, err := zipWriter.Create(filename)
	// body, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	return err
	// }
	// f.Write(body)

	_, err = io.Copy(f, src)
	if err != nil {
		return err
	}

	return nil
}
