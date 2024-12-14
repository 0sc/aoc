package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func init() {
	Solutions["day14"] = day14
}

var roboRegx = regexp.MustCompile("p=(.+),(.+) v=(.+),(.+)")

func day14(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	type robot struct{ p, v [2]int }
	var robots []robot

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		match := roboRegx.FindAllStringSubmatch(scanner.Text(), -1)
		robots = append(robots, robot{
			p: [2]int{toInt(match[0][1]), toInt(match[0][2])},
			v: [2]int{toInt(match[0][3]), toInt(match[0][4])},
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// X := 11
	// Y := 7
	X := 101
	Y := 103

	move := func(rb robot, t int) [2]int {
		x := (rb.v[0]*t + rb.p[0]) % X
		if x < 0 {
			x += X
		}

		y := (rb.v[1]*t + rb.p[1]) % Y
		if y < 0 {
			y += Y
		}

		return [2]int{x, y}
	}

	mX := X / 2
	mY := Y / 2

	var q1, q2, q3, q4 int
	for _, rb := range robots {
		xy := move(rb, 100)
		switch {
		case xy[0] < mX && xy[1] < mY:
			q1++
		case xy[0] > mX && xy[1] < mY:
			q2++
		case xy[0] < mX && xy[1] > mY:
			q3++
		case xy[0] > mX && xy[1] > mY:
			q4++
		}
	}

	partOne := q1 * q2 * q3 * q4

	var mu, partTwo int
	for i := 1000; i < 10000; i++ {
		u := map[[2]int]int{}
		for _, rb := range robots {
			pt := move(rb, i)
			u[pt] += 1
		}

		if us := len(u); us > mu { // partially safe heuristic
			drawRobots(u, X, Y)
			fmt.Println()

			mu = us
			partTwo = i
		}
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func uniquenessScore(pts [][2]int) int {
	dedup := map[[2]int]int{}
	for _, pt := range pts {
		dedup[pt] += 1
	}

	return len(dedup)
}

func drawRobots(pts map[[2]int]int, lx int, ly int) {
	grid := make([][]byte, ly)
	for pt := range pts {
		x, y := pt[0], pt[1]
		if len(grid[y]) == 0 {
			grid[y] = make([]byte, lx)
		}

		grid[y][x] = '#'
	}

	for _, r := range grid {
		if len(r) == 0 {
			r = make([]byte, lx)
		}

		for _, c := range r {
			if c != '#' {
				c = '.'
			}

			fmt.Print(string(c))
		}
		fmt.Println()
	}
}
