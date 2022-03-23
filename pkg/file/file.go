package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Put(data []byte, to string) error {
	err := ioutil.WriteFile(to, data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}

	return true
}

func FileNameWithoutExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
