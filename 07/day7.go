package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

type node struct {
	id       rune
	prereqs  []rune
	children []*node
}

func part1(nodes []*node) string {
	if len(nodes) == 0 {
		return ""
	}
	n, tail := nodes[0], nodes[1:]
	tail = updatePreReqs(n, tail)
	str := part1(tail)
	return string(n.id) + str
}

func updatePreReqs(n *node, list []*node) []*node {
	for _, c := range n.children {
		if isEmpty := c.removeFromPreReqs(n.id); isEmpty {
			list = addToSorted(list, c)
		}
	}
	return list
}

type worker struct {
	node    *node
	minDone int
}

func (w worker) idle() bool {
	return w.node == nil
}

func part2(nodes []*node, workers []*worker, minute int) int {
	// check workers are done
	nextMinute := math.MaxInt64
	idleWorkers := []*worker{}
	for _, w := range workers {
		if w.idle() {
			idleWorkers = append(idleWorkers, w)
			continue
		}
		if w.minDone == minute {
			nodes = updatePreReqs(w.node, nodes)
			w.node = nil
			idleWorkers = append(idleWorkers, w)
			continue
		}
		if w.minDone < nextMinute {
			nextMinute = w.minDone
		}
	}

	// divide new work among workers
	for _, w := range idleWorkers {
		if len(nodes) == 0 {
			break
		}
		work, tail := nodes[0], nodes[1:]
		nodes = tail
		w.node = work
		w.minDone = minute + int(work.id-64) + 60
		if w.minDone < nextMinute {
			nextMinute = w.minDone
		}
	}

	// last piece of work completed
	if nextMinute == math.MaxInt64 {
		return minute
	}
	return part2(nodes, workers, nextMinute)
}

func (n *node) removeFromPreReqs(id rune) (isEmpty bool) {
	if len(n.prereqs) == 1 {
		isEmpty = true
	}
	index := -1
	for i, r := range n.prereqs {
		if r == id {
			index = i
			break
		}
	}
	if index == len(n.prereqs)-1 {
		n.prereqs = n.prereqs[:index]
		return isEmpty
	}
	n.prereqs = append(n.prereqs[:index], n.prereqs[index+1:]...)
	return isEmpty
}

func addToSorted(nodes []*node, n *node) []*node {
	if len(nodes) == 0 {
		return []*node{n}
	}
	// find smallest index for which nodes[i].id > n.id
	index := sort.Search(len(nodes), func(i int) bool { return nodes[i].id > n.id })
	// not found: n.id is biggest so insert last
	if index == len(nodes) {
		return append(nodes, n)
	}
	return append(nodes[:index], append([]*node{n}, nodes[index:]...)...)
}

func parseInput(input []string) []*node {
	parseMap := map[rune]*node{}
	possibleRoots := map[rune]struct{}{}
	for _, s := range input {
		var prereq, id rune
		fmt.Sscanf(s, "Step %c must be finished before step %c can begin.", &prereq, &id)
		possibleRoots[prereq] = struct{}{}
		n, ok := parseMap[id]
		if !ok {
			parseMap[id] = &node{id: id, prereqs: []rune{prereq}}
			continue
		}
		n.prereqs = append(n.prereqs, prereq)
	}

	// find runes that have no prereqs
	roots := []*node{}
	for id, _ := range possibleRoots {
		if _, ok := parseMap[id]; !ok {
			n := &node{id: id}
			parseMap[id] = n
			roots = addToSorted(roots, n)
		}
	}

	// update children with info from parseMap
	for _, n := range parseMap {
		for _, prereq := range n.prereqs {
			parent := parseMap[prereq]
			parent.children = append(parent.children, n)
		}
	}
	return roots
}

func main() {
	input, err := ioutil.ReadFile("day7.input")
	if err != nil {
		panic(err)
	}
	splitinput := strings.Split(string(input), "\n")

	roots := parseInput(splitinput)
	fmt.Printf("Part 1: %s\n", part1(roots))

	// parse again because we build up a mutable graph
	roots = parseInput(splitinput)
	numWorkers := 5
	workers := make([]*worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = &worker{}
	}
	fmt.Printf("Part 2: %d\n", part2(roots, workers, 0))
}
