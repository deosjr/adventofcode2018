package main

import (
	"fmt"
	"math"
)

func part1(after int) int {
	recipes := map[int]int{}
	recipes[0] = 3
	recipes[1] = 7

	elf1 := 0
	elf2 := 1
	for len(recipes) < after+10 {
		sum := recipes[elf1] + recipes[elf2]
		if sum >= 10 {
			recipes[len(recipes)] = sum / 10
			sum = sum % 10
		}
		recipes[len(recipes)] = sum

		elf1 = (elf1 + 1 + recipes[elf1]) % len(recipes)
		elf2 = (elf2 + 1 + recipes[elf2]) % len(recipes)
	}
	ans := 0
	for i := after; i < after+10; i++ {
		ans += int(float64(recipes[i]) * math.Pow(10, float64(after+9-i)))
	}
	return ans
}

func part2(input int) int {
	recipes := map[int]int{
		0: 3, 1: 7, 2: 1, 3: 0, 4: 1, 5: 0,
	}
	s := 371010
	i := 0

	elf1 := 4
	elf2 := 3
	for {
		sum := recipes[elf1] + recipes[elf2]
		if sum >= 10 {
			recipes[len(recipes)] = sum / 10
			i++
			snew, end := shift(s, sum/10, input)
			if end {
				break
			}
			s = snew
			sum = sum % 10
		}
		i++
		recipes[len(recipes)] = sum
		snew, end := shift(s, sum, input)
		if end {
			break
		}
		s = snew

		elf1 = (elf1 + 1 + recipes[elf1]) % len(recipes)
		elf2 = (elf2 + 1 + recipes[elf2]) % len(recipes)
	}
	return i
}

func shift(s, current, input int) (int, bool) {
	s = s % 100000
	s = s * 10
	s += current
	if s == input {
		return s, true
	}
	return s, false
}

func main() {
	input := 147061
	out1 := part1(input)
	fmt.Printf("Part 1: %d\n", out1)
	out2 := part2(input)
	fmt.Printf("Part 2: %d\n", out2)
}
