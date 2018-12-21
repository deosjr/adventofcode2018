package main

import "fmt"

// this is the input code after compiling elfcode to go
// using the (half-finished) compiler from day 19 bonus

/*
0: r3 = 123
1: do {
	r3 &= 456
} while !(r3 == 72)
2: r3 = 0
3: r2 = r3 | 65536
4: r3 = 832312
5: r1 = r2 & 255
6: r3 += r1
7: r3 &= 16777215
8: r3 *= 65899
9: r3 &= 16777215
10: if 256 > r2 {
goto 19
}
11: r1 = 0
12: r4 = r1 + 1
13: r4 *= 256
14: if r4 > r2 {
goto 17
}
15: r1 += 1
16: goto 12
17: r2 = r1
18: goto 5
19: if r3 == r0 {
	goto 20
} else {
	goto 3
}
*/

// rewritten is the above as a function rewritten by hand in multiple steps
// line 0 and line 1 are discarded as they are just the bitwise-and check
// Basically the program generates a looping sequence of numbers for r3
// part 1 is get the first, part 2 is get the last number of the sequence
// this function is not actually used in the answer, see inner()
func rewritten(r0 int) {
	var r1, r2, r3, r4 int
	for {
		r2 = r3 | 65536
		r3 = 832312
		for {
			r1 = r2 & 255
			r3 += r1
			r3 &= 16777215
			r3 *= 65899
			r3 &= 16777215
			if 256 > r2 {
				break
			}
			r1 = 0
		L12:
			r4 = r1 + 1
			r4 *= 256
			if !(r4 > r2) {
				r1 += 1
				goto L12
			}
			r2 = r1
		}
		if r3 == r0 {
			return
		}
	}
}

// the inner loop of func rewritten
func inner(r0, r1, r2, r3, r4 int) (int, int, int, int, int) {
	for {
		r1 = r2 & 255
		r3 += r1
		r3 &= 16777215
		r3 *= 65899
		r3 &= 16777215
		if 256 > r2 {
			return r0, r1, r2, r3, r4
		}
		r2 = r2 / 256
	}
}

func part1() int {
	var r0, r1, r2, r3, r4 int
	r2 = r3 | 65536
	r3 = 832312
	r0, r1, r2, r3, r4 = inner(r0, r1, r2, r3, r4)
	return r3
}

func part2() int {
	answers := map[int]struct{}{}
	lastAnswer := 0
	var r0, r1, r2, r3, r4 int
	for {
		r2 = r3 | 65536
		r3 = 832312
		r0, r1, r2, r3, r4 = inner(r0, r1, r2, r3, r4)
		if _, ok := answers[r3]; ok {
			return lastAnswer
		}
		answers[r3] = struct{}{}
		lastAnswer = r3
	}
}

func main() {
	fmt.Printf("Part 1: %d\n", part1())
	fmt.Printf("Part 2: %d\n", part2())
}
