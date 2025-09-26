package models

// Command represents a single shell command with its arguments and I/O redirections.
type Command struct {
	Name   string   // Name is the executable or builtin command name.
	Args   []string // Args contains the arguments passed to the command.
	Stdin  string   // Stdin is the input file name; empty string means use default input.
	Stdout string   // Stdout is the output file name; empty string means use default output.
	Append bool     // Append indicates whether the output should be appended instead of truncated.
}

// Pipeline represents a sequence of commands connected by pipes.
type Pipeline []Command

// Operator represents a conditional operator between pipelines.
type Operator string

const (
	// And represents the "&&" operator, which runs the next pipeline only if the previous succeeded.
	And Operator = "&&"
	// Or represents the "||" operator, which runs the next pipeline only if the previous failed.
	Or Operator = "||"
)

// Job represents a sequence of pipelines combined with conditional operators.
type Job struct {
	Pipelines []Pipeline // Pipelines contains the list of pipelines to execute.
	CondAfter Operator   // CondAfter specifies the operator to apply after this job.
}

// Process represents a running process started by the shell.
type Process struct {
	PID int    // PID is the process ID of the running command.
	Cmd string // Cmd is the command string associated with the process.
}
