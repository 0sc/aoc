package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day2"] = day2
}

func day2(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lvls := strings.Split(scanner.Text(), " ")
		if isSafe(lvls) {
			partOne++
			partTwo++
			continue
		}

		if isSafeWithProblemDampener(lvls) {
			partTwo++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func isSafe(lvls []string) bool {
	var sign int

	for i := 1; i < len(lvls); i++ {
		curr, prev := toInt(lvls[i]), toInt(lvls[i-1])
		sn := cmp.Compare(curr, prev)

		if sign == 0 {
			sign = sn
		}

		if sn == 0 || sign != sn {
			return false
		}

		diff := abs(curr - prev)
		if diff > 3 {
			return false
		}
	}

	return true
}

func isSafeWithProblemDampener(lvls []string) bool {
	for i := 0; i < len(lvls); i++ {
		s := append([]string{}, lvls[0:i]...)
		s = append(s, lvls[i+1:]...)

		if isSafe(s) {
			return true
		}
	}

	return false
}
