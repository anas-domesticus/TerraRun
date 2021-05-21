package cmd

import (
	"fmt"
	"github.com/anas-domesticus/TerraRun/src/internal"
	"github.com/spf13/cobra"
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
			internal.Config{BaseDir: directory, Env: internal.Environment{Name: environment}},
			ListStacks)
	},
}

func ListStacks(config internal.Config, stack internal.TerraformStack) (internal.ExecuteOutput, error) {
	fmt.Println(stack.Path)
	return internal.ExecuteOutput{}, nil
}
