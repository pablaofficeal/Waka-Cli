package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	interval int
	once     bool
)

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Auto-format Go files and commit changes",
	Long: `Run Python formatter that:
- Formats all Go files using go fmt
- Commits and pushes changes if any
- Runs continuously with specified interval

Press Ctrl+C to stop.`,
	Run: runFormatter,
}

func init() {
	rootCmd.AddCommand(fmtCmd)
	fmtCmd.Flags().IntVarP(&interval, "interval", "i", 60, "Interval in seconds between runs")
	fmtCmd.Flags().BoolVarP(&once, "once", "o", false, "Run only once and exit")
}

func runFormatter(cmd *cobra.Command, args []string) {
	if once {
		runFormatOnce()
		return
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Printf("🔄 Starting formatter with %d second interval. Press Ctrl+C to stop.\n", interval)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Run immediately on start
	runFormatOnce()

	for {
		select {
		case <-ticker.C:
			runFormatOnce()
		case <-sigChan:
			fmt.Println("\n🛑 Shutting down formatter...")
			cancel()
			return
		case <-ctx.Done():
			return
		}
	}
}

func runFormatOnce() {
	fmt.Println("\n📁 Running formatter...")

	// Get the directory where the CLI is running from
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("❌ Error getting executable path: %v\n", err)
		return
	}

	// Assume formater.py is in the same directory as the CLI binary
	formaterPath := fmt.Sprintf("%s/formater.py", execPath[:len(execPath)-len("/fema")])

	// If not found, try current working directory
	if _, err := os.Stat(formaterPath); os.IsNotExist(err) {
		formaterPath = "./formater.py"
	}

	// Check if Python script exists
	if _, err := os.Stat(formaterPath); os.IsNotExist(err) {
		fmt.Printf("❌ formater.py not found at %s\n", formaterPath)
		return
	}

	// Run the Python script
	cmd := exec.Command("python3", formaterPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Error running formatter: %v\n", err)
		return
	}

	fmt.Println("✅ Format cycle completed")
}
