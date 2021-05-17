package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists eligible terraform stack directories",
	Long:  `Lists eligible terraform stack directories, searching from pwd`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ForAllStacks(
			internal.Config{BaseDir: "./", Env: internal.Environment{Name: "dev"}},
			func(config internal.Config, stack internal.TerraformStack) error {
				fmt.Printf("%s\n", stack.Path)
				return nil
			})
	},
}
