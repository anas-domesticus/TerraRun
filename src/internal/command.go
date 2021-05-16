package internal

import (
	"bytes"
	"os/exec"
)

type CommandIface interface {
	Execute(*Config, TerraformStack)
	ExecuteForStacks(*Config, []TerraformStack) []ExecuteOutput
}

type Command struct {
	Binary     string
	Parameters []Parameter
}

type ExecuteOutput struct {
	Stack  TerraformStack
	StdOut []byte
	StdErr []byte
	Error  error
}

//Get config, inc env & cmd
//Get stacks
//build & run command against stack
//capture output, squirt to stdout

func (c *Command) Execute(cfg *Config, stack TerraformStack) ExecuteOutput {
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
	return ExecuteOutput{
		Stack:  stack,
		StdOut: stdOut.Bytes(),
		StdErr: stdErr.Bytes(),
		Error:  err,
	}
}

func (c *Command) ExecuteForStacks(cfg *Config, stacks []TerraformStack) []ExecuteOutput {
	var outSlice []ExecuteOutput
	for i, _ := range stacks {
		outSlice = append(outSlice, c.Execute(cfg, stacks[i]))
	}
	return outSlice
}
