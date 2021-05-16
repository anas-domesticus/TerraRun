package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints Terrarun's version number",
	Long:  `Terrarun has a version number, this displays it, nothing more to say :)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Terrarun version")
	},
}
