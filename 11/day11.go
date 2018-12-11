package main

import (
	"fmt"
	"math"
)

const maxSize = 300

type coord struct {
	x, y int
}

func powerLevel(x, y, serial int) int {
	rackID := x + 10
	p := rackID * y
	p += serial
	p *= rackID
	p = (p / 100) % 10
	p -= 5
	return p
}

func maxPower(size, serial int) (coord, int) {
	intermediate := map[coord]int{}
	limit := maxSize - size + 1
	for y := 1; y <= limit; y++ {
		for x := 1; x <= maxSize; x++ {
			if y == 1 {
				sum := 0
				for yy := 1; yy <= size; yy++ {
					sum += powerLevel(x, yy, serial)
				}
				intermediate[coord{x, y}] = sum
				continue
			}
			sum := intermediate[coord{x, y - 1}] - powerLevel(x, y-1, serial) + powerLevel(x, y+size-1, serial)
			intermediate[coord{x, y}] = sum
		}
	}

	var xans, yans int
	max := math.MinInt64
	for y := 1; y <= limit; y++ {
		var sum int
		for x := 1; x <= limit; x++ {
			if x == 1 {
				for xx := 1; xx <= size; xx++ {
					sum += intermediate[coord{xx, y}]
				}
			} else {
				sum = sum - intermediate[coord{x - 1, y}] + intermediate[coord{x + size - 1, y}]
			}
			if sum > max {
				max = sum
				xans = x
				yans = y
			}
		}
	}
	return coord{xans, yans}, max
}

func part1(serial int) coord {
	c, _ := maxPower(3, serial)
	return c
}

type ans struct {
	c     coord
	size  int
	power int
}

func part2(serial int) (coord, int) {
	ch := make(chan ans, 1)
	for size := 1; size <= maxSize; size++ {
		go func(s int) {
			c, power := maxPower(s, serial)
			ch <- ans{c, s, power}
		}(size)
	}

	ansSoFar := ans{size: 0, power: math.MinInt64}
	for i := 1; i <= maxSize; i++ {
		answer := <-ch
		if answer.power > ansSoFar.power {
			ansSoFar = answer
		}
	}
	return ansSoFar.c, ansSoFar.size
}

func main() {
	input := 3628
	ans1 := part1(input)
	fmt.Printf("Part 1: %d,%d\n", ans1.x, ans1.y)

	ans2, size := part2(input)
	fmt.Printf("Part 2: %d,%d,%d\n", ans2.x, ans2.y, size)
}
