package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "terrarun",
	Short: "Terrarun is a helper for Terraform codebase",
	Long:  `Terrarun helps you with your CI with Terraform codebase`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Terrarun root")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("environment", "e", "", "name of current environment")
}
