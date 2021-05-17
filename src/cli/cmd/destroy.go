package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
)

func init() {
	rootCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Performs a terraform destroy against all the eligible stack directories",
	Long:  `Performs a terraform destroy against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ForAllStacks(
			internal.Config{BaseDir: "./", Env: internal.Environment{Name: "dev"}},
			func(config internal.Config, stack internal.TerraformStack) error {
				fmt.Printf("Destroying %s\n", stack.Path)
				return nil
			})
	},
}
