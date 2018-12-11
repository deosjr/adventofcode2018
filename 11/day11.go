package main

import (
	"fmt"
	"math"
)

type coord struct {
	x, y int
}

func powerLevel(c coord, serial int) int {
	rackID := c.x + 10
	p := rackID * c.y
	p += serial
	p *= rackID
	p = (p / 100) % 10
	p -= 5
	return p
}

func squareSum(minX, minY, size, serial int) int {
	sum := 0
	for y := minY; y < minY+size; y++ {
		for x := minX; x < minX+size; x++ {
			sum += powerLevel(coord{x, y}, serial)
		}
	}
	return sum
}

func maxPower(size, serial int) (coord, int) {
	var xans, yans int
	max := math.MinInt64
	limit := 300 - size + 2
	for y := 1; y < limit; y++ {
		for x := 1; x < limit; x++ {
			sum := squareSum(x, y, size, serial)
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
	for size := 1; size <= 300; size++ {
		go func(s int) {
			c, power := maxPower(size, serial)
			ch <- ans{c, s, power}
		}(size)
	}

	ansSoFar := ans{size: 0, power: math.MinInt64}
	for size := 1; size <= 300; size++ {
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
