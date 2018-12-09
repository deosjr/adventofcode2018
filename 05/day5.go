package main

import (
	"fmt"
	"io/ioutil"
	"math"
)

func main() {
	input, err := ioutil.ReadFile("day5.input")
	if err != nil {
		panic(err)
	}
	last := input[0]
	ans := []byte{last}
	for _, s := range input[1:] {
		diff := math.Abs(float64(last) - float64(s))
		if diff != 32 {
			ans = append(ans, s)
			last = s
			continue
		}
		ans = ans[:len(ans)-1]
		last = ans[len(ans)-1]
	}
	fmt.Printf("Part 1: %d\n", len(ans))

	shortest := 999999
	var c byte
	for c = 65; c <= 90; c++ {
		last := input[0]
		ans = []byte{}
		if last != c && last != c+32 {
			ans = append(ans, last)
		}
		for _, s := range input[1:] {
			if s == c || s == c+32 {
				continue
			}
			diff := math.Abs(float64(last) - float64(s))
			if diff != 32 {
				ans = append(ans, s)
				last = s
				continue
			}
			ans = ans[:len(ans)-1]
			last = ans[len(ans)-1]
		}
		if len(ans) < shortest {
			shortest = len(ans)
		}
	}

	fmt.Printf("Part 2: %d\n", shortest)
}
