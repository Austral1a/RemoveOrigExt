package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const OrigExt = ".orig"

var _FILES = make(chan string, 20)

func main() {
	var rootFlag = flag.String("root", ".", "relative root directory")

	flag.Parse()

	err := filepath.Walk(*rootFlag, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			go writeFilePathWithExt(path, OrigExt)
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

L:
	for {
		select {
		case file := <-_FILES:
			err := removeFile(file)

			sayFileRemoveSucceed(file)

			if err != nil {
				sayFileRemoveFailed(file, err)
			}
		default:
			break L
		}
	}
}

func removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func writeFilePathWithExt(fileName, ext string) {
	if strings.HasSuffix(fileName, ext) {
		_FILES <- fileName
	}
}

func sayFileRemoveSucceed(filePath string) {
	log.Println(fmt.Sprintf("`%s` has been deleted successfully", filePath))
}

func sayFileRemoveFailed(filePath string, err error) {
	log.Println(fmt.Sprintf("`%s` has not been deleted\nerror: %s", filePath, err))
}
