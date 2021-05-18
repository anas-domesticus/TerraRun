package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
	"io/ioutil"
	"os"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Performs a terraform validate against all the eligible stack directories",
	Long:  `Performs a terraform validate against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		err := CheckAllValidateOutputs(
			internal.Config{
				BaseDir: directory,
				Env:     internal.Environment{Name: environment},
			})
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
	},
}

func CheckAllValidateOutputs(config internal.Config) error {
	dir, err := ioutil.TempDir("", "tf-cache")
	if err != nil {
		fmt.Printf("Failed to create temporary directory")
		os.Exit(1)
	}
	config.TFPluginCacheDir = dir
	outputs, err := internal.ForAllStacks(
		config,
		internal.ValidateStack)
	if err != nil {
		return err
	}
	errOccurred := false
	for _, out := range outputs {
		if internal.ValidateWasSuccessful(out) {
			color.Green("Validation passed for %s\n", out.Stack.Path)
		} else {
			color.Red("Validation failed for %s\n", out.Stack.Path)
			color.Red(string(out.StdOut))
			color.Red(string(out.StdErr))
			errOccurred = true
		}
	}
	if errOccurred {
		return errors.New("one or more stacks failed validation\n")
	}
	return nil
}
