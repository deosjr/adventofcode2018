package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

// nice way: generate voronoi diagram, compute area of polygons
// naive / quick and dirty way: calculate manhattan distance to every node for every coordinate..
// never mind, for part 2 we need to calculate everything anyways!

type coord struct {
	x, y int
}

func manhattan(p, q coord) int {
	xDiff := math.Abs(float64(p.x) - float64(q.x))
	yDiff := math.Abs(float64(p.y) - float64(q.y))
	return int(xDiff + yDiff)
}

func main() {
	input, err := ioutil.ReadFile("day6.input")
	if err != nil {
		panic(err)
	}
	splitinput := strings.Split(string(input), "\n")

	var xMax, yMax int
	coords := make([]coord, len(splitinput))
	for i, s := range splitinput {
		split := strings.Split(s, ",")
		x, _ := strconv.Atoi(split[0])
		y, _ := strconv.Atoi(strings.TrimSpace(split[1]))
		coords[i] = coord{x, y}
		if x > xMax {
			xMax = x
		}
		if y > yMax {
			yMax = y
		}
	}

	// map of id -> areasize. size -1 means infinite
	sizes := map[int]int{}

	var c coord
	var safe, closest, numClosest int
	// loop over all coordinates from {0,0} to {xMax, yMax} inclusive
	for y := 0; y <= yMax; y++ {
		for x := 0; x <= xMax; x++ {
			c = coord{x, y}
			distance := 99999
			sum := 0
			for i, cc := range coords {
				m := manhattan(c, cc)
				sum += m
				switch {
				case m < distance:
					distance = m
					closest = i
					numClosest = 1
				case m == distance:
					numClosest += 1
				case m > distance:
					continue
				}
			}
			if sum < 10000 {
				safe += 1
			}
			if numClosest > 1 {
				continue
			}
			if v, ok := sizes[closest]; ok && v == -1 {
				continue
			}
			if x == 0 || x == xMax || y == 0 || y == yMax {
				sizes[closest] = -1
				continue
			}
			sizes[closest] += 1
		}
	}

	_, ans := maxValueWithKey(sizes)
	fmt.Printf("Part 1: %d\n", ans)

	fmt.Printf("Part 2: %d\n", safe)
}

// reused from day 3
func maxValueWithKey(m map[int]int) (int, int) {
	var key, max int
	for k, v := range m {
		if v >= max {
			max = v
			key = k
		}
	}
	return key, max
}
