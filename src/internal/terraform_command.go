package internal

import "os"

func NewTerraformCommand() Command {
	return Command{
		Binary:  "terraform",
		EnvVars: append(os.Environ(), "TF_IN_AUTOMATION=true"),
	}
}
