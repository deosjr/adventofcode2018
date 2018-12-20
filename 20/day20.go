package main

import (
	"fmt"
	"io/ioutil"
)

type coord struct {
	x, y int
}

type room struct {
	N, W, S, E *room
	distance   int
}

func (r *room) update(d int) {
	if r.distance <= d+1 {
		return
	}
	r.distance = d + 1
	if r.N != nil {
		r.N.update(d + 1)
	}
	if r.W != nil {
		r.W.update(d + 1)
	}
	if r.S != nil {
		r.S.update(d + 1)
	}
	if r.E != nil {
		r.E.update(d + 1)
	}
}

func part1(input string) int {
	currentPos := coord{0, 0}
	currentRoom := &room{distance: 0}
	rooms := map[coord]*room{
		currentPos: currentRoom,
	}
	stack := []coord{currentPos}
	for _, c := range input[1 : len(input)-1] {
		switch c {
		case 'N':
			newPos := coord{currentPos.x, currentPos.y - 1}
			r := newRoom(rooms, currentRoom, newPos)
			r.S = currentRoom
			currentRoom.N = r
			currentRoom = r
			currentPos = newPos
		case 'W':
			newPos := coord{currentPos.x - 1, currentPos.y}
			r := newRoom(rooms, currentRoom, newPos)
			r.E = currentRoom
			currentRoom.W = r
			currentRoom = r
			currentPos = newPos
		case 'S':
			newPos := coord{currentPos.x, currentPos.y + 1}
			r := newRoom(rooms, currentRoom, newPos)
			r.N = currentRoom
			currentRoom.S = r
			currentRoom = r
			currentPos = newPos
		case 'E':
			newPos := coord{currentPos.x + 1, currentPos.y}
			r := newRoom(rooms, currentRoom, newPos)
			r.W = currentRoom
			currentRoom.E = r
			currentRoom = r
			currentPos = newPos
		case '(':
			stack = append(stack, currentPos)
		case ')':
			stack = stack[:len(stack)-1]
		case '|':
			currentPos = stack[len(stack)-1]
			currentRoom = rooms[currentPos]
		}
	}
	max := 0
	for _, r := range rooms {
		if r.distance > max {
			max = r.distance
		}
	}
	return max
}

func newRoom(rooms map[coord]*room, currentRoom *room, newPos coord) *room {
	r, ok := rooms[newPos]
	if ok {
		currentRoom.update(r.distance)
		r.update(currentRoom.distance)
	}
	if !ok {
		r = &room{
			distance: currentRoom.distance + 1,
		}
		rooms[newPos] = r
	}
	return r
}

func main() {
	input, err := ioutil.ReadFile("day20.input")
	if err != nil {
		panic(err)
	}
	out := part1(string(input))
	fmt.Printf("Part 1: %d\n", out)
}
