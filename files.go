package main

import (
	"os"
	"path/filepath"
)



func getSongList(input string) ([]string, error) {

	result := make([]string, 0)
	addPath := func(path string, info os.FileInfo, err error) error {
		if (err != nil) {
			return err
		}

		if !info.IsDir() && Contains(supportedFormats, filepath.Ext(path)) {
			result = append(result, path)
		}
		return nil
	}
	err := filepath.Walk(input, addPath)

	return result, err

}

func Contains(arr []string ,input string) bool {
	for _, v := range arr {
		if v == input {
			return true
		}
	}
	return false
}
