package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day20"] = day20
}

func day20(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var grid []string
	var start [2]int
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if c := strings.Index(line, "S"); c > -1 {
			start = [2]int{len(grid), c}
		}

		grid = append(grid, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne := countRaceCheatSaves(2, 100, start, grid)
	partTwo := countRaceCheatSaves(20, 100, start, grid)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func countRaceCheatSaves(chts int, ms int, s [2]int, grid []string) int {
	dirs := [...][2]int{
		{-1, 0},
		{0, 1},
		{0, -1},
		{1, 0},
	}

	rev := func(d int) [2]int {
		dir := dirs[d]
		return [2]int{dir[0] * -1, dir[1] * -1}
	}

	p := struct{ r, c, d, t int }{r: s[0], c: s[1]}
	np := [][2]int{{p.r, p.c}}

	for grid[p.r][p.c] != 'E' {
		bd := rev(p.d)

		for d, dir := range dirs {
			if dir == bd {
				continue
			}

			nr, nc := p.r+dir[0], p.c+dir[1]
			if grid[nr][nc] == '#' {
				continue
			}

			p.r, p.c, p.d, p.t = nr, nc, d, p.t+1

			np = append(np, [2]int{p.r, p.c})

			break
		}
	}

	var saves int

	mi := len(np) - ms
	for i, p := range np[:mi] {
		for j, q := range np[i+ms:] {
			a := [2]int{p[0], p[1]}
			b := [2]int{q[0], q[1]}
			dst := manhattanDistance(a, b)

			if dst <= chts && dst <= j {
				saves++
			}
		}
	}

	return saves
}

func manhattanDistance(a [2]int, b [2]int) int {
	return abs(a[0]-b[0]) + abs(a[1]-b[1])
}
