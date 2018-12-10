package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

// quick/naive: loop over all seconds and visually inspect the answer		(done)
// better: detect the answer using the idea that the points will converge	(done)
// even better: dont loop over all seconds but do a gradient descent		(done)
// EVEN better: proper dynamic gamma so finding one is less volatile		(done)
// best(?): dont visually inspect but actually print the answer

type light struct {
	x, y   int
	vx, vy int
}

type boundingBox struct {
	lights     []light
	minX, maxX int
	minY, maxY int
	product    int
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

var evaluations = 0

// lightsAtSecond returns the positions of the lights at a certain second
// it also returns the product of the total bounding box dimensions
func lightsAtSecond(initialState []light, second int) boundingBox {
	evaluations++
	bb := boundingBox{
		lights: make([]light, len(initialState)),
		minX:   math.MaxInt64,
		maxX:   math.MinInt64,
		minY:   math.MaxInt64,
		maxY:   math.MinInt64,
	}
	for i, l := range initialState {
		x := l.x + second*l.vx
		y := l.y + second*l.vy
		bb.lights[i] = light{x: x, y: y}
		bb.minX = min(x, bb.minX)
		bb.maxX = max(x, bb.maxX)
		bb.minY = min(y, bb.minY)
		bb.maxY = max(y, bb.maxY)
	}
	w := bb.maxX - bb.minX
	h := bb.maxY - bb.minY
	bb.product = w * h
	return bb
}

// gradientDescent finds the minimum of the function lightsAtSecond on a certain second
// it returns the configuration of lights at that second and the second itself
// Note that this function (the boundingBox function) does not have a derivative,
// so we approximate it by evaluations second-1 and second+1, taking the diffs of
// current state to both and averaging the two.
// Since we know the function is convex and therefore has no local minima to get trapped in,
// we can just try again whenever we overshoot our goal.
func gradientDescent(initial []light, second int, gamma float64, precision, prevStepSize, prevSecond, prevValue int) (boundingBox, int) {
	nowBB := lightsAtSecond(initial, second)
	if prevStepSize <= precision {
		return nowBB, second
	}
	// dynamic gamma: if we overshot our goal, try again with lower gamma (by 15%)
	// otherwise, increase gamma by 5%
	if nowBB.product > prevValue {
		return gradientDescent(initial, prevSecond, gamma*0.85, precision, prevStepSize, prevSecond, prevValue)
	}

	prevBB := lightsAtSecond(initial, second-1)
	nextBB := lightsAtSecond(initial, second+1)
	prevDiff := nowBB.product - prevBB.product
	nextDiff := nextBB.product - nowBB.product
	df := float64(prevDiff+nextDiff) / 2.0

	newSecond := second + int(-gamma*df)
	stepSize := abs(second - newSecond)

	return gradientDescent(initial, newSecond, gamma*1.05, precision, stepSize, second, nowBB.product)
}

func main() {
	input, err := ioutil.ReadFile("day10.input")
	if err != nil {
		panic(err)
	}
	splitinput := strings.Split(string(input), "\n")

	lights := make([]light, len(splitinput))
	for i, s := range splitinput {
		var posX, posY, vX, vY int
		fmt.Sscanf(s, "position=<%d,%d> velocity=<%d,%d>", &posX, &posY, &vX, &vY)
		lights[i] = light{posX, posY, vX, vY}
	}

	// good gamma found by experimentation :)
	gamma := 0.004
	precision := 1
	bb, second := gradientDescent(lights, 0, gamma, precision, 2, 0, math.MaxInt64)

	fmt.Println("Part1:")
	for y := bb.minY; y <= bb.maxY; y++ {
	XLoop:
		for x := bb.minX; x <= bb.maxX; x++ {
			for _, l := range bb.lights {
				if l.x == x && l.y == y {
					fmt.Print("#")
					continue XLoop
				}
			}
			fmt.Print(".")
		}
		fmt.Println()
	}

	fmt.Printf("Part 2: %d\n", second)

	fmt.Printf("Bonus: evaluated lights state %d times\n", evaluations)
}
