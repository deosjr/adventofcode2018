package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type group struct {
	units        int
	hitPoints    int
	initiative   int
	attackDamage int
	attackType   string
	weaknesses   []string
	immunities   []string

	// false means infection
	immuneSystem bool
}

func (g *group) effectivePower() int {
	return g.units * g.attackDamage
}

func (g *group) wouldDealDamage(target *group) int {
	for _, t := range target.immunities {
		if t == g.attackType {
			return 0
		}
	}
	d := g.effectivePower()
	for _, t := range target.weaknesses {
		if t == g.attackType {
			return 2 * d
		}
	}
	return d
}

func (g *group) takeDamage(d int) {
	g.units -= d / g.hitPoints
}

func parse(input string) []*group {
	toplevelsplit := strings.Split(input, "\n\n")
	immuneStr, infectionStr := toplevelsplit[0], toplevelsplit[1]
	return append(parseGroups(immuneStr, true), parseGroups(infectionStr, false)...)
}

func parseGroups(str string, immuneSystem bool) []*group {
	split := strings.Split(str, "\n")[1:]
	groups := make([]*group, len(split))
	for i, s := range split {
		var units, hitPoints, attackDamage, initiative int
		var attackType string
		withAn := strings.Split(s, "with an")
		points := strings.Split(withAn[0], "points")
		fmt.Sscanf(points[0], "%d units each with %d hit", &units, &hitPoints)
		fmt.Sscanf(withAn[1], " attack that does %d %s damage at initiative %d", &attackDamage, &attackType, &initiative)
		g := &group{
			units:        units,
			hitPoints:    hitPoints,
			initiative:   initiative,
			attackDamage: attackDamage,
			attackType:   attackType,
			immuneSystem: immuneSystem,
		}

		// optional: immunities & weaknesses
		if points[1] != " " {
			for _, iwstr := range strings.Split(points[1][2:len(points[1])-2], "; ") {
				if strings.HasPrefix(iwstr, "immune to ") {
					g.immunities = strings.Split(iwstr[10:], ", ")
					continue
				}
				if strings.HasPrefix(iwstr, "weak to ") {
					g.weaknesses = strings.Split(iwstr[8:], ", ")
				}
			}
		}
		groups[i] = g
	}
	return groups
}

func targetSort(groups []*group) func(i, j int) bool {
	return func(i, j int) bool {
		gi, gj := groups[i], groups[j]
		// equal effective power, tiebreaker is initiative
		if gi.effectivePower() == gj.effectivePower() {
			return gi.initiative > gj.initiative
		}
		return gi.effectivePower() > gj.effectivePower()
	}
}

func attackSort(groups []*group) func(i, j int) bool {
	return func(i, j int) bool {
		return groups[i].initiative > groups[j].initiative
	}
}

func part1(groups []*group) int {
	groups = combat(groups)
	return sumUnits(groups)
}

func combat(groups []*group) []*group {
	for !combatEnds(groups) {
		sum := sumUnits(groups)
		groups = fight(groups)
		if sum == sumUnits(groups) {
			// stalemate detected
			break
		}
	}
	return groups
}

func combatEnds(groups []*group) bool {
	if len(groups) < 2 {
		return true
	}
	g := groups[0]
	for _, gg := range groups[1:] {
		if g.immuneSystem != gg.immuneSystem {
			return false
		}
	}
	return true
}

func sumUnits(groups []*group) int {
	sum := 0
	for _, g := range groups {
		sum += g.units
	}
	return sum
}

// fight simulates one round of combat
func fight(groups []*group) []*group {
	targets := map[*group]*group{}
	targeted := map[*group]struct{}{}

	sort.Slice(groups, targetSort(groups))
	for _, g := range groups {
		for _, e := range groups {
			// dont target own team
			if g.immuneSystem == e.immuneSystem {
				continue
			}
			// group can only be targeted once per fight
			if _, ok := targeted[e]; ok {
				continue
			}
			dmg := g.wouldDealDamage(e)
			// can't target if you can't deal damage
			if dmg == 0 {
				continue
			}
			prevTarget, ok := targets[g]
			// no target found yet? target e
			if !ok {
				targets[g] = e
				continue
			}
			dmgOld := g.wouldDealDamage(prevTarget)
			// if we deal less dmg to e than to previous target, target previous
			// if we deal more dmg to e than to previous target, target e
			if dmg < dmgOld {
				continue
			}
			if dmg > dmgOld {
				targets[g] = e
				continue
			}
			// dmg == dmgOld: equal damage to previous and to e
			// higher effective power is tiebreaker
			ep, prevp := e.effectivePower(), prevTarget.effectivePower()
			if ep < prevp {
				continue
			}
			if ep > prevp {
				targets[g] = e
				continue
			}
			// ep == prep: equal effective power
			// tiebreaker is higher initiative (which is unique)
			if e.initiative > prevTarget.initiative {
				targets[g] = e
			}
		}
		t, ok := targets[g]
		if ok {
			targeted[t] = struct{}{}
		}
	}

	sort.Slice(groups, attackSort(groups))
	for _, g := range groups {
		if g.units <= 0 {
			continue
		}
		target, ok := targets[g]
		if !ok {
			continue
		}
		target.takeDamage(g.wouldDealDamage(target))
	}

	remaining := []*group{}
	for _, g := range groups {
		if g.units <= 0 {
			continue
		}
		remaining = append(remaining, g)
	}
	return remaining
}

func part2(input string) int {
	var boost int
	var groups []*group
	for {
		groups = parse(input)
		boost++
		groups = combatWithBoost(groups, boost)
		if immuneSystemWins(groups) {
			break
		}
	}
	return sumUnits(groups)
}

func combatWithBoost(groups []*group, boost int) []*group {
	for _, g := range groups {
		if g.immuneSystem {
			g.attackDamage += boost
		}
	}
	return combat(groups)
}

func immuneSystemWins(groups []*group) bool {
	if len(groups) == 0 {
		return false
	}
	for _, g := range groups {
		if !g.immuneSystem {
			return false
		}
	}
	return true
}

func main() {
	input, err := ioutil.ReadFile("day24.input")
	if err != nil {
		panic(err)
	}
	groups := parse(string(input))
	fmt.Printf("Part 1: %d\n", part1(groups))
	fmt.Printf("Part 2: %d\n", part2(string(input)))
}
