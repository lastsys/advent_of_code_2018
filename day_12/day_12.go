package main

import (
	"bufio"
	"fmt"
	"golang.org/x/tools/container/intsets"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
)

type State map[int]bool
type StateList []State

func (sl StateList) Print() {
	min := intsets.MaxInt
	max := intsets.MinInt
	for _, state := range sl {
		for k := range state {
			if k < min {
				min = k
			}
			if k > max {
				max = k
			}
		}
	}
	min--
	max++
	for index, state := range sl {
		state.Print(index, min, max)
	}
}

func (s State) Get(i int) bool {
	if v, ok := s[i]; !ok {
		return false
	} else {
		return v
	}
}

func (s State) applyRule(rule *Rule, i int) (bool, bool) {
	if s.Get(i-2) == rule.pattern[0] &&
		s.Get(i-1) == rule.pattern[1] &&
		s.Get(i) == rule.pattern[2] &&
		s.Get(i+1) == rule.pattern[3] &&
		s.Get(i+2) == rule.pattern[4] {
		return rule.result, true
	}
	return false, false
}

func (s State) ApplyRuleSet(rules RuleSet) State {
	newState := make(State, len(s))
	min := intsets.MaxInt
	max := intsets.MinInt
	for k := range s {
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
	}
	for i := min - 2; i < max+2; i++ {
		for _, rule := range rules {
			result, ok := s.applyRule(rule, i)
			if ok {
				if result {
					newState[i] = result
				}
				break
			}
		}
	}
	return newState
}

func (s State) Print(index, min, max int) {
	for k := range s {
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
	}
	fmt.Printf("%5v: ", index)
	for i := min; i <= max; i++ {
		if s.Get(i) {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()
}

func (s State) Sum() int {
	sum := 0
	for k := range s {
		sum += k
	}
	return sum
}

type Rule struct {
	pattern [5]bool
	result  bool
}

type RuleSet []*Rule

func loadData(filename string) (RuleSet, State) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	initialStateRow, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	initialStatePattern, err := regexp.Compile(`[#\.]+`)
	if err != nil {
		log.Fatal(err)
	}
	initialState := []byte(initialStatePattern.FindString(initialStateRow))
	state := make(State)
	for i, s := range initialState {
		if s == '#' {
			state[i] = true
		}
	}

	_, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	rulePattern, err := regexp.Compile(`([#\.]{5}) => ([#\.])`)
	rules := make(RuleSet, 0)
	for {
		row, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		match := rulePattern.FindStringSubmatch(row)
		rule := &Rule{}
		for i, s := range []byte(match[1]) {
			if s == '#' {
				rule.pattern[i] = true
			}
		}
		if match[2] == "#" {
			rule.result = true
			rules = append(rules, rule)
		}
	}
	return rules, state
}

func part1(rules RuleSet, initialState State) {
	const steps = 20
	var s State
	history := make(StateList, 0, steps)
	history = append(history, initialState)
	s = initialState
	for i := 1; i <= steps; i++ {
		s = s.ApplyRuleSet(rules)
		history = append(history, s)
	}
	history.Print()
	fmt.Printf("Part 1: Sum = %v\n", history[20].Sum())
}

func part2(rules RuleSet, initialState State) {
	const steps = 120
	const maxSteps = 50000000000
	history := make(StateList, 0, steps)
	history = append(history, initialState)
	s := initialState
	for i := 1; i <= steps; i++ {
		s2 := s.ApplyRuleSet(rules)
		history = append(history, s2)
		s = s2
	}
	history.Print()

	// Takes too long to run all steps.
	// Just run until stable and then calculate an offset.

	a := make([]int, 0, len(s))
	for k := range s {
		a = append(a, k)
	}
	sort.Ints(a)
	fmt.Println(a)
	offset := maxSteps - steps
	sum := 0
	for _, v := range a {
		sum += v + offset
	}
	fmt.Printf("Part 2: Sum = %v\n", sum)
}

func main() {
	rules, state := loadData("input.txt")
	part1(rules, state)
	rules, state = loadData("input.txt")
	part2(rules, state)
}
