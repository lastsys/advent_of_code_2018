package main

import (
	"bufio"
	"container/list"
	"fmt"
	"golang.org/x/tools/container/intsets"
	"log"
	"os"
)

type Regex []rune

type Map map[int]map[int]rune

func (m Map) Set(x, y int, tile rune) {
	if _, ok := m[y]; !ok {
		m[y] = map[int]rune{}
	}
	m[y][x] = tile
}

func (m Map) Get(x, y int) rune {
	if _, ok := m[y]; !ok {
		return ' '
	}
	if _, ok := m[y][x]; !ok {
		return ' '
	}
	return m[y][x]
}

func (m Map) MakeRoom(rx, ry int) {
	xm := rx * 2
	ym := ry * 2
	m.Set(xm, ym, '.')
	m.Set(xm-1, ym-1, '#')
	m.Set(xm-1, ym+1, '#')
	m.Set(xm+1, ym+1, '#')
	m.Set(xm+1, ym-1, '#')
	if m.Get(xm, ym-1) == ' ' {
		m.Set(xm, ym-1, '?')
	}
	if m.Get(xm, ym+1) == ' ' {
		m.Set(xm, ym+1, '?')
	}
	if m.Get(xm-1, ym) == ' ' {
		m.Set(xm-1, ym, '?')
	}
	if m.Get(xm+1, ym) == ' ' {
		m.Set(xm+1, ym, '?')
	}
}

func (m Map) AddDoor(rx, ry, dx, dy int) {
	xm := rx * 2
	ym := ry * 2
	if dx != 0 && dy == 0 {
		m.Set(xm+dx, ym, '|')
		return
	}
	if dx == 0 && dy != 0 {
		m.Set(xm, ym+dy, '-')
		return
	}
	log.Fatal("Illegal door position.")
}

func (m Map) Finalize() {
	for y, xv := range m {
		for x, r := range xv {
			if r == '?' {
				m.Set(x, y, '#')
			}
		}
	}
}

func (m Map) Print() {
	xMin := intsets.MaxInt
	xMax := intsets.MinInt
	yMin := intsets.MaxInt
	yMax := intsets.MinInt
	for y, vx := range m {
		if y < yMin {
			yMin = y
		}
		if y > yMax {
			yMax = y
		}
		for x := range vx {
			if x < xMin {
				xMin = x
			}
			if x > xMax {
				xMax = x
			}
		}
	}
	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			fmt.Printf("%c", m.Get(x, y))
		}
		fmt.Println()
	}
}

type Stack struct {
	sp    int
	stack []*State
}

func NewStack(stackSize int) *Stack {
	return &Stack{0, make([]*State, stackSize)}
}

func (s *Stack) Push(state *State) {
	s.stack[s.sp] = state
	s.sp++
}

func (s *Stack) Pop() *State {
	s.sp--
	state := s.stack[s.sp]
	s.stack[s.sp] = nil
	return state
}

type State struct {
	x, y int
}

func GenerateMap(regex Regex) Map {
	x := 0
	y := 0
	m := Map{}
	stack := NewStack(1024)
	for _, r := range regex {
		switch r {
		case '^':
			m.MakeRoom(x, y)
			m.Set(0, 0, 'X')
		case 'N':
			m.AddDoor(x, y, 0, -1)
			y--
			m.MakeRoom(x, y)
		case 'E':
			m.AddDoor(x, y, 1, 0)
			x++
			m.MakeRoom(x, y)
		case 'S':
			m.AddDoor(x, y, 0, 1)
			y++
			m.MakeRoom(x, y)
		case 'W':
			m.AddDoor(x, y, -1, 0)
			x--
			m.MakeRoom(x, y)
		case '(':
			stack.Push(&State{x, y})
		case ')':
			state := stack.Pop()
			x = state.x
			y = state.y
		case '|':
			state := stack.Pop()
			x = state.x
			y = state.y
			stack.Push(state)
		case '$':
			m.Finalize()
		}
	}
	return m
}

type Vector2 [2]int

func (v *Vector2) Add(v2 *Vector2) Vector2 {
	return Vector2{v[0] + v2[0], v[1] + v2[1]}
}

func FindShortestPathWithMostDoors(m Map, limit int) (int, int) {
	deltas := []Vector2{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	start := Vector2{0, 0}
	queue := list.New()
	queue.PushBack(start)
	cameFrom := map[Vector2]Vector2{}
	pathLengths := make([]int, 0)
	roomCount := 0
	for queue.Len() > 0 {
		e := queue.Front()
		queue.Remove(e)
		v := e.Value.(Vector2)

		culDeSac := true
		for _, delta := range deltas {
			w := v.Add(&delta)
			if tile := m.Get(v[0]*2+delta[0], v[1]*2+delta[1]); tile == '-' || tile == '|' {
				if _, ok := cameFrom[w]; !ok {
					culDeSac = false
					queue.PushBack(w)
					cameFrom[w] = v
				}
			}
		}
		// Backtrack.
		l := 0
		c := v
		for c != start {
			l++
			c = cameFrom[c]
		}
		if l >= limit {
			roomCount++
		}
		if culDeSac && l > 0 {
			pathLengths = append(pathLengths, l)
		}
	}

	fmt.Println(pathLengths)
	max := intsets.MinInt
	for _, v := range pathLengths {
		if v > max {
			max = v
		}
	}
	return max, roomCount
}

func loadRegex(filename string) []rune {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	regex, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return []rune(regex)
}

func main() {
	regex := loadRegex("input.txt")
	m := GenerateMap(regex)
	l, n := FindShortestPathWithMostDoors(m, 1000)
	m.Print()
	fmt.Println("Longest path =", l)
	fmt.Println("Room count =", n)
}
