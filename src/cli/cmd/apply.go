package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
	"os"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Performs a terraform apply against all the eligible stack directories",
	Long:  `Performs a terraform apply against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := internal.ForAllStacks(
			internal.Config{BaseDir: "./", Env: internal.Environment{Name: "dev"}},
			ApplyTerraform)
		if err != nil {
			fmt.Printf("Error occurred applying Terraform: %s", err.Error())
			os.Exit(1)
		}
	},
}

func ApplyTerraform(config internal.Config, stack internal.TerraformStack) (internal.ExecuteOutput, error) {
	//command := internal.GetTerraformApply()
	return internal.ExecuteOutput{}, nil
}
