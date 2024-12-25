package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func init() {
	Solutions["day24"] = day24
}

type circuitGate struct {
	w1, op, w2, out string
	ready           bool
}

func day24(input string) {

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	gateRegx := regexp.MustCompile("^(.+) (AND|XOR|OR) (.*) -> (.*)$")
	toUint8 := func(s string) uint8 {
		v, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			panic(err)
		}

		return uint8(v)
	}
	output := map[string]uint8{}
	gates := map[string][]*circuitGate{}
	ready := []*circuitGate{}

	var maxz int
	var isGate bool

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			isGate = true
			continue
		}

		if !isGate {
			p := strings.Split(line, ": ")
			output[p[0]] = toUint8(p[1])
			continue
		}

		match := gateRegx.FindAllStringSubmatch(line, -1)[0]
		w1, op, w2, out := match[1], match[2], match[3], match[4]

		gate := &circuitGate{w1: w1, op: op, w2: w2, out: out}
		gates[w1] = append(gates[w1], gate)
		gates[w2] = append(gates[w2], gate)

		if gateIsReady(gate, output) {
			gate.ready = true
			ready = append(ready, gate)
		}

		if out[0] == 'z' {
			maxz = max(maxz, toInt(out[1:]))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var partOne int
	simulate(gates, ready, output)
	for i := 0; i <= maxz; i++ {
		z := wireLabel("z", i)
		if output[z] != 0 {
			partOne += 1 << i
		}
	}

	var partTwo string
	swaps := repairAdders(maxz, gates) // NOTE: doesn't work on all the example input
	slices.Sort(swaps)
	partTwo = strings.Join(swaps, ",")

	fmt.Printf("Part1: %d, Part2: %s\n", partOne, partTwo)
}

func wireLabel(w string, idx int) string { return fmt.Sprintf("%s%02d", w, idx) }

func simulate(
	gates map[string][]*circuitGate,
	ready []*circuitGate,
	output map[string]uint8,
) {
	type gateLogic func(a, b uint8) uint8

	logic := map[string]gateLogic{
		"AND": func(a, b uint8) uint8 { return a & b },
		"OR":  func(a, b uint8) uint8 { return a | b },
		"XOR": func(a, b uint8) uint8 { return a ^ b },
	}

	for i := 0; i < len(ready); i++ {
		gate := ready[i]

		w1, w2 := output[gate.w1], output[gate.w2]
		output[gate.out] = logic[gate.op](w1, w2)

		out := gate.out
		// queue all circuits dependent on you
		for _, cg := range gates[out] {
			if !cg.ready && gateIsReady(cg, output) {
				cg.ready = true
				ready = append(ready, cg)
			}
		}
	}
}

// https://content.instructables.com/F3D/2GZ2/KNVR5S0C/F3D2GZ2KNVR5S0C.png
func repairAdders(maxb int, gates map[string][]*circuitGate) []string {
	var swaps []string
	swapOut := func(a, b *circuitGate) {
		swaps = append(swaps, a.out, b.out)
		a.out, b.out = b.out, a.out
	}

	var cout string
	for i := 0; i < maxb; i++ {
		xl := wireLabel("x", i)
		yl := wireLabel("y", i)
		zl := wireLabel("z", i)

		// find the gates in the adder
		// the xor gate
		xor := findGate(xl, "XOR", yl, gates)
		// the and gate
		and := findGate(xl, "AND", yl, gates)

		if i == 0 {
			// special case for the first gate
			if xor.out != zl {
				panic(fmt.Sprintf("first gate out should be %s not %s", zl, xor.out))
			}

			cout = and.out
			continue
		}

		// the carry xor
		cxor := findGate(xor.out, "XOR", cout, gates)
		if cxor == nil {
			// xor is outputing to the wrong place
			cxor = findByOther(cout, "XOR", zl, gates)
			eout := cxor.w1
			if eout == cout {
				eout = cxor.w2
			}

			eo := findGateByOut(eout, gates)

			// swap with xor
			swapOut(xor, eo)
		}

		// expected output gate
		zg := findGateByOut(zl, gates)
		if cxor != zg {
			// we have an xor output that doesn't match expected
			swapOut(cxor, zg)
		}

		// the carry and
		cand := findGate(xor.out, "AND", cout, gates)

		// the next carry
		nc := findGate(and.out, "OR", cand.out, gates)
		if nc == nil {
			// and.out is outputting to the wrong place
			nc = findGateByOtherAndOp(cand.out, "OR", gates)
			oth := nc.w1
			if oth == cand.out {
				oth = nc.w2
			}

			eo := findGateByOut(oth, gates)
			swapOut(and, eo)
		}

		cout = nc.out
	}

	return swaps
}

func gateIsReady(gate *circuitGate, output map[string]uint8) bool {
	_, ok1 := output[gate.w1]
	_, ok2 := output[gate.w2]

	return ok1 && ok2
}

func findGate(lhs, op, rhs string, gates map[string][]*circuitGate) *circuitGate {
	for _, cg := range gates[lhs] {
		if cg.op == op &&
			(cg.w1 == lhs || cg.w1 == rhs) &&
			(cg.w2 == lhs || cg.w2 == rhs) {
			return cg
		}
	}

	return nil
}

func findByOther(in, op, out string, gates map[string][]*circuitGate) *circuitGate {
	for _, cg := range gates[in] {
		if cg.op == op && cg.out == out {
			return cg
		}
	}

	panic(fmt.Sprintf("op %s with out %s not found for %s\n", op, out, in))
}

func findGateByOtherAndOp(oth, op string, gates map[string][]*circuitGate) *circuitGate {
	for _, cg := range gates[oth] {
		if cg.op == op {
			return cg
		}
	}

	panic(fmt.Sprintf("op %s not found for %s\n", op, oth))
}

func findGateByOut(out string, gates map[string][]*circuitGate) *circuitGate {
	for _, gts := range gates {
		for _, gt := range gts {
			if gt.out == out {
				return gt
			}
		}
	}

	panic(fmt.Sprintf("couldn't find output to %s\n", out))
}
