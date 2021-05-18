package internal

import (
	"bytes"
	"os/exec"
)

type Command struct {
	Binary     string
	Parameters []Parameter
}

type ExecuteOutput struct {
	Stack   TerraformStack
	Command *Command
	StdOut  []byte
	StdErr  []byte
	Error   error
}

func (c *Command) Execute(cfg Config, stack TerraformStack) ExecuteOutput {
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
		Stack:   stack,
		Command: c,
		StdOut:  stdOut.Bytes(),
		StdErr:  stdErr.Bytes(),
		Error:   err,
	}
}
