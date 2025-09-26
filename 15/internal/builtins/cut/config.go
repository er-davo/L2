package cut

// Config for cut
type Config struct {
	Fields    string
	Delimiter string
	Separated bool
}

// ParseConfig parses the command line arguments
func ParseConfig(args ...string) Config {
	cfg := Config{Delimiter: "\t"}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f":
			i++
			cfg.Fields = args[i]
		case "-d":
			i++
			cfg.Delimiter = args[i]
		case "-s":
			cfg.Separated = true
		default:
		}
	}

	return cfg
}
