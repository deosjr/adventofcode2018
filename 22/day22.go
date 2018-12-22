package main

import (
	"fmt"
	"io/ioutil"
	"math"
)

type region int

const (
	rocky region = iota
	wet
	narrow
)

type coord struct {
	x, y int
}

// assumption from input: tx < ty
func part1(depth, tx, ty int) int {
	sum := 0
	m := map[coord]int{}
	for y := 0; y <= ty; y++ {
		for x := 0; x <= tx; x++ {
			c := coord{x, y}
			g := geologicIndex(m, x, y)
			if y == ty && x == tx {
				g = 0
			}
			e := erosionLevel(g, depth)
			sum += int(regionType(e))
			m[c] = e
		}
	}
	return sum
}

func geologicIndex(m map[coord]int, x, y int) int {
	// also catches the first case, {0,0}
	if y == 0 {
		return x * 16807
	}
	if x == 0 {
		return y * 48271
	}
	// guaranteed to exist due to traversal order
	return m[coord{x - 1, y}] * m[coord{x, y - 1}]
}

func erosionLevel(geo, depth int) int {
	return (geo + depth) % 20183
}

func regionType(erosion int) region {
	return region(erosion % 3)
}

// naive solution sketch for part 2:
// we start in state (0,0,torch) and win in (tx,ty,torch)
// find the minimal path to winning state
// each coord has two possible states
// each regiontype can only cross into one other regiontype
// ask every square to expand to its neighbours, starting at start
// up to a certain point, say manhattan distance to target times 2
// (naive upper bound).

// TODO: this first version is slowwww at ~50sec...

type tool uint8

const (
	torch tool = iota
	climbingGear
	neither
)

type posTool struct {
	x, y int
	tool tool
}

func geologicIndexPart2(m map[coord]int, depth int, pos, target coord) int {
	if pos.y == target.y && pos.x == target.x {
		return 0
	}
	// also catches the first case, {0,0}
	if pos.y == 0 {
		return pos.x * 16807
	}
	if pos.x == 0 {
		return pos.y * 48271
	}

	left := coord{pos.x - 1, pos.y}
	leftErosion, ok := m[left]
	if !ok {
		leftErosion = erosionLevelPart2(m, depth, left, target)
	}
	up := coord{pos.x, pos.y - 1}
	upErosion, ok := m[up]
	if !ok {
		upErosion = erosionLevelPart2(m, depth, up, target)
	}
	return leftErosion * upErosion
}

func erosionLevelPart2(m map[coord]int, depth int, pos, target coord) int {
	if v, ok := m[pos]; ok {
		return v
	}
	geo := geologicIndexPart2(m, depth, pos, target)
	erosion := (geo + depth) % 20183
	m[pos] = erosion
	return erosion
}

func regionTypePart2(m map[coord]int, depth int, pos, target coord) region {
	e := erosionLevelPart2(m, depth, pos, target)
	return regionType(e)
}

func part2(depth, tx, ty int) int {
	target := coord{tx, ty}
	threshold := 2 * (tx + ty)
	erosion := map[coord]int{}
	m := map[posTool]int{
		{0, 0, torch}: 0,
	}
	fringe := []posTool{{0, 0, torch}}
	var currentPos posTool
	for len(fringe) > 0 {
		if len(fringe) == 1 {
			currentPos = fringe[0]
			fringe = nil
		} else {
			currentPos, fringe = fringe[0], fringe[1:]
		}
		explored := explore(m, erosion, depth, threshold, currentPos, target)
		fringe = append(fringe, explored...)
	}

	min := math.MaxInt64
	for _, t := range []tool{torch, climbingGear, neither} {
		mins, ok := m[posTool{tx, ty, t}]
		if !ok {
			continue
		}
		if t != torch {
			mins += 7
		}
		if mins < min {
			min = mins
		}
	}
	return min
}

// given a position, check all its von neumann neighbours
// if we can travel there, compare travel times with tool needed for travel
// - it hasnt been explored with this tool: set travel time with tool
// - it has been explored with this tool: set travel time only if its lower
// add it to new fringe if we updated the travel time only
func explore(m map[posTool]int, erosion map[coord]int, depth, threshold int, current posTool, target coord) []posTool {
	currentPos := coord{current.x, current.y}
	targetPos := coord{target.x, target.y}
	currentMins := m[current]
	r := regionTypePart2(erosion, depth, currentPos, targetPos)

	newFringe := []posTool{}
	for _, n := range neighbours(currentPos) {
		nr := regionTypePart2(erosion, depth, n, targetPos)
		neededTool, ok := travel(r, nr, current.tool)
		if !ok {
			continue
		}
		npos := posTool{n.x, n.y, neededTool}
		newMins := currentMins + 1
		if neededTool != current.tool {
			newMins += 7
		}
		if newMins > threshold {
			continue
		}
		v, explored := m[npos]
		if !explored {
			m[npos] = newMins
			newFringe = append(newFringe, npos)
			continue
		}
		if v <= newMins {
			continue
		}
		m[npos] = newMins
		newFringe = append(newFringe, npos)
	}
	return newFringe
}

func neighbours(pos coord) []coord {
	n := []coord{{pos.x + 1, pos.y}, {pos.x, pos.y + 1}}
	if pos.x != 0 {
		n = append(n, coord{pos.x - 1, pos.y})
	}
	if pos.y != 0 {
		n = append(n, coord{pos.x, pos.y - 1})
	}
	return n
}

func travel(from, to region, t tool) (tool, bool) {
	switch from {
	case rocky:
		switch to {
		case rocky:
			return t, true
		case wet:
			return climbingGear, true
		case narrow:
			return torch, true
		}
	case wet:
		switch to {
		case rocky:
			return climbingGear, true
		case wet:
			return t, true
		case narrow:
			return neither, true
		}
	case narrow:
		switch to {
		case rocky:
			return torch, true
		case wet:
			return neither, true
		case narrow:
			return t, true
		}
	}
	return t, false
}

func main() {
	input, err := ioutil.ReadFile("day22.input")
	if err != nil {
		panic(err)
	}
	var depth, tx, ty int
	fmt.Sscanf(string(input), "depth: %d\ntarget: %d,%d", &depth, &tx, &ty)
	fmt.Printf("Part 1: %d\n", part1(depth, tx, ty))
	fmt.Printf("Part 2: %d\n", part2(depth, tx, ty))
}
