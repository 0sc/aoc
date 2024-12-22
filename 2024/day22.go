package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func init() {
	Solutions["day22"] = day22
}

func day22(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var seeds []int
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		seeds = append(seeds, toInt(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tally := map[[4]int]int{}
	for _, seed := range seeds {
		partOne += genSecretNum(seed, 2000, tally)
	}

	for _, v := range tally {
		partTwo = max(partTwo, v)
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func genSecretNum(seed int, iter int, tally map[[4]int]int) int {
	mixIn := func(num int, val int) int { return num ^ val }
	prune := func(num int) int { return num % 16777216 }
	price := func(num int) int { return num % 10 }

	lp := price(seed)

	dp := map[[4]int]int{}
	l4 := make([]int, 0, 4)
	track := func(p int) {
		l4 = append(l4, p-lp)
		if len(l4) < 4 {
			return
		}

		k := [4]int{l4[0], l4[1], l4[2], l4[3]}
		if _, ok := dp[k]; !ok {
			dp[k] = p
			tally[k] += p
		}

		l4 = l4[1:]
	}

	for i := 0; i < iter; i++ {
		seed = prune(mixIn(seed, seed*64))
		seed = prune(mixIn(seed, seed/32))
		seed = prune(mixIn(seed, seed*2048))

		p := price(seed)
		track(p)
		lp = p
	}

	return seed
}
