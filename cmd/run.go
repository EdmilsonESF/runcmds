package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var names []string

var runCmd = &cobra.Command{
	Use:   "run [commands...]",
	Short: "Run multiple shell commands concurrently",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCommands(args, names)
	},
}

func init() {
	runCmd.Flags().StringSliceVarP(&names, "name", "n", nil, "Names for each command (ex: --name api,test)")

	rootCmd.AddCommand(runCmd)
}

func runCommands(commands []string, names []string) {
	colors := []func(a ...interface{}) string{
		color.New(color.FgCyan).SprintFunc(),
		color.New(color.FgGreen).SprintFunc(),
		color.New(color.FgMagenta).SprintFunc(),
		color.New(color.FgYellow).SprintFunc(),
		color.New(color.FgBlue).SprintFunc(),
		color.New(color.FgHiRed).SprintFunc(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, len(commands))

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

			cmd := exec.CommandContext(ctx, "sh", "-c", command)
			stdout, _ := cmd.StdoutPipe()
			stderr, _ := cmd.StderrPipe()

			pty, err := pty.Start(cmd)
			if err != nil {
				fmt.Printf("%s Error starting command: %v\n", prefix, err)
				errChan <- err
				return
			}
			defer pty.Close()

			logPipe := func(pipe io.ReadCloser, prefix string) {
				defer pipe.Close()
				scanner := bufio.NewScanner(pipe)
				for scanner.Scan() {
					fmt.Printf("%s %s\n", prefix, scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					fmt.Printf("%s Error reading pipe: %v\n", prefix, err)
				}
			}

			go logPipe(stdout, prefix)
			go logPipe(stderr, prefix)

			if err := cmd.Wait(); err != nil {
				fmt.Printf("%s Command failed: %v\n", prefix, err)
				errChan <- err
			}
		}(cmdStr, name, colorFunc)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err, ok := <-errChan; ok {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
