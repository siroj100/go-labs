package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("please give 1 argument")
	}

	dirName := os.Args[1]
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		log.Fatalf("%s is not exist\n", dirName)
	}
	if !info.IsDir() {
		log.Fatalf("%s is not a directory\n", dirName)
	}

	fname := fmt.Sprintf("%s-%d.zip", dirName, time.Now().Unix())
	file, err := os.Create(fname)
	if err != nil {
		log.Fatalf("failed to open %s\n", fname)
	}
	defer file.Close()
	zipFile := zip.NewWriter(file)

	entries, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatalf("error reading %s, err %v\n", dirName, err)
	}
	dirHash := sha256.New()
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("skipping %s because it's a dir\n", entry.Name())
			continue
		}

		path := strings.Join([]string{dirName, entry.Name()}, string(os.PathSeparator))
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("error opening %s, err %v\n", path, err)
			continue
		}
		reader := bytes.NewReader(data)

		fhash := sha256.New()
		fileWritten, err := io.Copy(fhash, reader)
		if err != nil {
			log.Printf("error calculate fhash for %s, err %v\n", path, err)
			continue
		}
		reader.Seek(0, io.SeekStart)

		dirWritten, err := io.Copy(dirHash, reader)
		if err != nil {
			log.Printf("error calculate dir fhash for %s, err %v\n", path, err)
			continue
		}
		if fileWritten != dirWritten {
			log.Fatalf("%s fileWritten %d != dirWritten %d", path, fileWritten, dirWritten)
		}

		hash := fhash.Sum(nil)
		fmt.Printf("file %s, %d, %x\n", path, fileWritten, hash)
		finfo, _ := entry.Info()
		fh := zip.FileHeader{
			Name:     path,
			Comment:  fmt.Sprintf("sha256=%x", hash),
			Modified: finfo.ModTime(),
		}
		writer, err := zipFile.CreateHeader(&fh)
		if err != nil {
			log.Fatalf("error create header to zip: %+v\n", fh)
		}

		written, err := writer.Write(data)
		if written != int(fileWritten) || err != nil {
			log.Fatalf("error writing to zip file: %s, %d written, err %v\n", path, written, err)
		}
	}

	hash := dirHash.Sum(nil)
	fh := zip.FileHeader{
		Name:     dirName + "/",
		Comment:  fmt.Sprintf("sha256=%x", hash),
		Modified: info.ModTime(),
	}
	fmt.Printf("dir %s, %x\n", dirName, hash)
	_, err = zipFile.CreateHeader(&fh)
	if err != nil {
		log.Fatalf("error create header to zip: %+v\n", fh)
	}

	err = zipFile.Close()
	if err != nil {
		log.Fatalf("failed to close %s\n", fname)
	}
}
