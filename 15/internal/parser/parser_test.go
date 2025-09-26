package parser_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"minishell/internal/models"
	"minishell/internal/parser"
)

func TestParseSingleCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected models.Command
	}{
		{
			name:  "simple command",
			input: "ls -l -a",
			expected: models.Command{
				Name: "ls",
				Args: []string{"-l", "-a"},
			},
		},
		{
			name:  "output redirection",
			input: "echo hi > out.txt",
			expected: models.Command{
				Name:   "echo",
				Args:   []string{"hi"},
				Stdout: "out.txt",
				Append: false,
			},
		},
		{
			name:  "append redirection",
			input: "echo hi >> log.txt",
			expected: models.Command{
				Name:   "echo",
				Args:   []string{"hi"},
				Stdout: "log.txt",
				Append: true,
			},
		},
		{
			name:  "input redirection",
			input: "cat < file.txt",
			expected: models.Command{
				Name:  "cat",
				Args:  []string{},
				Stdin: "file.txt",
			},
		},
		{
			name:  "env substitution",
			input: "echo $GOPATH",
			expected: models.Command{
				Name: "echo",
				Args: []string{os.Getenv("GOPATH")},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: models.Command{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.ParseSingleCommand(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestSplitByOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		ops      []models.Operator
	}{
		{
			name:     "single command",
			input:    "ls -l",
			expected: []string{"ls -l"},
			ops:      []models.Operator{""},
		},
		{
			name:     "with &&",
			input:    "echo hi && ls",
			expected: []string{"echo hi", "ls"},
			ops:      []models.Operator{"", models.And},
		},
		{
			name:     "with ||",
			input:    "make || echo fail",
			expected: []string{"make", "echo fail"},
			ops:      []models.Operator{"", models.Or},
		},
		{
			name:     "mixed && and ||",
			input:    "cmd1 && cmd2 || cmd3",
			expected: []string{"cmd1", "cmd2", "cmd3"},
			ops:      []models.Operator{"", models.And, models.Or},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := parser.ParseCommand(tt.input)
			var cmds []string
			var ops []models.Operator
			for _, job := range tokens {
				for _, pipe := range job.Pipelines {
					for _, cmd := range pipe {
						cmds = append(cmds, stringifyCommand(cmd))
					}
				}
				ops = append(ops, job.CondAfter)
			}
			assert.Equal(t, tt.expected, cmds)
			assert.Equal(t, tt.ops, ops)
		})
	}
}

func stringifyCommand(cmd models.Command) string {
	if len(cmd.Args) == 0 {
		return cmd.Name
	}
	res := cmd.Name
	for _, a := range cmd.Args {
		res += " " + a
	}
	return res
}
