package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
)

func init() {
	Solutions["day21"] = day21
}

var (
	numPad = map[rune]map[rune]rune{
		'0': {'2': '^', 'A': '>'},
		'1': {'2': '>', '4': '^'},
		'2': {'0': 'v', '1': '<', '3': '>', '5': '^'},
		'3': {'A': 'v', '2': '<', '6': '^'},
		'4': {'1': 'v', '5': '>', '7': '^'},
		'5': {'2': 'v', '4': '<', '6': '>', '8': '^'},
		'6': {'3': 'v', '5': '<', '9': '^'},
		'7': {'4': 'v', '8': '>'},
		'8': {'5': 'v', '7': '<', '9': '>'},
		'9': {'6': 'v', '8': '<'},
		'A': {'0': '<', '3': '^'},
	}

	dirPad = map[rune]map[rune]rune{
		'^': {'A': '>', 'v': 'v'},
		'A': {'^': '<', '>': 'v'},
		'>': {'A': '^', 'v': '<'},
		'v': {'<': '<', '^': '^', '>': '>'},
		'<': {'v': '>'},
	}
)

func day21(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var codes []string
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	memo := make([]map[string]int, 26)
	for _, code := range codes {
		sl := shortestSeq(code, 2, numPad, memo)
		partOne += sl * toInt(code[:3])

		sl = shortestSeq(code, 25, numPad, memo)
		partTwo += sl * toInt(code[:3])
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func shortestSeq(
	code string,
	actor int,
	pad map[rune]map[rune]rune,
	memo []map[string]int,
) int {
	if memo[actor] == nil {
		memo[actor] = map[string]int{}
	}

	if m, ok := memo[actor][code]; ok {
		return m
	}

	var sl int
	fr := 'A'
	for _, r := range code {
		seqs := allSeqs(fr, r, pad)
		fr = r

		if actor == 0 {
			slices.SortFunc(seqs, func(a, b string) int { return cmp.Compare(len(a), len(b)) })
			sl += len(seqs[0])

			continue
		}

		msl := math.MaxInt
		for _, seq := range seqs {
			l := shortestSeq(seq, actor-1, dirPad, memo)
			msl = min(msl, l)
		}

		sl += msl
	}

	memo[actor][code] = sl
	return sl
}

func allSeqs(from rune, to rune, pad map[rune]map[rune]rune) []string {
	type item struct {
		val, dir rune
		seq      []rune
	}

	var seqs []string
	visited := map[rune]int{}
	q := []item{{val: from}}
	for i := 0; i < len(q); i++ {
		it := q[i]

		if it.val == to {
			it.seq = append(it.seq, 'A')
			seqs = append(seqs, string(it.seq))
			continue
		}

		for k, p := range pad[it.val] {
			nit := item{val: k, seq: append([]rune{}, it.seq...)}
			nit.seq = append(nit.seq, p)

			if p, ok := visited[nit.val]; !ok || len(nit.seq) <= p {
				visited[nit.val] = len(nit.seq)
				q = append(q, nit)
			}
		}
	}

	return seqs
}
