package cmd

import (
	"errors"
	"fmt"
	"github.com/anas-domesticus/TerraRun/src/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().BoolVarP(&outputPlanReport, "report", "r", false, "if set, plan will output an HTML plan report to plan.html")
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
	var stacksForReport []internal.TerraformStack
	errOccurred := false
	for _, out := range outputs {
		if internal.PlanWasSuccessful(out) {
			color.Green("Plan OK for %s\n", out.Stack.Path)
			if outputPlanReport {
				stacksForReport = append(stacksForReport, out.Stack)
			}
		} else {
			color.Red("Plan failed for %s\n", out.Stack.Path)
			color.Red(string(out.StdOut))
			color.Red(string(out.StdErr))
			errOccurred = true
		}
	}
	if errOccurred {
		return errors.New("one or more stacks failed to plan")
	}

	// Report generation
	if outputPlanReport {
		var reportSet internal.ShowOutputSet
		for _, stack := range stacksForReport {
			out, err := internal.GetShowOutput(config, stack)
			if err != nil {
				fmt.Printf("failed to get JSON data for: %s", stack.Path)
				continue
			}
			reportSet = append(reportSet, out)
		}
		err := ioutil.WriteFile("report.html", reportSet.GenerateHTMLReport(), 0644)
		if err != nil {
			fmt.Printf("Failed to write report: %s", err.Error())
		}
	}

	return nil
}
