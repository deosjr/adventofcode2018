:- module day1.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is cc_multi.
 
:- implementation.
:- import_module string, integer.
:- import_module list, map.
 
:- pred part1(list(string)::in, string::out) is det.
part1(List, Out) :-
  list.foldl(parse_and_add, List, integer.zero, OutInt),
  integer.to_string(OutInt) = Out.

:- pred parse_and_add(string::in, integer::in, integer::out) is det.
parse_and_add(Str, In, Out) :-
  integer.det_from_string(Str) = InParsed,
  Out = In + InParsed.

:- pred part2(list(string)::in, string::out) is multi.
part2(List, Out) :-
  map.init(Map),
  repeat_parse_and_add(List, List, Map, integer.zero, Out). 

:- pred repeat_parse_and_add(list(string)::in, list(string)::in, map(integer, int)::in, integer::in, string::out) is multi.
repeat_parse_and_add([], List, Map, X, Y) :-
  repeat_parse_and_add(List, List, Map, X, Y).
repeat_parse_and_add([H|T], List, Map, X, Y) :-
  integer.det_from_string(H) = HParsed,
  N = X + HParsed, 
  (
    map.contains(Map, N),
    integer.to_string(N) = Y
  ;
    map.det_insert(N, 1, Map, MapOut),
    repeat_parse_and_add(T, List, MapOut, N, Y)
  ).
 
main(!IO) :-
   io.open_input("day1.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
           ReadResult = ok(FileContents),
           string.split_at_string("\n", FileContents) = List,
           part1(List, Out1),
           io.write_string("Part1: " ++ Out1 ++ "\n", !IO),
           part2(List, Out2),
           io.write_string("Part2: " ++ Out2 ++ "\n", !IO)
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