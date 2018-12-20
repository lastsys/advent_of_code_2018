package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Tile rune

type Grid struct {
	Width  int
	Height int
	Tiles  [][]Tile
	Buffer [][]Tile
}

func (g *Grid) ResourceValue() int {
	wooded := 0
	lumberyard := 0
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			switch g.Tiles[y][x] {
			case '|':
				wooded++
			case '#':
				lumberyard++
			}
		}
	}
	return wooded * lumberyard
}

func (g *Grid) Equal(g2 *Grid) bool {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Tiles[y][x] != g2.Tiles[y][x] {
				return false
			}
		}
	}
	return true
}

func (g *Grid) Print() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			fmt.Printf("%c", g.Tiles[y][x])
		}
		fmt.Println()
	}
}

func (g *Grid) Step() {
	for y := 0; y < g.Height; y++ {
		g.Buffer[y] = make([]Tile, g.Width)
	}

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			switch g.Tiles[y][x] {
			case '.':
				if g.countAdjacent(x, y, '|') >= 3 {
					g.Buffer[y][x] = '|'
				} else {
					g.Buffer[y][x] = '.'
				}
			case '|':
				if g.countAdjacent(x, y, '#') >= 3 {
					g.Buffer[y][x] = '#'
				} else {
					g.Buffer[y][x] = '|'
				}
			case '#':
				if g.countAdjacent(x, y, '|') >= 1 && g.countAdjacent(x, y, '#') >= 1 {
					g.Buffer[y][x] = '#'
				} else {
					g.Buffer[y][x] = '.'
				}
			}
		}
	}
	// Double buffering to avoid allocating memory over and over again.
	g.Tiles, g.Buffer = g.Buffer, g.Tiles
}

func (g *Grid) countAdjacent(x, y int, tile Tile) int {
	xMin := x - 1
	xMax := x + 1
	yMin := y - 1
	yMax := y + 1
	if xMin < 0 {
		xMin = 0
	}
	if xMax >= g.Width {
		xMax = g.Width - 1
	}
	if yMin < 0 {
		yMin = 0
	}
	if yMax >= g.Height {
		yMax = g.Height - 1
	}
	count := 0
	for yy := yMin; yy <= yMax; yy++ {
		for xx := xMin; xx <= xMax; xx++ {
			if yy == y && xx == x {
				continue
			}
			if g.Tiles[yy][xx] == tile {
				count++
			}
		}
	}
	return count
}

func loadData(filename string) *Grid {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines := make([][]Tile, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, []Tile(scanner.Text()))
	}

	buffer := make([][]Tile, len(lines))
	for y := 0; y < len(lines); y++ {
		buffer[y] = make([]Tile, len(lines[0]))
	}

	return &Grid{len(lines[0]), len(lines), lines, buffer}
}

func main() {
	grid := loadData("input.txt")
	for i := 0; i < 10; i++ {
		grid.Step()
	}
	fmt.Printf("Resource value after 10 minutes = %v.\n", grid.ResourceValue())

	// Part 2:

	grid = loadData("input.txt")
	const iterations = 1000
	resourceHistory := make([]int, iterations+1)
	resourceHistory[0] = grid.ResourceValue()
	for i := 1; i <= iterations; i++ {
		grid.Step()
		resourceHistory[i] = grid.ResourceValue()
		//fmt.Printf("%4d : RV = %d\n", i, grid.ResourceValue())
	}
	fmt.Println(resourceHistory)

	cycleLength := 0
	value := resourceHistory[len(resourceHistory)-1]
	for i := len(resourceHistory) - 2; i >= 0; i-- {
		if value == resourceHistory[i] {
			cycleLength = len(resourceHistory) - i - 1
			break
		}
	}
	fmt.Println("Cycle length =", cycleLength)
	cycle := resourceHistory[len(resourceHistory)-cycleLength:]
	fmt.Println(cycle)

	// Distance from first element in cycle.
	offset := 1000000000 - len(resourceHistory) - len(cycle)
	offset %= len(cycle)

	fmt.Printf("Resource value after 1000000000 minutes = %v.\n", cycle[offset])
}
