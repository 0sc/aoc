package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day15"] = day15
}

func day15(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var grid [][]byte
	var moves string
	var isMoves bool
	var start [2]int

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			isMoves = true
			continue
		}

		if isMoves {
			moves += line
			continue
		}

		if c := strings.Index(line, "@"); c != -1 {
			start = [2]int{len(grid), c}
		}

		row := []byte(line)
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var start2 [2]int
	grid2 := make([][]byte, len(grid))
	for i, r := range grid {
		grid2[i] = make([]byte, 0, len(r)*2)
		for _, ch := range r {
			if ch == 'O' {
				grid2[i] = append(grid2[i], '[', ']')
				continue
			}

			if ch == '@' {
				start2 = [2]int{i, len(grid2[i])}
				grid2[i] = append(grid2[i], '@', '.')
				continue
			}

			grid2[i] = append(grid2[i], ch, ch)
		}
	}

	dirs := map[rune][2]int{
		'^': {-1, 0},
		'>': {0, 1},
		'v': {1, 0},
		'<': {0, -1},
	}

	var mvd bool
	for _, mv := range moves {
		dir := dirs[mv]
		start, _ = moveOne(start, dir, grid)

		g2 := cloneGrid(grid2)
		start2, mvd = moveTwo(start2, dir, grid2)
		if !mvd {
			grid2 = g2
		}
	}

	partOne := sumGPS(grid)
	partTwo := sumGPS(grid2)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func moveOne(from [2]int, by [2]int, grid [][]byte) ([2]int, bool) {
	r, c := from[0], from[1]

	nr, nc := r+by[0], c+by[1]
	if grid[nr][nc] == '#' {
		return from, false
	}

	if grid[nr][nc] == 'O' {
		// ask occupant to move
		_, mvd := moveOne([2]int{nr, nc}, by, grid)
		if !mvd {
			return from, false
		}
	}

	grid[nr][nc] = grid[r][c]
	grid[r][c] = '.'

	return [2]int{nr, nc}, true
}

func moveTwo(from [2]int, by [2]int, grid [][]byte) ([2]int, bool) {
	r, c := from[0], from[1]
	if grid[r][c] == '.' {
		return from, true
	}

	if grid[r][c] == '#' {
		return from, false
	}

	nr, nc := r+by[0], c+by[1]
	_, mvd := moveTwo([2]int{nr, nc}, by, grid)
	if !mvd {
		return from, false
	}

	// move if not vertical movement
	if by[0] == 0 || grid[r][c] == '@' {
		grid[nr][nc] = grid[r][c]
		grid[r][c] = '.'

		return [2]int{nr, nc}, true
	}

	// try move other component
	oc := c + 1
	if grid[r][c] == ']' {
		oc = c - 1
	}

	_, mvd = moveTwo([2]int{nr, oc}, by, grid)
	if !mvd {
		return from, false
	}

	// both parts can move
	grid[nr][nc] = grid[r][c]
	grid[r][c] = '.'

	grid[nr][oc] = grid[r][oc]
	grid[r][oc] = '.'

	return [2]int{nr, nc}, true
}

func cloneGrid[T any](grid [][]T) [][]T {
	cp := make([][]T, len(grid))
	for i, r := range grid {
		cp[i] = make([]T, len(r))
		copy(cp[i], r)
	}

	return cp
}

func sumGPS(grid [][]byte) int {
	var sum int
	for dt, r := range grid {
		for dl, ch := range r {
			if ch == 'O' || ch == '[' {
				sum += dt*100 + dl
			}
		}
	}

	return sum
}

func drawGrid(grid [][]byte) {
	for _, r := range grid {
		fmt.Println(string(r))
	}
	fmt.Println()
}
