package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type state struct {
	plants   map[int]bool
	min, max int
	shiftSum int
}

func parse(input string) (state, map[int]bool) {
	splitinput := strings.Split(input, "\n")
	initial := strings.Split(splitinput[0], ": ")
	initPlants := map[int]bool{}
	for i, c := range initial[1] {
		if c == '#' {
			initPlants[i] = true
		}
	}

	transitions := map[int]bool{}
	for _, s := range splitinput[2:] {
		// form is: [#.]{5} => [#.]
		sum := 0
		for i := 0; i < 5; i++ {
			if s[i] == '#' {
				sum += int(math.Exp2(float64(i)))
			}
		}
		if s[9] == '#' {
			transitions[sum] = true
		}
	}
	initState := state{
		plants: initPlants,
		min:    0,
		max:    len(initial[1]),
	}
	return initState, transitions
}

func generations(s state, transitions map[int]bool, g int) int {
	var genSimulated int
	// simulate until pattern no longer changes, only shifts
	for genSimulated = 0; genSimulated < g; genSimulated++ {
		newState := generation(s, transitions)
		if s.shiftSum == newState.shiftSum {
			break
		}
		s = newState
	}

	// include shifts in sum
	// assumption: pattern shifts by one every generation from now
	genRemaining := g - genSimulated
	sum := 0
	for k, _ := range s.plants {
		sum += k + genRemaining
	}
	return sum
}

func generation(s state, transitions map[int]bool) state {
	newPlants := map[int]bool{}
	min := s.max
	max := s.min
	checked := map[int]struct{}{}
	for k, _ := range s.plants {
		for i := k - 2; i <= k+2; i++ {
			if _, ok := checked[i]; ok {
				continue
			}
			checked[i] = struct{}{}
			if plantAt(i, s.plants, transitions) {
				newPlants[i] = true
				if i < min {
					min = i
				} else if i > max {
					max = i
				}
			}
		}
	}

	shiftSum := 0
	for k, _ := range newPlants {
		shiftSum += k - min
	}

	newState := state{
		plants:   newPlants,
		min:      min,
		max:      max,
		shiftSum: shiftSum,
	}
	return newState
}

func plantAt(index int, plants, transitions map[int]bool) bool {
	neighbourhood := 0
	for i := 0; i < 5; i++ {
		if plants[index-2+i] {
			neighbourhood += int(math.Exp2(float64(i)))
		}
	}
	return transitions[neighbourhood]
}

func main() {
	input, err := ioutil.ReadFile("day12.input")
	if err != nil {
		panic(err)
	}
	initial, transitions := parse(string(input))
	part1 := generations(initial, transitions, 20)
	fmt.Printf("Part 1: %d\n", part1)
	part2 := generations(initial, transitions, 50000000000)
	fmt.Printf("Part 2: %d\n", part2)
}
