package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type coord struct {
	x, y int
}

const overlap = -9999

func main() {
	input, err := ioutil.ReadFile("day3.input")
	if err != nil {
		panic(err)
	}

	m := map[coord]int{}
	ids := map[int]struct{}{}
	totalOverlap := 0
	splitinput := strings.Split(string(input), "\n")
	for i, claim := range splitinput {
		s := strings.Split(claim, " ")
		ss := strings.Split(s[2], ",")
		x, _ := strconv.Atoi(ss[0])
		y, _ := strconv.Atoi(ss[1][:len(ss[1])-1])
		dims := strings.Split(s[3], "x")
		w, _ := strconv.Atoi(dims[0])
		h, _ := strconv.Atoi(dims[1])

		id := i + 1

		for yy := y; yy < y+h; yy++ {
			for xx := x; xx < x+w; xx++ {
				c := coord{xx, yy}
				switch m[c] {
				case 0:
					m[c] = id
				case overlap:
					ids[id] = struct{}{}
					continue
				default:
					ids[m[c]] = struct{}{}
					ids[id] = struct{}{}
					m[c] = overlap
					totalOverlap += 1
				}
			}
		}
	}
	fmt.Printf("Part 1: %d\n", totalOverlap)
	for id := 1; id <= len(splitinput); id++ {
		if _, ok := ids[id]; !ok {
			fmt.Printf("Part 2: %d\n", id)
		}
	}
}
