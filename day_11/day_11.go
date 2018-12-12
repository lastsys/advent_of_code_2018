package main

import "fmt"

type Grid [][]int

func NewGrid(width, height int) Grid {
	grid := make([][]int, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]int, width)
	}
	return grid
}

func (g Grid) Initialize(serial int) {
	for y := 0; y < 300; y++ {
		for x := 0; x < 300; x++ {
			g[y][x] = powerLevel(x+1, y+1, serial)
		}
	}
}

func (g Grid) FindMaxSquare(width, height int) (int, int, int) {
	var max, mx, my int
	for y := 0; y < len(g)-height; y++ {
		for x := 0; x < len(g[0])-width; x++ {
			sum := 0
			for yy := y; yy < y+height; yy++ {
				for xx := x; xx < x+width; xx++ {
					sum += g[yy][xx]
				}
			}
			if sum > max {
				max = sum
				mx = x
				my = y
			}
		}
	}
	return mx + 1, my + 1, max
}

func (g Grid) FindMaxSquare2(maxSize int) (int, int, int, int) {
	var max, mx, my, msz int
	for sz := 1; sz < maxSize; sz++ {
		fmt.Printf("%v%%\r", 100*sz/maxSize)
		x, y, power := g.FindMaxSquare(sz, sz)
		if power > max {
			max = power
			mx = x
			my = y
			msz = sz
		}
	}
	fmt.Println()
	return mx, my, max, msz
}

func powerLevel(x, y, serial int) int {
	rackId := x + 10
	powerLevel := rackId*y + serial
	powerLevel *= rackId
	if powerLevel >= 100 {
		powerLevel = (powerLevel / 100) % 10
	}
	powerLevel -= 5
	return powerLevel
}

func part1(serial int) {
	grid := NewGrid(300, 300)
	grid.Initialize(serial)
	x, y, _ := grid.FindMaxSquare(3, 3)
	fmt.Println(x, y)
}

func part2(serial int) {
	grid := NewGrid(300, 300)
	grid.Initialize(serial)
	x, y, _, sz := grid.FindMaxSquare2(300)
	fmt.Println(x, y, sz)
}

func main() {
	serial := 7989
	part1(serial)
	part2(serial)
}
