package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type register [6]int

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
	opcode  opcode
	a, b, c int
}

var opcodes = map[string]opcode{"addr": addr, "addi": addi, "mulr": mulr, "muli": muli,
	"banr": banr, "bani": bani, "borr": borr, "bori": bori, "setr": setr, "seti": seti,
	"gtir": gtir, "gtri": gtri, "gtrr": gtrr, "eqir": eqir, "eqri": eqri, "eqrr": eqrr}

func parse(input string) (int, []instruction) {
	split := strings.Split(input, "\n")
	var insPtr int
	fmt.Sscanf(split[0], "#ip %d", &insPtr)
	program := make([]instruction, len(split[1:]))
	for i, s := range split[1:] {
		var opStr string
		var a, b, c int
		fmt.Sscanf(s, "%s %d %d %d", &opStr, &a, &b, &c)
		program[i] = instruction{opcodes[opStr], a, b, c}
	}
	return insPtr, program
}

func part1(insPtr int, program []instruction) int {
	var r register
	for r[insPtr] >= 0 && r[insPtr] < len(program) {
		ins := program[r[insPtr]]
		r = ins.opcode(r, ins.a, ins.b, ins.c)
		r[insPtr] = r[insPtr] + 1
	}
	return r[0]
}

func part2(insPtr int, program []instruction) int {
	r := [6]int{1, 0, 0, 0, 0, 0}
	for r[insPtr] >= 0 && r[insPtr] < len(program) {
		ins := program[r[insPtr]]
		r = ins.opcode(r, ins.a, ins.b, ins.c)
		r[insPtr] = r[insPtr] + 1
	}
	return r[0]
}

// running part2 will take an insane amount of time
// it spends most of its time in a loop
// this is the loop part optimised and written in Go
// (see annotated input file)
func part2_interpreted(r2 int) int {
	r0 := 0
	for r1 := 1; r1 <= r2; r1++ {
		if r2%r1 == 0 {
			r0 += r1
		}
	}
	return r0
}

func main() {
	input, err := ioutil.ReadFile("day19.input")
	if err != nil {
		panic(err)
	}
	insPtr, program := parse(string(input))
	reg0 := part1(insPtr, program)
	fmt.Printf("Part 1: %d\n", reg0)

	// reg0 = part2(insPtr, program)
	// fmt.Printf("Part 2: %d\n", reg0)
	fmt.Printf("Part 2: %d\n", part2_interpreted(10551376))
}
