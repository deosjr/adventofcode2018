:- module day4.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is cc_multi.
 
:- implementation.
:- import_module string, bool.
:- import_module list, map.

:- type status --->
  beginsShift ;
  fallsAsleep ;
  wakesUp.

:- type log ---> 
  log(
    datetime  :: string,
    minute    :: int,
    guardID   :: int,
    status    :: status
    ).
 
:- pred parse(list(string)::in, list(log)::in, list(log)::out, bool::out) is det.
parse([], A, A, yes).
parse([H|T], In, Out, Success) :-
  parse_line(H, Log, LineSuccess),
  (
    LineSuccess = yes,
    List = list.cons(Log, In),
    parse(T, List, Out, Success) 
  ;
    LineSuccess = no,
    Success = no,
    In = Out
  ).

:- pred parse_line(string::in, log::out, bool::out) is det.
parse_line(String, Log, Success) :-
  ( if [DateTimeRaw, StatusRaw] = string.split_at_char(']', String) then
    Success = yes,
    DateTime = string.remove_prefix_if_present("[", DateTimeRaw),
    string.split(DateTime, 14, _, MinsRaw),
    MinsStripped = string.remove_prefix_if_present("0", MinsRaw),
    Minute = string.det_to_int(MinsStripped),
    Log1 = log(DateTime, Minute, 0, _),
    % TODO switch status
    Log2 = Log1 ^ guardID := 1,
    Log = Log2 ^ status := wakesUp
  else
    Success = no,
    Log = log("fail", 1, 1, wakesUp)
  ).

%TODO
:- pred part1(list(log)::in, int::out) is det.
part1([], 42). 
part1([H|T], Int) :-
  Int = H ^ minute.

main(!IO) :-
   io.open_input("day4.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
           ReadResult = ok(FileContents),
           string.split_at_string("\n", FileContents) = List,
           parse(List, [], Parsed, Success),
           (
            Success = yes,
            part1(Parsed, Out1),
            S1 = string.format("Part1: %i\n", [i(Out1)]),
            io.write_string(S1, !IO)
           ;
            Success = no,
            io.write_string("Parse failed.", !IO)
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