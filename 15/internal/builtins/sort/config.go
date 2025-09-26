package sort

import (
	"strconv"
)

// Config of sorting
type Config struct {
	Column       int    // -k
	Delimiter    string // N
	Numeric      bool   // -n
	Reverse      bool   // -r
	Unique       bool   // -u
	Month        bool   // -M
	IgnoreBlanks bool   // -b
	CheckSorted  bool   // -c
	Human        bool   // -h
}

// ParseConfig parses command line arguments into Config
func ParseConfig(args ...string) Config {
	cfg := Config{
		Delimiter: " ", // по умолчанию пробел
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-k":
			i++
			if i < len(args) {
				n, err := strconv.Atoi(args[i])
				if err == nil {
					cfg.Column = n
				}
			}
		case "-t":
			i++
			if i < len(args) {
				cfg.Delimiter = args[i]
			}
		case "-n":
			i++
			cfg.Numeric = true
		case "-r":
			i++
			cfg.Reverse = true
		case "-u":
			i++
			cfg.Unique = true
		case "-M":
			i++
			cfg.Month = true
		case "-b":
			i++
			cfg.IgnoreBlanks = true
		case "-c":
			i++
			cfg.CheckSorted = true
		case "-h":
			i++
			cfg.Human = true
		default:
		}
	}
	return cfg
}
