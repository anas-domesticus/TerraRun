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
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Performs a terraform apply against all the eligible stack directories",
	Long:  `Performs a terraform apply against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		err := CheckAllApplyOutputs(BuildConfig())
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
	},
}

func CheckAllApplyOutputs(config internal.Config) error {
	config.TFPluginCacheDir = GetCacheDir()
	outputs, err := internal.ForAllStacks(
		config,
		internal.ApplyStack)
	if err != nil {
		return err
	}
	errOccurred := false
	for _, out := range outputs {
		if internal.ApplyWasSuccessful(out) {
			color.Green("Apply OK for %s\n", out.Stack.Path)
			if config.Debug {
				color.Green(string(out.StdOut))
			}
		} else {
			color.Red("Apply failed for %s\n", out.Stack.Path)
			color.Red(string(out.StdOut))
			color.Red(string(out.StdErr))
			errOccurred = true
		}
	}
	if errOccurred {
		return errors.New("one or more stacks failed to apply")
	}
	return nil
}
