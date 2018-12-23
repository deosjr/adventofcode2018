package main

import (
	"container/heap"
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

func part2(depth, tx, ty int) int {
	erosion := map[coord]int{}
	start := posTool{0, 0, torch}
	goal := posTool{tx, ty, torch}

	return findRoute(start, goal, erosion, depth)
}

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

func neighbours(pt, goal posTool, erosion map[coord]int, depth int) []posTool {
	pos := coord{pt.x, pt.y}
	ncoord := []coord{{pos.x + 1, pos.y}, {pos.x, pos.y + 1}}
	if pos.x != 0 {
		ncoord = append(ncoord, coord{pos.x - 1, pos.y})
	}
	if pos.y != 0 {
		ncoord = append(ncoord, coord{pos.x, pos.y - 1})
	}

	self := posTool{pt.x, pt.y, 0}
	goalPos := coord{goal.x, goal.y}
	r := regionTypePart2(erosion, depth, pos, goalPos)
	switch r {
	case rocky:
		if pt.tool == torch {
			self.tool = climbingGear
		} else {
			self.tool = torch
		}
	case wet:
		if pt.tool == climbingGear {
			self.tool = neither
		} else {
			self.tool = climbingGear
		}
	case narrow:
		if pt.tool == torch {
			self.tool = neither
		} else {
			self.tool = torch
		}
	}

	n := []posTool{self}
	for _, nc := range ncoord {
		nr := regionTypePart2(erosion, depth, nc, goalPos)
		neededTool, ok := travel(r, nr, pt.tool)
		if !ok {
			continue
		}
		n = append(n, posTool{nc.x, nc.y, neededTool})
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

func findRoute(start, goal posTool, erosion map[coord]int, depth int) int {
	openSet := map[posTool]bool{
		start: true,
	}
	cameFrom := map[posTool]posTool{}

	gScore := map[posTool]float64{}
	gScore[start] = 0

	fScore := map[posTool]float64{}
	fScore[start] = h(start, goal)

	pq := priorityQueue{
		&pqItem{
			posTool: start,
			fScore:  fScore[start],
			index:   0,
		},
	}
	heap.Init(&pq)

	goalScore := float64(math.MaxInt64)

	for pq.Len() != 0 {
		item := heap.Pop(&pq).(*pqItem)
		current := item.posTool
		if current == goal {
			goalScore = gScore[current]
		}

		delete(openSet, current)

		for _, n := range neighbours(current, goal, erosion, depth) {
			gCurrent := gScore[current]
			gn := gScore[n]

			tentativeGscore := gCurrent + g(current, n)
			f := tentativeGscore + h(n, goal)

			v, ok := fScore[n]
			if !openSet[n] && f < goalScore && (!ok || f < v) {
				openSet[n] = true
				item := &pqItem{
					posTool: n,
					fScore:  f,
				}
				heap.Push(&pq, item)
			} else if tentativeGscore >= gn {
				continue
			}
			cameFrom[n] = current
			gScore[n] = tentativeGscore
			fScore[n] = f
		}
	}

	return int(goalScore)
}

func g(p, q posTool) float64 {
	var min float64
	if !(p.x == q.x && p.y == q.y) {
		min += 1.0
	}
	if p.tool != q.tool {
		min += 7.0
	}
	return min
}

func h(p, q posTool) float64 {
	m := manhattan(p, q)
	if p.tool != q.tool {
		m += 7
	}
	return m
}

func manhattan(p, q posTool) float64 {
	dx := float64(q.x - p.x)
	dy := float64(q.y - p.y)
	return math.Abs(dx) + math.Abs(dy)
}

func reconstructPath(m map[posTool]posTool, current posTool) []posTool {
	path := []posTool{current}
	for {
		prev, ok := m[current]
		if !ok {
			break
		}
		current = prev
		path = append(path, current)
	}
	return path
}

type pqItem struct {
	posTool posTool
	fScore  float64
	index   int
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].fScore < pq[j].fScore
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
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
