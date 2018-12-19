package main

import (
	"fmt"
	"sort"
	"strings"
)

type Tile uint8

type Grid struct {
	width          int
	height         int
	tiles          [][]Tile
	units          UnitList
	elfAttackPower int
}

const (
	openTile Tile = 0
	wallTile Tile = 1
	unitTile Tile = 10 // Offset for UnitId.
)

func NewGrid(width, height int) *Grid {
	tiles := make([][]Tile, height)
	for i := 0; i < height; i++ {
		tiles[i] = make([]Tile, width)
	}
	units := make(UnitList, 0)
	return &Grid{width, height, tiles, units, 3}
}

func (g *Grid) GetUnitAt(x, y int) *Unit {
	if tile := g.tiles[y][x]; tile >= unitTile {
		return g.units[tile-unitTile]
	}
	return nil
}

func (g *Grid) GetUnit(id UnitId) *Unit {
	return g.units[id]
}

func (g *Grid) RemoveUnit(id UnitId) {
	p := g.units[id].position
	g.tiles[p[1]][p[0]] = openTile
	g.units[id] = nil
}

func (g *Grid) UnitReadOrder() UnitList {
	unitOrder := make(UnitList, 0, len(g.units))
	for _, unit := range g.units {
		if unit != nil {
			unitOrder = append(unitOrder, unit)
		}
	}
	sort.Sort(unitOrder)
	return unitOrder
}

func (g *Grid) MoveUnit(id UnitId, position Vector2) {
	unit := g.units[id]
	g.tiles[unit.position[1]][unit.position[0]] = openTile
	g.tiles[position[1]][position[0]] = unitTile + Tile(id)
	unit.position = position
}

func (g *Grid) Print() {
	for y, row := range g.tiles {
		unitsOnRow := make(UnitList, 0)
		for x, tile := range row {
			char := '.'
			if tile == wallTile {
				// A wall.
				char = '#'
			} else if unit := g.GetUnitAt(x, y); unit != nil {
				unitsOnRow = append(unitsOnRow, unit)
				// A unit.
				switch unit.race {
				case Goblin:
					char = 'G'
				case Elf:
					char = 'E'
				}
			}
			fmt.Printf("%c", char)
		}
		fmt.Print(strings.Repeat(" ", 3))
		for _, unit := range unitsOnRow {
			race := 'G'
			if unit.race == Elf {
				race = 'E'
			}
			fmt.Printf("%c(%3d), ", race, unit.hp)
		}
		fmt.Println()
	}
}
