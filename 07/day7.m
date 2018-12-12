:- module day7.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, int, char, bool.
:- import_module list, map.
:- import_module maybe.

:- pred parse(list(string)::in, list(char)::out, map(char,list(char))::out, map(char,list(char))::out) is semidet.
:- pred parse_lines(list(string)::in, map(char,bool)::out, map(char,list(char))::out) is semidet.
:- pred find_roots(list(char)::in, map(char,list(char))::in, list(char)::out) is det.
:- pred prereqs_to_children(char::in, list(char)::in, map(char,list(char))::in, map(char,list(char))::out) is det.
parse(Input, list.sort(Roots), PreReqs, Children) :-
	parse_lines(Input, PossibleRoots, PreReqs),
	find_roots(map.keys(PossibleRoots), PreReqs, Roots),
	map.foldl(prereqs_to_children, PreReqs, map.init, Children).

parse_lines([], map.init, map.init).
parse_lines([Line|T], PossibleRoots, PreReqs) :-
	string.words(Line) = [_, PreqStr, _, _, _, _, _, IDStr, _, _],
	Preq = string.det_index(PreqStr, 0),
	ID = string.det_index(IDStr, 0),
	parse_lines(T, PrevRoots, PrevPreqs),
	PossibleRoots = map.set(PrevRoots, Preq, yes),
	( if map.search(PrevPreqs, ID, Value) then
		PreReqs = map.update(PrevPreqs, ID, [Preq|Value])
	else
		PreReqs = map.set(PrevPreqs, ID, [Preq])
	).

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
	( if map.search(ChildMap, Preq, Value) then
		!:Children = map.det_update(ChildMap, Preq, [ID|Value])
	else
		!:Children = map.set(ChildMap, Preq, [ID])
	).

:- pred part1(list(char)::in, map(char,list(char))::in, map(char,list(char))::in, string::out) is det.
:- pred p1loop(list(char)::in, map(char,list(char))::in, map(char,list(char))::in, list(char)::out) is det.
:- pred update_prereqs(char::in, list(char)::in, map(char,list(char))::in, map(char,list(char))::out, list(char)::in, list(char)::out) is det.
:- pred remove_from_prereqs(char::in, char::in, map(char,list(char))::in, map(char,list(char))::out, list(char)::in, list(char)::out) is det.
part1(Roots, PreReqs, Children, Answer) :-
	p1loop(Roots, PreReqs, Children, CharList),
	Answer = string.from_char_list(CharList).

p1loop([], _ , _, []).
p1loop([Root|Roots], PreReqs, Children, [Root|Out]) :-
	( if map.search(Children, Root, ChildrenOfRoot) then
		update_prereqs(Root, ChildrenOfRoot, PreReqs, PreReqsUpdated, Roots, NewRoots)
	else
		Roots = NewRoots,
		PreReqsUpdated = PreReqs
	),
	p1loop(NewRoots, PreReqsUpdated, Children, Out).

update_prereqs(ID, Children, !PreReqs, !Roots) :-
	list.foldl2(remove_from_prereqs(ID), Children, !PreReqs, !.Roots, R),
	!:Roots = list.sort(R).

%% remove Parent from the list of PreReqs for Child
%% if Child now has no prereqs, add it to Roots and delete its entry in PreReqs
remove_from_prereqs(Parent, Child, !PreReqs, !Roots) :-
	ChildReqs = map.lookup(!.PreReqs, Child),
	NewReqs = list.delete_all(ChildReqs, Parent),
	(
		NewReqs = [],
		!:PreReqs = map.delete(!.PreReqs, Child),
		!:Roots = [Child|!.Roots]
	;
		NewReqs = [_|_],
		!:PreReqs = map.set(!.PreReqs, Child, NewReqs)
	).

%% a worker is a tuple of {nodeID, minute} where nodeID can be missing
:- type worker ---> {maybe(char), int}.
:- pred worker_idle(worker::in) is semidet.
worker_idle({no, _}).

:- pred part2(list(char)::in, map(char,list(char))::in, map(char,list(char))::in, int::in, int::out) is det.
:- pred init_workers(int::in, list(worker)::out) is det.
:- pred p2loop(list(char)::in, map(char,list(char))::in, map(char,list(char))::in, int::in, list(worker)::in, int::out) is det.
:- pred check_workers(int::in, map(char,list(char))::in, worker::in, worker::out, list(char)::in, list(char)::out, int::in, int::out, map(char,list(char))::in, map(char,list(char))::out) is det.
:- pred divide_work(int::in, worker::in, worker::out, list(char)::in, list(char)::out, int::in, int::out) is det.
part2(Roots, PreReqs, Children, NumWorkers, Answer) :-
	init_workers(NumWorkers, Workers),
	p2loop(Roots, PreReqs, Children, 0, Workers, Answer).

init_workers(N, Workers) :-
	( if N = 0 then
		Workers = []
	else
		init_workers(N-1, W),
		Workers = [{no, 0}|W]
	).

p2loop(Roots, PreReqs, Children, Minute, Workers, Answer) :-
	list.map_foldl3(check_workers(Minute, Children), Workers, W, Roots, R, int.max_int, M, PreReqs, NewPreReqs),
	list.map_foldl2(divide_work(Minute), W, NewWorkers, R, NewRoots, M, NextMinute),
	( if NextMinute = int.max_int then
		Answer = Minute
	else
		p2loop(NewRoots, NewPreReqs, Children, NextMinute, NewWorkers, Answer)
	).

check_workers(_, _, {no, Min}, {no, Min}, !Roots, !NextMinute, !PreReqs).
check_workers(Minute, Children, {yes(ID), MinDone}, NewWorker, !Roots, !NextMinute, !PreReqs) :-
	( if MinDone = Minute then
		( if map.search(Children, ID, IDChildren) then
			NewWorker = {no, 0},
			update_prereqs(ID, IDChildren, !PreReqs, !Roots)
		else 
			NewWorker = {no, 0}
		)
	else
		NewWorker = {yes(ID), MinDone},
		!:NextMinute = int.min(MinDone, !.NextMinute)
	).

divide_work(_, !Worker, [], [], !NextMinute).
divide_work(Minute, !Worker, [Work|R], RootsOut, !NextMinute) :-
	( if worker_idle(!.Worker) then
		RootsOut = R, 
		MinDone = Minute + char.to_int(Work) - 4,
		!:NextMinute = int.min(MinDone, !.NextMinute),
		!:Worker = {yes(Work), MinDone}
	else 
		!:Worker = !.Worker,
		RootsOut = [Work|R]
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
	        part1(Roots, PreReqs, Children, Out1),
	        S1 = string.format("Part 1: %s\n", [s(Out1)]),
	        io.write_string(S1, !IO),
	        part2(Roots, PreReqs, Children, 5, Out2),
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
