package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day25"] = day25
}

func day25(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	schematic := make([]string, 0, 7)
	var locks, keys [][5]int
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			schematic = append(schematic, scanner.Text())
			continue
		}

		kind, loh := convertSchematic(schematic)
		if kind == "lock" {
			locks = append(locks, loh)
		} else {
			keys = append(keys, loh)
		}

		schematic = make([]string, 0, 7)
	}

	kind, loh := convertSchematic(schematic)
	if kind == "lock" {
		locks = append(locks, loh)
	} else {
		keys = append(keys, loh)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne = countKeyLockFits(locks, keys)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func countKeyLockFits(locks [][5]int, keys [][5]int) int {
	var count int

	for _, lk := range locks {
		for _, k := range keys {
			var overlap bool
			for i := 0; i < 5; i++ {
				if lk[i]+k[i] > 5 {
					overlap = true
					break
				}
			}

			if !overlap {
				count++
			}
		}
	}

	return count
}

func convertSchematic(sch []string) (string, [5]int) {
	kind := "key"
	if strings.ContainsAny(sch[0], "#") {
		kind = "lock"
	}

	var height [5]int
	for i := 1; i < 6; i++ {
		for j := range sch[i] {
			if sch[i][j] == '#' {
				height[j]++
			}
		}
	}

	return kind, height
}
