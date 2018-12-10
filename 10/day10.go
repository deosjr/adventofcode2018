package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

// quick/naive: loop over all seconds and visually inspect the answer		(done)
// better: detect the answer using the idea that the points will converge	(done)
// even better: dont loop over all seconds but do a gradient descent		()
// best(?): dont visually inspect but actually print the answer				()

type light struct {
	x, y   int
	vx, vy int
}

func (l *light) update() (int, int) {
	l.x += l.vx
	l.y += l.vy
	return l.x, l.y
}

func main() {
	input, err := ioutil.ReadFile("day10.input")
	if err != nil {
		panic(err)
	}
	splitinput := strings.Split(string(input), "\n")

	lights := make([]*light, len(splitinput))
	for i, s := range splitinput {
		var posX, posY, vX, vY int
		fmt.Sscanf(s, "position=<%d,%d> velocity=<%d,%d>", &posX, &posY, &vX, &vY)
		lights[i] = &light{posX, posY, vX, vY}
	}

	var minX, maxX, minY, maxY int
	seconds := 0
	boundingBox := math.MaxInt64
Seconds:
	for {
		minX, maxX = 999999, -999999
		minY, maxY = 999999, -999999
		for _, l := range lights {
			x, y := l.update()
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
		w := maxX - minX
		h := maxY - minY
		currentBB := w * h
		if currentBB > boundingBox {
			// we went one too far, reverse
			for _, l := range lights {
				l.x -= l.vx
				l.y -= l.vy
			}
			break Seconds
		}
		boundingBox = currentBB
		seconds++
	}

	fmt.Println("Part1:")

	for y := minY; y <= maxY; y++ {
	XLoop:
		for x := minX; x <= maxX; x++ {
			for _, l := range lights {
				if l.x == x && l.y == y {
					fmt.Print("#")
					continue XLoop
				}
			}
			fmt.Print(".")
		}
		fmt.Println()
	}

	fmt.Printf("Part 2: %d\n", seconds)
}
