package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func init() {
	Solutions["day17"] = day17
}

type register struct{ a, b, c, p int }
type combo func(reg *register) int

func literal(v int) combo    { return func(_ *register) int { return v } }
func regA(reg *register) int { return reg.a }
func regB(reg *register) int { return reg.b }
func regC(reg *register) int { return reg.c }

type instruction func(int, *register, []combo)

func adv(op int, reg *register, cbs []combo) {
	num := reg.a
	den := 1 << cbs[op](reg)

	reg.a = num / den
}

func bxl(op int, reg *register, _ []combo)   { reg.b ^= op }
func bst(op int, reg *register, cbs []combo) { reg.b = cbs[op](reg) % 8 }
func jnz(op int, reg *register, _ []combo) {
	if reg.a != 0 {
		reg.p = op - 1
	}
}

func bxc(_ int, reg *register, _ []combo) { reg.b ^= reg.c }
func out(sb *strings.Builder) instruction {
	return func(op int, reg *register, cbs []combo) {
		sb.WriteString(fmt.Sprintf("%d,", cbs[op](reg)%8))
	}
}
func bdv(op int, reg *register, cbs []combo) {
	num := reg.a
	den := 1 << cbs[op](reg)

	reg.b = num / den
}

func cdv(op int, reg *register, cbs []combo) {
	num := reg.a
	den := 1 << cbs[op](reg)

	reg.c = num / den
}

func day17(input string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	regRegx := regexp.MustCompile("([A-C]): (\\d+)")
	r := make([]int, 3)
	var program string

	var isProgram bool
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			isProgram = true
			continue
		}

		if isProgram {
			program = strings.TrimPrefix(line, "Program: ")
			continue
		}

		match := regRegx.FindAllStringSubmatch(line, -1)[0]
		idx := int(match[1][0] - 'A')
		r[idx] = toInt(match[2])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	partOne := runProgram(program, &register{a: r[0], b: r[1], c: r[2]})
	partTwo := searchForA(program, strings.Split(program, ","))

	fmt.Printf("Part1: %s, Part2: %d\n", partOne, partTwo)
}

func runProgram(prg string, reg *register) string {
	cbs := []combo{literal(0), literal(1), literal(2), literal(3), regA, regB, regC}

	var sb strings.Builder
	ins := []instruction{adv, bxl, bst, jnz, bxc, out(&sb), bdv, cdv}

	sprg := strings.Split(prg, ",")
	for reg.p < len(sprg) {
		i := toInt(sprg[reg.p])
		reg.p++
		op := toInt(sprg[reg.p])

		ins[i](op, reg, cbs)
		reg.p++
	}

	return strings.TrimRight(sb.String(), ",")
}

func searchForA(prg string, target []string) int {
	var tmpA int
	if len(target) > 1 {
		tmpA = 8 * searchForA(prg, target[1:])
	}

	tgt := strings.Join(target, ",")
	for runProgram(prg, &register{a: tmpA}) != tgt {
		tmpA++
	}

	return tmpA
}
