package main

import (
	"fmt"
	"math"
	"time"
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
	var ans coord
	var ansSize int
	max := math.MinInt64
	for size := 1; size <= maxSize; size++ {
		c, power := maxPower(size, serial)
		if power > max {
			max = power
			ans = c
			ansSize = size
		}
	}
	return ans, ansSize
}

func part2_parallel(serial int) (coord, int) {
	numWorkers := 10
	inCh := make(chan int, numWorkers)
	outCh := make(chan ans, maxSize)
	defer close(outCh)
	for i := 0; i < numWorkers; i++ {
		go func(in chan int, out chan ans) {
			for size := range in {
				c, power := maxPower(size, serial)
				out <- ans{c, size, power}
			}
		}(inCh, outCh)
	}

	for i := 1; i <= maxSize; i++ {
		inCh <- i
	}
	close(inCh)

	ansSoFar := ans{size: 0, power: math.MinInt64}
	for i := 1; i <= maxSize; i++ {
		answer := <-outCh
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

	start := time.Now()
	ans2, size := part2(input)
	took := time.Now().Sub(start)
	fmt.Printf("Part 2: %d,%d,%d ; took %s\n", ans2.x, ans2.y, size, took.String())

	start = time.Now()
	ans2, size = part2_parallel(input)
	took = time.Now().Sub(start)
	fmt.Printf("Bonus : %d,%d,%d ; took %s\n", ans2.x, ans2.y, size, took.String())
}
