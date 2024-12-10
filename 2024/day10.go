package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func init() {
	Solutions["day10"] = day10
}

func day10(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var grid []string
	var ths [][2]int

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		for c := range line {
			if line[c] == '0' {
				ths = append(ths, [2]int{len(grid), c})
			}
		}

		grid = append(grid, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne, partTwo := scoreAndRateTrailheads(ths, grid)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func scoreAndRateTrailheads(ths [][2]int, grid []string) (int, int) {
	dirs := [...][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	dedup := make([]map[[2]int]int, len(ths))
	q := [][3]int{}
	for i, th := range ths {
		q = append(q, [3]int{th[0], th[1], i})
	}

	var rating int
	for i := 0; i < len(q); i++ {
		r, c, th := q[i][0], q[i][1], q[i][2]
		v := grid[r][c]
		if v == '9' {
			if dedup[th] == nil {
				dedup[th] = map[[2]int]int{}
			}

			dedup[th][[2]int{r, c}]++

			rating++
			continue
		}

		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if outOfBounds(nr, nc, len(grid), len(grid[0])) {
				continue
			}

			nv := grid[nr][nc]
			if v+1 == nv {
				q = append(q, [3]int{nr, nc, th})
			}
		}

	}

	var score int
	for _, th := range dedup {
		score += len(th)
	}

	return score, rating
}
