package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
)

func init() {
	Solutions["day9"] = day9
}

func day9(input string) {
	var partOne, partTwo int

	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		partOne = compactBlock([]byte(line))
		partTwo = compactFile([]byte(line))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part1: %d, Part2: %d\n", partOne, partTwo)
}

func compactBlock(layout []byte) int {
	var cs, bp int

	lfid, rfid := 0, len(layout)/2

	l, r := 0, len(layout)-1
	for l <= r {
		lbs := btoi(layout[l])
		cs, bp = updateCheckSum(cs, lfid, lbs, bp)
		lfid++

		l++
		fs := btoi(layout[l])
		for fs > 0 && r > l {
			if layout[r] == '0' {
				r -= 2
				rfid--

				continue
			}

			rbs := btoi(layout[r])
			ms := min(rbs, fs)
			cs, bp = updateCheckSum(cs, rfid, ms, bp)

			fs -= ms
			layout[r] -= byte(ms)
		}

		l++
	}

	return cs
}

// index, blocksize
type dfHeap [][2]int

func (h dfHeap) Len() int           { return len(h) }
func (h dfHeap) Less(i, j int) bool { return h[i][0] < h[j][0] }
func (h dfHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *dfHeap) Push(x any)        { *h = append(*h, x.([2]int)) }

func (h *dfHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func compactFile(layout []byte) int {
	dfh := &dfHeap{}

	var cs, bp int
	r := len(layout) - 1
	for j := 0; j < r; j += 2 {
		bp += btoi(layout[j])
		fs := btoi(layout[j+1])
		if fs == 0 {
			continue
		}

		heap.Push(dfh, [2]int{bp, fs})
		bp += fs
	}

	bp += btoi(layout[r])

	fid := len(layout) / 2
	for ; r >= 0; r -= 2 {
		// poll from heap till find a size that fits
		bs := btoi(layout[r])
		bp -= bs
		sp := [2]int{bp, bs}
		var xt [][2]int
		for dfh.Len() > 0 {
			s := heap.Pop(dfh).([2]int)
			if s[0] > sp[0] { // discard spaces that come after current pos
				continue
			}

			if s[1] >= bs {
				sp = s

				break
			}

			xt = append(xt, s)
		}

		// found a space for it
		cs, _ = updateCheckSum(cs, fid, bs, sp[0])

		// check if there's space left
		ls := sp[1] - bs
		if ls > 0 {
			heap.Push(dfh, [2]int{sp[0] + bs, ls})
		}

		for _, x := range xt {
			heap.Push(dfh, x)
		}

		if r > 0 {
			bp -= btoi(layout[r-1])
		}

		fid -= 1
	}

	return cs
}

func updateCheckSum(cs int, fid int, bs int, bp int) (int, int) {
	for i := 0; i < bs; i++ {
		cs += fid * bp
		bp++
	}

	return cs, bp
}

func btoi(b byte) int { return int(b - '0') }
