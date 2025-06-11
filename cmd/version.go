package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set during build time
var Version = "unknown"

// rootCmd represents the base command when called without any subcommands
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints version of buildx",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
