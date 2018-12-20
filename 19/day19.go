package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type register [6]int

type opcode interface {
	do(register, int, int, int) register
	// used for bonus
	translate(instr instruction, lineNumber int, instrPtr int) intermediateRepresentation
}

type addr struct{}

func (addr) do(r register, a, b, c int) register {
	r[c] = r[a] + r[b]
	return r
}

type addi struct{}

func (addi) do(r register, a, b, c int) register {
	r[c] = r[a] + b
	return r
}

type mulr struct{}

func (mulr) do(r register, a, b, c int) register {
	r[c] = r[a] * r[b]
	return r
}

type muli struct{}

func (muli) do(r register, a, b, c int) register {
	r[c] = r[a] * b
	return r
}

type banr struct{}

func (banr) do(r register, a, b, c int) register {
	r[c] = r[a] & r[b]
	return r
}

type bani struct{}

func (bani) do(r register, a, b, c int) register {
	r[c] = r[a] & b
	return r
}

type borr struct{}

func (borr) do(r register, a, b, c int) register {
	r[c] = r[a] | r[b]
	return r
}

type bori struct{}

func (bori) do(r register, a, b, c int) register {
	r[c] = r[a] | b
	return r
}

type setr struct{}

func (setr) do(r register, a, b, c int) register {
	r[c] = r[a]
	return r
}

type seti struct{}

func (seti) do(r register, a, b, c int) register {
	r[c] = a
	return r
}

type gtir struct{}

func (gtir) do(r register, a, b, c int) register {
	if a > r[b] {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

type gtri struct{}

func (gtri) do(r register, a, b, c int) register {
	if r[a] > b {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

type gtrr struct{}

func (gtrr) do(r register, a, b, c int) register {
	if r[a] > r[b] {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

type eqir struct{}

func (eqir) do(r register, a, b, c int) register {
	if a == r[b] {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

type eqri struct{}

func (eqri) do(r register, a, b, c int) register {
	if r[a] == b {
		r[c] = 1
		return r
	}
	r[c] = 0
	return r
}

type eqrr struct{}

func (eqrr) do(r register, a, b, c int) register {
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

func (i instruction) execute(r register) register {
	return i.opcode.do(r, i.a, i.b, i.c)
}

var opcodes = map[string]opcode{"addr": addr{}, "addi": addi{}, "mulr": mulr{}, "muli": muli{},
	"banr": banr{}, "bani": bani{}, "borr": borr{}, "bori": bori{}, "setr": setr{}, "seti": seti{},
	"gtir": gtir{}, "gtri": gtri{}, "gtrr": gtrr{}, "eqir": eqir{}, "eqri": eqri{}, "eqrr": eqrr{}}

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
		r = ins.execute(r)
		r[insPtr] = r[insPtr] + 1
	}
	return r[0]
}

// running part2 like part 1 will take an insane amount of time
// it spends most of its time in a loop
// this is the loop part optimised and written in Go
// (see bonus for a code version)
func part2(r2 int) int {
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
	fmt.Printf("Part 2: %d\n", part2(10551376))
	fmt.Printf("Bonus: %s\n", bonus(insPtr, program))
}
