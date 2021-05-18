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
	//command := internal.GetTerraformValidate()
	return internal.ExecuteOutput{}, nil
}
