package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
	"os"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Performs a terraform validate against all the eligible stack directories",
	Long:  `Performs a terraform validate against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := internal.ForAllStacks(
			internal.Config{BaseDir: "./", Env: internal.Environment{Name: "dev"}},
			ValidateTerraform)
		if err != nil {
			fmt.Printf("Error occurred validating Terraform: %s", err.Error())
			os.Exit(1)
		}
	},
}

func ValidateTerraform(config internal.Config, stack internal.TerraformStack) (internal.ExecuteOutput, error) {
	initCmd := internal.GetTerraformInit()
	validateCmd := internal.GetTerraformValidate()
	output := initCmd.Execute(config, stack)
	if output.Error != nil {
		return internal.ExecuteOutput{}, output.Error
	}

	output = validateCmd.Execute(config, stack)
	return output, nil
}
