package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type nanobot struct {
	pos coord
	r   int
}

type coord struct {
	x, y, z int
}

func manhattan(p, q coord) int {
	dx := float64(p.x - q.x)
	dy := float64(p.y - q.y)
	dz := float64(p.z - q.z)
	return int(math.Abs(dx) + math.Abs(dy) + math.Abs(dz))
}

func (n nanobot) inRange(p coord) bool {
	return manhattan(n.pos, p) <= n.r
}

func parse(input string) []nanobot {
	split := strings.Split(input, "\n")
	nanobots := make([]nanobot, len(split))
	for i, s := range split {
		var x, y, z, r int
		fmt.Sscanf(s, "pos=<%d,%d,%d>, r=%d", &x, &y, &z, &r)
		nanobots[i] = nanobot{coord{x, y, z}, r}
	}
	return nanobots
}

func part1(nanobots []nanobot) int {
	var largest nanobot
	for _, n := range nanobots {
		if n.r > largest.r {
			largest = n
		}
	}

	sum := 0
	for _, n := range nanobots {
		if largest.inRange(n.pos) {
			sum++
		}
	}
	return sum
}

func part2(nanobots []nanobot) int {
	xmin, xmax := -10, 10
	ymin, ymax := -10, 10
	zmin, zmax := -10, 10
	var best []coord
	var bestSum int

	for exp := 7; exp >= 0; exp-- {
		zoom := int(math.Pow(10, float64(exp)))
		zoomedBots := make([]nanobot, len(nanobots))
		for i, n := range nanobots {
			zoomedBots[i] = nanobot{coord{n.pos.x / zoom, n.pos.y / zoom, n.pos.z / zoom}, n.r / zoom}
		}

		best = nil
		bestSum = 0

		for z := zmin; z <= zmax; z++ {
			for y := ymin; y <= ymax; y++ {
				for x := xmin; x <= xmax; x++ {
					c := coord{x, y, z}
					sum := 0
					for _, n := range zoomedBots {
						if n.inRange(c) {
							sum++
						}
					}
					if sum > bestSum {
						best = []coord{c}
						bestSum = sum
						continue
					}
					if sum == bestSum {
						best = append(best, c)
					}
				}
			}
		}

		// assumption: len(best)==1
		b := best[0]
		xmin, xmax = (b.x-1)*10, (b.x+1)*10
		ymin, ymax = (b.y-1)*10, (b.y+1)*10
		zmin, zmax = (b.z-1)*10, (b.z+1)*10
	}
	return manhattan(coord{0, 0, 0}, best[0])
}

func main() {
	input, err := ioutil.ReadFile("day23.input")
	if err != nil {
		panic(err)
	}
	nanobots := parse(string(input))
	fmt.Printf("Part 1: %d\n", part1(nanobots))
	fmt.Printf("Part 2: %d\n", part2(nanobots))
}
