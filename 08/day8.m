:- module day8.
:- interface.
 
:- import_module io.
:- pred main(io::di, io::uo) is det.
 
:- implementation.
:- import_module string, int.
:- import_module list, map.

:- pred part1(list(int)::in, list(int)::out, int::out) is det.
:- pred part1_children(int::in, list(int)::in, list(int)::out, int::out) is det.
:- pred part1_metadata(int::in, list(int)::in, list(int)::out, int::out) is det.
part1(!List, Answer) :-
  NumChildren = list.det_head(!.List),
  TempList = list.det_tail(!.List),
  NumMetadata = list.det_head(TempList),
  Tail = list.det_tail(TempList),
  part1_children(NumChildren, Tail, TailOut, ChildAns),
  part1_metadata(NumMetadata, TailOut, !:List, MetaAns),
  Answer = ChildAns + MetaAns.

part1_children(RangeCounter, !List, Answer) :-
  ( if RangeCounter = 0 then
    Answer = 0
  else
    part1(!.List, Tail, Ans),
    MinOne = RangeCounter - 1,
    part1_children(MinOne, Tail, !:List, ChildAns),
    Answer = Ans + ChildAns
  ).

part1_metadata(RangeCounter, !List, Answer) :-
  ( if RangeCounter = 0 then
    Answer = 0
  else
    Meta = list.det_head(!.List),
    Tail = list.det_tail(!.List),
    MinOne = RangeCounter - 1,
    part1_metadata(MinOne, Tail, !:List, Ans),
    Answer = Meta + Ans
  ).

:- pred part2(list(int)::in, list(int)::out, int::out) is det.
:- pred part2_children(int::in, int::in, list(int)::in, list(int)::out, map(int,int)::in, map(int,int)::out) is det.
:- pred metadata_no_children(int::in, list(int)::in, list(int)::out, int::out) is det.
:- pred metadata_with_children(int::in, int::in, list(int)::in, list(int)::out, map(int,int)::in, int::out) is det.
part2(!List, Answer) :-
  NumChildren = list.det_head(!.List),
  TempList = list.det_tail(!.List),
  NumMetadata = list.det_head(TempList),
  Tail = list.det_tail(TempList),
  ( if NumChildren = 0 then
    metadata_no_children(NumMetadata, Tail, !:List, Answer)
  else
    map.init(EmptyMap),
    part2_children(0, NumChildren, Tail, TailOut, EmptyMap, Map),
    metadata_with_children(NumMetadata, NumChildren, TailOut, !:List, Map, Answer)
  ).

part2_children(RangeCounter, NumChildren, !List, !Map) :-
  ( if RangeCounter = NumChildren then
    !.List = !:List
  else
    part2(!.List, Tail, ChildValue),
    Next = RangeCounter + 1,
    map.det_insert(Next, ChildValue, !.Map, M),
    part2_children(Next, NumChildren, Tail, !:List, M, !:Map)
  ).

% exactly the same predicate
metadata_no_children(RangeCounter, !List, Answer) :-
  part1_metadata(RangeCounter, !List, Answer).

metadata_with_children(RangeCounter, NumChildren, !List, Map, Answer) :-
  ( if RangeCounter = 0 then
    Answer = 0
  else
    Meta = list.det_head(!.List),
    Tail = list.det_tail(!.List),
    MinOne = RangeCounter - 1,
    ( if Meta > NumChildren then
      metadata_with_children(MinOne, NumChildren, Tail, !:List, Map, Answer)
    else
      metadata_with_children(MinOne, NumChildren, Tail, !:List, Map, Ans),
      map.lookup(Map, Meta, Value),
      Answer = Ans + Value
    )
  ).

:- pred parse(string::in, list(int)::out) is det.
parse(String, IntList) :-
  string.words(String) = List,
  list.map((pred(Str::in, Int::out) is det :- Int = string.det_to_int(Str)), List, IntList).

main(!IO) :-
   io.open_input("day8.input", OpenResult, !IO),
   (
      OpenResult = ok(File),
      io.read_file_as_string(File, ReadResult, !IO),
      (
        ReadResult = ok(FileContents),
        parse(FileContents, Input),
        part1(Input, _, Out1),
        S1 = string.format("Part1: %i\n", [i(Out1)]),
        io.write_string(S1, !IO),
        part2(Input, _, Out2),
        S2 = string.format("Part2: %i\n", [i(Out2)]),
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