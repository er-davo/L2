package shell

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"minishell/internal/builtins"
	"minishell/internal/models"
	"minishell/internal/parser"
)

// Colors
var (
	red     = "\u001b[31m"
	blue    = "\u001b[34m"
	green   = "\u001b[32m"
	yellow  = "\u001b[33m"
	magenta = "\u001b[35m"
	reset   = "\u001b[0m"
)

// Shell represents a shell with builtins commands
type Shell struct {
	builtin builtins.Builtins

	mu             sync.Mutex
	currentProcess *os.Process
}

// Run runs the shell
func (s *Shell) Run() error {
	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		for range sigChan {
			s.mu.Lock()
			proc := s.currentProcess
			s.mu.Unlock()
			if proc != nil {
				proc.Kill()
				fmt.Println()
				s.printConsoleLine()
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		if err := s.printConsoleLine(); err != nil {
			return err
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { // Ctrl+D
				break
			}
			readingLineError(err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			break
		}

		jobs := parser.ParseCommand(line)
		for _, job := range jobs {
			out, err := s.runJob(job)
			if err != nil {
				if errors.Is(err, builtins.ErrCmdNotFound) {
					commandNotFound(job.Pipelines[0][0].Name)
				} else {
					handleError(err)
				}
			} else if out != "" {
				fmt.Println(out)
			}
		}
	}
	return nil
}

func (s *Shell) runJob(job models.Job) (string, error) {
	var result string
	var err error
	for _, pipeline := range job.Pipelines {
		result, err = s.runPipeline(pipeline)
		if err != nil {
			break
		}
	}
	return result, err
}

func (s *Shell) runPipeline(pipeline models.Pipeline) (string, error) {
	var currentInput string

	for _, cmd := range pipeline {
		var output bytes.Buffer

		outStr, runErr := s.builtin.Run(cmd.Name, currentInput, cmd.Args...)
		if runErr == nil {
			currentInput = outStr
			continue
		} else if !errors.Is(runErr, builtins.ErrCmdNotFound) {
			return "", runErr
		}

		execCmd := exec.Command(cmd.Name, cmd.Args...)

		if cmd.Stdin != "" {
			data, readErr := os.ReadFile(cmd.Stdin)
			if readErr != nil {
				return "", fmt.Errorf("cannot read file %s: %w", cmd.Stdin, readErr)
			}
			execCmd.Stdin = bytes.NewReader(data)
		} else if currentInput != "" {
			execCmd.Stdin = strings.NewReader(currentInput)
		} else {
			execCmd.Stdin = os.Stdin
		}

		execCmd.Stdout = &output
		execCmd.Stderr = os.Stderr

		if err := execCmd.Start(); err != nil {
			return "", err
		}

		s.mu.Lock()
		s.currentProcess = execCmd.Process
		s.mu.Unlock()
		s.builtin.NewProcess(execCmd.Process.Pid, cmd.Name)

		err := execCmd.Wait()

		s.mu.Lock()
		s.currentProcess = nil
		s.mu.Unlock()
		s.builtin.RemoveProcess(execCmd.Process.Pid)

		if err != nil {
			return "", err
		}

		currentInput = output.String()

		if cmd.Stdout != "" {
			flags := os.O_CREATE | os.O_WRONLY
			if cmd.Append {
				flags |= os.O_APPEND
			} else {
				flags |= os.O_TRUNC
			}
			if writeErr := os.WriteFile(cmd.Stdout, []byte(currentInput), 0644); writeErr != nil {
				return "", fmt.Errorf("cannot write to file %s: %w", cmd.Stdout, writeErr)
			}
			currentInput = ""
		}
	}

	return currentInput, nil
}

func (s *Shell) printConsoleLine() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Print(magenta + "minishell" + reset + " " + blue + wd + reset + " $ ")
	return nil
}

func commandNotFound(cmd string) {
	fmt.Fprintf(os.Stderr, "%s command not found\n", cmd)
}

func readingLineError(err error) {
	fmt.Fprintln(os.Stderr, "error reading line:", err)
}

func handleError(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
}
