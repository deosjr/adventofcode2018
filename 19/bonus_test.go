package main

import (
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	for i, tt := range []struct {
		instr instruction
		line  int
		ptr   int
		want  intermediateRepresentation
	}{
		{
			instr: instruction{addr{}, 1, 1, 1},
			line:  1,
			ptr:   4,
			want: assignment{
				assignee: r{1},
				value:    op{"+", r{1}, r{1}},
			},
		},
		{
			instr: instruction{seti{}, 1, 1, 4},
			line:  15,
			ptr:   4,
			want:  goTo{lineNumber: 2},
		},
		{
			instr: instruction{addi{}, 4, 1, 4},
			line:  15,
			ptr:   4,
			want:  goTo{lineNumber: 17},
		},
	} {
		got := tt.instr.opcode.translate(tt.instr, tt.line, tt.ptr)
		if got != tt.want {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}

func TestIfElsePattern(t *testing.T) {
	for i, tt := range []struct {
		list []intermediateRepresentation
		want []intermediateRepresentation
		ptr  int
	}{
		{
			list: []intermediateRepresentation{
				comparison{cond: op{operator: "==", x: r{}, y: r{}}, r: r{}},
				assignmentOp{operator: "+", assignee: r{}, value: r{}},
				goTo{lineNumber: 10},
			},
			want: []intermediateRepresentation{
				ifGotos{cond: op{operator: "==", x: r{}, y: r{}}, ifg: goTo{1}, elseg: goTo{8}},
			},
			ptr: 0,
		},
		{
			list: []intermediateRepresentation{
				goTo{lineNumber: 2},
				goTo{lineNumber: 12},
				assignmentOp{operator: "+", assignee: r{}, value: r{}},
				comparison{cond: op{operator: "==", x: r{}, y: r{}}, r: r{}},
				assignmentOp{operator: "+", assignee: r{}, value: r{}},
				goTo{lineNumber: 1},
				goTo{lineNumber: 15},
			},
			want: []intermediateRepresentation{
				goTo{lineNumber: 2},
				goTo{lineNumber: 10},
				assignmentOp{operator: "+", assignee: r{}, value: r{}},
				ifGotos{cond: op{operator: "==", x: r{}, y: r{}}, ifg: goTo{4}, elseg: goTo{1}},
				goTo{lineNumber: 13},
			},
			ptr: 0,
		},
	} {
		got, _ := ifelsePattern(tt.list, tt.ptr)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got\n%s --want \n%s", i, print(got), print(tt.want))
		}
	}
}

func TestSingleIf(t *testing.T) {
	for i, tt := range []struct {
		list []intermediateRepresentation
		want []intermediateRepresentation
		ptr  int
	}{
		{
			list: []intermediateRepresentation{
				ifGotos{cond: op{operator: "==", x: r{}, y: r{}}, ifg: goTo{1}, elseg: goTo{2}},
				assignmentOp{operator: "+", assignee: r{1}, value: v{5}},
				assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
			},
			want: []intermediateRepresentation{
				ifstatement{cond: op{operator: "==", x: r{}, y: r{}}, irs: []intermediateRepresentation{
					assignmentOp{operator: "+", assignee: r{1}, value: v{5}},
				}},
				assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
			},
			ptr: 0,
		},
		{
			list: []intermediateRepresentation{
				goTo{lineNumber: 5},
				ifGotos{cond: op{operator: "==", x: r{}, y: r{}}, ifg: goTo{2}, elseg: goTo{3}},
				assignmentOp{operator: "+", assignee: r{1}, value: v{5}},
				assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
				goTo{lineNumber: 1},
				goTo{lineNumber: 10},
			},
			want: []intermediateRepresentation{
				goTo{lineNumber: 4},
				ifstatement{cond: op{operator: "==", x: r{}, y: r{}}, irs: []intermediateRepresentation{
					assignmentOp{operator: "+", assignee: r{1}, value: v{5}},
				}},
				assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
				goTo{lineNumber: 1},
				goTo{lineNumber: 9},
			},
			ptr: 0,
		},
	} {
		got, _ := singleIf(tt.list, tt.ptr)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got\n%s --want \n%s", i, print(got), print(tt.want))
		}
	}
}

func TestDoWhile(t *testing.T) {
	for i, tt := range []struct {
		list []intermediateRepresentation
		want []intermediateRepresentation
		ptr  int
	}{
		{
			list: []intermediateRepresentation{
				assignmentOp{operator: "+", assignee: r{1}, value: v{5}},
				assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
				ifGotos{cond: op{operator: "==", x: r{}, y: r{}}, ifg: goTo{3}, elseg: goTo{0}},
			},
			want: []intermediateRepresentation{
				doWhile{
					cond: op{operator: "==", x: r{}, y: r{}},
					irs: []intermediateRepresentation{
						assignmentOp{operator: "+", assignee: r{1}, value: v{5}},
						assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
					},
				},
			},
			ptr: 0,
		},
		{
			list: []intermediateRepresentation{
				goTo{2}, // goto to inside dowhile is a problem!!
				assignmentOp{operator: "+", assignee: r{1}, value: v{5}},
				assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
				ifGotos{cond: op{operator: "==", x: r{}, y: r{}}, ifg: goTo{4}, elseg: goTo{1}},
				goTo{10},
			},
			want: []intermediateRepresentation{
				doWhile{
					cond: op{operator: "==", x: r{}, y: r{}},
					irs: []intermediateRepresentation{
						assignmentOp{operator: "+", assignee: r{1}, value: v{5}},
						assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
					},
				},
			},
			ptr: 0,
		},
	} {
		got, _ := dowhile(tt.list, tt.ptr)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got\n%s --want \n%s", i, print(got), print(tt.want))
		}
	}
}

func TestForLoop(t *testing.T) {
	for i, tt := range []struct {
		list []intermediateRepresentation
		want []intermediateRepresentation
		ptr  int
	}{
		{
			list: []intermediateRepresentation{
				assignment{assignee: r{1}, value: v{5}},
				doWhile{
					cond: op{operator: "==", x: r{}, y: r{}},
					irs: []intermediateRepresentation{
						assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
						assignmentOp{operator: "+", assignee: r{1}, value: v{1}},
					},
				},
			},
			want: []intermediateRepresentation{
				forLoop{
					loopvar: r{1},
					init:    v{5},
					cond:    op{operator: "==", x: r{}, y: r{}},
					irs: []intermediateRepresentation{
						assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
					},
				},
			},
			ptr: 0,
		},
	} {
		got, _ := forloop(tt.list, tt.ptr)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got\n%s --want \n%s", i, print(got), print(tt.want))
		}
	}
}

func TestRewritePointer(t *testing.T) {
	for i, tt := range []struct {
		list []intermediateRepresentation
		want []intermediateRepresentation
		ptr  int
	}{
		{
			list: []intermediateRepresentation{
				goTo{lineNumber: 1},
				goTo{lineNumber: 9},
				assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
				assignmentOp{operator: "+", assignee: r{4}, value: r{3}},
				goTo{lineNumber: 1},
				goTo{lineNumber: 9},
			},
			want: []intermediateRepresentation{
				goTo{lineNumber: 1},
				goTo{lineNumber: 10},
				assignmentOp{operator: "+", assignee: r{2}, value: v{6}},
				comparison{cond: op{operator: "==", x: r{3}, y: v{1}}, r: r{3}},
				assignmentOp{operator: "+", assignee: r{4}, value: r{3}},
				goTo{lineNumber: 1},
				goTo{lineNumber: 10},
			},
			ptr: 4,
		},
	} {
		got := rewritePointer(tt.list, tt.ptr)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got\n%s --want \n%s", i, print(got), print(tt.want))
		}
	}
}
