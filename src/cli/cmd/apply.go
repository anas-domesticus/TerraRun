package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Performs a terraform apply against all the eligible stack directories",
	Long:  `Performs a terraform apply against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ForAllStacks(
			internal.Config{BaseDir: "./", Env: internal.Environment{Name: "dev"}},
			func(config internal.Config, stack internal.TerraformStack) error {
				fmt.Printf("Applying %s\n", stack.Path)
				return nil
			})
	},
}
