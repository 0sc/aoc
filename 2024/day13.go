package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func init() {
	Solutions["day13"] = day13
}

var macRegx = regexp.MustCompile("(\\d+)\\D+(\\d+)")

func day13(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var machines [][][2]int
	mac := make([][2]int, 0, 3)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			machines = append(machines, mac)
			mac = make([][2]int, 0, 3)
			continue
		}

		match := macRegx.FindAllStringSubmatch(line, -1)
		x := toInt(match[0][1])
		y := toInt(match[0][2])
		mac = append(mac, [2]int{x, y})
	}

	machines = append(machines, mac)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	inc := 10000000000000
	for _, mac = range machines {
		partOne += fewestTokenToPrice(mac[0], mac[1], mac[2])

		m2 := [2]int{mac[2][0] + inc, mac[2][1] + inc}
		partTwo += fewestTokenToPrice(mac[0], mac[1], m2)
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

// Xx1 + Yy1 = z1
// Xx2 + Yy2 = z2
func fewestTokenToPrice(a [2]int, b [2]int, p [2]int) int {
	x1, y1 := a[0], b[0]
	x2, y2 := a[1], b[1]
	z1, z2 := p[0], p[1]

	y := ((x1 * z2) - (x2 * z1)) / ((x1 * y2) - (x2 * y1))
	x := (z1 - (y * y1)) / x1

	if (x1*x+(y1*y)) == z1 && (x2*x+(y2*y)) == z2 {
		return x*3 + y
	}

	return 0
}
