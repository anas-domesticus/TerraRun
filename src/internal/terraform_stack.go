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
	Path   string
	config StackConfig
}

func NewTerraformStack(path string) (TerraformStack, error) {
	config, err := getStackConfig(path)
	if err != nil {
		return TerraformStack{}, err
	}
	return TerraformStack{
		Path:   path,
		config: config,
	}, nil
}

func getStackConfig(path string) (StackConfig, error) {
	fullPath := filepath.Join(path, "terrarun.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return StackConfig{}, err // Directory doesn't exist, return error
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return StackConfig{}, nil // File doesn't exist, so provide empty config & throw away error
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

func FindAllStacks(path string) (map[int]TerraformStack, error) {
	stacks := make(map[int]TerraformStack)

	count := 0
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
				stack, err := NewTerraformStack(file.Name())
				if err != nil {
					return err
				}
				stacks[count] = stack
				count = count + 1
			}
			return nil
		})

	// Here we should check dependencies

	return stacks, err
}

func dependencyInSlice(stacks map[int]TerraformStack, path Dependency) bool {
	for _, stack := range stacks {
		if stack.Path == string(path) {
			return true
		}
	}
	return false
}

func pathToKey(stacks map[int]TerraformStack, path Dependency) int {
	for i, stack := range stacks {
		if stack.Path == string(path) {
			return i
		}
	}
	return -1
}

func CheckDependencies(stacks map[int]TerraformStack) error {
	for i := range stacks {
		// Loop through each of the dependency paths
		for _, v := range stacks[i].config.Depends {
			// Checking for the presence of the path in the stacks slice
			if !dependencyInSlice(stacks, v) {
				// It's not there!!
			}
			// Need to check for dep loops here
		}
	}
	return nil
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
