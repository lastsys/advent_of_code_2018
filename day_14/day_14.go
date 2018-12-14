package main

import (
	"fmt"
)

type Score uint8

type State struct {
	elfPositions []int
	scoreboard   []Score
}

func NewState(elfCount int, scoreboard []Score) *State {
	elfPositions := make([]int, elfCount)
	for i := 0; i < elfCount; i++ {
		elfPositions[i] = i
	}
	return &State{elfPositions, scoreboard}
}

func (s *State) Step() {
	// Generate new recipes.
	sum := 0
	for _, elfPosition := range s.elfPositions {
		sum += int(s.scoreboard[elfPosition])
	}

	// Special case when the sum is zero.
	if sum == 0 {
		s.scoreboard = append(s.scoreboard, 0)
	} else {
		// Extract all new recipes.
		newScores := make([]Score, 0, 4)
		for {
			// Extract last digit in number.
			score := Score(sum % 10)

			newScores = append(newScores, score)

			// Shift down one digit.
			sum /= 10

			// If we have zero we are done.
			if sum == 0 {
				break
			}
		}
		// Reverse the order of the recipes.
		for i, j := 0, len(newScores)-1; i < j; i, j = i+1, j-1 {
			newScores[i], newScores[j] = newScores[j], newScores[i]
		}
		s.scoreboard = append(s.scoreboard, newScores...)
	}
	// Move elves.
	for elfId, elfPosition := range s.elfPositions {
		s.elfPositions[elfId] += 1 + int(s.scoreboard[elfPosition])
		s.elfPositions[elfId] %= len(s.scoreboard)
	}
	// Done.
}

func (s *State) Print() {
	var startChar, endChar rune
	for i, score := range s.scoreboard {
		startChar, endChar = ' ', ' '
		for elfId, elfPosition := range s.elfPositions {
			if elfPosition == i {
				switch elfId {
				case 0:
					startChar, endChar = '(', ')'
				case 1:
					startChar, endChar = '[', ']'
				default:
					startChar, endChar = '*', '*'
				}
				break
			}
		}
		fmt.Printf("%c%d%c", startChar, score, endChar)
	}
	fmt.Println()
}

func part1(scoreboard []Score, recipesUntilAnswer int) {
	state := NewState(2, scoreboard)
	for {
		state.Step()
		if len(state.scoreboard) >= recipesUntilAnswer+10 {
			break
		}
	}
	the10 := state.scoreboard[recipesUntilAnswer : recipesUntilAnswer+10]
	fmt.Print("Part 1: ")
	for _, v := range the10 {
		fmt.Printf("%v", v)
	}
	fmt.Println()
}

func (s *State) PatternMatch(pattern []Score, offset int) bool {
	for i, j := offset, 0; i < offset+len(pattern); i, j = i+1, j+1 {
		if s.scoreboard[i] != pattern[j] {
			return false
		}
	}
	return true
}

func part2(scoreboard []Score, pattern []Score) {
	state := NewState(2, scoreboard)
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

	fmt.Printf("Part 2: %v", offset)
}

func main() {
	scoreboard := []Score{3, 7}
	part1(scoreboard, 880751)
	part2(scoreboard, []Score{8, 8, 0, 7, 5, 1})
}
