package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	Solutions["day11"] = day11
}

func day11(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var line string
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne, partTwo := blink(line)

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func blink(str string) (int, int) {
	dedup := map[int]int{}
	for _, s := range strings.Split(str, " ") {
		dedup[toInt(s)] = 1
	}

	count := func() int {
		var qty int
		for _, q := range dedup {
			qty += q
		}

		return qty
	}

	var at25 int

	for i := 0; i < 75; i++ {
		if i == 25 {
			at25 = count()
		}

		ndedup := map[int]int{}

		for st, qty := range dedup {
			if st == 0 {
				ndedup[1] += qty
				continue
			}

			s := fmt.Sprint(st)
			if l := len(s); l&1 == 0 {
				m := l / 2
				ndedup[toInt(s[:m])] += qty
				ndedup[toInt(s[m:])] += qty

				continue
			}

			ndedup[st*2024] += qty
		}

		dedup = ndedup
	}

	at75 := count()

	return at25, at75
}
