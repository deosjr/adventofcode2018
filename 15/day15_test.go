package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestPart1(t *testing.T) {
	for i, tt := range []struct {
		input string
		want  int
	}{
		{
			input: `#######
					#.G...#
					#...EG#
					#.#.#G#
					#..G#E#
					#.....#
					#######`,
			want: 27730,
		},
		{
			input: `#######
					#G..#E#
					#E#E.E#
					#G.##.#
					#...#E#
					#...E.#
					#######`,
			want: 36334,
		},
		{
			input: `#######
					#E..EG#
					#.#G.E#
					#E.##E#
					#G..#.#
					#..E#.#
					#######`,
			want: 39514,
		},
		{
			input: `#######
					#E.G#.#
					#.#G..#
					#G.#.G#
					#G..#.#
					#...E.#
					#######`,
			want: 27755,
		},
		{
			input: `#######
					#.E...#
					#.#..G#
					#.###.#
					#E#G#G#
					#...#G#
					#######`,
			want: 28944,
		},
		{
			input: `#########
					#G......#
					#.E.#...#
					#..##..G#
					#...##..#
					#...#...#
					#.G...G.#
					#.....G.#
					#########`,
			want: 18740,
		},
	} {
		in := strings.Replace(tt.input, "\t", "", -1)
		game := parse(in)
		got := part1(game)
		if got != tt.want {
			t.Errorf("%d): got %d want %d", i, got, tt.want)
		}
	}
}

func TestAdjacentTiles(t *testing.T) {
	e := &elf{&unit{pos: coord{2, 2}}}
	g := &goblin{&unit{pos: coord{2, 3}}}
	for i, tt := range []struct {
		tiles map[coord]tile
		p     coord
		want  []tile
	}{
		{
			tiles: map[coord]tile{
				{2, 2}: tile{unit: e},
				{1, 2}: tile{},
				{3, 2}: tile{},
				{2, 1}: tile{},
				{2, 3}: tile{},
			},
			p:    coord{2, 2},
			want: []tile{{}, {}, {}, {}},
		},
		{
			tiles: map[coord]tile{
				{2, 2}: tile{unit: e},
				{1, 2}: tile{isWall: true},
				{3, 2}: tile{},
				{2, 1}: tile{isWall: true},
				{2, 3}: tile{unit: g},
			},
			p:    coord{2, 2},
			want: []tile{{}, {unit: g}},
		},
	} {
		game := gameState{tiles: tt.tiles}
		got := game.adjacentTiles(tt.p)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}

func TestAdjacentEmptyCoords(t *testing.T) {
	e := &elf{&unit{pos: coord{2, 2}}}
	g := &goblin{&unit{pos: coord{2, 3}}}
	for i, tt := range []struct {
		tiles map[coord]tile
		p     coord
		want  []coord
	}{
		{
			tiles: map[coord]tile{
				{2, 2}: tile{unit: e},
				{1, 2}: tile{},
				{3, 2}: tile{},
				{2, 1}: tile{},
				{2, 3}: tile{},
			},
			p:    coord{2, 2},
			want: []coord{{2, 1}, {1, 2}, {3, 2}, {2, 3}},
		},
		{
			tiles: map[coord]tile{
				{2, 2}: tile{unit: e},
				{1, 2}: tile{isWall: true},
				{3, 2}: tile{},
				{2, 1}: tile{isWall: true},
				{2, 3}: tile{unit: g},
			},
			p:    coord{2, 2},
			want: []coord{{3, 2}},
		},
	} {
		game := gameState{tiles: tt.tiles}
		got := game.adjacentEmptyCoords(tt.p)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}

func TestMoves(t *testing.T) {
	for i, tt := range []struct {
		input string
		want  string
	}{
		{
			input: `#########
					#G..G..G#
					#.......#
					#.......#
					#G..E..G#
					#.......#
					#.......#
					#G..G..G#
					#########`,
			want: `#########
					#.G...G.#
					#...G...#
					#...E..G#
					#.G.....#
					#.......#
					#G..G..G#
					#.......#
					#########`,
		},
		{
			input: `#########
					#.G...G.#
					#...G...#
					#...E..G#
					#.G.....#
					#.......#
					#G..G..G#
					#.......#
					#########`,
			want: `#########
					#..G.G..#
					#...G...#
					#.G.E.G.#
					#.......#
					#G..G..G#
					#.......#
					#.......#
					#########`,
		},
		{
			input: `#########
					#..G.G..#
					#...G...#
					#.G.E.G.#
					#.......#
					#G..G..G#
					#.......#
					#.......#
					#########`,
			want: `#########
					#.......#
					#..GGG..#
					#..GEG..#
					#G..G...#
					#......G#
					#.......#
					#.......#
					#########`,
		},
		{
			input: `#########
					#.......#
					#..GGG..#
					#..GEG..#
					#G..G...#
					#......G#
					#.......#
					#.......#
					#########`,
			want: `#########
					#.......#
					#..GGG..#
					#..GEG..#
					#G..G...#
					#......G#
					#.......#
					#.......#
					#########`,
		},
		{
			input: `E..#....E
					##..G####`,
			want: `.E.#...E.
					##.G.####`,
		},
	} {
		in := strings.Replace(tt.input, "\t", "", -1)
		game := parse(in)
		round(game)
		split := strings.Split(in, "\n")
		got := game.testPrint(len(split[0]), len(split), false)
		want := strings.Replace(tt.want, "\t", "", -1)
		if got != want {
			t.Errorf("%d): \ngot \n%v\nwant \n%v\n", i, got, want)
		}
	}
}

func TestPossibleTargets(t *testing.T) {
	for i, tt := range []struct {
		input string
		unit  int
		want  []coord
	}{
		{
			input: `#####
					#.G.#
					#...#
					#...#
					#####`,
			unit: 0,
			want: []coord{},
		},
		{
			input: `#####
					#.G.#
					#####
					#.E.#
					#####`,
			unit: 0,
			want: []coord{{1, 3}, {3, 3}},
		},
		{
			input: `#####
					#.G.#
					#...#
					#.E.#
					#####`,
			unit: 1,
			want: []coord{{1, 1}, {3, 1}, {2, 2}},
		},
		{
			input: `#######
					#.G...#
					#...EG#
					#.#.#G#
					#..G#E#
					#.....#
					#######`,
			unit: 1,
			want: []coord{{1, 1}, {3, 1}, {2, 2}, {5, 1}, {3, 3}, {2, 4}, {3, 5}},
		},
	} {
		in := strings.Replace(tt.input, "\t", "", -1)
		game := parse(in)
		unit := game.units[tt.unit]
		got := game.possibleTargets(unit, map[Unit]struct{}{})
		wantMap := map[coord]struct{}{}
		for _, k := range tt.want {
			wantMap[k] = struct{}{}
		}
		if !reflect.DeepEqual(got, wantMap) {
			t.Errorf("%d): got %v want %v", i, got, wantMap)
		}
	}
}

func TestFloodFill(t *testing.T) {
	for i, tt := range []struct {
		input      string
		unit       int
		wantFound  []coord
		wantString string
		wantFirst  coord
		wantEnd    *coord
	}{
		{
			input: `#####
					#.G.#
					#####
					#.E.#
					#####`,
			unit:      1,
			wantFound: []coord{},
			wantString: `#####
						#.G.#
						#####
						#101#
						#####`,
		},
		{
			input: `#####
					#.G.#
					#...#
					#.E.#
					#####`,
			unit:      0,
			wantFound: []coord{{2, 2}},
			wantString: `#####
						#101#
						#.1.#
						#.E.#
						#####`,
			wantFirst: coord{2, 2},
		},
		{
			input: `#######
					#.E...#
					#.....#
					#...G.#
					#######`,
			unit:      0,
			wantFound: []coord{{4, 2}, {3, 3}},
			wantString: `#######
						#10123#
						#2123.#
						#323G.#
						#######`,
			wantFirst: coord{3, 1},
		},
		{
			input: `#########
					#G..G..G#
					#.......#
					#.......#
					#G..E..G#
					#.......#
					#.......#
					#G..G..G#
					#########`,
			unit:      4,
			wantFound: []coord{{4, 2}, {2, 4}, {6, 4}, {4, 6}},
			wantString: `#########
						#G..G..G#
						#...2...#
						#..212..#
						#G21012G#
						#..212..#
						#...2...#
						#G..G..G#
						#########`,
			wantFirst: coord{4, 3},
		},
		{
			input: `######
					#E...#
					#..#.#
					#.##.#
					#....#
					#..#G#
					######`,
			unit:      0,
			wantFound: []coord{{4, 4}},
			wantString: `######
						#0123#
						#12#4#
						#2##5#
						#3456#
						#45#G#
						######`,
			wantFirst: coord{2, 1},
		},
		{
			input: `######
					#.G..#
					#....#
					#.##.#
					#.##.#
					#..E.#
					######`,
			unit:      1,
			wantFound: []coord{{1, 1}, {2, 2}, {3, 1}},
			wantString: `######
						#6G65#
						#5654#
						#4##3#
						#3##2#
						#2101#
						######`,
			wantFirst: coord{2, 5},
		},
		{
			input: `#######
					##E...#
					#..##.#
					#.#.#.#
					#.....#
					#..#G##
					#######`,
			unit:      0,
			wantFound: []coord{{4, 4}},
			wantString: `#######
						##0123#
						#21##4#
						#3#7#5#
						#45676#
						#56#G##
						#######`,
			wantFirst: coord{3, 1},
		},
		{
			input: `#########
					#.......#
					#..GGG..#
					#..GEG..#
					#.......#
					#G..G..G#
					#.......#
					#.......#
					#########`,
			unit:      6,
			wantFound: []coord{{4, 4}},
			wantString: `#########
						#4......#
						#34GGG..#
						#23GEG..#
						#1234...#
						#012G..G#
						#1234...#
						#234....#
						#########`,
			wantFirst: coord{1, 4},
		},
		{
			input: `.G..
					#..#
					...E
					E...`,
			unit:      0,
			wantFound: []coord{{2, 2}, {0, 2}, {1, 3}},
			wantString: `1012
						#12#
						323E
						E3..`,
			wantFirst: coord{1, 1},
			wantEnd:   &coord{0, 2},
		},
		{
			input: `..G...#.##.#
					#....GGE...#
					#.......#..#
					#..........#
					#....E.....#
					#...EGE.#..#
					.E.........#`,
			unit:      0,
			wantFound: []coord{{1, 5}, {5, 3}, {4, 4}, {3, 5}, {2, 6}},
			wantString: `210123#.##.#
						#2123GGE...#
						#323456.#..#
						#43456.....#
						#5456E.....#
						#656EGE.#..#
						.E6........#`,
			wantFirst: coord{3, 0},
		},
	} {
		in := strings.Replace(tt.input, "\t", "", -1)
		game := parse(in)
		unit := game.units[tt.unit]
		gotFound, gotMap := game.floodFill(unit, game.possibleTargets(unit, map[Unit]struct{}{}))
		if !reflect.DeepEqual(gotFound, tt.wantFound) {
			t.Errorf("%d): got %v want %v", i, gotFound, tt.wantFound)
		}
		// parse wantString into floodfillmap
		wantMap := map[coord]int{}
		split := strings.Split(tt.wantString, "\n")
		for y := 0; y < len(split); y++ {
			s := strings.Replace(split[y], "\t", "", -1)
			for x := 0; x < len(s); x++ {
				i, err := strconv.Atoi(string(s[x]))
				if err != nil {
					continue
				}
				wantMap[coord{x, y}] = i
			}
		}
		if !reflect.DeepEqual(gotMap, wantMap) {
			t.Errorf("%d): got %v want %v", i, gotMap, wantMap)
		}
		if len(tt.wantFound) == 0 {
			continue
		}
		gotEnd := firstInReadingOrder(gotFound)
		if tt.wantEnd != nil {
			if gotEnd != *tt.wantEnd {
				t.Errorf("%d): got %v want %v", i, gotEnd, *tt.wantEnd)
			}
		}
		gotFirst := game.findFirstStep(gotEnd, gotMap)
		if gotFirst != tt.wantFirst {
			t.Errorf("%d): got %v want %v", i, gotFirst, tt.wantFirst)
		}
	}
}

func TestFirstInReadingOrder(t *testing.T) {
	for i, tt := range []struct {
		coords []coord
		want   coord
	}{
		{
			coords: []coord{{0, 0}},
			want:   coord{0, 0},
		},
		{
			coords: []coord{{4, 2}, {1, 0}, {2, 3}},
			want:   coord{1, 0},
		},
		{
			coords: []coord{{4, 2}, {3, 3}},
			want:   coord{4, 2},
		},
	} {
		got := firstInReadingOrder(tt.coords)
		if got != tt.want {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}

func TestRoundWithCombat(t *testing.T) {
	for i, tt := range []struct {
		input string
		setHP map[int]int
		want  string
	}{
		{
			input: `#######
					#.G...#
					#...EG#
					#.#.#G#
					#..G#E#
					#.....#
					#######`,
			want: `#######
					#..G..#   G(200)
					#...EG#   E(197), G(197)
					#.#G#G#   G(200), G(197)
					#...#E#   E(197)
					#.....#
					#######`,
		},
		{
			input: `#######
					#G..#E#
					#E#E.E#
					#G.##.#
					#...#E#
					#...E.#
					#######`,
			want: `#######
					#G.E#E#   G(197), E(200), E(200)
					#E#..E#   E(194), E(200)
					#G.##.#   G(200)
					#...#E#   E(200)
					#..E..#   E(200)
					#######`,
		},
		{
			input: `#######
					#G.E#E#
					#E#..E#
					#G.##.#
					#...#E#
					#..E..#
					#######`,
			want: `#######
					#GE.#E#   G(194), E(200), E(200)
					#E#..E#   E(194), E(200)
					#G.##.#   G(200)
					#..E#E#   E(200), E(200)
					#.....#
					#######`,
		},
		{
			input: `#######
					#GE.#E#
					#E#..E#
					#G.##.#
					#..E#E#
					#.....#
					#######`,
			want: `#######
					#GE.#E#   G(194), E(200), E(200)
					#E#..E#   E(191), E(200)
					#G.##.#   G(200)
					#.E.#.#   E(200)
					#....E#   E(200)
					#######`,
			setHP: map[int]int{3: 197},
		},
		{
			input: `#######
					#G....#
					#..G..#
					#..EG.#
					#..G..#
					#...G.#
					#######`,
			want: `#######
					#.G...#   G(9)
					#..G..#   G(4)
					#..E..#   E(194)
					#..GG.#   G(2), G(1)
					#.....#
					#######`,
			setHP: map[int]int{0: 9, 1: 4, 3: 2, 4: 2, 5: 1},
		},
		{
			input: `.G....#
					GEG...#
					.G.#..#
					.....G#
					.......
					#.....G
					......E`,
			want: `.G....#   G(200)
					...G..#   G(200)
					G.G#..#   G(143), G(200)
					......#
					.....G.   G(200)
					#.....G   G(197)
					......E   E(182)`,
			setHP: map[int]int{1: 143, 2: 2, 7: 185},
		},
	} {
		in := strings.Replace(tt.input, "\t", "", -1)
		game := parse(in)
		for k, v := range tt.setHP {
			u := game.units[k]
			if e, ok := u.(*elf); ok {
				e.hp = v
				continue
			}
			u.(*goblin).hp = v
		}
		round(game)
		split := strings.Split(in, "\n")
		got := game.testPrint(len(split[0]), len(split), true)
		want := strings.Replace(tt.want, "\t", "", -1)
		if got != want {
			t.Errorf("%d): \ngot \n%v\nwant \n%v\n", i, got, want)
		}
	}
}
