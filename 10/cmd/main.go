package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"mysort/internal/config"
	"mysort/internal/sort"
)

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.Column, "k", 0, "sort by column number (tab-separated)")
	flag.StringVar(&cfg.Delimiter, "N", "\t", "use DELIM instead of TAB for field delimiter")
	flag.BoolVar(&cfg.Numeric, "n", false, "compare according to string numerical value")
	flag.BoolVar(&cfg.Reverse, "r", false, "reverse the result of comparisons")
	flag.BoolVar(&cfg.Unique, "u", false, "output only the first of an equal run")
	flag.BoolVar(&cfg.Month, "M", false, "compare by month name")
	flag.BoolVar(&cfg.IgnoreBlanks, "b", false, "ignore trailing blanks")
	flag.BoolVar(&cfg.CheckSorted, "c", false, "check whether input is sorted")
	flag.BoolVar(&cfg.Human, "h", false, "compare human readable numbers (e.g. 2K, 1G)")

	flag.Parse()

	args := flag.Args()
	var r io.Reader

	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		r = f
	} else {
		r = os.Stdin
	}

	lines, err := readLines(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines, err = sort.Sort(lines, cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
