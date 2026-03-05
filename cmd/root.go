package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fema",
	Short: "Fema DevOps CLI",
	Long:  "A lightweight DevOps CLI tool written in Go for logs, scanning, sysinfo and networking.",
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
