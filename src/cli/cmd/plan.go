package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/lewisedginton/aws_common/terraform_wrapper/src/internal"
)

func init() {
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Performs a terraform plan against all the eligible stack directories",
	Long:  `Performs a terraform plan against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ForAllStacks(
			internal.Config{BaseDir: "./", Env: internal.Environment{Name: "dev"}},
			PlanTerraform)
	},
}

func PlanTerraform(config internal.Config, stack internal.TerraformStack) (internal.ExecuteOutput, error) {
	//command := internal.Command{
	//	Binary: "terraform",
	//	Parameters: []internal.Parameter{
	//		&internal.SimpleParameter{Value: "plan"},
	//	},
	//}
	return internal.ExecuteOutput{}, nil
}
