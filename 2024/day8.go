package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func init() {
	Solutions["day8"] = day8
}

func day8(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var maxRow, maxCol int
	freqs := map[rune][][2]int{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		maxCol = len(line)
		for col, b := range line {
			if b != '.' {
				freqs[b] = append(freqs[b], [2]int{maxRow, col})
			}
		}

		maxRow++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	oob := func(pos [2]int) bool {
		r, c := pos[0], pos[1]
		return r < 0 || r >= maxRow || c < 0 || c >= maxCol
	}

	partOne, partTwo := uniqueAntiNodeLocations(freqs, oob)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func uniqueAntiNodeLocations(freqs map[rune][][2]int, oob func([2]int) bool) (int, int) {
	locs := map[[2]int]bool{}
	locs2 := map[[2]int]bool{}

	for _, grp := range freqs {
		// match each group with each other
		l := len(grp) - 1
		for i := 0; i < l; i++ {
			for j := i + 1; j <= l; j++ {
				l1, l2 := antiNodeLocations(grp[i], grp[j])

				if !oob(l1) {
					locs[l1] = true
				}

				if !oob(l2) {
					locs[l2] = true
				}

				lhs := antiNodeLocationsWithHarmonics(grp[i], grp[j], oob)
				for _, l := range lhs {
					locs2[l] = true
				}
			}
		}
	}

	return len(locs), len(locs2)
}

func antiNodeLocations(a [2]int, b [2]int) ([2]int, [2]int) {
	dy := b[0] - a[0]
	dx := b[1] - a[1]

	la := [2]int{a[0] - dy, a[1] - dx}
	lb := [2]int{b[0] + dy, b[1] + dx}

	return la, lb
}

func antiNodeLocationsWithHarmonics(a [2]int, b [2]int, oob func([2]int) bool) [][2]int {
	dy := b[0] - a[0]
	dx := b[1] - a[1]

	locs := [][2]int{}
	// assumes a comes before b in the grid
	n := a
	for !oob(n) {
		locs = append(locs, n)

		n[0] -= dy
		n[1] -= dx
	}

	n = b
	for !oob(n) {
		locs = append(locs, n)

		n[0] += dy
		n[1] += dx
	}

	return locs
}
