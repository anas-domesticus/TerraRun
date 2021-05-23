package internal

import (
	"fmt"
	"os"
)

func GetTerraformInit(cacheDir string) Command {
	cmd := NewTerraformCommand()
	envVars := append(os.Environ(), fmt.Sprintf("TF_PLUGIN_CACHE_DIR=%s", cacheDir))
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "init"},
		&SimpleParameter{Value: "-input=false"},
	}...)
	cmd.EnvVars = envVars
	return cmd
}

func InitStack(config Config, stack TerraformStack) (ExecuteOutput, error) {
	initCmd := GetTerraformInit(config.TFPluginCacheDir)
	output, err := initCmd.Execute(config, stack)
	if output.Error != nil {
		return output, output.Error
	}
	return output, err
}
