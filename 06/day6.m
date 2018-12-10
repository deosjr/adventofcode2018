:- module day6.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, int.
:- import_module list, map.

:- type size ---> s(int) ; infinite.

:- pred parse(list(string)::in, list({int,int})::out, int::out, int::out, int::out, int::out) is det.
parse([], [], int.max_int, int.min_int, int.max_int, int.min_int).
parse([Line|T], Coordinates, Xmin, Xmax, Ymin, Ymax) :-
	( if [XStr, YStr] = string.split_at_char(',', Line) then
		parse(T, Coords, PrevXmin, PrevXmax, PrevYmin, PrevYmax),
		X = string.det_to_int(XStr),
		Y = string.det_to_int(string.remove_prefix_if_present(" ", YStr)),
		Coordinates = [{X, Y}|Coords],
		Xmin = int.min(X, PrevXmin), Xmax = int.max(X, PrevXmax), 
		Ymin = int.min(Y, PrevYmin), Ymax = int.max(Y, PrevYmax)
	else
		% error/1 does not work locally?
		% should just error out here
		Coordinates = [], Xmin = 0, Xmax = 0, Ymin = 0, Ymax = 0
	).

:- pred yloop(list({int,int})::in, int::in, {int,int}::in, {int,int}::in, map(int,size)::in, map(int,size)::out, int::out) is det.
:- pred xloop(list({int,int})::in, {int,int}::in, {int,int}::in, {int,int}::in, map(int,size)::in, map(int,size)::out, int::out) is det.
:- pred coordinate_loop({int,int}::in, {int,int}::in, int::in, int::out, int::in, int::out, int::in, int::out, int::in, int::out, int::in, int::out) is det.
:- pred update_map(map(int,size)::in, map(int,size)::out, int::in, int::in, {int,int}::in, {int,int}::in, {int,int}::in) is det.

yloop(Coordinates, Y, Xrange, Yrange, !SizeMap, NumSafe) :-
	Xrange = {Xmin,_}, Yrange = {_,Ymax},
	( if Y = Ymax then
		NumSafe = 0
	else
		yloop(Coordinates, Y+1, Xrange, Yrange, !.SizeMap, Map, YSafe),
		xloop(Coordinates, {Xmin,Y}, Xrange, Yrange, Map, !:SizeMap, XSafe),
		NumSafe = XSafe + YSafe
	).

xloop(Coordinates, {X,Y}, Xrange, Yrange, !SizeMap, NumSafe) :-
	Xrange = {_,Xmax},
	( if X = Xmax then
		NumSafe = 0
	else
		xloop(Coordinates, {X+1,Y}, Xrange, Yrange, !.SizeMap, Map, PrevSafe),
		list.foldl5(coordinate_loop({X,Y}), Coordinates, 0, _, 0, Sum, int.max_int, _, 0, Closest, 0, NumClosest),
		( if Sum < 10000 then
			NumSafe = PrevSafe + 1
		else
			NumSafe = PrevSafe
		),
		update_map(Map, !:SizeMap, Closest, NumClosest, {X,Y}, Xrange, Yrange)
	).

coordinate_loop(C1, C2, !ID, !Safe, !Distance, !Closest, !NumClosest) :-
	manhattan(C1, C2, M),
	!:Safe = !.Safe + M,
	!:ID = !.ID + 1,
	( if M < !.Distance then
		!:Distance = M,
		!:Closest = !.ID,
		!:NumClosest = 1
	else 
		( if M = !.Distance then
			!:NumClosest = !.NumClosest
		else
			% M > !.Distance
			!:Distance = !.Distance
		)
	).

update_map(!Map, Closest, NumClosest, {X,Y}, {Xmin,Xmax}, {Ymin,Ymax}) :-
	( if NumClosest > 1 then
		!.Map = !:Map
	else
		( if map.search(!.Map, Closest, Value), Value = infinite then
			!.Map = !:Map
		else
			( if X = Xmin; X = Xmax; Y = Ymin; Y = Ymax then
				map.set(Closest, infinite, !.Map, !:Map)
			else
				% set map(closest) to value+1
				( if map.search(!.Map, Closest, s(N)) then
					map.set(Closest, s(N+1), !.Map, !:Map)
				else
					map.set(Closest, s(1), !.Map, !:Map)
				)
			)
		)
	).


:- pred manhattan({int,int}::in, {int,int}::in, int::out) is det.
:- pred max_value(map(int,size)::in, size::out) is det.
:- pred max_size(size::in, size::in, size::out) is det.

manhattan({PX, PY}, {QX, QY}, X) :-
	X = int.abs(PX - QX) + int.abs(PY - QY).

max_value(Map, Max) :-
	list.foldl(max_size, map.values(Map), infinite, Max).
max_size(infinite, Size, Size).
max_size(s(X), S2, Max) :-
	( 
		S2 = infinite,
		Max = s(X)
	;
		S2 = s(Y),
		N = int.max(X, Y),
		Max = s(N)
	).

main(!IO) :-
   io.open_input("day6.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
        ReadResult = ok(FileContents),
        string.split_at_string("\n", FileContents) = List,
        parse(List, Input, Xmin, Xmax, Ymin, Ymax),
        map.init(MapIn),
        yloop(Input, Ymin, {Xmin,Xmax+1}, {Ymin,Ymax+1}, MapIn, Map, Safe),
        max_value(Map, MaxSize),
        ( 
        	MaxSize = infinite,
        	io.write_string("Part1: no non-infinite areas found!\n", !IO)
        ;
        	MaxSize = s(N),
        	S1 = string.format("Part1: %i\n", [i(N)]),
        	io.write_string(S1, !IO)
        ),
        S2 = string.format("Part2: %i\n", [i(Safe)]),
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