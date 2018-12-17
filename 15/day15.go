package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type coord struct {
	x, y int
}

type tile struct {
	isWall bool
	unit   Unit
}

type Unit interface {
	HP() int
	Attack() int
	Pos() coord
	Damage(int) bool
	Enemy(Unit) bool
	MoveTo(coord)
}

type unit struct {
	hp     int
	attack int
	pos    coord
}

func (u *unit) HP() int {
	return u.hp
}
func (u *unit) Attack() int {
	return u.attack
}
func (u *unit) Pos() coord {
	return u.pos
}
func (u *unit) Damage(d int) bool {
	u.hp -= d
	return u.hp <= 0
}
func (u *unit) MoveTo(c coord) {
	u.pos = c
}

type elf struct {
	*unit
}

func (*elf) Enemy(u Unit) bool {
	_, ok := u.(*goblin)
	return ok
}

type goblin struct {
	*unit
}

func (*goblin) Enemy(u Unit) bool {
	_, ok := u.(*elf)
	return ok
}

type gameState struct {
	tiles map[coord]tile
	units []Unit // sorted by pos
}

func parse(input string) *gameState {
	tiles := map[coord]tile{}
	units := []Unit{}
	for y, line := range strings.Split(string(input), "\n") {
		for x, c := range line {
			pos := coord{x, y}
			switch c {
			case '.':
				tiles[pos] = tile{}
			case '#':
				tiles[pos] = tile{isWall: true}
			case 'E':
				e := &elf{&unit{hp: 200, pos: pos, attack: 3}}
				units = addToSorted(units, e)
				tiles[pos] = tile{unit: e}
			case 'G':
				g := &goblin{&unit{hp: 200, pos: pos, attack: 3}}
				units = addToSorted(units, g)
				tiles[pos] = tile{unit: g}
			}
		}
	}
	return &gameState{
		tiles: tiles,
		units: units,
	}
}

func round(game *gameState) ([]Unit, bool) {
	casualties := map[Unit]struct{}{}
	end := false
	for _, u := range game.units {
		if _, ok := casualties[u]; ok {
			continue
		}
		if game.noTargets(casualties) {
			end = true
			break
		}
		t := game.target(u)
		if t == nil {
			// move
			game.move(u, casualties)

			// retarget
			t = game.target(u)
			if t == nil {
				continue
			}
		}
		//combat
		if lethal := t.Damage(u.Attack()); lethal {
			casualties[t] = struct{}{}
			game.tiles[t.Pos()] = tile{}
		}
	}
	newUnits := []Unit{}
	for _, u := range game.units {
		if _, ok := casualties[u]; ok {
			continue
		}
		newUnits = addToSorted(newUnits, u)
	}
	return newUnits, end
}

// adjacent should only be called for non-walls
// !ok should therefore never occur
func (game *gameState) adjacentTiles(p coord) []tile {
	adj := []tile{}
	neighbours := []coord{{p.x, p.y - 1}, {p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y + 1}}
	for _, c := range neighbours {
		t, ok := game.tiles[c]
		if !ok || t.isWall {
			continue
		}
		adj = append(adj, t)
	}
	return adj
}

func (game *gameState) adjacentEmptyCoords(p coord) []coord {
	coords := []coord{}
	neighbours := []coord{{p.x, p.y - 1}, {p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y + 1}}
	for _, c := range neighbours {
		t, ok := game.tiles[c]
		if !ok || t.isWall || t.unit != nil {
			continue
		}
		coords = append(coords, c)
	}
	return coords
}

func (game *gameState) possibleTargets(u Unit, casualties map[Unit]struct{}) map[coord]struct{} {
	targetCoords := map[coord]struct{}{}
	for _, unit := range game.units {
		if unit == u {
			continue
		}
		if _, ok := casualties[unit]; ok {
			continue
		}
		if !u.Enemy(unit) {
			continue
		}
		for _, c := range game.adjacentEmptyCoords(unit.Pos()) {
			targetCoords[c] = struct{}{}
		}
	}
	return targetCoords
}

func (game *gameState) floodFill(u Unit, targetCoords map[coord]struct{}) (found []coord, floodFillMap map[coord]int) {
	floodFillMap = map[coord]int{u.Pos(): 0}
	found = []coord{}
	i := 0
	fringe := []coord{u.Pos()}
	for len(fringe) != 0 {
		i++
		newFringe := []coord{}

		for _, f := range fringe {
			for _, c := range game.adjacentEmptyCoords(f) {
				if _, ok := floodFillMap[c]; ok {
					continue
				}
				floodFillMap[c] = i
				newFringe = append(newFringe, c)
				if _, ok := targetCoords[c]; ok {
					found = append(found, c)
				}
			}
		}

		if len(found) > 0 {
			break
		}
		fringe = newFringe
	}
	return found, floodFillMap
}

func (game *gameState) move(u Unit, casualties map[Unit]struct{}) {
	// find possible target locations
	targetCoords := game.possibleTargets(u, casualties)
	if len(targetCoords) == 0 {
		return
	}
	// expand to all tiles with increasing stepsize
	found, floodFillMap := game.floodFill(u, targetCoords)
	if len(found) == 0 {
		return
	}
	// find path to closest target and pick step 1 on the path
	end := firstInReadingOrder(found)
	first := game.findFirstStep(end, floodFillMap)
	// move u to coord first
	game.step(u, first)
}

// find ALL shortest paths from end coord to start (score 0)
// out of those shortest paths, take the first step at reading order
func (game *gameState) findFirstStep(end coord, ffm map[coord]int) coord {
	length := ffm[end]
	paths := map[int]map[coord]struct{}{length: map[coord]struct{}{end: {}}}
	for i := length - 1; i > 0; i-- {
		prev := paths[i+1]
		paths[i] = map[coord]struct{}{}
		for p, _ := range prev {
			for _, next := range game.adjacentEmptyCoords(p) {
				if ffm[next] == i {
					paths[i][next] = struct{}{}
				}
			}
		}
	}
	firsts := []coord{}
	for k, _ := range paths[1] {
		firsts = append(firsts, k)
	}
	return firstInReadingOrder(firsts)
}

func (game *gameState) step(u Unit, c coord) {
	old := game.tiles[c]
	old.unit = u
	game.tiles[c] = old
	told := game.tiles[u.Pos()]
	told.unit = nil
	game.tiles[u.Pos()] = told
	u.MoveTo(c)
}

func (game *gameState) target(u Unit) Unit {
	var t Unit = nil
	for _, tile := range game.adjacentTiles(u.Pos()) {
		if tile.unit == nil {
			continue
		}
		if !u.Enemy(tile.unit) {
			continue
		}
		if t == nil {
			t = tile.unit
			continue
		}
		if tile.unit.HP() < t.HP() {
			t = tile.unit
			continue
		}
		if order2d(tile.unit.Pos(), t.Pos()) {
			t = tile.unit
		}
	}
	return t
}

func (game *gameState) noTargets(casualties map[Unit]struct{}) bool {
	var first Unit = nil
	for _, u := range game.units {
		if _, ok := casualties[u]; ok {
			continue
		}
		if first == nil {
			first = u
			continue
		}
		if first.Enemy(u) {
			return false
		}
	}
	return true
}

func (game *gameState) outcome(rounds int) int {
	sum := 0
	for _, u := range game.units {
		sum += u.HP()
	}
	return rounds * sum
}

// adapted from day 13
func addToSorted(list []Unit, u Unit) []Unit {
	if len(list) == 0 {
		return []Unit{u}
	}
	// find smallest index for which list[i].pos comes after u.pos in 2d order
	index := sort.Search(len(list), func(i int) bool { return order2d(u.Pos(), list[i].Pos()) })
	// not found: m.id is biggest so insert last
	if index == len(list) {
		return append(list, u)
	}
	return append(list[:index], append([]Unit{u}, list[index:]...)...)
}

func removeFromSorted(list []unit, u unit) []unit {
	toRemove := -1
	for i, v := range list {
		if v.pos == u.pos {
			toRemove = i
		}
	}
	if toRemove == -1 {
		fmt.Println("LIST DID NOT CONTAIN U")
		return list
	}
	if toRemove == 0 {
		return list[1:]
	}
	if toRemove == len(list)-1 {
		return list[:len(list)-1]
	}
	return append(list[:toRemove], list[toRemove+1:]...)
}

func firstInReadingOrder(found []coord) coord {
	t := found[0]
	// if multiple closest found: pick the closest target in reading order
	if len(found) > 1 {
		for _, f := range found[1:] {
			if order2d(f, t) {
				t = f
			}
		}
	}
	return t
}

// order2d returns true if q comes after p in 2d grid order
func order2d(p, q coord) bool {
	if q.y > p.y {
		return true
	}
	if q.y < p.y {
		return false
	}
	if q.x > p.x {
		return true
	}
	return false
}

func part1(game *gameState) int {
	rounds := 0
	for {
		units, noTargetsFound := round(game)
		game.units = units
		if noTargetsFound {
			return game.outcome(rounds)
		}
		rounds++
	}
}

func part2(input string) int {
	attackPower := 3
	for {
		attackPower++
		game := parse(input)
		elfCount := 0
		for _, u := range game.units {
			if e, ok := u.(*elf); ok {
				elfCount++
				e.unit.attack = attackPower
			}
		}
		outcome := part1(game)
		for _, u := range game.units {
			if _, ok := u.(*elf); ok {
				elfCount--
			}
		}
		if elfCount == 0 {
			return outcome
		}
	}
}

func main() {
	input, err := ioutil.ReadFile("day15.input")
	if err != nil {
		panic(err)
	}
	game := parse(string(input))
	fmt.Printf("Part 1: %d\n", part1(game))
	fmt.Printf("Part 2: %d\n", part2(string(input)))
}

func (game *gameState) testPrint(xMax, yMax int, printHealth bool) string {
	s := ""
	for y := 0; y < yMax; y++ {
		units := []string{}
		for x := 0; x < xMax; x++ {
			t := game.tiles[coord{x, y}]
			if t.unit == nil {
				if t.isWall {
					s += "#"
					continue
				}
				s += "."
				continue
			}
			if _, ok := t.unit.(*elf); ok {
				s += "E"
				units = append(units, fmt.Sprintf("E(%d)", t.unit.HP()))
				continue
			}
			s += "G"
			units = append(units, fmt.Sprintf("G(%d)", t.unit.HP()))
		}
		if printHealth && len(units) > 0 {
			s += "   "
			s += strings.Join(units, ", ")
		}
		s += "\n"
	}
	return s[:len(s)-1]
}
