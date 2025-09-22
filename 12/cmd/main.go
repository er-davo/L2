package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"mygrep/internal/config"
	"mygrep/internal/grep"
)

func main() {
	cfg := config.Config{}

	flag.StringVar(&cfg.Pattern, "pattern", "", "search pattern (required)")
	flag.IntVar(&cfg.After, "A", 0, "print N lines after match")
	flag.IntVar(&cfg.Before, "B", 0, "print N lines before match")
	flag.IntVar(&cfg.Context, "C", 0, "print N lines around match (before and after)")
	flag.BoolVar(&cfg.CountOnly, "c", false, "print only the count of matching lines")
	flag.BoolVar(&cfg.IgnoreCase, "i", false, "ignore case distinctions")
	flag.BoolVar(&cfg.InvertMatch, "v", false, "invert match: select non-matching lines")
	flag.BoolVar(&cfg.Fixed, "F", false, "treat pattern as a fixed string (not regex)")
	flag.BoolVar(&cfg.LineNum, "n", false, "print line number with output lines")
	flag.Parse()

	var lines []string
	var filename string

	if flag.NArg() > 0 {
		filename = flag.Arg(0)
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("failed to read file: %v", err)
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			log.Fatal("file or stdin input required")
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("failed to read stdin: %v", err)
		}
	}

	result := grep.Grep(lines, cfg)

	fmt.Println(strings.Join(result, "\n"))
}
