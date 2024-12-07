package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func init() {
	Solutions["day3"] = day3
}

var mulRegx = regexp.MustCompile("(mul\\(\\d{1,3},\\d{1,3}\\)|do\\(\\)|don't\\(\\))")

func day3(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var skip bool
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		matches := mulRegx.FindAllString(line, -1)

		for _, op := range matches {
			if strings.HasPrefix(op, "don") {
				skip = true
				continue
			}

			if strings.HasPrefix(op, "do") {
				skip = false
				continue
			}

			mul := doMul(op)
			partOne += mul

			if !skip {
				partTwo += mul
			}

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func doMul(str string) int {
	substr := str[4 : len(str)-1]
	ops := strings.Split(substr, ",")

	return toInt(ops[0]) * toInt(ops[1])
}
