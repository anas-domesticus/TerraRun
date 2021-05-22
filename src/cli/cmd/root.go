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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "name of current environment")
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "./", "directory within which to look for stacks")
	rootCmd.PersistentFlags().BoolVar(&debugLogging, "debug", false, "whether to output debug logging")
}
