:- module day5.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, char, int.
:- import_module list.

:- pred parse(string::in, list(int)::out) is det.
parse(String, List) :-
	Chars = string.to_char_list(String),
	list.map((pred(C::in, X::out) is det :- X = char.to_int(C)), Chars, List).

:- pred part1(list(int)::in, list(int)::in, int::in, int::out) is det.
part1([], _, A, A).
part1([H|T], Stack, !Answer) :-
	(
		Stack = [],
		part1(T, [H], !.Answer + 1, !:Answer)
	;
		Stack = [SH|ST],
		(if react(H, SH) then
			part1(T, ST, !.Answer - 1, !:Answer)
		else
			part1(T, [H|Stack], !.Answer + 1, !:Answer)
		)
	).

:- pred part2(list(int)::in, int::out) is det.
:- pred unit_types(int::in, int::in, list(int)::out) is det.
:- pred part1_with_skip(list(int)::in, int::in, list(int)::in, int::in, int::out) is det.
part2(Input, Answer) :-
	unit_types(65, 90, UnitTypes),
	list.map((pred(T::in, A::out) is det :- part1_with_skip(Input, T, [], 0, A)), UnitTypes, Lengths),
	list.foldl(int.min, Lengths, 999999, Answer).
	
unit_types(Current, End, List) :-
	( if Current = End then
		List = [End]
	else
		unit_types(Current+1, End, L),
		List = [Current|L]
	).

part1_with_skip([], _, _, A, A).
part1_with_skip([H|T], Skip, Stack, !Answer) :-
	( if skip(Skip, H) then
		part1_with_skip(T, Skip, Stack, !Answer)
	else
		(
			Stack = [],
			part1_with_skip(T, Skip, [H], !.Answer + 1, !:Answer)
		;
			Stack = [SH|ST],
			(if react(H, SH) then
				part1_with_skip(T, Skip, ST, !.Answer - 1, !:Answer)
			else
				part1_with_skip(T, Skip, [H|Stack], !.Answer + 1, !:Answer)
			)
		)
	).

:- pred react(int::in, int::in) is semidet.
:- pred skip(int::in, int::in) is semidet.
react(X, Y) :-
	Diff = X - Y,
	int.abs(Diff) = 32.
skip(Skip, X) :-
	X = Skip ; X = Skip + 32.

main(!IO) :-
   io.open_input("day5.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
        ReadResult = ok(FileContents),
        parse(FileContents, Input),
        part1(Input, [], 0, Out1),
        S1 = string.format("Part 1: %i\n", [i(Out1)]),
        io.write_string(S1, !IO),
        part2(Input, Out2),
        S2 = string.format("Part 2: %i\n", [i(Out2)]),
        io.write_string(S2, !IO)
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
