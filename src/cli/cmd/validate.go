package cmd

import (
	"errors"
	"fmt"
	"github.com/anas-domesticus/TerraRun/src/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
	config.TFPluginCacheDir = GetCacheDir()
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
		return errors.New("one or more stacks failed validation")
	}
	return nil
}
