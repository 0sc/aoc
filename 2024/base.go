package main

import "strconv"

var Solutions = map[string]func(string){}

func toInt(v string) int {
	num, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return num
}

func abs(num int) int {
	if num < 0 {
		num = 0 - num
	}

	return num
}
