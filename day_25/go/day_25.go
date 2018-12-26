package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

type Vector4 [4]int

func (v *Vector4) ManhattanDistance(v2 *Vector4) int {
	return Abs(v[0]-v2[0]) + Abs(v[1]-v2[1]) + Abs(v[2]-v2[2]) + Abs(v[3]-v2[3])
}

type Graph map[Vector4]map[Vector4]bool

func NewGraph(v []Vector4) Graph {
	g := Graph{}
	for _, v1 := range v {
		for _, v2 := range v {
			if v1.ManhattanDistance(&v2) <= 3 {
				g.add(v1, v2)
			}
		}
	}
	return g
}

func (g Graph) add(v1 Vector4, v2 Vector4) {
	if _, ok := g[v1]; !ok {
		g[v1] = make(map[Vector4]bool)
	}
	g[v1][v2] = true
	if _, ok := g[v2]; !ok {
		g[v2] = make(map[Vector4]bool)
	}
	g[v2][v1] = true
}

func (g Graph) countIslands() int {
	visited := map[Vector4]bool{}
	count := 0
	for v := range g {
		if visited[v] {
			continue
		}
		queue := list.New()
		queue.PushBack(v)
		for queue.Len() > 0 {
			e := queue.Front()
			queue.Remove(e)
			w1 := e.Value.(Vector4)
			visited[w1] = true
			for w2 := range g[w1] {
				if visited[w2] {
					continue
				}
				queue.PushBack(w2)
			}
			if queue.Len() == 0 {
				count++
			}
		}
	}
	return count
}

func loadData(filename string) []Vector4 {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	points := make([]Vector4, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p := Vector4{}
		i, err := fmt.Sscanf(scanner.Text(), "%d,%d,%d,%d\n",
			&p[0], &p[1], &p[2], &p[3])
		if err != nil {
			log.Fatal(i, err)
		}
		points = append(points, p)
	}

	return points
}

func main() {
	points := loadData("input.txt")
	g := NewGraph(points)
	fmt.Println(g.countIslands())
}
