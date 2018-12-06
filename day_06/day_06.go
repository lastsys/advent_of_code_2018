package main

import (
	"bufio"
	"fmt"
	"golang.org/x/tools/container/intsets"
	"log"
	"os"
)

type Point struct {
	id int
	x  int
	y  int
}

type PointList []*Point

const occupied = -1

type Location struct {
	closestPointId int
	totalDistance  int
}

type Area struct {
	maxX     int
	maxY     int
	location [][]Location
}

func NewArea(points PointList) *Area {
	maxX := intsets.MinInt
	maxY := intsets.MinInt
	for _, p := range points {
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	maxX += 5
	maxY += 5
	locations := make([][]Location, maxX)
	for i := 0; i < maxX; i++ {
		locations[i] = make([]Location, maxY)
	}
	return &Area{maxX, maxY, locations}
}

type Distances []int

func (v Distances) Min() (min, count, id int) {
	min = intsets.MaxInt
	count = 0
	id = 0
	for i, v := range v {
		if min == v {
			count++
		} else if min > v {
			min = v
			count = 1
			id = i
		}
	}
	return min, count, id + 1
}

func (a *Area) Fill(points PointList) {
	distances := make(Distances, len(points))
	for x := 0; x < a.maxX; x++ {
		for y := 0; y < a.maxY; y++ {
			testPoint := Point{0, x, y}
			for i, p := range points {
				d := testPoint.ManhattanDistance(p)
				distances[i] = d
			}
			_, count, id := distances.Min()
			if count > 1 {
				a.location[x][y].closestPointId = occupied
			} else {
				a.location[x][y].closestPointId = id
			}
		}
	}
}

func (a *Area) LargestBoundedArea(points PointList) (int, int) {
	const unbounded = -1
	areas := make([]int, len(points))
	var id int
	for x := 0; x < a.maxX; x++ {
		for y := 0; y < a.maxY; y++ {
			id = a.location[x][y].closestPointId
			if id == -1 {
				continue
			}
			if x == 0 || y == 0 || x == a.maxX-1 || y == a.maxY-1 {
				areas[id-1] = unbounded
				continue
			}
			if areas[id-1] != unbounded {
				areas[id-1]++
			}
		}
	}
	maxArea := 0
	maxId := 0
	for i, p := range points {
		if areas[i] > maxArea {
			maxArea = areas[i]
			maxId = p.id
		}
	}
	return maxId, maxArea
}

func (a *Area) Print() {
	for y := 0; y < a.maxY; y++ {
		for x := 0; x < a.maxX; x++ {
			fmt.Printf("%4d", a.location[x][y].closestPointId)
		}
		fmt.Println()
	}
}

func (p1 *Point) ManhattanDistance(p2 *Point) int {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func loadData(filename string) PointList {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var x, y int
	count := 1
	points := make(PointList, 0)
	for scanner.Scan() {
		if _, err := fmt.Sscanf(scanner.Text(), "%d, %d\n", &x, &y); err != nil {
			log.Fatal(err)
		}
		points = append(points, &Point{count, x, y})
		count++
	}
	return points
}

func part1(points PointList) {
	area := NewArea(points)
	area.Fill(points)
	area.Print()
	id, a := area.LargestBoundedArea(points)
	fmt.Println(id, a)
}

func (a *Area) Fill2(points PointList) {
	for x := 0; x < a.maxX; x++ {
		for y := 0; y < a.maxY; y++ {
			d := 0
			for _, p := range points {
				d += p.ManhattanDistance(&Point{0, x, y})
			}
			a.location[x][y].totalDistance = d
		}
	}
}

func (a *Area) Print2() {
	for y := 0; y < a.maxY; y++ {
		for x := 0; x < a.maxX; x++ {
			fmt.Printf("%6d", a.location[x][y].totalDistance)
		}
		fmt.Println()
	}
}

func (a *Area) LocationsWithTotalDistanceLessThan(threshold int) int {
	count := 0
	for x := 0; x < a.maxX; x++ {
		for y := 0; y < a.maxY; y++ {
			if a.location[x][y].totalDistance < threshold {
				count++
			}
		}
	}
	return count
}

func part2(points PointList) {
	area := NewArea(points)
	area.Fill2(points)
	fmt.Println("Locations with total distance less than 10000:", area.LocationsWithTotalDistanceLessThan(10000))
}

func main() {
	points := loadData("input.txt")
	part1(points)
	part2(points)
}
