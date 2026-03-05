package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping [url]",
	Short: "Check HTTP endpoint availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		client := http.Client{Timeout: 3 * time.Second}
		res, err := client.Get(url)

		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}

		fmt.Printf("Status: %s (%d)\n", res.Status, res.StatusCode)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
