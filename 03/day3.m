:- module day3.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is cc_multi.
 
:- implementation.
:- import_module int, string, char, bool.
:- import_module list, map.

:- type fabric ---> id(int) ; overlap.
 
:- pred part1(list(string)::in, string::out, map(int,bool)::out, bool::out) is cc_multi.
part1(List, Out, IDMapOut, Success) :-
  % foldl function over List, aggregating giant map:
  % - parse line
  % - for x for y in range coordinates
  % - set id or count overlap in giant map(tuple->overlap_state)
  map.init(MapIn),
  map.init(IDMapIn),
  list.foldl4(do_per_line, List, MapIn, _, IDMapIn, IDMapOut, 0, TotalOverlap, yes, ParseSucceed),
  ( 
    ParseSucceed = yes,
    Success = yes,
    Out = string.from_int(TotalOverlap)
  ;
    ParseSucceed = no,
    Success = no,
    Out = "FAIL"
  ).

:- pred do_per_line(string::in, map({int,int},fabric)::in, map({int,int},fabric)::out, map(int,bool)::in, map(int,bool)::out, int::in, int::out, bool::in, bool::out) is cc_multi.
do_per_line(String, !Map, !IDMap, !TotalOverlap, !Success) :-
  ( 
    !.Success = no,
    !:Success = no
  ;
    !.Success = yes,
    parse_line(String, Tuple, LineSuccess),
    (
      LineSuccess = no,
      !:Success = no
    ;
      LineSuccess = yes,
      Tuple = {ID,X,Y,W,H},
      XMax = X + W,
      YMax = Y + H,
      nested_for_loop(ID, X, X, XMax, Y, YMax, !Map, !IDMap, !TotalOverlap),
      !:Success = yes
    )
  ).

% ranges module isnt working; would've loved to just foldl over range
:- pred nested_for_loop(int::in, int::in, int::in, int::in, int::in, int::in, map({int,int},fabric)::in, map({int,int},fabric)::out, map(int,bool)::in, map(int,bool)::out, int::in, int::out) is multi.
nested_for_loop(ID, XBase, X, XMax, Y, YMax, !Map, !IDMap, !TotalOverlap) :-
  ( 
    % end of both loops
    Y = YMax
  ;
    (
      % end of x loop, repeat
      X = XMax,
      YNew = Y + 1,
      nested_for_loop(ID, XBase, XBase, XMax, YNew, YMax, !Map, !IDMap, !TotalOverlap)
    ;
      update_maps(ID, {X,Y}, !.Map, NewMap, !.IDMap, NewIDMap, !.TotalOverlap, TONew),
      XNew = X + 1,
      nested_for_loop(ID, XBase, XNew, XMax, Y, YMax, NewMap, !:Map, NewIDMap, !:IDMap, TONew, !:TotalOverlap)
    )
  ).

:- pred update_maps(int::in, {int,int}::in, map({int,int},fabric)::in, map({int,int},fabric)::out, map(int,bool)::in, map(int,bool)::out, int::in, int::out) is det.
update_maps(ID, Coord, !Map, !IDMap, !TotalOverlap) :-
  ( if map.search(!.Map, Coord, Value) then
    (
      % first time overlap
      Value = id(ID2),
      !:TotalOverlap = !.TotalOverlap + 1,
      map.det_update(Coord, overlap, !Map),
      update_id_map(ID, !.IDMap, IDMap1),
      update_id_map(ID2, IDMap1, !:IDMap)
    ;
      % has already overlapped
      Value = overlap,
      update_id_map(ID, !IDMap)
    )
  else
    % new in map
    map.det_insert(Coord, id(ID), !Map)
  ).

:- pred update_id_map(int::in, map(int,bool)::in, map(int,bool)::out) is det.
update_id_map(ID, !Map) :-
  ( if map.contains(!.Map, ID) then
    !.Map = !:Map
  else
    map.det_insert(ID, yes, !Map)
  ).

% tuple: id, x, y, w, h
:- pred parse_line(string::in, {int,int,int,int,int}::out, bool::out) is det.
parse_line(String, Tuple, Success) :-
  Words = string.words(String),
  ( if [IDStr, _, XY, WH] = Words then
    parse_id(IDStr, ID),
    parse_xy(XY, ',', X, Y, SuccessXY),
    (
      SuccessXY = no,
      Success = no,
      Tuple = {1,1,1,1,1}
    ;
      SuccessXY = yes,
      parse_xy(WH, 'x', W, H, SuccessWH),
      (
        SuccessWH = no,
        Success = no,
        Tuple = {1,1,1,1,1}
      ;
        SuccessWH = yes,
        Success = yes,
        Tuple = {ID, X, Y, W, H}
      )
    )
  else
    Success = no,
    Tuple = {1,1,1,1,1}
  ).

:- pred parse_id(string::in, int::out) is det.
parse_id(String, ID) :-
  IDFixed = string.remove_prefix_if_present("#", String),
  ID = string.det_to_int(IDFixed).

:- pred parse_xy(string::in, char::in, int::out, int::out, bool::out) is det.
parse_xy(String, Char, X, Y, Success) :-
  ( if [XStr, YStr] = string.split_at_char(Char, String) then
    Success = yes,
    X = string.det_to_int(XStr),
    YFixed = string.remove_suffix_if_present(":", YStr),
    Y = string.det_to_int(YFixed)
  else
    Success = no,
    X = 1,
    Y = 1
  ).

:- pred part2(map(int,bool)::in, int::in, int::in, int::out, bool::out) is det.
part2(Map, X, Len, Int, Success) :-
  ( if X = Len then
    Success = no,
    Int = 0
  else
    ( if map.contains(Map, X) then
      Y = X + 1,
      part2(Map, Y, Len, Int, Success)
    else
      Success = yes,
      Int = X
    )
  ).
  
main(!IO) :-
   io.open_input("day3.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
           ReadResult = ok(FileContents),
           string.split_at_string("\n", FileContents) = List,
           part1(List, Out1, IDMap, Success),
           (
            Success = yes,
            S1 = string.format("Part1: %s\n", [s(Out1)]),
            io.write_string(S1, !IO),
            Len = list.length(List),
            Limit = Len + 1,
            part2(IDMap, 1, Limit, Out2, Success2),
            (
              Success2 = yes,
              S2 = string.format("Part2: %i\n", [i(Out2)]),
              io.write_string(S2, !IO)
            ;
              Success2 = no,
              io.write_string("Part 2 failed.", !IO)
            )
           ;
            Success = no,
            io.write_string("Part 1 failed.", !IO)
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