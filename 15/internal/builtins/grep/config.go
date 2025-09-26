package grep

import "strconv"

// Config for grep
type Config struct {
	After       int
	Before      int
	Context     int
	CountOnly   bool
	IgnoreCase  bool
	InvertMatch bool
	Fixed       bool
	LineNum     bool
	Pattern     string
	File        string
}

// ParseConfig parses command line arguments into Config
func ParseConfig(args ...string) Config {
	cfg := Config{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-A":
			i++
			n, _ := strconv.Atoi(args[i])
			cfg.After = n
		case "-B":
			i++
			n, _ := strconv.Atoi(args[i])
			cfg.Before = n
		case "-C":
			i++
			n, _ := strconv.Atoi(args[i])
			cfg.Context = n
		case "-c":
			cfg.CountOnly = true
		case "-i":
			cfg.IgnoreCase = true
		case "-v":
			cfg.InvertMatch = true
		case "-F":
			cfg.Fixed = true
		case "-n":
			cfg.LineNum = true
		default:
			if cfg.Pattern == "" {
				cfg.Pattern = args[i]
			} else {
				cfg.File = args[i]
			}
		}
	}

	return cfg
}
