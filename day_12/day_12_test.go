package main

import "testing"

func TestRules(t *testing.T) {
	rules, initialState := loadData("test_input.txt")
	const steps = 20
	var s State
	history := make(StateList, 0, steps)
	history = append(history, initialState)
	s = initialState
	for i := 1; i <= steps; i++ {
		s = s.ApplyRuleSet(rules)
		history = append(history, s)
	}
	if sum := history[steps].Sum(); sum != 325 {
		t.Errorf("Got sum %v, expected 325.", sum)
	}
}
