package genenv

import (
	"log"
	"os"
	"path/filepath"
)

func GetLatestTestDirectory(rootPath string) string {
	configsDir := filepath.Join(rootPath, "inabox", "testdata")
	files, err := os.ReadDir(configsDir)
	if err != nil {
		panic(err)
		// return err
	}
	if len(files) == 0 {
		log.Panicf("no default experiment available")
	}
	return files[len(files)-1].Name()
}
