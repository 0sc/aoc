package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day18"] = day18
}

type memHeapItem struct {
	x, y, s int
	p       [][2]int
}
type memHeap []memHeapItem // x, y, steps

func (h memHeap) Len() int           { return len(h) }
func (h memHeap) Less(i, j int) bool { return h[i].s < h[j].s }
func (h memHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *memHeap) Push(x any) { *h = append(*h, x.(memHeapItem)) }

func (h *memHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type memHeap2 struct{ memHeap }

func (h memHeap2) Less(i, j int) bool { return h.memHeap[i].s > h.memHeap[j].s }

func day18(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var coords [][2]int

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		xy := strings.Split(scanner.Text(), ",")
		x, y := toInt(xy[0]), toInt(xy[1])
		coords = append(coords, [2]int{x, y})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// mx, my, p1 := 7, 7, 12
	mx, my, p1 := 71, 71, 1024
	grid := make([][]byte, my)
	for y := range grid {
		grid[y] = make([]byte, mx)
	}

	var i int
	for ; i < p1; i++ {
		x, y := coords[i][0], coords[i][1]
		grid[y][x] = '#'
	}

	partOne, mh := stepsToExit(grid, &memHeap{{}})

	var partTwo string
	var s int
	for ; i < len(coords); i++ {
		x, y := coords[i][0], coords[i][1]
		grid[y][x] = '#'

		mh = ensureValidHeapItems(grid, mh)
		s, mh = stepsToExit(grid, mh)
		if s == -1 {
			partTwo = fmt.Sprintf("%d,%d", x, y)
			break
		}
	}

	fmt.Printf("Part1: %d, Part2: %s\n", partOne, partTwo)
}

func stepsToExit(grid [][]byte, mh heap.Interface) (int, heap.Interface) {
	my, mx := len(grid), len(grid[0])

	dirs := [...][2]int{
		{0, 1},
		{-1, 0},
		{0, -1},
		{1, 0},
	}

	sh := &memHeap2{}
	visited := map[[2]int]bool{}

	for mh.Len() > 0 {
		mhi := heap.Pop(mh).(memHeapItem)

		k := [2]int{mhi.x, mhi.y}
		if visited[k] {
			heap.Push(sh, mhi)
			continue
		}

		visited[k] = true

		if mhi.x == mx-1 && mhi.y == my-1 {
			heap.Push(sh, mhi)
			mvHeapItems(mh, sh)

			return mhi.s, sh
		}

		for _, dir := range dirs {
			nx, ny := mhi.x+dir[0], mhi.y+dir[1]
			if outOfBounds(ny, nx, my, mx) || grid[ny][nx] == '#' {
				continue
			}

			nmhi := memHeapItem{x: nx, y: ny, s: mhi.s + 1}
			nmhi.p = append(nmhi.p, mhi.p...)
			nmhi.p = append(nmhi.p, [2]int{mhi.x, mhi.y})

			heap.Push(mh, nmhi)
		}
	}

	return -1, mh
}

func ensureValidHeapItems(grid [][]byte, mh heap.Interface) heap.Interface {
	vh := &memHeap2{}

	for mh.Len() > 0 {
		mhi := heap.Pop(mh).(memHeapItem)

		var invalid bool
		for _, pos := range mhi.p {
			x, y := pos[0], pos[1]
			if grid[y][x] == '#' {
				invalid = true
				break
			}
		}

		if !invalid {
			heap.Push(vh, mhi)
		}

	}

	return vh
}

func mvHeapItems(from heap.Interface, to heap.Interface) {
	for from.Len() > 0 {
		heap.Push(to, heap.Pop(from))
	}
}
