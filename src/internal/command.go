package internal

import (
	"bytes"
	"os/exec"
)

type CommandIface interface {
	Execute() error
}

type Command struct {
	Binary     string
	Parameters []Parameter
}

//Get config, inc env & cmd
//Get stacks
//build & run command against stack
//capture output, squirt to stdout

func (c *Command) Execute(cfg *Config, stack *TerraformStack) ([]byte, []byte, error) {
	// First we populate a slice of standard placeholders for this run
	stdPlaceholders := stack.GetStackPlaceholders()
	stdPlaceholders = append(stdPlaceholders, cfg.Env.GetEnvPlaceholder())

	// Getting the params we'll actually be using
	var params []string
	for i, v := range c.Parameters {
		c.Parameters[i].AddPlaceholders(stdPlaceholders)
		params = append(params, v.GetValue())
	}

	// Building the command
	cmd := exec.Command(c.Binary, params...)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	cmd.Dir = stack.Path
	err := cmd.Run()
	return stdOut.Bytes(), stdErr.Bytes(), err
}
