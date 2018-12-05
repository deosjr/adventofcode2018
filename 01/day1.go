package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input, err := ioutil.ReadFile("day1.input")
	if err != nil {
		panic(err)
	}
	inputList := strings.Split(string(input), "\n")

	i := 0
	for _, s := range inputList {
		n, _ := strconv.Atoi(s)
		i += n
	}
	fmt.Printf("Part 1: %d\n", i)

	i = 0
	m := map[int]struct{}{0: struct{}{}}
Loop:
	for {
		for _, s := range inputList {
			n, _ := strconv.Atoi(s)
			i += n
			if _, ok := m[i]; ok {
				fmt.Printf("Part 2: %d\n", i)
				break Loop
			}
			m[i] = struct{}{}
		}
	}
}
