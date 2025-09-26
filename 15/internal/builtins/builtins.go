package builtins

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"minishell/internal/builtins/cut"
	"minishell/internal/builtins/grep"
	"minishell/internal/builtins/sort"
	"minishell/internal/models"
)

// ErrCmdNotFound is returned when a command is not found
var ErrCmdNotFound = fmt.Errorf("builtins: command not found")

// Builtins is a struct that represents the built-in commands
type Builtins struct {
	processes []models.Process
}

// Run executes a builtin command and returns its output
func (b *Builtins) Run(cmd string, input string, args ...string) (string, error) {
	switch cmd {
	case "cd":
		return "", b.Cd(args...)
	case "pwd":
		return b.Pwd()
	case "echo":
		return b.Echo(args...)
	case "kill":
		return "", b.Kill(args...)
	case "ps":
		return b.Ps(), nil
	case "grep":
		return b.Grep(input, args...), nil
	case "cut":
		return b.Cut(input, args...)
	case "sort":
		return b.Sort(input, args...)
	default:
		return "", ErrCmdNotFound
	}
}

// Cd changes the current directory
func (b *Builtins) Cd(args ...string) error {
	path := ""
	if len(args) == 0 {
		path = os.Getenv("HOME")
	} else {
		path = args[0]
	}
	return os.Chdir(path)
}

// Pwd prints the current directory
func (b *Builtins) Pwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("pwd error: %w", err)
	}
	return dir, nil
}

// Echo prints the arguments
func (b *Builtins) Echo(args ...string) (string, error) {
	noNewline := false
	interpretEscapes := false
	outArgs := []string{}

	for i, arg := range args {
		if i == 0 && arg == "-n" {
			noNewline = true
			continue
		}
		if i == 0 && arg == "-e" {
			interpretEscapes = true
			continue
		}
		outArgs = append(outArgs, arg)
	}

	out := strings.Join(outArgs, " ")
	if interpretEscapes {
		out = strings.ReplaceAll(out, `\n`, "\n")
		out = strings.ReplaceAll(out, `\t`, "\t")
		out = strings.ReplaceAll(out, `\r`, "\r")
	}

	if !noNewline {
		out += "\n"
	}

	return out, nil
}

// Kill sends a signal to a process
func (b *Builtins) Kill(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: kill [-SIGNAL] pid")
	}

	sig := syscall.SIGTERM
	pidArg := args[0]

	if strings.HasPrefix(pidArg, "-") && len(args) > 1 {
		switch pidArg {
		case "-9":
			sig = syscall.SIGKILL
		case "-15":
			sig = syscall.SIGTERM
		default:
			return fmt.Errorf("unsupported signal: %s", pidArg)
		}
		pidArg = args[1]
	}

	pid, err := strconv.Atoi(pidArg)
	if err != nil {
		return fmt.Errorf("invalid pid: %s", pidArg)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(sig)
}

// NewProcess adds a new process to the list of processes
func (b *Builtins) NewProcess(pid int, cmd string) {
	b.processes = append(b.processes, models.Process{PID: pid, Cmd: cmd})
}

// RemoveProcess removes a process from the list of processes
func (b *Builtins) RemoveProcess(pid int) {
	for i, p := range b.processes {
		if p.PID == pid {
			b.processes = append(b.processes[:i], b.processes[i+1:]...)
			return
		}
	}
}

// Ps returns the list of processes
func (b *Builtins) Ps() string {
	output := "PID    CMD"
	for _, p := range b.processes {
		output += "\n" + fmt.Sprintf("%-6d %s", p.PID, p.Cmd)
	}
	return output
}

// Grep filters lines
func (b *Builtins) Grep(lines string, args ...string) string {
	result := grep.Grep(strings.Split(lines, "\n"), grep.ParseConfig(args...))
	return strings.Join(result, "\n")
}

// Cut extracts columns from lines
func (b *Builtins) Cut(lines string, args ...string) (string, error) {
	result, err := cut.Cut(strings.Split(lines, "\n"), cut.ParseConfig(args...))
	return strings.Join(result, "\n"), err
}

// Sort sorts lines
func (b *Builtins) Sort(lines string, args ...string) (string, error) {
	result, err := sort.Sort(strings.Split(lines, "\n"), sort.ParseConfig(args...))
	return strings.Join(result, "\n"), err
}
