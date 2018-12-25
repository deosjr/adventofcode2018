package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type point struct {
	a, b, c, d int
}

func manhattan(p, q point) int {
	da := float64(p.a - q.a)
	db := float64(p.b - q.b)
	dc := float64(p.c - q.c)
	dd := float64(p.d - q.d)
	return int(math.Abs(da) + math.Abs(db) + math.Abs(dc) + math.Abs(dd))
}

func main() {
	input, err := ioutil.ReadFile("day25.input")
	if err != nil {
		panic(err)
	}
	split := strings.Split(string(input), "\n")
	constellations := [][]point{}

	for _, s := range split {
		var a, b, c, d int
		fmt.Sscanf(s, "%d,%d,%d,%d", &a, &b, &c, &d)
		p := point{a, b, c, d}

		matched := map[int]struct{}{}
	Constellations:
		for i, c := range constellations {
			for _, cc := range c {
				if manhattan(p, cc) <= 3 {
					matched[i] = struct{}{}
					continue Constellations
				}
			}
		}
		newCons := [][]point{}
		for i, c := range constellations {
			if _, ok := matched[i]; ok {
				continue
			}
			newCons = append(newCons, c)
		}
		if len(matched) == 0 {
			constellations = append(constellations, []point{p})
			continue
		}
		combined := []point{p}
		for i, _ := range matched {
			combined = append(combined, constellations[i]...)
		}
		if len(combined) > 0 {
			newCons = append(newCons, combined)
		}
		constellations = newCons
	}
	fmt.Printf("Part 1: %d\n", len(constellations))
}
