:- module day7.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, int, char, bool.
:- import_module list, map.

:- pred parse(list(string)::in, list(char)::out, map(char,list(char))::out, map(char,list(char))::out) is semidet.
:- pred parse_lines(list(string)::in, map(char,bool)::out, map(char,list(char))::out) is semidet.
:- pred find_roots(list(char)::in, map(char,list(char))::in, list(char)::out) is det.
:- pred prereqs_to_children(char::in, list(char)::in, map(char,list(char))::in, map(char,list(char))::out) is det.
parse(Input, Roots, PreReqs, Children) :-
	parse_lines(Input, PossibleRoots, PreReqs),
	find_roots(map.keys(PossibleRoots), PreReqs, Roots),
	map.foldl(prereqs_to_children, PreReqs, map.init, Children).

parse_lines([], map.init, map.init).
parse_lines([Line|T], PossibleRoots, PreReqs) :-
	string.words(Line) = [_, PreqStr, _, _, _, _, _, IDStr, _, _],
	Preq = string.det_index(PreqStr, 0),
	ID = string.det_index(IDStr, 0),
	parse_lines(T, PrevRoots, PrevPreqs),
	PossibleRoots = map.set(PrevRoots, ID, yes),
	( if map.search(PrevPreqs, ID, Value) then
		PreReqs = map.update(PrevPreqs, ID, [Preq|Value])
	else
		PreReqs = map.set(PrevPreqs, ID, [Preq])
	).

% TODO: Roots should be sorted list
find_roots([], _, []).
find_roots([PRoot|PRoots], PreReqs, Roots) :-
	find_roots(PRoots, PreReqs, R),
	(if map.contains(PreReqs, PRoot) then 
		Roots = R
	else
		Roots = [PRoot | R]
	).

prereqs_to_children(_, [], !Children).
prereqs_to_children(ID, [Preq|PreReqs], !Children) :-
	prereqs_to_children(ID, PreReqs, !.Children, ChildMap),
	( if map.search(ChildMap, ID, Value) then
		!:Children = map.det_update(ChildMap, ID, [Preq|Value])
	else
		!:Children = map.set(ChildMap, ID, [Preq])
	).

main(!IO) :-
   io.open_input("day7.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
        ReadResult = ok(FileContents),
        string.split_at_string("\n", FileContents) = List,
        (if parse(List, Roots, PreReqs, Children) then
	        Out1 = 42, % TODO
	        Out2 = 42, % TODO
	        S1 = string.format("Part 1: %i\n", [i(Out1)]),
	        io.write_string(S1, !IO),
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
