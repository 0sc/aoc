package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	Solutions["day7"] = day7
}

func day7(input string) {
	var partOne, partTwo int64

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		pOne, pTwo := calibrate(scanner.Text())

		partOne += pOne
		partTwo += pTwo
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func calibrate(s string) (int64, int64) {
	ss := strings.Split(s, ": ")

	res := toInt64(ss[0])

	nums := strings.Split(ss[1], " ")
	num := toInt64(nums[0])
	if canApplyOps(nums[1:], num, res) {
		return res, res
	}

	if canApplyOpsWithMerge(nums[1:], num, res) {
		return 0, res
	}

	return 0, 0
}

func canApplyOps(nums []string, acc int64, target int64) bool {
	if len(nums) == 0 {
		return acc == target
	}

	num := toInt64(nums[0])
	if canApplyOps(nums[1:], acc+num, target) {
		return true
	}

	return canApplyOps(nums[1:], acc*num, target)
}

func canApplyOpsWithMerge(nums []string, acc int64, target int64) bool {
	if len(nums) == 0 {
		return acc == target
	}

	num := toInt64(nums[0])
	if canApplyOpsWithMerge(nums[1:], acc+num, target) {
		return true
	}

	if canApplyOpsWithMerge(nums[1:], acc*num, target) {
		return true
	}

	acc = (acc * pow10(len(nums[0]))) + num

	return canApplyOpsWithMerge(nums[1:], acc, target)
}

func toInt64(v string) int64 {
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(err)
	}

	return n
}

func pow10(v int) int64 {
	var r int64 = 1
	for i := 0; i < v; i++ {
		r *= 10
	}

	return r
}
