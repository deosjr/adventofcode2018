package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type coord struct {
	x, y int
}

type acre uint8

const (
	open acre = iota
	wooded
	lumberyard
)

func parse(input string) map[coord]acre {
	m := map[coord]acre{}
	for y, s := range strings.Split(input, "\n") {
		for x, c := range s {
			switch c {
			case '.':
				m[coord{x, y}] = open
			case '|':
				m[coord{x, y}] = wooded
			case '#':
				m[coord{x, y}] = lumberyard
			}
		}
	}
	return m
}

func part1(grid map[coord]acre) int {
	// print(grid)
	for minute := 0; minute < 10; minute++ {
		grid = tick(grid)
	}
	var wood, yards int
	for _, v := range grid {
		if v == wooded {
			wood++
			continue
		}
		if v == lumberyard {
			yards++
		}
	}
	return wood * yards
}

func tick(grid map[coord]acre) map[coord]acre {
	newGrid := map[coord]acre{}
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			c := coord{x, y}
			numWood, numYard := neighbours(c, grid)
			switch grid[c] {
			case open:
				if numWood >= 3 {
					newGrid[c] = wooded
					continue
				}
				newGrid[c] = open
			case wooded:
				if numYard >= 3 {
					newGrid[c] = lumberyard
					continue
				}
				newGrid[c] = wooded
			case lumberyard:
				if numWood >= 1 && numYard >= 1 {
					newGrid[c] = lumberyard
					continue
				}
				newGrid[c] = open
			}
		}
	}
	return newGrid
}

func neighbours(c coord, grid map[coord]acre) (int, int) {
	var numWood, numYard int
	ncs := []coord{{c.x - 1, c.y - 1}, {c.x, c.y - 1}, {c.x + 1, c.y - 1}, {c.x - 1, c.y}, {c.x + 1, c.y}, {c.x - 1, c.y + 1}, {c.x, c.y + 1}, {c.x + 1, c.y + 1}}
	for _, n := range ncs {
		v, ok := grid[n]
		if !ok {
			continue
		}
		switch v {
		case wooded:
			numWood++
		case lumberyard:
			numYard++
		}
	}
	return numWood, numYard
}

func print(m map[coord]acre) {
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			switch m[coord{x, y}] {
			case open:
				fmt.Print(".")
			case wooded:
				fmt.Print("|")
			case lumberyard:
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

func main() {
	input, err := ioutil.ReadFile("day18.input")
	if err != nil {
		panic(err)
	}
	grid := parse(string(input))
	out := part1(grid)
	fmt.Printf("Part 1: %d\n", out)
}
