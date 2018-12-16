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
	r[c] = r[a]
	return r
}

func seti(r register, a, b, c int) register {
	r[c] = a
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

type instruction struct {
	opcode  int
	a, b, c int
}

type sample struct {
	before      register
	instruction instruction
	after       register
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
			before:      [4]int{b0, b1, b2, b3},
			instruction: instruction{i0, i1, i2, i3},
			after:       [4]int{a0, a1, a2, a3},
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
			after := f(s.before, s.instruction.a, s.instruction.b, s.instruction.c)
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

func parseProgram(input string) []instruction {
	split := strings.Split(input, "\n")
	program := make([]instruction, len(split))
	for i, s := range split {
		var opcode, a, b, c int
		fmt.Sscanf(s, "%d %d %d %d", &opcode, &a, &b, &c)
		program[i] = instruction{opcode, a, b, c}
	}
	return program
}

func determineOpcodes(samples []sample) map[int]opcode {
	opcodes := []opcode{addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}
	found := map[int]struct{}{}
	mapping := map[int]opcode{}
	for len(mapping) < len(opcodes) {
	Samples:
		for _, s := range samples {
			if _, ok := mapping[s.instruction.opcode]; ok {
				continue
			}
			behavesLike := -1
			for i, f := range opcodes {
				if _, ok := found[i]; ok {
					continue
				}
				after := f(s.before, s.instruction.a, s.instruction.b, s.instruction.c)
				if after == s.after {
					if behavesLike != -1 {
						continue Samples
					}
					behavesLike = i
				}
			}
			mapping[s.instruction.opcode] = opcodes[behavesLike]
			found[behavesLike] = struct{}{}
		}
	}
	return mapping
}

func part2(samples []sample, program []instruction) int {
	opcodes := determineOpcodes(samples)
	r := [4]int{0, 0, 0, 0}
	for _, ins := range program {
		f := opcodes[ins.opcode]
		r = f(r, ins.a, ins.b, ins.c)
	}
	return r[0]
}

func main() {
	input, err := ioutil.ReadFile("day16.input")
	if err != nil {
		panic(err)
	}
	twoParts := strings.Split(string(input), "\n\n\n\n")
	part1input, part2input := twoParts[0], twoParts[1]
	samples := parseSamples(part1input)
	fmt.Printf("Part 1: %d\n", part1(samples))
	program := parseProgram(part2input)
	fmt.Printf("Part 2: %d\n", part2(samples, program))
}
