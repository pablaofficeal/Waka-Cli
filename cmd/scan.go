package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan /var/log directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scanning /var/log ...")

		entries, _ := os.ReadDir("/var/log")
		for _, e := range entries {
			info, _ := e.Info()
			fmt.Println(info.Name(), info.Size(), "bytes", info.ModTime().Format(time.RFC822))
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
