package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var sysCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Show system information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Go Version:", runtime.Version())
		fmt.Println("OS:", runtime.GOOS)
		fmt.Println("ARCH:", runtime.GOARCH)
		fmt.Println("CPUs:", runtime.NumCPU())
		fmt.Println("Goroutines:", runtime.NumGoroutine())
	},
}

func init() {
	rootCmd.AddCommand(sysCmd)
}
