package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	input, err := ioutil.ReadFile("day2.input")
	if err != nil {
		panic(err)
	}
	inputList := strings.Split(string(input), "\n")
	twos := 0
	threes := 0

	for _, s := range inputList {
		m := map[rune]int{}
		for _, r := range s {
			m[r] += 1
		}
		var setTwos, setThrees bool
		for _, v := range m {
			if !setTwos && v == 2 {
				twos += 1
				setTwos = true
			}
			if !setThrees && v == 3 {
				threes += 1
				setThrees = true
			}
		}
	}
	fmt.Printf("Part 1: %d\n", twos*threes)

	var id1, id2 string

	for j, s := range inputList[:len(inputList)] {
		for _, ss := range inputList[j+1:] {
			diff := 0
			for i := 0; i < len(s); i++ {
				if s[i] != ss[i] {
					diff += 1
				}
			}
			if diff <= 1 {
				id1, id2 = s, ss
			}
		}
	}

	s := ""
	for i := 0; i < len(id1); i++ {
		if id1[i] == id2[i] {
			s += string(id1[i])
		}
	}

	fmt.Printf("Part 2: %s\n", s)
}
