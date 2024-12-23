package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func init() {
	Solutions["day23"] = day23
}

func day23(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	net := map[string][]string{}
	ts := map[string]bool{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		cs := strings.Split(scanner.Text(), "-")
		c1, c2 := cs[0], cs[1]

		net[c1] = append(net[c1], c2)
		net[c2] = append(net[c2], c1)

		if c1[0] == 't' {
			ts[c1] = true
		}

		if c2[0] == 't' {
			ts[c2] = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	dedup := map[string]bool{}

	for c1 := range ts {
		for _, c2 := range net[c1] {
			for _, c3 := range net[c2] {
				if slices.Contains(net[c1], c3) {
					grp := make([]string, 0, 3)
					grp = append(grp, c1, c2, c3)
					slices.Sort(grp)
					dedup[strings.Join(grp, ",")] = true
				}
			}
		}
	}

	partOne := len(dedup)
	cmps := make([]string, 0, len(net))
	for cmp := range net {
		cmps = append(cmps, cmp)
	}

	cliques := bronKerbosch1([]string{}, cmps, []string{}, net)
	lc := slices.MaxFunc(cliques, func(a []string, b []string) int { return cmp.Compare(len(a), len(b)) })

	slices.Sort(lc)
	partTwo := strings.Join(lc, ",")

	fmt.Printf("Part1: %d, Part2: %s\n", partOne, partTwo)
}

// https://towardsdatascience.com/graphs-paths-bron-kerbosch-maximal-cliques-e6cab843bc2c
// R := is the set of nodes of a maximal clique.
// P := is the set of possible nodes in a maximal clique.
// X := is the set of nodes that are excluded.
func bronKerbosch1(
	r []string,
	p []string,
	x []string,
	net map[string][]string,
) [][]string {
	if len(p) == 0 && len(x) == 0 {
		return [][]string{r}
	}

	var cliques [][]string
	for i := 0; i < len(p); i++ {
		v := p[i]
		// R ⋃ {v}
		nr := append([]string{v}, r...)

		var np, nx []string
		for _, n := range net[v] {
			// P ⋂ N(v)
			if slices.Contains(p, n) {
				np = append(np, n)
			}

			// X ⋂ N(v)
			if slices.Contains(x, n) {
				nx = append(nx, n)
			}
		}

		clique := bronKerbosch1(nr, np, nx, net)
		cliques = append(cliques, clique...)

		// P := P \ {v}
		p = slices.Delete(p, i, i+1)
		i--

		// X := X ⋃ {v}
		x = append(x, v)
	}

	return cliques
}
