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
		for i := 4; i >= 0; i-- {
			if s[i] == '#' {
				sum += int(math.Exp2(float64(-i + 4)))
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

func part1(s state, transitions map[int]bool) int {
	for i := 0; i < 20; i++ {
		s = generation(s, transitions)
	}
	sum := 0
	for k, _ := range s.plants {
		sum += k
	}
	return sum
}

// start from ....# which is two left of the leftmost plant
// then loop over all pots until two right of rightmost plant
// calculate neighbourhood by using bitshift
func generation(s state, transitions map[int]bool) state {
	newPlants := map[int]bool{}
	neighbourhood := 1
	min := s.max
	max := s.min
	for i := s.min - 2; i <= s.max+2; i++ {
		if i != s.min-2 {
			if neighbourhood >= 16 {
				neighbourhood -= 16
			}
			neighbourhood *= 2
			if s.plants[i+2] {
				neighbourhood++
			}
		}
		if transitions[neighbourhood] {
			newPlants[i] = true
			if i < min {
				min = i
			}
			if i > max {
				max = i
			}
		}
	}
	return state{min: min, max: max, plants: newPlants}
}

func main() {
	input, err := ioutil.ReadFile("day12.input")
	if err != nil {
		panic(err)
	}
	initial, transitions := parse(string(input))
	out := part1(initial, transitions)
	fmt.Printf("Part 1: %d\n", out)
}
