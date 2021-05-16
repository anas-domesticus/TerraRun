package cmd

import (
	"fmt"

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
		fmt.Println("Apply!")
	},
}
