:- module day9.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, int.
:- import_module list, map.

:- pred marble_game(int::in, int::in, int::out) is det.
:- pred game(int::in, int::in, int::in, int::in, map(int,{int,int})::in, map(int,int)::in, map(int,int)::out) is det.
:- pred init_player_map(int::in, map(int,int)::out) is det.
marble_game(NumPlayers, NumMarbles, WinningScore) :-
	MarbleMap = map.singleton(0, {0,0}),
	init_player_map(NumPlayers, PlayerInit),
	game(NumPlayers, 1, NumMarbles, 0, MarbleMap, PlayerInit, PlayerMap),
	list.foldl(int.max, map.values(PlayerMap), int.min_int, WinningScore).

%% MarbleMap stores neighbours for all marbles: id -> {clockwise, counterclockwise}
game(NumPlayers, N, NumMarbles, CurrentId, MarbleMap, !PlayerMap) :-
	(if N = NumMarbles+1 then 
		!:PlayerMap = !.PlayerMap
	else
		( if N mod 23 = 0 then
			remove(CurrentId, NewCurrent, Score, MarbleMap, NewMarbleMap),
			PlayerId = ((N-1) mod NumPlayers) + 1,
			map.lookup(!.PlayerMap, PlayerId, Value),
			PM = map.set(!.PlayerMap, PlayerId, Value+Score),
			game(NumPlayers, N+1, NumMarbles, NewCurrent, NewMarbleMap, PM, !:PlayerMap)
		else
			place(CurrentId, N, MarbleMap, NewMarbleMap),
			game(NumPlayers, N+1, NumMarbles, N, NewMarbleMap, !PlayerMap)
		)
	).

init_player_map(NumPlayers, Map) :-
	(if NumPlayers = 0 then
		Map = map.init
	else
		init_player_map(NumPlayers-1, M),
		Map = map.set(M, NumPlayers, 0)
	).

:- pred place(int::in, int::in, map(int,{int,int})::in, map(int,{int,int})::out) is det.
:- pred remove(int::in, int::out, int::out, map(int,{int,int})::in, map(int,{int,int})::out) is det.
:- pred clockwise(map(int,{int,int})::in, int::in, int::out) is det.
:- pred counter_clockwise(map(int,{int,int})::in, int::in, int::out) is det.
:- pred counter_clockwise(int::in, int::in, map(int,{int,int})::in, int::out) is det.
place(CurrentId, NewId, !Map) :- 
	clockwise(!.Map, CurrentId, Current1C),
	clockwise(!.Map, Current1C, Current2C),
	clockwise(!.Map, Current2C, Current3C),
	(if Current1C = Current2C then
		M = map.det_update(!.Map, Current1C, {NewId, NewId}),
		!:Map = map.det_insert(M, NewId, {Current1C, Current1C})
	else
		M1 = map.det_update(!.Map, Current1C, {NewId, CurrentId}),
		M2 = map.det_update(M1, Current2C, {Current3C, NewId}),
		!:Map = map.det_insert(M2, NewId, {Current2C, Current1C})
	).

remove(CurrentId, NewCurrent, Score, !Map) :-
	counter_clockwise(6, CurrentId, !.Map, Current6CC),
	{Current5CC, Current7CC} = map.lookup(!.Map, Current6CC),
	counter_clockwise(!.Map, Current7CC, Current8CC),
	counter_clockwise(!.Map, Current8CC, Current9CC),
	M1 = map.det_update(!.Map, Current6CC, {Current5CC, Current8CC}),
	M2 = map.det_update(M1, Current8CC, {Current6CC, Current9CC}), 
	!:Map = map.delete(M2, Current7CC),
	Score = CurrentId + 1 + Current7CC,
	NewCurrent = Current6CC.

clockwise(Map, Id, Clockwise) :-
	{Clockwise, _} = map.lookup(Map, Id).

counter_clockwise(Map, Id, CounterClockwise) :-
	{_, CounterClockwise} = map.lookup(Map, Id).

counter_clockwise(Steps, Id, Map, CC) :-
	( if Steps = 0 then
		CC = Id
	else
		map.lookup(Map, Id, {_, CounterClockwise}),
		counter_clockwise(Steps-1, CounterClockwise, Map, CC)
	).

main(!IO) :-
   io.open_input("day9.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
        ReadResult = ok(FileContents),
        (if [NumPlayersStr, _,_,_,_,_, NumMarblesStr, _] = string.words(FileContents) then
        	NumPlayers = string.det_to_int(NumPlayersStr),
        	NumMarbles = string.det_to_int(NumMarblesStr),
	        marble_game(NumPlayers, NumMarbles, Out1),
	        S1 = string.format("Part 1: %i\n", [i(Out1)]),
	        io.write_string(S1, !IO),
	        marble_game(NumPlayers, NumMarbles*100, Out2),
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
