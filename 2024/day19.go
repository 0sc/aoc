package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day19"] = day19
}

func day19(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var isDesign bool
	var patterns, designs []string

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			isDesign = true
			continue
		}

		if isDesign {
			designs = append(designs, line)
			continue
		}

		patterns = strings.Split(line, ", ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne, partTwo := solveTowelDesigns(designs, patterns)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

type trie struct {
	end  bool
	next [26]*trie
	memo map[string]int
}

func solveTowelDesigns(designs []string, patterns []string) (int, int) {
	root := &trie{memo: map[string]int{}}
	for _, pt := range patterns {
		node := root
		for _, r := range pt {
			i := int(r - 'a')
			if node.next[i] == nil {
				node.next[i] = &trie{memo: map[string]int{}}
			}

			node = node.next[i]
		}

		node.end = true
	}

	var count, ways int

	for _, des := range designs {
		if w := isPossibleDesign(des, root, root); w > 0 {
			ways += w
			count++
		}
	}

	return count, ways
}

func isPossibleDesign(design string, node *trie, root *trie) int {
	if node == nil {
		return 0
	}

	if v, ok := node.memo[design]; ok {
		return v
	}

	if design == "" {
		var c int
		if node.end {
			c = 1
		}

		return c
	}

	i := int(design[0] - 'a')

	var c int
	c += isPossibleDesign(design[1:], node.next[i], root)

	if node.end {
		c += isPossibleDesign(design[1:], root.next[i], root)
	}

	node.memo[design] = c

	return c
}
