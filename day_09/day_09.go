package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

type Game struct {
	players          int
	lastMarblePoints int
	highScore        int
}

func insertMarbleIntoCircle(circle *list.List, index *list.Element, marble int) *list.Element {
	newIndex := index.Next()
	if newIndex == nil {
		newIndex = circle.Front()
	}
	circle.InsertAfter(marble, newIndex)
	return newIndex.Next()
}

func handleMultipleOf23(circle *list.List, index *list.Element, marble int) (int, *list.Element) {
	var newIndex = index
	for i := 0; i < 7; i++ {
		newIndex = newIndex.Prev()
		if newIndex == nil {
			newIndex = circle.Back()
		}
	}
	score := newIndex.Value.(int)
	nextIndex := newIndex.Next()
	circle.Remove(newIndex)
	score += marble
	return score, nextIndex
}

func loadData(filename string, test bool) []Game {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	games := make([]Game, 0)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		game := Game{}
		if test {
			_, err := fmt.Sscanf(scanner.Text(), "%d players; last marble is worth %d points: high score is %d",
				&game.players, &game.lastMarblePoints, &game.highScore)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err := fmt.Sscanf(scanner.Text(), "%d players; last marble is worth %d points",
				&game.players, &game.lastMarblePoints)
			if err != nil {
				log.Fatal(err)
			}
		}
		games = append(games, game)
	}
	return games
}

func play(game *Game) int {
	circle := list.New()
	circle.PushBack(0)
	index := circle.Back()

	//fmt.Print("[---] ")
	//circle.Print(0)
	score := make([]int, game.players)
	player := 0
	var tempScore int
	for i := 1; i <= game.lastMarblePoints; i++ {
		// Special rule if marble number is a multiple of 23.
		if i%23 == 0 {
			tempScore, index = handleMultipleOf23(circle, index, i)
			score[player] += tempScore
		} else {
			index = insertMarbleIntoCircle(circle, index, i)
		}
		player = (player + 1) % game.players
	}
	highScore := 0
	for _, s := range score {
		if s > highScore {
			highScore = s
		}
	}
	return highScore
}

func part1(games []Game) {
	score := play(&games[0])
	fmt.Printf("Part 1 score: %v\n", score)
}

func part2(games []Game) {
	largerGame := &Game{
		games[0].players,
		games[0].lastMarblePoints * 100,
		0,
	}
	score := play(largerGame)
	fmt.Printf("Part 2 score: %v\n", score)
}

func main() {
	games := loadData("input.txt", false)
	part1(games)
	part2(games)
}
