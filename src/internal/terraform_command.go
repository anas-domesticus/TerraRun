package internal

import "os"

func NewTerraformCommand() Command {
	return Command{
		Binary:  "terraform",
		EnvVars: os.Environ(),
	}
}
