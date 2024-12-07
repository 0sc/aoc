package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func init() {
	Solutions["day4"] = day4
}

func day4(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	var xs, as [][2]int
	var lines []string
	var row int
	for scanner.Scan() {
		line := scanner.Text()
		for i := range line {
			if line[i] == 'X' {
				xs = append(xs, [2]int{row, i})
				continue
			}

			if line[i] == 'A' {
				as = append(as, [2]int{row, i})
			}
		}

		lines = append(lines, line)
		row++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, x := range xs {
		partOne += countXmasFrom(x[0], x[1], lines)
	}

	for _, a := range as {
		partTwo += countXMasFrom(a[0], a[1], lines)
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func countXmasFrom(r, c int, lines []string) int {
	dirs := [...][2]int{
		// horizontal
		{0, 1},  // l -> r
		{0, -1}, // l <- r

		// vertical
		{1, 0},  // t -> b
		{-1, 0}, // t <- b

		// diagonal
		{1, 1},   // tl -> br
		{-1, -1}, // tl <- br
		{-1, 1},  // bl -> tr
		{1, -1},  // bl <- tr
	}

	var count int
	var x, m, a, s [2]int
	for _, dir := range dirs {
		x = [...]int{r, c}
		m = [...]int{r + dir[0], c + dir[1]}
		a = [...]int{r + dir[0]*2, c + dir[1]*2}
		s = [...]int{r + dir[0]*3, c + dir[1]*3}

		if isXmas(x, m, a, s, lines) {
			count++
		}
	}

	return count
}

func isXmas(x, m, a, s [2]int, lines []string) bool {
	mr, mc := len(lines), len(lines[0])
	if outOfBounds(x[0], x[1], mr, mc) || lines[x[0]][x[1]] != 'X' {
		return false
	}

	if outOfBounds(m[0], m[1], mr, mc) || lines[m[0]][m[1]] != 'M' {
		return false
	}

	if outOfBounds(a[0], a[1], mr, mc) || lines[a[0]][a[1]] != 'A' {
		return false
	}

	if outOfBounds(s[0], s[1], mr, mc) || lines[s[0]][s[1]] != 'S' {
		return false
	}

	return true
}

func countXMasFrom(r, c int, lines []string) int {
	tl := [2]int{r - 1, c - 1}
	tr := [2]int{r - 1, c + 1}
	bl := [2]int{r + 1, c - 1}
	br := [2]int{r + 1, c + 1}
	p := [2]int{r, c}

	if isMas(tl, p, br, lines) && isMas(tr, p, bl, lines) {
		return 1
	}

	return 0
}

func isMas(m, a, s [2]int, lines []string) bool {
	mr, mc := len(lines), len(lines[0])
	if outOfBounds(m[0], m[1], mr, mc) || outOfBounds(a[0], a[1], mr, mc) || outOfBounds(s[0], s[1], mr, mc) {
		return false
	}

	mv := lines[m[0]][m[1]]
	av := lines[a[0]][a[1]]
	sv := lines[s[0]][s[1]]

	return (mv == 'M' && av == 'A' && sv == 'S') || (mv == 'S' && av == 'A' && sv == 'M')
}

func outOfBounds(row, col, maxRow, maxCol int) bool {
	return row < 0 || row >= maxRow || col < 0 || col >= maxCol
}
