package main

import (
	"testing"
)

func expectScore(t *testing.T, initial []Score, iterations int, expected []Score) {
	state := NewState(2, initial)
	for {
		state.Step()
		if len(state.scoreboard) >= iterations+10 {
			break
		}
	}
	the10 := state.scoreboard[iterations : iterations+10]
	equal := true
	for i := 0; i < 10; i++ {
		if the10[i] != expected[i] {
			equal = false
			break
		}
	}
	if !equal {
		t.Error(state.scoreboard)
		t.Errorf("The 10 %v not equal to %v", the10, expected)
	}
}

func expectLen(t *testing.T, initial, pattern []Score, expected int) {
	state := NewState(2, initial)
	offset := 0
outside:
	for {
		lengthBefore := len(state.scoreboard)
		state.Step()
		lengthAfter := len(state.scoreboard)
		if lengthAfter <= len(pattern) {
			continue
		}
		for i := lengthBefore; i < lengthAfter; i++ {
			if state.PatternMatch(pattern, i-len(pattern)+1) {
				offset = i - len(pattern) + 1
				break outside
			}
		}
	}

	if offset != expected {
		t.Errorf("Got %v expected %v", offset, expected)
	}
}

func TestSteps(t *testing.T) {
	expectScore(t, []Score{3, 7}, 9, []Score{5, 1, 5, 8, 9, 1, 6, 7, 7, 9})
	expectScore(t, []Score{3, 7}, 5, []Score{0, 1, 2, 4, 5, 1, 5, 8, 9, 1})
	expectScore(t, []Score{3, 7}, 18, []Score{9, 2, 5, 1, 0, 7, 1, 0, 8, 5})
	expectScore(t, []Score{3, 7}, 2018, []Score{5, 9, 4, 1, 4, 2, 9, 8, 8, 2})

	expectLen(t, []Score{3, 7}, []Score{5, 1, 5, 8, 9}, 9)
	expectLen(t, []Score{3, 7}, []Score{0, 1, 2, 4, 5}, 5)
	expectLen(t, []Score{3, 7}, []Score{9, 2, 5, 1, 0}, 18)
	expectLen(t, []Score{3, 7}, []Score{5, 9, 4, 1, 4}, 2018)
}
