package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func anagramm(words []string) map[string][]string {
	angrms := map[string][]string{}

	for _, word := range words {
		word = strings.ToLower(word)
		runes := []rune(word)

		slices.Sort(runes)

		angrms[string(runes)] = append(angrms[string(runes)], word)
	}

	for k, v := range angrms {
		if len(v) == 1 {
			delete(angrms, k)
		}
	}

	finalAnagrams := make(map[string][]string, len(angrms))

	// skip one words, change key to first word and sorts
	for _, v := range angrms {
		if len(v) > 1 {
			finalAnagrams[v[0]] = v
			slices.Sort(finalAnagrams[v[0]])
		}
	}

	return finalAnagrams
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var s string

	s, _ = reader.ReadString('\n')

	lines := strings.Fields(s)

	for word, angrms := range anagramm(lines) {
		fmt.Printf("%s: %s\n", word, strings.Join(angrms, " "))
	}
}
