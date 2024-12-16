package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func init() {
	Solutions["day16"] = day16
}

type mazePos struct {
	r, c, s, d int
	path       [][2]int
}
type mazeHeap []mazePos

func (h mazeHeap) Len() int           { return len(h) }
func (h mazeHeap) Less(i, j int) bool { return h[i].s < h[j].s }
func (h mazeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *mazeHeap) Push(x any) { *h = append(*h, x.(mazePos)) }

func (h *mazeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func day16(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var start [2]int
	var grid []string

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		if c := strings.Index(line, "S"); c != -1 {
			start = [2]int{len(grid), c}
		}

		grid = append(grid, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne, partTwo := findLowestScore(start, grid)
	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func findLowestScore(s [2]int, grid []string) (int, int) {
	mr, mc := len(grid), len(grid[0])
	dirs := [...][2]int{
		{0, 1},
		{-1, 0},
		{0, -1},
		{1, 0},
	}

	isTurn := func(f, t [2]int) bool {
		return (f[0] == 0 && t[1] == 0) || (f[1] == 0 && t[0] == 0)
	}

	mh := &mazeHeap{
		{
			r:    s[0],
			c:    s[1],
			path: [][2]int{{s[0], s[1]}},
		},
	}

	visited := map[[3]int]int{{s[0], s[1]}: 0}
	var bests []mazePos
	ms := math.MaxInt
	for mh.Len() > 0 {
		pos := heap.Pop(mh).(mazePos)
		if pos.s > ms {
			break
		}

		k := [3]int{pos.r, pos.c, pos.d}
		if vs, ok := visited[k]; ok && pos.s > vs {
			continue
		}

		visited[k] = pos.s

		if grid[pos.r][pos.c] == 'E' {
			ms = pos.s
			bests = append(bests, pos)
			continue
		}

		pd := dirs[pos.d]
		for di, dir := range dirs {
			nr, nc := pos.r+dir[0], pos.c+dir[1]

			if outOfBounds(nr, nc, mr, mc) || grid[nr][nc] == '#' {
				continue
			}
			s := pos.s + 1
			if isTurn(pd, dir) {
				s += 1000
			}

			mp := mazePos{r: nr, c: nc, d: di, s: s}
			mp.path = append(mp.path, pos.path...)
			mp.path = append(mp.path, [2]int{nr, nc})

			heap.Push(mh, mp)
		}
	}

	dedup := map[[2]int]bool{}
	for _, pos := range bests {
		for _, k := range pos.path {
			dedup[k] = true
		}
	}

	// plotPath(grid, dedup)

	return bests[0].s, len(dedup)
}

func plotPath(grid []string, path map[[2]int]bool) {
	plot := make([][]byte, 0, len(grid))
	for _, r := range grid {
		plot = append(plot, []byte(r))
	}

	for p := range path {
		r, c := p[0], p[1]
		plot[r][c] = '0'
	}

	for _, r := range plot {
		fmt.Println(string(r))
	}

	fmt.Println()
}
