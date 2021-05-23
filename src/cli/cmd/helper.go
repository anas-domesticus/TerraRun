package cmd

import (
	"fmt"
	"github.com/anas-domesticus/TerraRun/src/internal"
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

func BuildConfig() internal.Config {
	return internal.Config{
		BaseDir: directory,
		Env:     internal.Environment{Name: environment},
		Debug:   debugLogging,
	}
}
