:- module day16.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, int.
:- import_module list, map.

:- type register ---> {int,int,int,int}.
:- type instruction ---> {int,int,int,int}.
:- type sample ---> {register,instruction,register}.

:- pred parse(string::in, list(sample)::out, list(instruction)::out) is semidet.
:- pred parseSamples(list(string)::in, list(sample)::out) is semidet.
:- pred parseProgram(list(string)::in, list(instruction)::out) is semidet.
:- pred parseRegister(string::in, register::out) is semidet.
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
	Register = {X1, X2, X3, X4}.

parseInstruction(Line, Instruction) :-
	[OP, A, B, C] = string.words(Line),
	Instruction = {string.det_to_int(OP), string.det_to_int(A), string.det_to_int(B), string.det_to_int(C)}.

:- pred part1(list(sample)::in, int::out) is det.
:- pred part2(list(sample)::in, list(instruction)::in, int::out) is det.

part1(Samples, Out) :- Out=42.

part2(Samples, Program, Out) :- Out=42.

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
