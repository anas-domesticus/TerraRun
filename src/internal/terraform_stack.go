package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type TerraformStack struct {
	Path string
}

func (tfs *TerraformStack) GetAbsPath() string {
	absPath, err := filepath.Abs(filepath.Join("./", tfs.Path))
	if err != nil {
		return ""
	}
	return absPath
}

func (tfs *TerraformStack) ShouldRunForEnv(env Environment) bool {
	if env.Name == "" {
		return true
	}
	files, err := ioutil.ReadDir(tfs.Path)
	if err != nil {
		return false
	}
	for _, f := range files {
		if strings.ToLower(f.Name()) == fmt.Sprintf("env-%s.tfvars", strings.ToLower(env.Name)) {
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

func (tfs *TerraformStack) GetStackConfig() (StackConfig, error) {
	fullPath := filepath.Join(tfs.Path, "terrarun.yaml")
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return StackConfig{}, err
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return StackConfig{}, err
	}
	stackConfig := StackConfig{}
	err = yaml.Unmarshal(data, &stackConfig)
	if err != nil {
		return StackConfig{}, err
	}
	return stackConfig, nil
}

func FindAllStacks(path string) ([]TerraformStack, error) {
	var stacks []TerraformStack

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, iErr error) error {
			if iErr != nil {
				return iErr
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
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

func ForAllStacks(cfg Config, fn func(Config, TerraformStack) (ExecuteOutput, error)) ([]ExecuteOutput, error) {
	stacks, err := FindAllStacks(cfg.BaseDir)
	if err != nil {
		return nil, err
	}
	var outputs []ExecuteOutput
	for _, s := range stacks {
		if s.ShouldRunForEnv(cfg.Env) {
			out, err := fn(cfg, s)
			if err != nil {
				return outputs, err
			}
			outputs = append(outputs, out)
		}
	}
	return outputs, nil
}
