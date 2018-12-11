:- module day4.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, bool, int.
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
  FailLog = log("fail", 1, 1, wakesUp),
  ( if [DateTimeRaw, StatusRaw] = string.split_at_char(']', String) then
    DateTime = string.remove_prefix_if_present("[", DateTimeRaw),
    string.split(DateTime, 14, _, MinsRaw),
    MinsStripped = string.remove_prefix_if_present("0", MinsRaw),
    Minute = string.det_to_int(MinsStripped),
    Log1 = log(DateTime, Minute, 0, _),
    parse_status(StatusRaw, Status, GuardID, StatusSuccess),
    ( if StatusSuccess = yes then
      Success = yes,
      Log2 = Log1 ^ status := Status,
      Log = Log2 ^ guardID := GuardID
    else
      Success = no,
      Log = FailLog
    )
  else
    Success = no,
    Log = FailLog
  ).

:- pred parse_status(string::in, status::out, int::out, bool::out) is det.
parse_status(String, Status, GuardID, Success) :-
  ( if String = " wakes up" then
    Success = yes,
    Status = wakesUp,
    GuardID = 0
  else ( if String = " falls asleep" then
    Success = yes,
    Status = fallsAsleep,
    GuardID = 0
  % " Guard #2273 begins shift"
  else ( if [_, IDHash|_] = string.words(String) then 
    Success = yes,
    Status = beginsShift,
    IDStr = string.remove_prefix_if_present("#", IDHash),
    GuardID = string.det_to_int(IDStr)
  else
    Status = beginsShift,
    GuardID = 0,
    Success = no
    ))
  ).

:- pred aggregate_data(list(log)::in, map(int,int)::out, map(int,map(int,int))::out) is det.
aggregate_data(LogList, TotalSleepMapOut, MinsPerGuardMapOut) :-
  map.init(TotalSleepMapIn),
  map.init(MinsPerGuardMapIn),
  list.foldl4(aggregate_loop, LogList, TotalSleepMapIn, TotalSleepMapOut, MinsPerGuardMapIn, MinsPerGuardMapOut, 0, _, 0, _).

:- pred aggregate_loop(log::in, map(int,int)::in, map(int,int)::out, map(int,map(int,int))::in, map(int,map(int,int))::out, int::in, int::out, int::in, int::out) is det.
aggregate_loop(Log, !TotalSleepMap, !MinsPerGuardMap, !ID, !Min) :-
  (
    Log ^ status = beginsShift,
    !:ID = Log ^ guardID
  ;
    Log ^ status = fallsAsleep,
    !:Min = Log ^ minute
  ;
    Log ^ status = wakesUp,
    Minute = Log ^ minute,
    Diff = Minute - !.Min,
    update_sleep_map(!.ID, Diff, !TotalSleepMap),
    update_mins_map(!.ID, !.Min, Minute, !MinsPerGuardMap)
  ).

:- pred update_sleep_map(int::in, int::in,  map(int,int)::in, map(int,int)::out) is det.
:- pred update_mins_map(int::in, int::in, int::in,  map(int,map(int,int))::in, map(int,map(int,int))::out) is det.

update_sleep_map(Key, AddValue, !Map) :-
  ( if map.search(!.Map, Key, Value) then
    Added = Value + AddValue,
    map.det_update(Key, Added, !Map)
  else
    map.det_insert(Key, AddValue, !Map)
  ).

update_mins_map(ID, Current, End, !Map) :-
  ( if End = Current then 
    !.Map = !:Map
  else
    ( if map.search(!.Map, ID, NestedMap) then
      ( if map.search(NestedMap, Current, Value) then
        Added = Value + 1,
        map.det_update(Current, Added, NestedMap, NewNestedMap),
        map.det_update(ID, NewNestedMap, !Map)
      else
        map.det_insert(Current, 1, NestedMap, NewNestedMap),
        map.det_update(ID, NewNestedMap, !Map)
      )
    else
      map.init(NewMap),
      map.det_insert(ID, NewMap, !Map)
    ),
    Next = Current + 1,
    update_mins_map(ID, Next, End, !Map)
  ).

:- pred max_value_with_key(map(int,int)::in, int::out, int::out) is det.
:- pred max_value(int::in, int::in, int::in, int::out, int::in, int::out) is det.
max_value_with_key(Map, K, V) :-
  map.foldl2(max_value, Map, 0, K, 0, V).

max_value(Key, Value, !KeyAtMax, !MaxValue) :-
  ( if Value > !.MaxValue then
    !:KeyAtMax = Key,
    !:MaxValue = Value
  else
    !.MaxValue = !:MaxValue
  ).

:- pred part1(map(int,int)::in, map(int,map(int,int))::in, int::out) is det.
part1(TotalSleepMap, MinsPerGuardMap, Answer) :-
  max_value_with_key(TotalSleepMap, Guard, _),
  GuardMap = map.lookup(MinsPerGuardMap, Guard),
  max_value_with_key(GuardMap, Minute, _),
  Answer = Guard * Minute.

:- pred part2(map(int,map(int,int))::in, int::out) is det.
part2(MinsPerGuardMap, Answer) :-
  map.foldl3(map_loop, MinsPerGuardMap, 0, _, 0, Guard, 0, Minute),
  Answer = Guard * Minute.

:- pred map_loop(int::in, map(int,int)::in, int::in, int::out, int::in, int::out, int::in, int::out) is det.
map_loop(ID, Map, !MaxMinute, !Guard, !Minute) :-
  max_value_with_key(Map, Key, Value),
  ( if Value >= !.MaxMinute then
    !:MaxMinute = Value,
    !:Guard = ID,
    !:Minute = Key
  else
    % implied, but I dont know how to set up if without else?
    !:MaxMinute = !.MaxMinute
  ).

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
            aggregate_data(list.sort(Parsed), TotalSleepMap, MinsPerGuardMap),
            part1(TotalSleepMap, MinsPerGuardMap, Out1),
            S1 = string.format("Part 1: %i\n", [i(Out1)]),
            io.write_string(S1, !IO),
            part2(MinsPerGuardMap, Out2),
            S2 = string.format("Part 2: %i\n", [i(Out2)]),
            io.write_string(S2, !IO)
           ;
            Success = no,
            io.write_string("Parse failed.\n", !IO)
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
