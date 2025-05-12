package cmd

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [commands...]",
	Short: "Run multiple shell commands concurrently",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCommands(args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runCommands(commands []string) {
	var wg sync.WaitGroup

	for i, cmdStr := range commands {
		wg.Add(1)

		go func(index int, command string) {
			defer wg.Done()

			name := fmt.Sprintf("cmd%d", index+1)
			fmt.Printf("[%s] Starting: %s\n", name, command)

			cmd := exec.Command("sh", "-c", command)
			output, err := cmd.CombinedOutput()

			fmt.Printf("[%s] Output:\n%s\n", name, string(output))
			if err != nil {
				fmt.Printf("[%s] Error: %v\n", name, err)
			}
		}(i, cmdStr)
	}

	wg.Wait()
}
