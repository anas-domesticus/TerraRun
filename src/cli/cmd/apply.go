package cmd

import (
	"fmt"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Performs a terraform apply against all the eligible stack directories",
	Long:  `Performs a terraform apply against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		stacks, err := internal.FindAllStacks("./")
		if err != nil {
			fmt.Printf("Error occurred: %s\n", err.Error())
			os.Exit(1)
		}
		command := internal.Command{
			Binary: "terraform",
			Parameters: []internal.Parameter{
				&internal.SimpleParameter{Value: "apply"},
				&internal.SimpleParameter{Value: "plan.tfplan"},
			},
		}
	},
}
