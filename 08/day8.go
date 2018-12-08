package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func part1(input []int) (output []int, answer int) {
	numChildren := input[0]
	numMetadata := input[1]
	tail := input[2:]
	for i := 0; i < numChildren; i++ {
		t, ans := part1(tail)
		tail = t
		answer += ans
	}
	for i := 0; i < numMetadata; i++ {
		answer += tail[i]
	}
	tail = tail[numMetadata:]
	return tail, answer
}

func part2(input []int) (output []int, answer int) {
	numChildren := input[0]
	numMetadata := input[1]
	tail := input[2:]
	childValues := map[int]int{}
	for i := 0; i < numChildren; i++ {
		t, ans := part2(tail)
		tail = t
		childValues[i] = ans
	}
	for i := 0; i < numMetadata; i++ {
		m := tail[i]
		if numChildren == 0 {
			answer += m
			continue
		}
		if m > numChildren {
			continue
		}
		answer += childValues[m-1]
	}
	tail = tail[numMetadata:]
	return tail, answer
}

func main() {
	input, err := ioutil.ReadFile("day8.input")
	if err != nil {
		panic(err)
	}
	splitinput := strings.Split(string(input), " ")
	inputInts := make([]int, len(splitinput))
	for i, s := range splitinput {
		parsed, _ := strconv.Atoi(s)
		inputInts[i] = parsed
	}

	_, answer1 := part1(inputInts)
	fmt.Printf("Part 1: %d\n", answer1)

	_, answer2 := part2(inputInts)
	fmt.Printf("Part 2: %d\n", answer2)
}
