package parser

import (
	"os"
	"strings"

	"minishell/internal/models"
)

// ParseSingleCommand parses a command string into a Command struct
func ParseSingleCommand(input string) models.Command {
	tokens := strings.Fields(input)

	cmd := models.Command{}
	if len(tokens) == 0 {
		return cmd
	}

	cmd.Name = tokens[0]
	args := []string{}

	for i := 1; i < len(tokens); i++ {
		switch tokens[i] {
		case ">":
			if i+1 < len(tokens) {
				cmd.Stdout = tokens[i+1]
				cmd.Append = false
				i++
			}
		case ">>":
			if i+1 < len(tokens) {
				cmd.Stdout = tokens[i+1]
				cmd.Append = true
				i++
			}
		case "<":
			if i+1 < len(tokens) {
				cmd.Stdin = tokens[i+1]
				i++
			}
		default:
			args = append(args, tokens[i])
		}
	}

	cmd.Args = substituteEnv(args)
	return cmd
}

// ParseCommand parses a full input string into Jobs
func ParseCommand(input string) []models.Job {
	tokens := splitByOperators(input)

	jobs := []models.Job{}
	for _, t := range tokens {
		pipeParts := strings.Split(strings.TrimSpace(t.Cmd), "|")
		pipeline := make(models.Pipeline, 0, len(pipeParts))

		for _, cmdStr := range pipeParts {
			cmdStr = strings.TrimSpace(cmdStr)
			if cmdStr == "" {
				continue
			}
			pipeline = append(pipeline, ParseSingleCommand(cmdStr))
		}

		job := models.Job{
			Pipelines: []models.Pipeline{pipeline},
			CondAfter: t.Op,
		}
		jobs = append(jobs, job)
	}
	return jobs
}

type token struct {
	Cmd string
	Op  models.Operator
}

func splitByOperators(input string) []token {
	result := []token{}
	cur := ""
	lastOp := models.Operator("")

	i := 0
	for i < len(input) {
		if strings.HasPrefix(input[i:], string(models.And)) {
			result = append(result, token{Cmd: strings.TrimSpace(cur), Op: lastOp})
			cur = ""
			lastOp = models.And
			i += len(models.And)
		} else if strings.HasPrefix(input[i:], string(models.Or)) {
			result = append(result, token{Cmd: strings.TrimSpace(cur), Op: lastOp})
			cur = ""
			lastOp = models.Or
			i += len(models.Or)
		} else {
			cur += string(input[i])
			i++
		}
	}

	if strings.TrimSpace(cur) != "" {
		result = append(result, token{Cmd: strings.TrimSpace(cur), Op: lastOp})
	}

	return result
}

func substituteEnv(tokens []string) []string {
	result := make([]string, len(tokens))
	for i, tok := range tokens {
		if strings.HasPrefix(tok, "$") && len(tok) > 1 {
			val := os.Getenv(tok[1:]) // срезаем $
			result[i] = val
		} else {
			if strings.Contains(tok, "$") {
				result[i] = os.ExpandEnv(tok)
			} else {
				result[i] = tok
			}
		}
	}
	return result
}
