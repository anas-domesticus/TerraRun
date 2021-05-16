package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
	"os"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists eligible terraform stack directories",
	Long:  `Lists eligible terraform stack directories, searching from pwd`,
	Run: func(cmd *cobra.Command, args []string) {
		stacks, err := internal.FindAllStacks("./")
		if err != nil {
			fmt.Printf("Error occurred: %s\n", err.Error())
			os.Exit(1)
		}
		for _, s := range stacks {
			fmt.Printf("%s\n", s.Path)
		}
	},
}
