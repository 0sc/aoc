package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day5"] = day5
}

func day5(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var isUpdate bool
	rules := map[string][]string{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			isUpdate = true
			continue
		}

		if !isUpdate {
			pages := strings.Split(line, "|")
			rules[pages[0]] = append(rules[pages[0]], pages[1])

			continue
		}

		if v, ok := verify(line, rules); ok {
			partOne += v

			continue
		}

		partTwo += reorder(line, rules)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func verify(line string, rules map[string][]string) (int, bool) {
	pages := strings.Split(line, ",")
	seen := map[string]bool{}

	for _, page := range pages {
		for _, fp := range rules[page] {
			if seen[fp] {
				return 0, false
			}
		}

		seen[page] = true
	}

	m := len(pages) / 2
	return toInt(pages[m]), true
}

func reorder(line string, rules map[string][]string) int {
	pages := strings.Split(line, ",")

	seen := map[string]bool{}
	for _, page := range pages {
		seen[page] = true
	}

	wait := map[string]int{}
	after := map[string][]string{}
	for page := range seen {
		for _, fp := range rules[page] {
			if seen[fp] {
				wait[fp]++
				after[page] = append(after[page], fp)
			}
		}
	}

	var ready []string
	for _, page := range pages {
		if wait[page] == 0 {
			ready = append(ready, page)
		}
	}

	m := len(pages) / 2
	for i := 0; i < len(ready); i++ {
		page := ready[i]

		if i == m {
			return toInt(page)
		}

		for _, np := range after[page] {
			wait[np]--
			if wait[np] == 0 {
				ready = append(ready, np)
			}
		}
	}

	panic("reorder failed") // should not happen
}
