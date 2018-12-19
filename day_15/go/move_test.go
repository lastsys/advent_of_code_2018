package main

import (
	"testing"
)

func TestFindTargetPositionsInRange(t *testing.T) {
	grid := loadData("test_target.txt")
	targets := findTargetPositionsInRange(grid, 0)
	targetSet := make(map[Vector2]bool, len(targets))
	for _, target := range targets {
		targetSet[target] = true
	}

	expectedTargets := []Vector2{{3, 1}, {5, 1}, {2, 2}, {5, 2}, {1, 3}, {3, 3}}

	if len(expectedTargets) != len(targetSet) {
		t.Errorf("Expected to have %v targets, got %v.", len(expectedTargets), len(targetSet))
	}

	for _, expectedTarget := range expectedTargets {
		if _, ok := targetSet[expectedTarget]; !ok {
			t.Errorf("Expected to have %v in set.", expectedTarget)
		}
	}
}

func TestFindReachableTargets(t *testing.T) {
	grid := loadData("test_target.txt")
	targetsInRange := findTargetPositionsInRange(grid, 0)
	reachableTargets := filterReachable(grid, grid.units[0].position, targetsInRange)
	targetSet := make(map[Vector2]bool, len(reachableTargets))
	for _, target := range reachableTargets {
		targetSet[target] = true
	}

	expectedTargets := []Vector2{{3, 1}, {2, 2}, {1, 3}, {3, 3}}

	if len(expectedTargets) != len(targetSet) {
		t.Errorf("Expected to have %v targets, got %v.", len(expectedTargets), len(targetSet))
	}

	for _, expectedTarget := range expectedTargets {
		if _, ok := targetSet[expectedTarget]; !ok {
			t.Errorf("Expected to have %v in set.", expectedTarget)
		}
	}
}

func TestShortestDistance(t *testing.T) {
	grid := loadData("test_target.txt")
	if d, _ := shortestDistance(grid, Vector2{1, 2}, Vector2{3, 3}); d != 3 {
		t.Errorf("Expected distance 3, got %v.", d)
	}
}

func TestNearestTargets(t *testing.T) {
	grid := loadData("test_target.txt")
	targetsInRange := findTargetPositionsInRange(grid, 0)
	reachableTargets := filterReachable(grid, grid.units[0].position, targetsInRange)
	nearestTargets := findNearestTargets(grid, Vector2{1, 1}, reachableTargets)
	if l := len(nearestTargets); l != 3 {
		t.Errorf("Expected 3 targets, got %v.", l)
	}
	targetSet := make(map[Vector2]bool, len(nearestTargets))
	for _, target := range nearestTargets {
		targetSet[target] = true
	}

	expectedTargets := []Vector2{{3, 1}, {2, 2}, {1, 3}}

	for _, expectedTarget := range expectedTargets {
		if _, ok := targetSet[expectedTarget]; !ok {
			t.Errorf("Expected to have %v in set.", expectedTarget)
		}
	}
}

func TestNextPosition(t *testing.T) {
	grid := loadData("test_target.txt")
	nextPos := nextPosition(grid, 0)
	expectedPos := Vector2{2, 1}
	if nextPos != expectedPos {
		t.Errorf("Expected (2,1) got %v.", nextPos)
	}
}

func TestMove(t *testing.T) {
	grid := loadData("test_movement_0.txt")
	grid.Print()
	Step(grid)
	grid.Print()
	Step(grid)
	grid.Print()
	Step(grid)
	grid.Print()
}

func expectCombat(t *testing.T, filename string, expectedOutcome int) {
	grid := loadData(filename)
	if outcome := FullCombat(grid); outcome != expectedOutcome {
		t.Errorf("For %v we expected %v, but got %v.", filename, expectedOutcome, outcome)
	}
}

func TestSummarizedCombats(t *testing.T) {
	expectCombat(t, "combat_start_1.txt", 27730)
	expectCombat(t, "summarized_combat_1.txt", 36334)
	expectCombat(t, "summarized_combat_2.txt", 39514)
	expectCombat(t, "summarized_combat_3.txt", 27755)
	expectCombat(t, "summarized_combat_4.txt", 28944)
	expectCombat(t, "summarized_combat_5.txt", 18740)
}
