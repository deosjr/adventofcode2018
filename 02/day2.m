:- module day2.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is cc_multi.
 
:- implementation.
:- import_module string, int, char.
:- import_module list, map.

:- pred part1(list(string)::in, string::out) is cc_multi.
part1(List, Out) :-
  list.map(set_2s_and_3s, List, ListOut),
  list.foldl2(count_2s_and_3s, ListOut, 0, Twos, 0, Threes),
  OutInt = Twos * Threes,
  Out = string.from_int(OutInt).

:- pred set_2s_and_3s(string::in, {int, int}::out) is cc_multi.
set_2s_and_3s(String, Tuple) :-
  map.init(Map),
  string.foldl(char_count, String, Map, MapOut),
  map.foldl(check_map, MapOut, {0,0}, Tuple).

:- pred char_count(char::in, map(char, int)::in, map(char, int)::out) is det.
char_count(Char, Map, MapOut) :-
  ( if map.search(Map, Char, Value) then
    NewValue = Value + 1,
    map.det_update(Char, NewValue, Map, MapOut)
  else
    map.det_insert(Char, 1, Map, MapOut)
  ).

:- pred check_map(char::in, int::in, {int, int}::in, {int, int}::out) is cc_multi.
check_map(_, 2, {TwoIn, ThreeIn}, {TwoOut, ThreeIn}) :-
  TwoOut = TwoIn + 1.
check_map(_, 3, {TwoIn, ThreeIn}, {TwoIn, ThreeOut}) :-
  ThreeOut = ThreeIn + 1.
check_map(_, _, A, A).

:- pred count_2s_and_3s({int, int}::in, int::in, int::out, int::in, int::out) is det.
count_2s_and_3s({Twos, Threes}, In2, Out2, In3, Out3) :-
  zero_or_one(Twos, ZOOTwos),
  Out2 = In2 + ZOOTwos,
  zero_or_one(Threes, ZOOThrees),
  Out3 = In3 + ZOOThrees.

:- pred zero_or_one(int::in, int::out) is det.
zero_or_one(X, Y) :-
  ( if X = 0 then
    Y = 0
  else
    Y = 1
  ).

:- pred part2(list(string)::in, string::out) is det.
part2(List, Overlap) :-
  find_diff_by1(List, String1, String2),
  ( if String1 = "Not", String2 = "Found" then
    Overlap = "No solution found"
  else
    L1 = string.to_char_list(String1),
    L2 = string.to_char_list(String2),
    overlap(L1, L2, LOverlap),
    Overlap = string.from_char_list(LOverlap)
  ).

:- pred find_diff_by1(list(string)::in, string::out, string::out) is det.
find_diff_by1([], "Not", "Found").
find_diff_by1([H|T], S1, S2) :-
  diff_by1(H, T, SS),
  ( if SS = "NotFound" then
    find_diff_by1(T, S1, S2)
  else
    H = S1,
    SS = S2
  ).

:- pred diff_by1(string::in, list(string)::in, string::out) is det.
diff_by1(_, [], "NotFound").
diff_by1(String, [H|T], S) :-
  ( if diff_strings(String, H, 1) then
    S = H
  else
    diff_by1(String, T, S)
  ).

:- pred diff_strings(string::in, string::in, int::out) is det.
diff_strings(S1, S2, X) :-
  L1 = string.to_char_list(S1),
  L2 = string.to_char_list(S2),
  diff_lists(L1, L2, X).

:- pred diff_lists(list(T)::in, list(T)::in, int::out) is det.
diff_lists([], B, C) :-
  ( if B = [] then
    C = 0
  else
    C = list.length(B)
  ).
diff_lists(A, B, C) :-
  A = [H1|T1],
  (
    B = [],
    C = list.length(A)
  ;
    B = [H2|T2],
    ( if H1 = H2 then
      diff_lists(T1, T2, C)
    else
      C is X + 1,
      diff_lists(T1, T2, X)
    )
  ).

:- pred overlap(list(T)::in, list(T)::in, list(T)::out) is det.
overlap([], A, A).
overlap(A, B, C) :-
  A = [H1|T1],
  (
    B = [],
    A = C
  ;
    B = [H2|T2],
    overlap(T1, T2, L),
    ( if H1 = H2 then
      C = [H1|L]
    else
      C = L
    )
  ).

main(!IO) :-
   io.open_input("day2.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
           ReadResult = ok(FileContents),
           string.split_at_string("\n", FileContents) = List,
           part1(List, Out1),
           io.write_string("Part 1: " ++ Out1 ++ "\n", !IO),
           part2(List, Out2),
           io.write_string("Part 2: " ++ Out2 ++ "\n", !IO)
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
