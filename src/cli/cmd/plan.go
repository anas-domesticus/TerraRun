package cmd

import (
	"fmt"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Performs a terraform plan against all the eligible stack directories",
	Long:  `Performs a terraform plan against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Plan!")
		command := internal.Command{
			Binary:     "terraform",
			Parameters: []internal.Parameter{&internal.SimpleParameter{Value: "plan"}},
		}
	},
}
