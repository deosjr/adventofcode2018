package main

import (
	"fmt"
	"io/ioutil"
)

type marble struct {
	number           int
	clockwise        *marble
	counterclockwise *marble
}

func marbleGame(numPlayers, numMarbles int) int {
	current := &marble{number: 0}
	current.clockwise = current
	current.counterclockwise = current

	scores := map[int]int{}
	player := 0
	for n := 1; n <= numMarbles; n++ {
		if n%23 == 0 {
			playerId := player + 1
			scores[playerId] += n
			for i := 0; i < 6; i++ {
				current = current.counterclockwise
			}

			oneCounter := current.counterclockwise
			twoCounter := oneCounter.counterclockwise

			scores[playerId] += oneCounter.number
			current.counterclockwise = twoCounter
			twoCounter.clockwise = current

			player = (player + 1) % numPlayers
			continue
		}
		newMarble := &marble{number: n}

		oneClockwise := current.clockwise
		twoClockwise := oneClockwise.clockwise

		oneClockwise.clockwise = newMarble
		twoClockwise.counterclockwise = newMarble
		newMarble.clockwise = twoClockwise
		newMarble.counterclockwise = oneClockwise

		current = newMarble
		player = (player + 1) % numPlayers
	}
	max := 0
	for _, v := range scores {
		if v > max {
			max = v
		}
	}
	return max
}

func main() {
	input, err := ioutil.ReadFile("day9.input")
	if err != nil {
		panic(err)
	}
	var numPlayers, numMarbles int
	fmt.Sscanf(string(input), "%d players; last marble is worth %d points", &numPlayers, &numMarbles)
	max := marbleGame(numPlayers, numMarbles)
	fmt.Printf("Part 1: %d\n", max)
	max = marbleGame(numPlayers, numMarbles*100)
	fmt.Printf("Part 2: %d\n", max)
}
