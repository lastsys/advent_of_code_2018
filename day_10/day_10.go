package main

import (
	"bufio"
	"fmt"
	"golang.org/x/tools/container/intsets"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Point struct {
	x  int
	y  int
	vx int
	vy int
}

type PointList []*Point

func (p PointList) Step() {
	for _, point := range p {
		point.x += point.vx
		point.y += point.vy
	}
}

func (p PointList) Spread() float64 {
	ax := 0.0
	ay := 0.0

	// Average.
	for _, point := range p {
		ax += float64(point.x)
		ay += float64(point.y)
	}
	ax /= float64(len(p))
	ay /= float64(len(p))

	sx := 0.0
	sy := 0.0

	// Squared error.
	for _, point := range p {
		sx += math.Pow(float64(point.x)-ax, 2.0)
		sy += math.Pow(float64(point.y)-ay, 2.0)
	}

	return math.Sqrt(sx + sy)
}

func (p PointList) Limits() []int {
	minX := intsets.MaxInt
	maxX := intsets.MinInt
	minY := intsets.MaxInt
	maxY := intsets.MinInt

	for _, point := range p {
		if point.x < minX {
			minX = point.x
		} else if point.x > maxX {
			maxX = point.x
		}
		if point.y < minY {
			minY = point.y
		} else if point.y > maxY {
			maxY = point.y
		}
	}

	return []int{minX, maxX, minY, maxY}
}

type Grid struct {
	x1   int
	x2   int
	y1   int
	y2   int
	grid [][]bool
}

func NewGrid(x1, x2, y1, y2 int) *Grid {
	grid := make([][]bool, y2-y1)
	for i := 0; i < y2-y1; i++ {
		grid[i] = make([]bool, x2-x1)
	}
	return &Grid{x1, x2, y1, y2, grid}
}

func (g *Grid) Clear() {
	for y := 0; y < g.y2-g.y1; y++ {
		for x := 0; x < g.x2-g.x1; x++ {
			g.grid[y][x] = false
		}
	}
}

func (g *Grid) Render(pointList PointList) {
	for _, p := range pointList {
		if p.x >= g.x1 && p.x < g.x2 &&
			p.y >= g.y1 && p.y < g.y2 {
			g.grid[p.y-g.y1][p.x-g.x1] = true
		}
	}
}

func (g *Grid) Print() {
	for _, row := range g.grid {
		for _, isStar := range row {
			if isStar {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (g *Grid) HasVisiblePoint() bool {
	atLeastOnePointVisible := false
	for _, row := range g.grid {
		for _, isStar := range row {
			atLeastOnePointVisible = atLeastOnePointVisible || isStar
		}
	}
	return atLeastOnePointVisible
}

func loadData(filename string) PointList {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	points := make(PointList, 0)

	pattern, err := regexp.Compile(`position=<\s?([-\d]+),\s+([-\d]+)>\svelocity=<\s?([-\d]+),\s+([-\d]+)>`)
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		matches := pattern.FindStringSubmatch(scanner.Text())
		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])
		vx, _ := strconv.Atoi(matches[3])
		vy, _ := strconv.Atoi(matches[4])
		p := &Point{x, y, vx, vy}
		points = append(points, p)
	}
	return points
}

func part1(points PointList) {
	grid := NewGrid(180, 260, 180, 210)

	var spread float64
	i := 0
	for {
		points.Step()
		i++
		spread = points.Spread()
		if spread < 1000.0 {
			fmt.Println(spread)
		}
		if spread < 352.0 {
			fmt.Println(spread)
			fmt.Printf("Seconds: %v", i)
			break
		}
	}
	fmt.Println(points.Limits())
	grid.Clear()
	grid.Render(points)
	grid.Print()
}

func main() {
	points := loadData("input.txt")
	part1(points)
}
