package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func loadData(filename string) *Grid {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read lines.
	reader := bufio.NewReader(f)
	lines := make([]string, 0)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		lines = append(lines, strings.TrimSpace(line))
	}

	// Read information.
	width := len(lines[0])
	height := len(lines)
	grid := NewGrid(width, height)
	unitId := UnitId(0)
	for y, row := range lines {
		for x, tile := range []rune(row) {
			switch tile {
			case 'G':
				grid.units = append(grid.units, NewUnit(unitId, Goblin, x, y))
				grid.tiles[y][x] = unitTile + Tile(unitId)
				unitId++
			case 'E':
				grid.units = append(grid.units, NewUnit(unitId, Elf, x, y))
				grid.tiles[y][x] = unitTile + Tile(unitId)
				unitId++
			case '.':
				grid.tiles[y][x] = openTile
			case '#':
				grid.tiles[y][x] = wallTile
			default:
				log.Fatalf("unknown tile type '%v'", string(tile))
			}
		}
	}

	return grid
}
