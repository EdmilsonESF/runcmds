package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var names []string
var detach bool

var runCmd = &cobra.Command{
	Use:   "run [commands...]",
	Short: "Run multiple shell commands concurrently",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCommands(args, names, detach)
	},
}

func init() {
	runCmd.Flags().StringSliceVarP(&names, "name", "n", nil, "Names for each command (ex: --name api,test)")
	runCmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run commands in background")

	rootCmd.AddCommand(runCmd)
}

func runCommands(commands []string, names []string, detach bool) {
	colors := []func(a ...interface{}) string{
		color.New(color.FgCyan).SprintFunc(),
		color.New(color.FgGreen).SprintFunc(),
		color.New(color.FgMagenta).SprintFunc(),
		color.New(color.FgYellow).SprintFunc(),
		color.New(color.FgBlue).SprintFunc(),
		color.New(color.FgHiRed).SprintFunc(),
	}

	var wg sync.WaitGroup

	for i, cmdStr := range commands {
		wg.Add(1)

		name := fmt.Sprintf("cmd%d", i+1)
		if i < len(names) {
			name = names[i]
		}
		colorFunc := colors[i%len(colors)]

		go func(command, name string, colorize func(a ...interface{}) string) {
			defer wg.Done()
			prefix := colorize("[" + name + "]")

			fmt.Printf("%s Starting: %s\n", prefix, command)

			cmd := exec.Command("sh", "-c", command)
			output, err := cmd.CombinedOutput()

			for _, line := range strings.Split(string(output), "\n") {
				if strings.TrimSpace(line) != "" {
					fmt.Printf("%s %s\n", prefix, line)
				}
			}

			if err != nil {
				fmt.Printf("%s Erro: %v\n", prefix, err)
			}
		}(cmdStr, name, colorFunc)
	}

	if !detach {
		wg.Wait()
	}
}
