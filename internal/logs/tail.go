package logs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func Tail() {
	files, _ := filepath.Glob("/var/log/*.log")

	for _, file := range files {
		fmt.Printf("\n---- %s ----\n", file)

		f, err := os.Open(file)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		defer f.Close()

		var lines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		// вывести последние 20
		start := len(lines) - 20
		if start < 0 {
			start = 0
		}

		for _, line := range lines[start:] {
			fmt.Println(line)
		}
	}
}
