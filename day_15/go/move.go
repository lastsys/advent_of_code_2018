package main

import (
	"container/list"
	"errors"
	"golang.org/x/tools/container/intsets"
	"sort"
)

var deltas = []Vector2{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

func FullCombat(grid *Grid) int {
	i := 0
	for {
		i++
		incomplete := Step(grid)
		if incomplete {
			i--
		}
		win, _, hp := grid.units.WinCondition()
		if win {
			return i * hp
		}
	}
}

// Take a step forward in time.
// Return true if incomplete due to winning condition.
func Step(grid *Grid) bool {
	moveOrder := grid.UnitReadOrder()
	for _, unit := range moveOrder {
		if win, _, _ := grid.units.WinCondition(); win {
			return true
		}
		// Do not handle units which were eliminated during the step.
		if unit.hp <= 0 {
			continue
		}
		if !combat(grid, unit.id) {
			next := nextPosition(grid, unit.id)
			grid.MoveUnit(unit.id, next)
			// May attack after a move.
			combat(grid, unit.id)
		}
	}
	return false
}

// Find the shortest distance from a starting position among a set of other locations.
// If more than one location is closest then choose the first in given reading order.
// If more than one direction from the starting position has the same distance to the
// nearest target the direction which has the first reading order should be selected.
func nextPosition(grid *Grid, id UnitId) Vector2 {
	// Find targets in range.
	targetsInRange := findTargetPositionsInRange(grid, id)
	if len(targetsInRange) == 0 {
		// Remain still if no targets are in range.
		return grid.units[id].position
	}
	// Find out which of these targets are reachable.
	reachableTargets := filterReachable(grid, grid.units[id].position, targetsInRange)
	if len(reachableTargets) == 0 {
		return grid.units[id].position
	}
	// Find out which of the reachable targets are nearest.
	nearestTargets := findNearestTargets(grid, grid.units[id].position, reachableTargets)
	// Get the target which is first in the reading order.
	sort.Sort(nearestTargets)
	nearestTarget := nearestTargets[0]
	// Find out which direction to go by checking the distance between all four possible
	// directions to walk in and choose the shortest one. If there are several which
	// are shortest we choose the first in reading order.
	startLocations := make(VectorList, 0)
	for _, delta := range deltas {
		location := grid.units[id].position.Add(delta)
		if grid.tiles[location[1]][location[0]] != openTile {
			continue
		}
		startLocations = append(startLocations, location)
	}
	reachableStartLocations := filterReachable(grid, nearestTarget, startLocations)
	bestDirections := findNearestTargets(grid, nearestTarget, reachableStartLocations)
	sort.Sort(bestDirections)

	return bestDirections[0]
}

// Find all potential targets which are located on an open tile.
// The targets are not necessarily reachable.
func findTargetPositionsInRange(grid *Grid, id UnitId) VectorList {
	currentUnit := grid.units[id]
	targetUnits := make(UnitList, 0)
	for _, unit := range grid.units {
		if unit != nil && unit.race != currentUnit.race {
			targetUnits = append(targetUnits, unit)
		}
	}

	targetPositions := make(VectorList, 0)
	for _, unit := range targetUnits {
		for _, delta := range deltas {
			potentialTargetPosition := unit.position.Add(delta)
			if grid.tiles[potentialTargetPosition[1]][potentialTargetPosition[0]] == openTile {
				targetPositions = append(targetPositions, potentialTargetPosition)
			}
		}
	}

	return targetPositions
}

// Return reachable destinations.
func filterReachable(grid *Grid, start Vector2, targets VectorList) VectorList {
	frontier := list.New()
	frontier.PushBack(start)

	visited := make(map[Vector2]bool, grid.width*grid.height)

	reachableTargets := make(VectorList, 0)

	targetSet := make(map[Vector2]bool, len(targets))
	for _, target := range targets {
		targetSet[target] = true
	}

	for {
		if frontier.Len() == 0 {
			break
		}
		element := frontier.Front()
		frontier.Remove(element)
		currentLocation := element.Value.(Vector2)
		// Check if this is one of the targets.
		if _, ok := targetSet[currentLocation]; ok {
			reachableTargets = append(reachableTargets, currentLocation)
		}
		for _, delta := range deltas {
			neighbor := currentLocation.Add(delta)
			if grid.tiles[neighbor[1]][neighbor[0]] != openTile {
				continue
			}
			if _, ok := visited[neighbor]; !ok {
				frontier.PushBack(neighbor)
				visited[neighbor] = true
			}
		}
	}

	return reachableTargets
}

// Find the shortest distance between two locations.
func shortestDistance(grid *Grid, start Vector2, end Vector2) (int, error) {
	frontier := list.New()
	frontier.PushBack(start)
	cameFrom := make(map[Vector2]Vector2)

	for {
		if frontier.Len() == 0 {
			break
		}
		element := frontier.Front()
		frontier.Remove(element)
		currentLocation := element.Value.(Vector2)
		if currentLocation == end {
			distance := 0
			for {
				if currentLocation == start {
					return distance, nil
				}
				currentLocation = cameFrom[currentLocation]
				distance++
			}
		}
		for _, delta := range deltas {
			neighbor := currentLocation.Add(delta)
			if grid.tiles[neighbor[1]][neighbor[0]] != openTile {
				continue
			}
			if _, ok := cameFrom[neighbor]; !ok {
				frontier.PushBack(neighbor)
				cameFrom[neighbor] = currentLocation
			}
		}
	}

	return 0, errors.New("not reachable")
}

// Find nearest targets among a set of targets.
func findNearestTargets(grid *Grid, start Vector2, targets VectorList) VectorList {
	minDistance := intsets.MaxInt
	targetDistances := make([]int, len(targets))
	for i, target := range targets {
		d, err := shortestDistance(grid, start, target)
		if err != nil {
			continue
		}
		if d < minDistance {
			minDistance = d
		}
		targetDistances[i] = d
	}
	nearestTargets := make(VectorList, 0)
	for i, distance := range targetDistances {
		if minDistance == distance {
			nearestTargets = append(nearestTargets, targets[i])
		}
	}
	return nearestTargets
}

// Do combat if possible.
func combat(grid *Grid, id UnitId) bool {
	unit := grid.units[id]
	enemies := make(HpUnitList, 0)
	for _, delta := range deltas {
		neighbor := unit.position.Add(delta)
		for i := 0; i < len(grid.units); i++ {
			otherUnit := grid.units[i]
			if otherUnit != nil && otherUnit.race != unit.race && otherUnit.position == neighbor {
				enemies = append(enemies, otherUnit)
			}
		}
	}
	if len(enemies) == 0 {
		return false
	}
	sort.Sort(enemies)
	attackedEnemy := enemies[0]
	if unit.race == Elf {
		attackedEnemy.hp -= grid.elfAttackPower
	} else {
		attackedEnemy.hp -= 3
	}
	if attackedEnemy.hp <= 0 {
		// Enemy is dead. Remove.
		grid.RemoveUnit(attackedEnemy.id)
	}
	return true
}
