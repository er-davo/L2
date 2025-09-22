package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"mycut/internal/config"
	"mycut/internal/cut"
)

func main() {
	cfg := config.Config{}
	flag.StringVar(&cfg.Fields, "f", "", "fields to select (e.g. 1,3-5)")
	flag.StringVar(&cfg.Delimiter, "d", "\t", "delimiter (default tab)")
	flag.BoolVar(&cfg.Separated, "s", false, "only print lines with delimiter")

	flag.Parse()

	if cfg.Fields == "" {
		fmt.Fprintln(os.Stderr, "error: -f flag is required")
		os.Exit(1)
	}

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
		os.Exit(1)
	}

	result, err := cut.Cut(lines, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cut error:", err)
		os.Exit(1)
	}

	for _, line := range result {
		fmt.Println(line)
	}
}
