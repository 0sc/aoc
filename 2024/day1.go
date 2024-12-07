package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func init() {
	Solutions["day1"] = day1
}

const SEP = "   "

func day1(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var left, right []int
	freq := map[int]int{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		lists := strings.Split(line, SEP)

		left = append(left, toInt(lists[0]))

		b := toInt(lists[1])
		right = append(right, b)
		freq[b]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	slices.Sort(left)
	slices.Sort(right)

	for i := range left {
		partOne += abs(left[i] - right[i])
		partTwo += left[i] * freq[left[i]]
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}
