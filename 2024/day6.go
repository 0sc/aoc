package main

import (
	"bufio"
	"container/ring"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day6"] = day6
}

func day6(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var grid [][]byte
	var sr, sc int

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []byte(line))
		if c := strings.Index(line, "^"); c != -1 {
			sr = len(grid) - 1
			sc = c
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne, partTwo := simulateDay6(sr, sc, grid)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func simulateDay6(row, col int, grid [][]byte) (int, int) {
	steps, _, visited := traverse(row, col, grid)

	var cycles int
	st := [2]int{row, col}
	for pos := range visited {
		if pos == st {
			continue
		}

		grid[pos[0]][pos[1]] = '#'
		_, cycle, _ := traverse(row, col, grid)
		if cycle {
			cycles++
		}

		grid[pos[0]][pos[1]] = '.'
	}

	return steps, cycles
}

var (
	dirs = [...][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	reps = [...]byte{'N', 'E', 'S', 'W'}
)

func traverse(row, col int, grid [][]byte) (int, bool, map[[2]int]byte) {
	iter := ring.New(4)
	for i := 0; i < 4; i++ {
		iter.Value = i
		iter = iter.Next()
	}

	iv := func() int { return iter.Value.(int) }
	dir := dirs[iv()]
	rep := reps[iv()]
	move := func(inc int) {
		row += dir[0] * inc
		col += dir[1] * inc
	}

	visited := map[[2]int]byte{{row, col}: rep}
	steps := 1
	move(1)
	for row >= 0 && row < len(grid) && col >= 0 && col < len(grid[0]) {
		k := [2]int{row, col}
		if grid[row][col] == '#' {
			// time to rotate
			move(-1)
			iter = iter.Next()

			dir = dirs[iv()]
			rep = reps[iv()]
		} else if _, ok := visited[k]; !ok {
			steps++
			visited[[2]int{row, col}] = rep
		} else if visited[k] == rep {
			return steps, true, visited
		}

		move(1)
	}

	return steps, false, visited
}
