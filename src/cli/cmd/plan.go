package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
	"os"
)

func init() {
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Performs a terraform plan against all the eligible stack directories",
	Long:  `Performs a terraform plan against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		err := CheckAllPlanOutputs(
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

func CheckAllPlanOutputs(config internal.Config) error {
	config.TFPluginCacheDir = GetCacheDir()
	outputs, err := internal.ForAllStacks(
		config,
		internal.PlanStack)
	if err != nil {
		return err
	}
	errOccurred := false
	for _, out := range outputs {
		if internal.PlanWasSuccessful(out) {
			color.Green("Plan OK for %s\n", out.Stack.Path)
		} else {
			color.Red("Plan failed for %s\n", out.Stack.Path)
			color.Red(string(out.StdOut))
			color.Red(string(out.StdErr))
			errOccurred = true
		}
	}
	if errOccurred {
		return errors.New("one or more stacks failed to plan\n")
	}
	return nil
}
