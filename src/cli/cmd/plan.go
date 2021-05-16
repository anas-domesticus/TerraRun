package cmd

import (
	"fmt"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
	"os"

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
		stacks, err := internal.FindAllStacks("./")
		if err != nil {
			fmt.Printf("Error occurred: %s\n", err.Error())
			os.Exit(1)
		}
		command := internal.Command{
			Binary: "terraform",
			Parameters: []internal.Parameter{
				&internal.SimpleParameter{Value: "plan"},
				&internal.SimpleParameter{Value: "-out=plan.tfplan"},
			},
		}
	},
}
