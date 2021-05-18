package cmd

import (
	"errors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
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
		err := CheckAllValidateOutputs(internal.Config{BaseDir: "./", Env: internal.Environment{Name: "dev"}})
		if err != nil {
			os.Exit(1)
		}
	},
}

func CheckAllValidateOutputs(config internal.Config) error {
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
