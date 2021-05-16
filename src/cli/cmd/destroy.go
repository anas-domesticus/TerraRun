package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Performs a terraform destroy against all the eligible stack directories",
	Long:  `Performs a terraform destroy against all the eligible stack directories`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Destroy!")
	},
}
