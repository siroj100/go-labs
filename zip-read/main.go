package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("please give 1 argument")
	}

	reader, err := zip.OpenReader(os.Args[1])
	if err != nil {
		log.Fatalf("can't open %s, %v\n", os.Args[1], err)
	}
	defer reader.Close()
	fmt.Println("file", os.Args[1], ", comment:", reader.Comment)

	for i, file := range reader.File {
		fmt.Printf("%d, %s, %s, %v\n", i, file.Name, file.Comment, file.Modified)
	}
}
