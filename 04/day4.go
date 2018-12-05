package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type state uint8

const (
	beginsShift state = iota
	fallsAsleep
	wakesUp
)

type log struct {
	date    string
	hour    string
	minute  int
	guardID int
	status  state
}

type byDate []log

func (s byDate) Len() int {
	return len(s)
}
func (s byDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byDate) Less(i, j int) bool {
	if s[i].date == s[j].date {
		if s[i].hour == s[j].hour {
			return s[i].minute < s[j].minute
		}
		return s[i].hour < s[j].hour
	}
	return s[i].date < s[j].date
}

func main() {
	input, err := ioutil.ReadFile("day4.input")
	if err != nil {
		panic(err)
	}

	splitinput := strings.Split(string(input), "\n")
	logs := make([]log, len(splitinput))
	for i, v := range splitinput {
		s := strings.Split(v, "]")
		datetime := strings.Split(s[0][1:], " ")
		date := datetime[0]
		clock := strings.Split(datetime[1], ":")
		hour, min := clock[0], clock[1]
		minute, _ := strconv.Atoi(min)
		l := log{
			date:   date,
			hour:   hour,
			minute: minute,
		}
		switch s[1] {
		case " falls asleep":
			l.status = fallsAsleep
		case " wakes up":
			l.status = wakesUp
		default:
			l.status = beginsShift
			split := strings.Split(s[1], "#")
			furthersplit := strings.Split(split[1], " ")
			l.guardID, _ = strconv.Atoi(furthersplit[0])
		}
		logs[i] = l
	}
	sort.Sort(byDate(logs))

	totalSleep := map[int]int{}
	sleepMinsPerGuard := map[int]map[int]int{}
	var id, min int
	for _, log := range logs {
		switch log.status {
		case beginsShift:
			id = log.guardID
		case fallsAsleep:
			min = log.minute
		case wakesUp:
			diff := log.minute - min
			totalSleep[id] += diff
			for i := min; i < log.minute; i++ {
				if _, ok := sleepMinsPerGuard[id]; !ok {
					sleepMinsPerGuard[id] = map[int]int{}
				}
				sleepMinsPerGuard[id][i] += 1
			}
		}
	}

	chosenGuard, _ := maxValueWithKey(totalSleep)
	chosenMinute, _ := maxValueWithKey(sleepMinsPerGuard[chosenGuard])
	fmt.Printf("Part 1: %d\n", chosenGuard*chosenMinute)

	var maxMinute int
	for id, m := range sleepMinsPerGuard {
		k, v := maxValueWithKey(m)
		if v >= maxMinute {
			maxMinute = v
			chosenGuard = id
			chosenMinute = k
		}
	}
	fmt.Printf("Part 2: %d\n", chosenGuard*chosenMinute)
}

func maxValueWithKey(m map[int]int) (int, int) {
	var key, max int
	for k, v := range m {
		if v >= max {
			max = v
			key = k
		}
	}
	return key, max
}
