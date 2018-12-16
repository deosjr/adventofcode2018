package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type register [4]int

type opcode func(register, int, int, int) register

func addr(r register, a, b, c int) register {
	r[c] = r[a] + r[b]
	return r
}

func addi(r register, a, b, c int) register {
	r[c] = r[a] + b
	return r
}

func mulr(r register, a, b, c int) register {
	r[c] = r[a] * r[b]
	return r
}

func muli(r register, a, b, c int) register {
	r[c] = r[a] * b
	return r
}

func banr(r register, a, b, c int) register {
	r[c] = r[a] & r[b]
	return r
}

func bani(r register, a, b, c int) register {
	r[c] = r[a] & b
	return r
}

func borr(r register, a, b, c int) register {
	r[c] = r[a] | r[b]
	return r
}

func bori(r register, a, b, c int) register {
	r[c] = r[a] | b
	return r
}

func setr(r register, a, b, c int) register {
	r[c] = a
	return r
}

func seti(r register, a, b, c int) register {
	r[c] = r[a] * b
	return r
}

func gtir(r register, a, b, c int) register {
	if a > r[b] {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

func gtri(r register, a, b, c int) register {
	if r[a] > b {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

func gtrr(r register, a, b, c int) register {
	if r[a] > r[b] {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

func eqir(r register, a, b, c int) register {
	if a == r[b] {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

func eqri(r register, a, b, c int) register {
	if r[a] == b {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

func eqrr(r register, a, b, c int) register {
	if r[a] == r[b] {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

type sample struct {
	before          register
	opcode, a, b, c int
	after           register
}

func parseSamples(input string) []sample {
	samples := []sample{}
	for _, lines := range strings.Split(input, "\n\n") {
		var b0, b1, b2, b3 int
		var i0, i1, i2, i3 int
		var a0, a1, a2, a3 int
		fmt.Sscanf(lines, `Before: [%d, %d, %d, %d]
%d %d %d %d
After:  [%d, %d, %d, %d]`, &b0, &b1, &b2, &b3, &i0, &i1, &i2, &i3, &a0, &a1, &a2, &a3)
		s := sample{
			before: [4]int{b0, b1, b2, b3},
			opcode: i0,
			a:      i1,
			b:      i2,
			c:      i3,
			after:  [4]int{a0, a1, a2, a3},
		}
		samples = append(samples, s)
	}
	return samples
}

func part1(samples []sample) int {
	opcodes := []opcode{addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}
	amount := 0
	for _, s := range samples {
		behavesLike := 0
		for _, f := range opcodes {
			after := f(s.before, s.a, s.b, s.c)
			if after == s.after {
				behavesLike++
			}
			if behavesLike >= 3 {
				amount++
				break
			}
		}
	}
	return amount
}

func main() {
	input, err := ioutil.ReadFile("day16.input")
	if err != nil {
		panic(err)
	}
	twoParts := strings.Split(string(input), "\n\n\n\n")
	part1input, _ := twoParts[0], twoParts[1]
	samples := parseSamples(part1input)
	out1 := part1(samples)
	fmt.Printf("Part 1: %d\n", out1)
}
