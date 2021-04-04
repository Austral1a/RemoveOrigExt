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

const NodeModulesDirName = "node_modules"
const RootFolderPath = "./"

func main() {
	var rootFlag = flag.String("root", RootFolderPath, "relative root directory")
	var omitDirsFlag = flag.String("omitDirs", NodeModulesDirName, "directories through comma, where files with `.orig` ext won't be removed\ne.g --omitDirs='node_modules,someOtherDir,AndSoOnDir'\n")

	flag.Parse()

	omitDirs := strings.ReplaceAll(*omitDirsFlag, " ", "")

	err := filepath.WalkDir(*rootFlag, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			for _, omitDirName := range strings.Split(omitDirs, ",") {
				if d.Name() == omitDirName {
					return filepath.SkipDir
				}
			}
		}

		if d.Type().IsRegular() {
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
