:- module day16.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, int.
:- import_module list, map.

%% I can't define a fixed length immediately, but I can do it iteratively..
:- inst register0 ---> [].
:- inst register1 ---> [ground|register0].
:- inst register2 ---> [ground|register1].
:- inst register3 ---> [ground|register2].
:- inst register4 ---> [ground|register3].

:- inst sample ---> {register4,ground,register4}.
:- inst sample_list ---> [] ; [sample|sample_list].

:- type register == list(int).
:- type instruction ---> {int,int,int,int}.
:- type sample ---> {register,instruction,register}.

:- type opcode == pred(register, int, int, int, register).

%% NOTE: position index for det_replace_nth starts at 1
:- pred addr(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred addi(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred mulr(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred muli(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred banr(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred bani(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred borr(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred bori(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred setr(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred seti(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred gtir(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred gtri(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred gtrr(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred eqir(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred eqri(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.
:- pred eqrr(register::in(register4), int::in, int::in, int::in, register::out(register4)) is det.

addr(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	RB = det_index0(RIn, B),
	ROut = det_replace_nth(RIn, C+1, RA+RB).

addi(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	ROut = det_replace_nth(RIn, C+1, RA+B).

mulr(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	RB = det_index0(RIn, B),
	ROut = det_replace_nth(RIn, C+1, RA*RB).

muli(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	ROut = det_replace_nth(RIn, C+1, RA*B).

banr(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	RB = det_index0(RIn, B),
	ROut = det_replace_nth(RIn, C+1, RA/\RB).

bani(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	ROut = det_replace_nth(RIn, C+1, RA/\B).

borr(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	RB = det_index0(RIn, B),
	ROut = det_replace_nth(RIn, C+1, RA\/RB).

bori(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	ROut = det_replace_nth(RIn, C+1, RA\/B).

setr(RIn, A, _, C, ROut) :-
	RA = det_index0(RIn, A),
	ROut = det_replace_nth(RIn, C+1, RA).

seti(RIn, A, _, C, ROut) :-
	ROut = det_replace_nth(RIn, C+1, A).

gtir(RIn, A, B, C, ROut) :-
	RB = det_index0(RIn, B),
	(if A > RB then
		ROut = det_replace_nth(RIn, C+1, 1)
	else
		ROut = det_replace_nth(RIn, C+1, 0)
	).

gtri(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	(if RA > B then
		ROut = det_replace_nth(RIn, C+1, 1)
	else
		ROut = det_replace_nth(RIn, C+1, 0)
	).

gtrr(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	RB = det_index0(RIn, B),
	(if RA > RB then
		ROut = det_replace_nth(RIn, C+1, 1)
	else
		ROut = det_replace_nth(RIn, C+1, 0)
	).

eqir(RIn, A, B, C, ROut) :-
	RB = det_index0(RIn, B),
	(if A = RB then
		ROut = det_replace_nth(RIn, C+1, 1)
	else
		ROut = det_replace_nth(RIn, C+1, 0)
	).

eqri(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	(if RA = B then
		ROut = det_replace_nth(RIn, C+1, 1)
	else
		ROut = det_replace_nth(RIn, C+1, 0)
	).

eqrr(RIn, A, B, C, ROut) :-
	RA = det_index0(RIn, A),
	RB = det_index0(RIn, B),
	%% NOTE: this + 0 forces RB to typecheck to int
	%% how otherwise to resolve ambiguity between pred(list,int) and int?
	(if RA = RB + 0 then
		ROut = det_replace_nth(RIn, C+1, 1)
	else
		ROut = det_replace_nth(RIn, C+1, 0)
	).

:- pred parse(string::in, list(sample)::out(sample_list), list(instruction)::out) is semidet.
:- pred parseSamples(list(string)::in, list(sample)::out(sample_list)) is semidet.
:- pred parseProgram(list(string)::in, list(instruction)::out) is semidet.
:- pred parseRegister(string::in, register::out(register4)) is semidet.
:- pred parseInstruction(string::in, instruction::out) is semidet.

parse(Input, Samples, Program) :-
	[Part1, Part2] = string.split_at_string("\n\n\n\n", Input),
	parseSamples(string.split_at_string("\n\n", Part1), Samples),
	parseProgram(string.split_at_string("\n", Part2), Program).

parseSamples([], []).
parseSamples([H|T], [Sample|Samples]) :-
	[BeforeStr, InstructionStr, AfterStr] = string.split_at_string("\n", H),
	parseRegister(BeforeStr, Before),
	parseInstruction(InstructionStr, Instruction),
	parseRegister(AfterStr, After),
	Sample = {Before, Instruction, After},
	parseSamples(T, Samples).

parseProgram([], []).
parseProgram([H|T], [Instruction|Program]) :-
	parseInstruction(H, Instruction),
	parseProgram(T, Program).

parseRegister(Line, Register) :-
	[S1, S2, S3, S4] = string.split_at_string(", ", Line),
	[_, SS1] = string.words(S1),
	X1 = string.det_to_int(string.remove_prefix_if_present("[", SS1)),
	X2 = string.det_to_int(S2),
	X3 = string.det_to_int(S3),
	X4 = string.det_to_int(string.remove_suffix_if_present("]",S4)),
	Register = [X1, X2, X3, X4].

parseInstruction(Line, Instruction) :-
	[OP, A, B, C] = string.words(Line),
	Instruction = {string.det_to_int(OP), string.det_to_int(A), string.det_to_int(B), string.det_to_int(C)}.

:- pred part1(list(sample)::in(sample_list), int::out) is det.
:- pred behaves_likeN(list(opcode), sample, int).
:- mode behaves_likeN(in(list(pred(in(register4), in, in, in, out(register4)) is det)), in(sample), out) is det.
:- pred behaves_like(sample, opcode).
:- mode behaves_like(in(sample), pred(in(register4), in, in, in, out(register4)) is det) is semidet.

part1(Samples, Sum) :-
	Opcodes = [addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr],
	map(behaves_likeN(Opcodes), Samples, Counts),
	foldl((pred(X::in, A::in, O::out) is det :- (if X >= 3 then O=A+1 else O=A)), Counts, 0, Sum).

behaves_likeN([], _, 0).
behaves_likeN([Opcode|Tail], Sample, N) :-
	behaves_likeN(Tail, Sample, NOld),
	(if behaves_like(Sample, Opcode) then
		N = NOld + 1
	else 
		N = NOld
	).

behaves_like(Sample, Opcode) :-
	Sample = {Before, {_, A, B, C}, After},
	call(Opcode, Before, A, B, C, After).

:- pred part2(list(sample)::in, list(instruction)::in, int::out) is det.

part2(Samples, Program, Out) :- Out=42.

%% redefine head predicate to enforce correct inst on sample and the registers within
:- pred head(list(sample)::in, sample::out(sample)) is semidet.
head([H|_], H).

main(!IO) :-
   io.open_input("day16.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
        ReadResult = ok(FileContents),
        (if parse(FileContents, Samples, Program) then
	        part1(Samples, Out1),
	        S1 = string.format("Part 1: %i\n", [i(Out1)]),
	        io.write_string(S1, !IO),
	        part2(Samples, Program, Out2),
	        S2 = string.format("Part 2: %i\n", [i(Out2)]),
	        io.write_string(S2, !IO)
	    else
	    	io.write_string("Parse failed\n", !IO)
	    )
      ;
        ReadResult = error(_, IO_Error),
        io.stderr_stream(Stderr, !IO),
        io.write_string(Stderr, io.error_message(IO_Error) ++ "\n", !IO)
      )        
   ;
      OpenResult = error(IO_Error),
      io.stderr_stream(Stderr, !IO),
      io.write_string(Stderr, io.error_message(IO_Error) ++ "\n", !IO)
   ).
