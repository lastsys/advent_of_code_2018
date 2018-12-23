package main

import (
	"bufio"
	"fmt"
	"golang.org/x/tools/container/intsets"
	"log"
	"os"
	"strings"
)

func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

type Vector3 [3]int

func (v *Vector3) ManhattanDistance(v2 *Vector3) int {
	return Abs(v[0]-v2[0]) + Abs(v[1]-v2[1]) + Abs(v[2]-v2[2])
}

type Nanobot struct {
	Position Vector3
	Radius   int
}

type Nanobots []*Nanobot

func (b Nanobots) MaxRadius() *Nanobot {
	maxRadius := intsets.MinInt
	var maxBot *Nanobot
	for _, bot := range b {
		if bot.Radius > maxRadius {
			maxRadius = bot.Radius
			maxBot = bot
		}
	}
	return maxBot
}

func (b Nanobots) Span() (Vector3, Vector3) {
	minX := intsets.MaxInt
	minY := intsets.MaxInt
	minZ := intsets.MaxInt
	maxX := intsets.MinInt
	maxY := intsets.MinInt
	maxZ := intsets.MinInt
	for _, bot := range b {
		if bot.Position[0] < minX {
			minX = bot.Position[0]
		}
		if bot.Position[0] > maxX {
			maxX = bot.Position[0]
		}
		if bot.Position[1] < minY {
			minY = bot.Position[1]
		}
		if bot.Position[1] > maxY {
			maxY = bot.Position[1]
		}
		if bot.Position[2] < minZ {
			minZ = bot.Position[2]
		}
		if bot.Position[2] > maxZ {
			maxZ = bot.Position[2]
		}
	}

	return Vector3{minX, minY, minZ}, Vector3{maxX, maxY, maxZ}
}

func (b Nanobots) InRangeOf(bot *Nanobot) Nanobots {
	botsInRange := Nanobots{}
	for _, x := range b {
		if bot.Position.ManhattanDistance(&x.Position) <= bot.Radius {
			botsInRange = append(botsInRange, x)
		}
	}
	return botsInRange
}

func loadData(filename string) Nanobots {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	bots := Nanobots{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		bot := Nanobot{}
		n, err := fmt.Sscanf(line, "pos=<%d,%d,%d>, r=%d",
			&bot.Position[0], &bot.Position[1], &bot.Position[2], &bot.Radius)
		if err != nil {
			log.Println(line)
			log.Fatal(n, err)
		}
		bots = append(bots, &bot)
	}
	return bots
}

type IntArray []int

func (a IntArray) Max() int {
	max := intsets.MinInt
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

func (a IntArray) Min() int {
	min := intsets.MaxInt
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}

func part2(bots Nanobots) {
	xs := make(IntArray, 0)
	ys := make(IntArray, 0)
	zs := make(IntArray, 0)

	for _, bot := range bots {
		xs = append(xs, bot.Position[0])
		ys = append(ys, bot.Position[1])
		zs = append(zs, bot.Position[2])
	}

	dist := 1
	for dist < xs.Max()-xs.Min() {
		dist *= 2
	}

	for {
		targetCount := 0
		var best Vector3
		var bestVal int
		for x := xs.Min(); x <= xs.Max(); x += dist {
			for y := ys.Min(); y <= ys.Max(); y += dist {
				for z := zs.Min(); z <= zs.Max(); z += dist {
					count := 0
					for _, bot := range bots {
						v := Vector3{x, y, z}
						d := v.ManhattanDistance(&bot.Position)
						if (d-bot.Radius)/dist <= 0 {
							count++
						}
					}
					if count > targetCount {
						targetCount = count
						bestVal = Abs(x) + Abs(y) + Abs(z)
						best = Vector3{x, y, z}
					} else if count == targetCount {
						if Abs(x)+Abs(y)+Abs(z) < bestVal {
							bestVal = Abs(x) + Abs(y) + Abs(z)
							best = Vector3{x, y, z}
						}
					}
				}
			}
		}
		if dist == 1 {
			fmt.Println(bestVal)
			fmt.Println(best)
			break
		} else {
			xs = IntArray{best[0] - dist, best[0] + dist}
			ys = IntArray{best[1] - dist, best[1] + dist}
			zs = IntArray{best[2] - dist, best[2] + dist}
			dist /= 2
		}
	}
}

func main() {
	bots := loadData("input.txt")
	maxBot := bots.MaxRadius()
	inRange := bots.InRangeOf(maxBot)
	fmt.Printf("Bot with range %v has %v bots in its range.\n", maxBot.Radius, len(inRange))

	part2(bots)
}
