package main

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/tools/container/intsets"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Tile rune
type Grid struct {
	MinY  int
	MaxY  int
	MinX  int
	MaxX  int
	Tiles map[int]map[int]Tile
}

func NewGrid() Grid {
	grid := Grid{}
	grid.Tiles = map[int]map[int]Tile{}
	grid.Set(500, 0, '+')
	return grid
}

func (g Grid) Set(x, y int, tile Tile) {
	if tile == '.' {
		if _, ok := g.Tiles[y]; ok {
			delete(g.Tiles[y], x)
			if len(g.Tiles[y]) == 0 {
				delete(g.Tiles, y)
			}
		}
	} else {
		if _, ok := g.Tiles[y]; !ok {
			g.Tiles[y] = map[int]Tile{}
		}
		g.Tiles[y][x] = tile
	}
}

func (g Grid) Get(x, y int) Tile {
	if _, ok := g.Tiles[y]; !ok {
		return '.'
	}
	if _, ok := g.Tiles[y][x]; !ok {
		return '.'
	}
	return g.Tiles[y][x]
}

func (g Grid) Print() {
	for y := g.MinY - 1; y <= g.MaxY+1; y++ {
		fmt.Printf("%4d ", y)
		for x := g.MinX - 1; x <= g.MaxX+1; x++ {
			fmt.Printf("%c", g.Get(x, y))
		}
		fmt.Println()
	}
}

var SearchError = errors.New("search error")

func (g Grid) EnclosedRange(x, y int) (int, int, error) {
	xMin := x
	xMax := x
	for xMin >= g.MinX && xMax <= g.MaxX {
		// Not enclosed if there is no floor.
		if g.Get(xMin, y+1) == '.' || g.Get(xMax, y+1) == '.' {
			return 0, 0, SearchError
		}
		// Check for walls.
		if g.Get(xMin, y) == '#' && g.Get(xMax, y) == '#' {
			// We must have a minimum separation to be able to have
			// an enclosed volume.
			if xMax-xMin < 2 {
				return 0, 0, SearchError
			}
			return xMin + 1, xMax - 1, nil
		}
		// Search on.
		if g.Get(xMin, y) != '#' {
			xMin--
		}
		if g.Get(xMax, y) != '#' {
			xMax++
		}
	}
	return 0, 0, SearchError
}

func (g Grid) OpenRange(x, y int) (int, int, error) {
	xMin := x
	xMax := x
	for xMin >= g.MinX && xMax <= g.MaxX {
		if g.Get(xMin, y+1) == '.' ||
			g.Get(xMax, y+1) == '.' ||
			g.Get(xMin, y+1) == '|' ||
			g.Get(xMax, y+1) == '|' {
			return xMin, xMax, nil
		}
		if g.Get(xMin-1, y) != '#' {
			xMin--
		}
		if g.Get(xMax+1, y) != '#' {
			xMax++
		}
	}
	return 0, 0, SearchError
}

func (g Grid) Step() {
	for y := g.MaxY; y >= 0; y-- {
		for x, tile := range g.Tiles[y] {
			switch tile {
			case '+':
				// Start flow.
				if g.Get(x, y+1) != '|' {
					g.Set(x, y+1, '|')
				}
			case '|':
				switch g.Get(x, y+1) {
				case '.':
					// Just flow down.
					g.Set(x, y+1, '|')
				case '#', '~':
					// We hit a floor or resting water. Try to expand.
					if x1, x2, err := g.EnclosedRange(x, y); err == nil {
						for xx := x1; xx <= x2; xx++ {
							g.Set(xx, y, '~')
						}
					} else {
						// Flow sideways until we have no more floor.
						if x1, x2, err := g.OpenRange(x, y); err == nil {
							for xx := x1; xx <= x2; xx++ {
								g.Set(xx, y, '|')
							}
						}
					}
				}
			}
		}
	}
}

func (g Grid) WaterCount(startY int) int {
	water := 0
	for y := startY; y <= g.MaxY; y++ {
		for _, tile := range g.Tiles[y] {
			switch tile {
			case '|', '~':
				water++
			}
		}
	}
	return water
}

func (g Grid) StillWaterCount(startY int) int {
	water := 0
	for y := startY; y <= g.MaxY; y++ {
		for _, tile := range g.Tiles[y] {
			switch tile {
			case '~':
				water++
			}
		}
	}
	return water
}

func (g Grid) FinalWaterCount(allWater bool) (int, error) {
	// Find y-coordinate with first #.
	for y := 0; y <= g.MaxY; y++ {
		for _, tile := range g.Tiles[y] {
			if tile == '#' {
				if allWater {
					return g.WaterCount(y), nil
				} else {
					return g.StillWaterCount(y), nil
				}
			}
		}
	}
	return 0, errors.New("could not find any clay")
}

func loadData(filename string) Grid {
	verticalPattern, _ := regexp.Compile(`x=(\d+),\s+y=(\d+)\.\.(\d+)`)
	horizontalPattern, _ := regexp.Compile(`y=(\d+),\s+x=(\d+)\.\.(\d+)`)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	grid := NewGrid()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		row := strings.TrimSpace(scanner.Text())

		match := verticalPattern.FindStringSubmatch(row)
		if len(match) == 0 {
			match = horizontalPattern.FindStringSubmatch(row)
			y, _ := strconv.Atoi(match[1])
			x1, _ := strconv.Atoi(match[2])
			x2, _ := strconv.Atoi(match[3])
			for x := x1; x <= x2; x++ {
				grid.Set(x, y, '#')
			}
		} else {
			x, _ := strconv.Atoi(match[1])
			y1, _ := strconv.Atoi(match[2])
			y2, _ := strconv.Atoi(match[3])
			for y := y1; y <= y2; y++ {
				grid.Set(x, y, '#')
			}
		}
	}

	grid.MinY = 1
	grid.MaxY = intsets.MinInt

	grid.MinX = intsets.MaxInt
	grid.MaxX = intsets.MinInt
	for y, xv := range grid.Tiles {
		if y > grid.MaxY {
			grid.MaxY = y
		}
		for x := range xv {
			if x < grid.MinX {
				grid.MinX = x
			}
			if x > grid.MaxX {
				grid.MaxX = x
			}
		}
	}

	grid.MinX--
	grid.MaxX++

	return grid
}

func main() {
	grid := loadData("input.txt")

	lastWaterCount := grid.WaterCount(1)
	for i := 0; ; i++ {
		grid.Step()
		//if i % 1000 == 0 {
		//	fmt.Println(strings.Repeat("-", 40))
		//	fmt.Println(i)
		//	grid.Print()
		//}
		waterCount := grid.WaterCount(1)
		if waterCount == lastWaterCount {
			fmt.Println(strings.Repeat("-", 40))
			fmt.Printf("Finished after %v steps.\n", i)
			break
		}
		lastWaterCount = waterCount
	}
	fmt.Println()
	grid.Print()
	if wc, err := grid.FinalWaterCount(true); err == nil {
		fmt.Println("Water Count =", wc)
	} else {
		log.Fatal(err)
	}

	if wc, err := grid.FinalWaterCount(false); err == nil {
		fmt.Println("Still Water Count =", wc)
	} else {
		log.Fatal(err)
	}
}
