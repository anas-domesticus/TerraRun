package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GetCacheDir() string {
	dir, err := ioutil.TempDir("", "tf-cache")
	if err != nil {
		fmt.Printf("Failed to create temporary directory")
		os.Exit(1)
	}
	return dir
}
