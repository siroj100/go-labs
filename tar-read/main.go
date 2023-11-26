package main

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("please give 1 argument")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("can't open %s, %v\n", os.Args[1], err)
	}

	tarFile := tar.NewReader(file)
	for {
		hdr, err := tarFile.Next()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		fmt.Printf("header: %v, %v, %v, %v, %v, %v, %v, %v\n", hdr.Typeflag, hdr.Name, hdr.Linkname, hdr.Size, hdr.Mode, hdr.ModTime, hdr.Format, hdr.PAXRecords)
		//fmt.Printf("header: %+v\n", hdr)
	}
}
