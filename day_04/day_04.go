package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type EventType int

const (
	BeginShift EventType = iota
	FallAsleep
	WakeUp
)

const minutes = 60

type Event struct {
	EventType EventType
	TimeStamp time.Time
	GuardId   int
}

type EventList []*Event

func (e EventList) Len() int {
	return len(e)
}

func (e EventList) Less(i, j int) bool {
	return e[j].TimeStamp.After(e[i].TimeStamp)
}

func (e EventList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func loadFile(filename string) EventList {
	var eventPattern, err = regexp.Compile(`\[(\d{4})-(\d{2})-(\d{2})\s(\d{2}):(\d{2})\]\s(Guard #(\d+) begins shift|wakes up|falls asleep)`)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	events := make(EventList, 0)

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		trimmedLine := strings.TrimSpace(string(line))
		events = append(events, parseLine(eventPattern, trimmedLine))
	}
	sort.Sort(events)

	var guardId int
	for _, event := range events {
		if event.GuardId > 0 {
			guardId = event.GuardId
		} else {
			event.GuardId = guardId
		}
	}

	return events
}

func parseLine(p *regexp.Regexp, s string) *Event {
	matches := p.FindStringSubmatch(s)
	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])
	hour, _ := strconv.Atoi(matches[4])
	minute, _ := strconv.Atoi(matches[5])
	event := BeginShift
	guardId, err := strconv.Atoi(matches[7])
	if err != nil {
		if matches[6] == "wakes up" {
			event = WakeUp
		} else if matches[6] == "falls asleep" {
			event = FallAsleep
		}
	}
	return &Event{event,
		time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC),
		guardId}
}

type Schedule struct {
	GuardId int
	Minutes []int
}

func (s *Schedule) wakeUp(minute int) {
	s.Minutes[minute] = -1
}

func (s *Schedule) fallAsleep(minute int) {
	s.Minutes[minute] = 1
}

func (s *Schedule) finalize() {
	for i := 1; i < minutes; i++ {
		s.Minutes[i] = s.Minutes[i-1] + s.Minutes[i]
	}
}

func MakeSchedule(guardId int) *Schedule {
	return &Schedule{guardId, make([]int, minutes)}
}

func part1(events EventList) {
	grid := initializeGrid(events)
	guardIds := uniqueGuardIds(events)
	findGuardMostAsleep(guardIds, grid)
	//for _, id := range guardIds {
	//	sum := sumGuard(id, grid)
	//	fmt.Printf("%v : %v\n", id, sum)
	//}
}

func part2(events EventList) {
	grid := initializeGrid(events)
	guardIds := uniqueGuardIds(events)
	guardMostAsleep := 0
	maxSleep := 0
	maxMinute := 0
	for _, id := range guardIds {
		sum := sumGuard(id, grid)
		for m, v := range sum {
			if v > maxSleep {
				guardMostAsleep = id
				maxSleep = v
				maxMinute = m
			}
		}
	}
	fmt.Printf("GuardId = %[1]v, Minute = %[2]v (%[4]v), %[1]v * %[2]v = %[3]v\n",
		guardMostAsleep, maxMinute, guardMostAsleep*maxMinute, maxSleep)
}

func findGuardMostAsleep(guardIds []int, grid []*Schedule) (guardId, frequency, minute, total int) {
	type Result struct {
		index     int
		frequency int
		total     int
	}
	result := map[int]Result{}

	for _, id := range guardIds {
		sum := sumGuard(id, grid)
		maxIdx := 0
		maxFreq := 0
		total := 0
		for i := 0; i < minutes; i++ {
			total += sum[i]
			if sum[i] > maxFreq {
				maxIdx = i
				maxFreq = sum[i]
			}
		}
		result[id] = Result{maxIdx, maxFreq, total}
	}

	max := 0
	maxId := 0
	var maxResult Result
	for k, r := range result {
		if r.total > max {
			max = r.total
			maxResult = r
			maxId = k
		}
	}

	fmt.Println(maxResult)
	fmt.Printf("id = %[1]v, minute = %[2]v, %[1]v * %[2]v = %[3]v\n", maxId, maxResult.index,
		maxId*maxResult.index)

	return 0, 0, 0, 0
}

func sumGuard(guardId int, grid []*Schedule) []int {
	sum := make([]int, minutes)

	for _, s := range grid {
		if s.GuardId == guardId {
			for i := 0; i < minutes; i++ {
				sum[i] += s.Minutes[i]
			}
		}
	}

	return sum
}

func uniqueGuardIds(events EventList) []int {
	ids := map[int]bool{}
	for _, event := range events {
		ids[event.GuardId] = true
	}
	keys := make([]int, 0, len(ids))
	for k := range ids {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func initializeGrid(events EventList) []*Schedule {
	// Initialize data structure.
	shiftStarts := 0
	for _, event := range events {
		if event.EventType == BeginShift {
			shiftStarts++
		}
	}
	grid := make([]*Schedule, shiftStarts)
	row := -1
	for _, event := range events {
		// Always awake to start with.
		switch event.EventType {
		case BeginShift:
			if row > -1 {
				grid[row].finalize()
			}
			row++
			grid[row] = MakeSchedule(event.GuardId)
		case WakeUp:
			grid[row].wakeUp(event.TimeStamp.Minute())
		case FallAsleep:
			grid[row].fallAsleep(event.TimeStamp.Minute())
		}
	}
	grid[row].finalize()
	return grid
}

func main() {
	events := loadFile("input.txt")
	part1(events)
	part2(events)
}
