package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type coord struct {
	x, y int
}

type square uint8

const (
	sand square = iota
	clay
	waterFlowing
	waterStanding
)

type slice struct {
	squares    map[coord]square
	xMin, xMax int
	yMin, yMax int
}

func parse(input string) slice {
	squares := map[coord]square{}
	txMin, txMax := math.MaxInt64, math.MinInt64
	tyMin, tyMax := math.MaxInt64, math.MinInt64
	for _, s := range strings.Split(input, "\n") {
		if s[0] == 'x' {
			var x, yMin, yMax int
			fmt.Sscanf(s, "x=%d, y=%d..%d", &x, &yMin, &yMax)
			for y := yMin; y <= yMax; y++ {
				squares[coord{x, y}] = clay
			}
			if x < txMin {
				txMin = x
			}
			if x > txMax {
				txMax = x
			}
			if yMin < tyMin {
				tyMin = yMin
			}
			if yMax > tyMax {
				tyMax = yMax
			}
			continue
		}
		var y, xMin, xMax int
		fmt.Sscanf(s, "y=%d, x=%d..%d", &y, &xMin, &xMax)
		for x := xMin; x <= xMax; x++ {
			squares[coord{x, y}] = clay
		}
		if y < tyMin {
			tyMin = y
		}
		if y > tyMax {
			tyMax = y
		}
		if xMin < txMin {
			txMin = xMin
		}
		if xMax > txMax {
			txMax = xMax
		}
	}
	return slice{
		squares: squares,
		xMin:    txMin,
		xMax:    txMax,
		yMin:    tyMin,
		yMax:    tyMax,
	}
}

func flow(s slice, spring coord) {
	sources := []coord{spring}
	checked := map[coord]struct{}{}
	for len(sources) > 0 {
		var source coord
		source, sources = sources[0], sources[1:]
		if _, ok := checked[source]; ok {
			continue
		}
		checked[source] = struct{}{}
		landing, outOfBounds := s.flowDown(source)
		if outOfBounds {

			continue
		}
		newSources := s.fillUp(landing)
		sources = append(sources, newSources...)
	}
}

func (s slice) numWater() int {
	sum := 0
	for k, v := range s.squares {
		if k.y < s.yMin {
			continue
		}
		if v == waterFlowing || v == waterStanding {
			sum++
		}
	}
	return sum
}

func (s slice) numWaterStanding() int {
	sum := 0
	for k, v := range s.squares {
		if k.y < s.yMin {
			continue
		}
		if v == waterStanding {
			sum++
		}
	}
	return sum
}

// water flows down until it hits clay or water
// bool returns true if we flow out of bounds
func (s slice) flowDown(source coord) (coord, bool) {
	if source.y+1 > s.yMax {
		return coord{}, true
	}
	down := coord{source.x, source.y + 1}
	downSquare := s.squares[down]
	if downSquare == sand || downSquare == waterFlowing {
		s.squares[down] = waterFlowing
		return s.flowDown(down)
	}
	return source, false
}

// water flows to the left and right.
// if it hits walls on both ends, it does the same one level up
// if it hits a gap, create a source
func (s slice) fillUp(landing coord) []coord {
	s.squares[landing] = waterFlowing
	filledLeft, sourceLeft := s.fillLeft(landing)
	filledRight, sourceRight := s.fillRight(landing)
	// hit walls on both ends
	if sourceLeft == nil && sourceRight == nil {
		for _, c := range filledLeft {
			s.squares[c] = waterStanding
		}
		for _, c := range filledRight {
			s.squares[c] = waterStanding
		}
		s.squares[landing] = waterStanding
		up := coord{landing.x, landing.y - 1}
		return s.fillUp(up)
	}
	// hit a gap on at least one end
	sources := []coord{}
	if sourceLeft != nil {
		sources = append(sources, *sourceLeft)
	}
	if sourceRight != nil {
		sources = append(sources, *sourceRight)
	}
	return sources
}

func (s slice) fillLeft(landing coord) ([]coord, *coord) {
	filled := []coord{}
	for {
		left := coord{landing.x - 1, landing.y}
		if s.squares[left] == clay {
			return filled, nil
		}
		s.squares[left] = waterFlowing
		leftDown := coord{left.x, left.y + 1}
		squareLeftDown := s.squares[leftDown]
		// found a gap, left is new source
		if squareLeftDown == sand || squareLeftDown == waterFlowing {
			return filled, &left
		}
		filled = append(filled, left)
		landing = left
	}
}

func (s slice) fillRight(landing coord) ([]coord, *coord) {
	filled := []coord{}
	for {
		right := coord{landing.x + 1, landing.y}
		if s.squares[right] == clay {
			return filled, nil
		}
		s.squares[right] = waterFlowing
		rightDown := coord{right.x, right.y + 1}
		squarerightDown := s.squares[rightDown]
		// found a gap, right is new source
		if squarerightDown == sand || squarerightDown == waterFlowing {
			return filled, &right
		}
		filled = append(filled, right)
		landing = right
	}
}

// NOTE: catch 1: water can go beyond x bounds and still count
func (s slice) PrintSelf() string {
	list := []string{}
	for y := s.yMin; y <= s.yMax; y++ {
		out := ""
		for x := s.xMin - 1; x <= s.xMax+1; x++ {
			sq := s.squares[coord{x, y}]
			switch sq {
			case sand:
				out += "."
			case clay:
				out += "#"
			case waterFlowing:
				out += "|"
			case waterStanding:
				out += "~"
			}
		}
		list = append(list, out)
	}
	return strings.Join(list, "\n")
}

func main() {
	input, err := ioutil.ReadFile("day17.input")
	if err != nil {
		panic(err)
	}
	s := parse(string(input))
	spring := coord{500, 0}
	flow(s, spring)
	fmt.Printf("Part 1: %d\n", s.numWater())
	fmt.Printf("Part 2: %d\n", s.numWaterStanding())
}
