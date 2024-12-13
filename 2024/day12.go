package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func init() {
	Solutions["day12"] = day12
}

type unionFind struct {
	region []int
	size   []int
	count  int
}

func newUnionFind(count int) *unionFind {
	region := make([]int, count)
	size := make([]int, count)
	for i := 0; i < count; i++ {
		region[i] = i
		size[i] = 1
	}
	return &unionFind{
		region: region,
		size:   size,
		count:  count,
	}
}

func (uf *unionFind) find(plt int) int {
	for plt != uf.region[plt] {
		uf.region[plt] = uf.region[uf.region[plt]]
		plt = uf.region[plt]
	}

	return plt
}

func (uf *unionFind) union(plt1, plt2 int) {
	region1 := uf.find(plt1)
	region2 := uf.find(plt2)

	if region1 == region2 {
		return
	}

	if uf.size[region1] < uf.size[region2] {
		region1, region2 = region2, region1
	}

	uf.region[region2] = region1
	uf.size[region1] += uf.size[region2]

	uf.count -= 1
}

func day12(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var garden []string
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		garden = append(garden, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne, partTwo = fencingPrice(garden)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func fencingPrice(garden []string) (int, int) {
	lr := len(garden)
	lc := len(garden[0])
	oob := func(r, c int) bool { return outOfBounds(r, c, lr, lc) }

	dirs := [...][2]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}

	uf := newUnionFind(lr * lc)
	visited := make([]bool, lr*lc)
	perimeter := make([]int, lr*lc)
	key := func(r, c int) int { return r*lc + c }

	visited[0] = true
	q := make([][2]int, 0, lr*lc) // row, col, region
	q = append(q, [2]int{0, 0})
	for i := 0; i < len(q); i++ {
		r, c := q[i][0], q[i][1]
		k := key(r, c)

		for _, dir := range dirs {
			nr, nc := r+dir[0], c+dir[1]
			if oob(nr, nc) {
				perimeter[k]++
				continue
			}

			nk := key(nr, nc)
			if garden[r][c] == garden[nr][nc] {
				uf.union(k, nk)
			} else {
				perimeter[k]++
			}

			if visited[nk] {
				continue
			}

			visited[nk] = true
			q = append(q, [2]int{nr, nc})
		}
	}

	regs := map[int][][3]int{}
	for r := range garden {
		for c := range garden[r] {
			k := key(r, c)
			reg := uf.find(k)
			regs[reg] = append(regs[reg], [3]int{r, c, k})
		}
	}

	var cto = map[[2]int][][2]int{ // magic
		{1, 1}:   {{1, 0}, {0, 1}},
		{-1, 1}:  {{-1, 0}, {0, 1}},
		{1, -1}:  {{1, 0}, {0, -1}},
		{-1, -1}: {{-1, 0}, {0, -1}},
	}

	match := func(r, c, rr, cc int) bool {
		if oob(r, c) && oob(rr, cc) {
			return true
		}

		if !oob(r, c) && !oob(rr, cc) {
			return garden[r][c] == garden[rr][cc]
		}

		return false
	}

	var p1Price, p2Price int
	for reg, plts := range regs {
		var sides, pmt int

		for _, plt := range plts {
			r, c, k := plt[0], plt[1], plt[2]
			pmt += perimeter[k]

			for co, pa := range cto {
				nr := [2]int{r + co[0], c + co[1]}
				i1 := [2]int{r + pa[0][0], c + pa[0][1]}
				i2 := [2]int{r + pa[1][0], c + pa[1][1]}

				if !match(i1[0], i1[1], r, c) && !match(i2[0], i2[1], r, c) {
					sides++
				}

				if match(i1[0], i1[1], r, c) && match(i2[0], i2[1], r, c) && !match(nr[0], nr[1], r, c) {
					sides++
				}

			}
		}

		p1Price += uf.size[reg] * pmt
		p2Price += sides * uf.size[reg]
	}

	return p1Price, p2Price
}
