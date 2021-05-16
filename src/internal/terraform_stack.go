package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type TerraformStack struct {
	Path string
}

func (tfs *TerraformStack) ShouldRunForEnv(env Environment) bool {
	files, err := ioutil.ReadDir(tfs.Path)
	if err != nil {
		return false
	}
	for _, f := range files {
		if strings.ToLower(f.Name()) == fmt.Sprintf("%s.tfvars", strings.ToLower(env.Name)) {
			return true
		}
	}
	return false
}

func (tfs *TerraformStack) GetStackPlaceholders() []Placeholder {
	AbsPath, err := filepath.Abs(tfs.Path)
	if err != nil {
		return nil
	}
	RelPath, err := filepath.Rel("./", tfs.Path)
	if err != nil {
		return nil
	}
	return []Placeholder{{
		Before: "StackAbsPath",
		After:  AbsPath,
	}, {
		Before: "StackRelPath",
		After:  RelPath,
	}}
}

func FindAllStacks(path string) ([]TerraformStack, error) {
	var stacks []TerraformStack

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			file, err := os.Open(path)
			defer file.Close()
			if err != nil {
				return err
			}
			fileInfo, err := file.Stat()
			if err != nil {
				return err
			}
			if fileInfo.IsDir() && IsTerraformStack(file.Name()) {
				stacks = append(stacks, TerraformStack{Path: file.Name()})
			}
			return nil
		})
	return stacks, err
}

func IsTerraformStack(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".tf") {
			return true
		}
	}
	return false
}
