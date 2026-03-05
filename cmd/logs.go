package cmd

import (
	"fema-cli/internal/logs"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Tail /var/log/*.log",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Tail()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
