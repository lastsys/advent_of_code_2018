package main

import (
	"container/list"
	"fmt"
	"golang.org/x/tools/container/intsets"
	"strings"
)

type Map struct {
	Depth   int
	Target  Position
	Erosion map[int]map[int]int
}

func NewMap(depth int, target Position) *Map {
	return &Map{
		depth,
		target,
		map[int]map[int]int{},
	}
}

func (m *Map) ErosionLevel(x, y int) int {
	if _, ok := m.Erosion[y]; !ok {
		m.Erosion[y] = map[int]int{}
	}
	if _, ok := m.Erosion[y][x]; !ok {
		geo := 0
		if x == m.Target.x && y == m.Target.y {
			geo = 0
		} else if y == 0 {
			geo = x * 16807
		} else if x == 0 {
			geo = y * 48271
		} else {
			geo = m.ErosionLevel(x-1, y) * m.ErosionLevel(x, y-1)
		}
		m.Erosion[y][x] = (geo + m.Depth) % 20183
	}
	return m.Erosion[y][x]
}

func (m *Map) Tile(x, y int) rune {
	var char rune
	switch m.ErosionLevel(x, y) % 3 {
	case 0:
		char = '.'
	case 1:
		char = '='
	case 2:
		char = '|'
	}
	return char
}

func (m *Map) Print() {
	var char rune
	maxY := intsets.MinInt
	for y := range m.Erosion {
		if y > maxY {
			maxY = y
		}
	}

	for y := 0; y <= maxY; y++ {
		maxX := intsets.MinInt
		for x := range m.Erosion[y] {
			if x > maxX {
				maxX = x
			}
		}
		for x := 0; x < maxX; x++ {
			if x == 0 && y == 0 {
				char = 'M'
			} else if x == m.Target.x && y == m.Target.y {
				char = 'T'
			} else {
				char = m.Tile(x, y)
			}
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
}

func (m *Map) RiskLevel() int {
	risk := 0
	for y := 0; y <= m.Target.y; y++ {
		for x := 0; x <= m.Target.x; x++ {
			risk += m.ErosionLevel(x, y) % 3
		}
	}
	return risk
}

type Tool int

const (
	neither Tool = iota
	climbingGear
	torch
)

func (t Tool) String() string {
	switch t {
	case neither:
		return "Neither"
	case climbingGear:
		return "Climbing Gear"
	case torch:
		return "Torch"
	}
	return ""
}

type Position struct {
	x    int
	y    int
	tool Tool
}

func (p *Position) Add(p2 *Position) Position {
	return Position{
		p.x + p2.x,
		p.y + p2.y,
		p2.tool,
	}
}

func (p *Position) Distance(p2 *Position) int {
	distance := 0
	// Move.
	if p.x != p2.x || p.y != p2.y {
		distance++
	}
	// Tool switch.
	if p.tool != p2.tool {
		distance += 7
	}
	return distance
}

func findPath(m *Map) {
	deltas := []Position{
		{-1, 0, neither},
		{-1, 0, climbingGear},
		{-1, 0, torch},

		{1, 0, neither},
		{1, 0, climbingGear},
		{1, 0, torch},

		{0, -1, neither},
		{0, -1, climbingGear},
		{0, -1, torch},

		{0, 1, neither},
		{0, 1, climbingGear},
		{0, 1, torch},
	}

	start := Position{0, 0, torch}
	queue := list.New()
	queue.PushBack(start)
	distanceFromSource := map[Position]int{start: 0}
	//cameFrom := map[Position]Position{}
	for queue.Len() > 0 {
		minDistance := intsets.MaxInt
		var minElement *list.Element
		for e := queue.Front(); e != nil; e = e.Next() {
			pos := e.Value.(Position)
			if d := distanceFromSource[pos]; d < minDistance {
				minDistance = d
				minElement = e
			}
		}
		queue.Remove(minElement)
		p := minElement.Value.(Position)
		if p == m.Target {
			break
		}

		for _, delta := range deltas {
			p2 := p.Add(&delta)
			if !canMove(m, &p, &p2) {
				continue
			}

			d2 := p.Distance(&p2)
			if ds, ok := distanceFromSource[p2]; !ok {
				distanceFromSource[p2] = distanceFromSource[p] + d2
				//cameFrom[p2] = p
				fmt.Printf("(%v, %v, %s) (%c) -> (%v, %v, %s) (%c) = (%v, %v)\n", p.x, p.y, p.tool.String(), m.Tile(p.x, p.y), p2.x, p2.y, p2.tool.String(), m.Tile(p2.x, p2.y), d2, distanceFromSource[p2])
				queue.PushBack(p2)
			} else if ds > distanceFromSource[p]+d2 {
				distanceFromSource[p2] = distanceFromSource[p] + d2
				//cameFrom[p2] = p
				fmt.Printf("(%v, %v, %s) (%c) -> (%v, %v, %s) (%c) = (%v, %v)\n", p.x, p.y, p.tool.String(), m.Tile(p.x, p.y), p2.x, p2.y, p2.tool.String(), m.Tile(p2.x, p2.y), d2, distanceFromSource[p2])
				queue.PushBack(p2)
			}
		}
		fmt.Println(strings.Repeat("-", 80))
	}

	fmt.Println(distanceFromSource[m.Target])
}

func canMove(m *Map, p1 *Position, p2 *Position) bool {
	// Rules:

	// Cannot move outside map range.
	if p2.x < 0 || p2.y < 0 {
		return false
	}
	if p2.x > m.Target.x*3 {
		return false
	}
	t1 := m.Tile(p1.x, p1.y)
	t2 := m.Tile(p2.x, p2.y)
	// Cannot use neither in rocky region.
	if t1 == '.' && p2.tool == neither {
		return false
	}
	if t2 == '.' && p2.tool == neither {
		return false
	}
	// Cannot use torch in wet region.
	if t1 == '=' && p2.tool == torch {
		return false
	}
	if t2 == '=' && p2.tool == torch {
		return false
	}
	// Cannot use climbing gear in narrow region.
	if t1 == '|' && p2.tool == climbingGear {
		return false
	}
	if t2 == '|' && p2.tool == climbingGear {
		return false
	}
	return true
}

func main() {
	//target := Position{10, 10, torch}
	//m := NewMap(510, target)
	target := Position{13, 743, torch}
	m := NewMap(8112, target)
	fmt.Println("Risk Level =", m.RiskLevel())

	findPath(m)
	//m.Print()
}
