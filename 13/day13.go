package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type track uint8

const (
	HORIZONTAL track = iota
	VERTICAL
	SLASH
	BACKSLASH
	CROSSING
)

type heading uint8

const (
	LEFT heading = iota
	RIGHT
	UP
	DOWN
	STRAIGHT
)

type coord struct {
	x, y int
}

type minecart struct {
	pos        coord
	heading    heading
	nextSwitch heading
}

func parse(filename string) (map[coord]track, []minecart) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	splitinput := strings.Split(string(input), "\n")

	tracks := map[coord]track{}
	minecarts := []minecart{}
	for y, s := range splitinput {
		for x, c := range s {
			pos := coord{x, y}
			switch c {
			// tracks
			case '-':
				tracks[pos] = HORIZONTAL
			case '|':
				tracks[pos] = VERTICAL
			case '/':
				tracks[pos] = SLASH
			case '\\':
				tracks[pos] = BACKSLASH
			case '+':
				tracks[pos] = CROSSING
			// carts
			case '<':
				tracks[pos] = HORIZONTAL
				minecarts = append(minecarts, minecart{pos: pos, heading: LEFT, nextSwitch: LEFT})
			case '>':
				tracks[pos] = HORIZONTAL
				minecarts = append(minecarts, minecart{pos: pos, heading: RIGHT, nextSwitch: LEFT})
			case '^':
				tracks[pos] = VERTICAL
				minecarts = append(minecarts, minecart{pos: pos, heading: UP, nextSwitch: LEFT})
			case 'v':
				tracks[pos] = VERTICAL
				minecarts = append(minecarts, minecart{pos: pos, heading: DOWN, nextSwitch: LEFT})
			}
		}
	}
	return tracks, minecarts
}

func newPos(pos coord, h heading) coord {
	switch h {
	case LEFT:
		return coord{pos.x - 1, pos.y}
	case RIGHT:
		return coord{pos.x + 1, pos.y}
	case UP:
		return coord{pos.x, pos.y - 1}
	case DOWN:
		return coord{pos.x, pos.y + 1}
	}
	panic("INCORRECT HEADING")
	return coord{0, 0}
}

// this can probably be done better by picking smart values for heading enum
func direction(t track, h, s heading) (heading, heading) {
	switch t {
	case HORIZONTAL, VERTICAL:
		return h, s
	case SLASH:
		switch h {
		case LEFT:
			return DOWN, s
		case RIGHT:
			return UP, s
		case UP:
			return RIGHT, s
		case DOWN:
			return LEFT, s
		}
	case BACKSLASH:
		switch h {
		case LEFT:
			return UP, s
		case RIGHT:
			return DOWN, s
		case UP:
			return LEFT, s
		case DOWN:
			return RIGHT, s
		}
	case CROSSING:
		switch s {
		case LEFT:
			switch h {
			case LEFT:
				return DOWN, STRAIGHT
			case RIGHT:
				return UP, STRAIGHT
			case UP:
				return LEFT, STRAIGHT
			case DOWN:
				return RIGHT, STRAIGHT
			}
		case STRAIGHT:
			return h, RIGHT
		case RIGHT:
			switch h {
			case LEFT:
				return UP, LEFT
			case RIGHT:
				return DOWN, LEFT
			case UP:
				return RIGHT, LEFT
			case DOWN:
				return LEFT, LEFT
			}
		}
	}
	panic("INCORRECT HEADING")
	return 0, 0
}

func addToSorted(list []minecart, m minecart) []minecart {
	if len(list) == 0 {
		return []minecart{m}
	}
	// find smallest index for which list[i].pos comes after m.pos in 2d order
	index := sort.Search(len(list), func(i int) bool { return order2d(m.pos, list[i].pos) })
	// not found: m.id is biggest so insert last
	if index == len(list) {
		return append(list, m)
	}
	return append(list[:index], append([]minecart{m}, list[index:]...)...)
}

func removeFromSorted(list []minecart, m minecart) []minecart {
	toRemove := -1
	for i, v := range list {
		if v.pos == m.pos {
			toRemove = i
		}
	}
	if toRemove == -1 {
		fmt.Println("LIST DID NOT CONTAIN M")
		return list
	}
	if toRemove == 0 {
		return list[1:]
	}
	if toRemove == len(list)-1 {
		return list[:len(list)-1]
	}
	return append(list[:toRemove], list[toRemove+1:]...)
}

// order2d returns true if q comes after p in 2d grid order
func order2d(p, q coord) bool {
	if q.y > p.y {
		return true
	}
	if q.y < p.y {
		return false
	}
	if q.x > p.x {
		return true
	}
	return false
}

func part1(tracks map[coord]track, minecarts []minecart) coord {
	for {
		newMinecarts := make([]minecart, 0, len(minecarts))
		for i, m := range minecarts {
			// move the minecart
			p := newPos(m.pos, m.heading)

			// check collision
			for _, other := range newMinecarts {
				if other.pos == p {
					return p
				}
			}
			if i != len(minecarts)-1 {
				for _, other := range minecarts[i+1:] {
					if other.pos == p {
						return p
					}
				}
			}

			// set new heading/switch based on track
			newHeading, newSwitch := direction(tracks[p], m.heading, m.nextSwitch)

			newMinecart := minecart{p, newHeading, newSwitch}
			newMinecarts = addToSorted(newMinecarts, newMinecart)
		}
		minecarts = newMinecarts
	}
}

func part2(tracks map[coord]track, minecarts []minecart) coord {
	for {
		newMinecarts := tick(tracks, minecarts)
		if len(newMinecarts) == 1 {
			return newMinecarts[0].pos
		}
		minecarts = newMinecarts
	}
}

func tick(tracks map[coord]track, minecarts []minecart) []minecart {
	collided := map[int]struct{}{}
	newMinecarts := make([]minecart, 0, len(minecarts))
Minecarts:
	for i, m := range minecarts {
		// minecart has already collided
		if _, ok := collided[i]; ok {
			continue
		}

		// move the minecart
		p := newPos(m.pos, m.heading)

		// check collision
		for _, other := range newMinecarts {
			if other.pos == p {
				newMinecarts = removeFromSorted(newMinecarts, other)
				continue Minecarts
			}
		}
		if i != len(minecarts)-1 {
			for j, other := range minecarts[i+1:] {
				if _, ok := collided[j+i+1]; ok {
					continue
				}
				if other.pos == p {
					collided[j+i+1] = struct{}{}
					continue Minecarts
				}
			}
		}

		// set new heading/switch based on track
		newHeading, newSwitch := direction(tracks[p], m.heading, m.nextSwitch)

		newMinecart := minecart{p, newHeading, newSwitch}
		newMinecarts = addToSorted(newMinecarts, newMinecart)
	}
	return newMinecarts
}

func main() {
	tracks, minecarts := parse("day13.input")

	out1 := part1(tracks, minecarts)
	fmt.Printf("Part 1: %d,%d\n", out1.x, out1.y)

	out2 := part2(tracks, minecarts)
	fmt.Printf("Part 1: %d,%d\n", out2.x, out2.y)
}
