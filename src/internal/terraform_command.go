package internal

import "os"

func NewTerraformCommand() Command {
	return Command{
		Binary:  "terraform",
		EnvVars: os.Environ(),
	}
}

func GetTerraformInit() Command {
	cmd := NewTerraformCommand()
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "init"},
		&SimpleParameter{Value: "-input=false"},
	}...)
	return cmd
}
