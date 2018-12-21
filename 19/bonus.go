package main

import "fmt"

// bonus: do the optimisation itself in code
// - rewrite program to list of optimised instructions
// - execute optimised program
// option 1: run the intermediate representation in Go (interpreter)
// option 2: write the equivalent Go code, go build, go run (compiler)
func bonus(insPtr int, program []instruction) string {
	irlist := make([]intermediateRepresentation, len(program))
	for i, p := range program {
		irlist[i] = p.opcode.translate(p, i, insPtr)
	}
	irlist = rewriteOperators(irlist, insPtr)
	irlist = rewritePointer(irlist, insPtr)
	irlist = compile(irlist, insPtr)
	// fmt.Println(print(irlist))
	return "unfinished"
}

func compile(irlist []intermediateRepresentation, insPtr int) []intermediateRepresentation {
	// fmt.Println(print(irlist))
	// fmt.Println("========================")
	rewrites := []rewriteFunc{ifelsePattern, singleIf, forloop, dowhile}
	for {
		for _, f := range rewrites {
			newProgram, found := f(irlist, insPtr)
			if found {
				return compile(newProgram, insPtr)
			}
		}
		break
	}
	return irlist
}

func (addr) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	return assignment{
		assignee: r{instr.c},
		value:    op{"+", r{instr.a}, r{instr.b}},
	}
}

func (addi) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.a == insPtr && instr.c == insPtr {
		return goTo{lineNumber + 1 + instr.b}
	}
	return assignment{
		assignee: r{instr.c},
		value:    op{"+", r{instr.a}, v{instr.b}},
	}
}

func (mulr) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		// TODO: this is not always valid.
		// Good enough for now though
		return end{}
	}
	return assignment{
		assignee: r{instr.c},
		value:    op{"*", r{instr.a}, r{instr.b}},
	}
}

func (muli) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return assignment{
		assignee: r{instr.c},
		value:    op{"*", r{instr.a}, v{instr.b}},
	}
}

func (banr) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return assignment{
		assignee: r{instr.c},
		value:    op{"&", r{instr.a}, r{instr.b}},
	}
}

func (bani) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return assignment{
		assignee: r{instr.c},
		value:    op{"&", r{instr.a}, v{instr.b}},
	}
}

func (borr) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return assignment{
		assignee: r{instr.c},
		value:    op{"|", r{instr.a}, r{instr.b}},
	}
}

func (bori) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return assignment{
		assignee: r{instr.c},
		value:    op{"|", r{instr.a}, v{instr.b}},
	}
}

func (setr) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return assignment{
		assignee: r{instr.c},
		value:    r{instr.a},
	}
}

func (seti) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		return goTo{instr.a + 1}
	}
	return assignment{
		assignee: r{instr.c},
		value:    v{instr.a},
	}
}

func (gtir) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return comparison{
		cond: op{">", v{instr.a}, r{instr.b}},
		r:    r{instr.c},
	}
}

func (gtri) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return comparison{
		cond: op{">", r{instr.a}, v{instr.b}},
		r:    r{instr.c},
	}
}

func (gtrr) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return comparison{
		cond: op{">", r{instr.a}, r{instr.b}},
		r:    r{instr.c},
	}
}

func (eqir) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return comparison{
		cond: op{"==", v{instr.a}, r{instr.b}},
		r:    r{instr.c},
	}
}

func (eqri) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return comparison{
		cond: op{"==", r{instr.a}, v{instr.b}},
		r:    r{instr.c},
	}
}

func (eqrr) translate(instr instruction, lineNumber, insPtr int) intermediateRepresentation {
	if instr.c == insPtr {
		panic("unimplemented")
	}
	return comparison{
		cond: op{"==", r{instr.a}, r{instr.b}},
		r:    r{instr.c},
	}
}

// intermediateRepresentation interface covers simple statements
// but also more complex control structures such as for loops
type intermediateRepresentation interface {
	// execute(register) register
	print(indent int) string
}

type rOrv interface {
	print() string
}

// register
type r struct {
	i int
}

func (r r) print() string {
	return fmt.Sprintf("r%d", r.i)
}

// value
type v struct {
	i int
}

func (v v) print() string {
	return fmt.Sprintf("%d", v.i)
}

// arity-2 infix operator
type op struct {
	operator string
	x        rOrv
	y        rOrv
}

func (op op) print() string {
	return fmt.Sprintf("%s %s %s", op.x.print(), op.operator, op.y.print())
}

func indent(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "\t"
	}
	return s
}

type assignment struct {
	assignee r
	value    rOrv
}

func (a assignment) print(n int) string {
	s := indent(n)
	s += fmt.Sprintf("%s = %s", a.assignee.print(), a.value.print())
	return s
}

type assignmentOp struct {
	operator string
	assignee r
	value    rOrv
}

func (a assignmentOp) print(n int) string {
	s := indent(n)
	s += fmt.Sprintf("%s %s= %s", a.assignee.print(), a.operator, a.value.print())
	return s
}

type goTo struct {
	lineNumber int
}

func (g goTo) print(n int) string {
	return fmt.Sprintf("goto %d", g.lineNumber)
}

type end struct {
	lineNumber int
}

func (e end) print(n int) string {
	return "end"
}

type comparison struct {
	cond op
	r    r
}

func (c comparison) print(n int) string {
	s := indent(n)
	s = fmt.Sprintf("if %s {\n%s\t%s = 1\n%s} else {\n%s\t%s = 0\n%s}", c.cond.print(), s, c.r.print(), s, s, c.r.print(), s)
	return s
}

type ifGotos struct {
	cond  op
	ifg   goTo
	elseg goTo
}

func (i ifGotos) print(n int) string {
	s := indent(n)
	s = fmt.Sprintf("if %s {\n%s\t%s\n%s} else {\n%s\t%s\n%s}", i.cond.print(), s, i.ifg.print(0), s, s, i.elseg.print(0), s)
	return s
}

type ifstatement struct {
	cond op
	irs  []intermediateRepresentation
}

func (i ifstatement) print(n int) string {
	s := indent(n)
	irStr := ""
	for _, ir := range i.irs {
		irStr += ir.print(n+1) + "\n"
	}
	s = fmt.Sprintf("%sif %s {\n%s%s}", s, i.cond.print(), irStr, s)
	return s
}

type doWhile struct {
	cond op
	irs  []intermediateRepresentation
}

func (dw doWhile) print(n int) string {
	s := indent(n)
	irStr := ""
	for _, ir := range dw.irs {
		irStr += ir.print(n+1) + "\n"
	}
	s = fmt.Sprintf("%sdo {\n%s%s} while !(%s)", s, irStr, s, dw.cond.print())
	return s
}

type forLoop struct {
	loopvar r
	init    rOrv
	cond    op
	irs     []intermediateRepresentation
}

func (fl forLoop) print(n int) string {
	s := indent(n)
	irStr := ""
	for _, ir := range fl.irs {
		irStr += ir.print(n+1) + "\n"
	}
	s = fmt.Sprintf("%sfor %s=%s;%s;%s++ {\n%s%s}", s, fl.loopvar.print(), fl.init.print(), fl.cond.print(), fl.loopvar.print(), irStr, s)
	return s
}

type rewriteFunc func([]intermediateRepresentation, int) ([]intermediateRepresentation, bool)

// TODO: for the case rP += rX we need to use an assumption:
// rX at this point can only be 0 or 1
// this happens in _our_ case but will not extend to the general case at all
// However, _if_ this is the case, then addr P X P (rP += rX) 'simplifies' to:
// eqri X 1 X
// addr P X P
func rewritePointer(irlist []intermediateRepresentation, insPtr int) []intermediateRepresentation {
	updates := map[int]int{}
	newList := []intermediateRepresentation{irlist[0]}
	for line, ir := range irlist[1:] {
		a, ok := ir.(assignmentOp)
		if !ok {
			newList = append(newList, ir)
			continue
		}
		if _, ok := irlist[line].(comparison); ok {
			newList = append(newList, ir)
			continue
		}
		if a.operator != "+" {
			newList = append(newList, ir)
			continue
		}
		if a.assignee.i != insPtr {
			newList = append(newList, ir)
			continue
		}
		r, ok := a.value.(r)
		if !ok {
			newList = append(newList, ir)
			continue
		}
		newList = append(newList, eqri{}.translate(instruction{eqri{}, r.i, 1, r.i}, line, insPtr))
		newList = append(newList, assignmentOp{"+", a.assignee, r})
		updates[line] = 1
	}
	for k, v := range updates {
		newList = updateGotos(newList, k, v)
	}
	return newList
}

// TODO: assignments to 0*X or rX+0 can be simplified
func rewriteOperators(irlist []intermediateRepresentation, insPtr int) []intermediateRepresentation {
	canBeDeleted := map[int]struct{}{}
	currentValue := map[r]rOrv{}
	newList := make([]intermediateRepresentation, len(irlist))
	for i, ir := range irlist {
		a, ok := ir.(assignment)
		if !ok {
			// no assignment, we dont know how values change. start again
			currentValue = map[r]rOrv{}
			newList[i] = ir
			continue
		}
		op, ok := a.value.(op)
		if !ok {
			newList[i] = ir
			// rX = Y when rX already is equal to Y
			if a.value == currentValue[a.assignee] {
				canBeDeleted[i] = struct{}{}
			}
			currentValue[a.assignee] = a.value
			continue
		}
		if a.assignee == op.x {
			newList[i] = assignmentOp{op.operator, a.assignee, op.y}
			continue
		}
		if a.assignee == op.y {
			newList[i] = assignmentOp{op.operator, a.assignee, op.x}
			continue
		}
		newList[i] = ir
	}
	cleanList := []intermediateRepresentation{}
	for i, ir := range newList {
		if _, ok := canBeDeleted[i]; ok {
			continue
		}
		cleanList = append(cleanList, ir)
	}
	return cleanList
}

// N:   rX = comparison (so rX is either 1 or 0)
// N+1: rP += rX (pointer register +1 depending on comparison)
// N+2: goto Y
// N+3: some other statement
// ------------------------
// if comparison {
// 	goto N+3
// }
// goto Y
func ifelsePattern(irlist []intermediateRepresentation, insPtr int) ([]intermediateRepresentation, bool) {
	for line, ir := range irlist[:len(irlist)-2] {
		c, ok := ir.(comparison)
		if !ok {
			continue
		}
		a, ok := irlist[line+1].(assignmentOp)
		if !ok {
			continue
		}
		if a.operator != "+" {
			continue
		}
		if a.assignee.i != insPtr {
			continue
		}
		if c.r != a.value {
			continue
		}
		g, ok := irlist[line+2].(goTo)
		if !ok {
			continue
		}
		newIR := ifGotos{
			cond:  c.cond,
			ifg:   goTo{line + 3},
			elseg: g,
		}
		return updateList(newIR, irlist, line, 2), true
	}
	return irlist, false
}

// N:   if condition goto N+1 else goto N+2
// N+1: statement 1
// N+2: statement 2
// --------------------
// if condition {
// 	statement 1
// }
// statement 2
func singleIf(irlist []intermediateRepresentation, insPtr int) ([]intermediateRepresentation, bool) {
	for line, ir := range irlist[:len(irlist)-1] {
		igt, ok := ir.(ifGotos)
		if !ok {
			continue
		}
		if igt.ifg.lineNumber == line+1 && igt.elseg.lineNumber == line+2 {
			newIR := ifstatement{
				cond: igt.cond,
				irs:  []intermediateRepresentation{irlist[line+1]},
			}
			return updateList(newIR, irlist, line, 1), true
		}
	}
	return irlist, false
}

// ---> do - while loop:
// N:  statement
// (N+1).... only statements, no gotos
// M:  if condition { goto M+1 } else { goto N }
// ---------------------
// do (N - M-1 statements) while (not condition)
func dowhile(irlist []intermediateRepresentation, insPtr int) ([]intermediateRepresentation, bool) {
ListLoop:
	for i, ir := range irlist[1:] {
		igt, ok := ir.(ifGotos)
		if !ok {
			continue
		}
		line := i + 1
		if igt.ifg.lineNumber == line+1 && igt.elseg.lineNumber < line {
			innerIrs := []intermediateRepresentation{}
			// TODO: one problem is gotos that reach within the control struct
			// for now, dont allow this rewrite if such a goto exists in the program
			for _, irr := range irlist {
				if g, ok := irr.(goTo); ok {
					if g.lineNumber > igt.elseg.lineNumber && g.lineNumber < line+1 {
						continue ListLoop
					}
				}
			}
			for _, previr := range irlist[igt.elseg.lineNumber:line] {
				if _, ok := previr.(goTo); ok {
					continue ListLoop
				}
				if _, ok := previr.(ifGotos); ok {
					continue ListLoop
				}
				innerIrs = append(innerIrs, previr)
			}
			newIR := doWhile{
				cond: igt.cond,
				irs:  innerIrs,
			}
			return updateList(newIR, irlist, igt.elseg.lineNumber, len(innerIrs)), true
		}
	}
	return irlist, false
}

// ---> simple iterative for loop
// N:   rX = Y
// N+1: do {
// ...		statements not assigning to rX (no gotos guaranteed)
// M-1:	rX++
// M:   } while condition
// ---> go code:
// for rX = Y; condition; rX++ {
// 	statements not assigning to rX
// }
// ---------------
func forloop(irlist []intermediateRepresentation, insPtr int) ([]intermediateRepresentation, bool) {
ListLoop:
	for line, ir := range irlist[:len(irlist)-1] {
		a, ok := ir.(assignment)
		if !ok {
			continue
		}
		dw, ok := irlist[line+1].(doWhile)
		if !ok {
			continue
		}
		for _, irr := range dw.irs {
			if _, ok := irr.(goTo); ok {
				continue ListLoop
			}
			if _, ok := irr.(ifGotos); ok {
				continue ListLoop
			}
		}
		dwa, ok := dw.irs[len(dw.irs)-1].(assignmentOp)
		if !ok {
			continue
		}
		if dwa.assignee != a.assignee {
			continue
		}
		v, ok := dwa.value.(v)
		if !ok {
			continue
		}
		if v.i != 1 {
			continue
		}
		newIR := forLoop{
			loopvar: a.assignee,
			init:    a.value,
			cond:    dw.cond,
			irs:     dw.irs[:len(dw.irs)-1],
		}
		return updateList(newIR, irlist, line, 1), true
	}
	return irlist, false
}

func updateList(ir intermediateRepresentation, irlist []intermediateRepresentation, line, toDelete int) []intermediateRepresentation {
	newList := make([]intermediateRepresentation, len(irlist)-toDelete)
	for i, v := range irlist[:line] {
		newList[i] = v
	}
	newList[line] = ir
	if line == len(irlist)-1-toDelete {
		return updateGotos(newList, line, -toDelete)
	}
	for i, v := range irlist[line+1+toDelete:] {
		newList[i+line+1] = v
	}
	return updateGotos(newList, line, -toDelete)
}

func updateGotos(irlist []intermediateRepresentation, line, n int) []intermediateRepresentation {
	newList := make([]intermediateRepresentation, len(irlist))
	for i, ir := range irlist {
		g, ok := ir.(goTo)
		if ok {
			newList[i] = updateGoto(g, line, n)
			continue
		}
		ifgotos, ok := ir.(ifGotos)
		if ok {
			newList[i] = ifGotos{
				cond:  ifgotos.cond,
				ifg:   updateGoto(ifgotos.ifg, line, n),
				elseg: updateGoto(ifgotos.elseg, line, n),
			}
			continue
		}
		ifst, ok := ir.(ifstatement)
		if ok {
			newList[i] = ifstatement{
				cond: ifst.cond,
				irs:  updateGotos(ifst.irs, line, n),
			}
			continue
		}
		newList[i] = ir
	}
	return newList
}

func updateGoto(g goTo, line, n int) goTo {
	if g.lineNumber <= line {
		return g
	}
	return goTo{g.lineNumber + n}
}

func print(irlist []intermediateRepresentation) string {
	s := ""
	for i, ir := range irlist {
		s += fmt.Sprintf("%d: %s\n", i, ir.print(0))
	}
	return s
}
