package main

import (
	"math"
	"strings"
	"testing"
)

func TestFlow(t *testing.T) {
	for i, tt := range []struct {
		input      string
		wantString string
		want       int
	}{
		{
			input: `......+.......
					............#.
					.#..#.......#.
					.#..#..#......
					.#..#..#......
					.#.....#......
					.#.....#......
					.#######......
					..............
					..............
					....#.....#...
					....#.....#...
					....#.....#...
					....#######...`,
			wantString: `......|.....#.
						.#..#||||...#.
						.#..#~~#|.....
						.#..#~~#|.....
						.#~~~~~#|.....
						.#~~~~~#|.....
						.#######|.....
						........|.....
						...|||||||||..
						...|#~~~~~#|..
						...|#~~~~~#|..
						...|#~~~~~#|..
						...|#######|..`,
			want: 57,
		},
		{
			input: `......+.......
					............#.
					............#.
					...#.....#..#.
					...#.###.#....
					...#.#.#.#....
					...#.###.#....
					...#.....#....
					...#######....
					..............
					..............
					.#..........#.
					.#..........#.
					.#..........#.
					.############.`,
			wantString: `......|.....#.
						..|||||||||.#.
						..|#~~~~~#|.#.
						..|#~###~#|...
						..|#~#.#~#|...
						..|#~###~#|...
						..|#~~~~~#|...
						..|#######|...
						..|.......|...
						||||||||||||||
						|#~~~~~~~~~~#|
						|#~~~~~~~~~~#|
						|#~~~~~~~~~~#|
						|############|`,
			want: 92,
		},
	} {
		s, spring := testParse(tt.input)
		got := flow(s, spring)
		gotString := s.PrintSelf()
		if got != tt.want {
			t.Errorf("%d): got %d want %d\n", i, got, tt.want)
		}
		str := strings.Replace(tt.wantString, "\t", "", -1)
		if gotString != str {
			t.Errorf("%d): got\n%s want\n%s\n", i, gotString, str)
		}
	}
}

func testParse(input string) (slice, coord) {
	var spring coord
	m := map[coord]square{}
	strip := strings.Replace(input, "\t", "", -1)
	split := strings.Split(strip, "\n")

	xMin, xMax := math.MaxInt64, math.MinInt64
	yMin, yMax := math.MaxInt64, math.MinInt64
	for y := 0; y < len(split); y++ {
		for x := 0; x < len(split[0]); x++ {
			if split[y][x] == '#' {
				m[coord{x, y}] = clay
				if x < xMin {
					xMin = x
				}
				if x > xMax {
					xMax = x
				}
				if y < yMin {
					yMin = y
				}
				if y > yMax {
					yMax = y
				}
				continue
			}
			if split[y][x] == '+' {
				spring = coord{x, y}
			}
		}
	}
	s := slice{
		squares: m,
		xMin:    xMin,
		xMax:    xMax,
		yMin:    yMin,
		yMax:    yMax,
	}
	return s, spring
}
