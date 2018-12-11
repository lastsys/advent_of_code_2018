package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
)

type Node rune

type DirectedEdge struct {
	start Node
	end   Node
}

type NodeList []Node

func (n NodeList) Len() int {
	return len(n)
}

func (n NodeList) Less(i, j int) bool {
	return n[i] < n[j]
}

func (n NodeList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n NodeList) String() string {
	r := make([]rune, 0, len(n))
	for _, node := range n {
		r = append(r, rune(node))
	}
	return string(r)
}

type EdgeList []DirectedEdge

func (e EdgeList) Len() int {
	return len(e)
}

func (e EdgeList) Less(i, j int) bool {
	return e[i].end < e[j].end
}

func (e EdgeList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type Graph struct {
	nodes NodeList
	edges EdgeList
}

func loadData(filename string) *Graph {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var v1, v2 Node
	edges := make(EdgeList, 0)
	nodes := make(map[Node]bool)
	for scanner.Scan() {
		if _, err := fmt.Sscanf(scanner.Text(),
			"Step %c must be finished before step %c can begin.", &v1, &v2); err != nil {
			log.Fatal(err)
		}
		edges = append(edges, DirectedEdge{v1, v2})
		nodes[v1] = true
		nodes[v2] = true
	}
	nodeList := make(NodeList, 0)
	for k := range nodes {
		nodeList = append(nodeList, k)
	}
	return &Graph{nodeList, edges}
}

// Prioritized topological sort.
func topologicalSort(graph *Graph) (NodeList, error) {
	inDegree := make(map[Node]int, len(graph.nodes))
	for _, n := range graph.nodes {
		inDegree[n] = 0
	}

	// Store lists of edges nodes in a map.
	adjacency := make(map[Node]NodeList, len(graph.nodes))

	for _, e := range graph.edges {
		if _, ok := adjacency[e.start]; !ok {
			adjacency[e.start] = NodeList{e.end}
		} else {
			adjacency[e.start] = append(adjacency[e.start], e.end)
		}
	}

	// Initialize in-degree for each node.
	for _, v := range adjacency {
		for _, v2 := range v {
			inDegree[v2]++
		}
	}

	// Enqueue all vertices with in-degree of 0.
	queue := make(NodeList, 0)
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
		}
	}

	// Counter of visited vertices.
	count := 0

	// Store result.
	order := make(NodeList, 0)

	var u Node
	for {
		if len(queue) == 0 {
			break
		}
		sort.Sort(queue)
		u, queue = queue[0], queue[1:]
		order = append(order, u)
		for _, n := range adjacency[u] {
			inDegree[n]--
			if inDegree[n] == 0 {
				queue = append(queue, n)
			}
		}
		count++
	}

	if count != len(graph.nodes) {
		return nil, errors.New("cycle detected")
	}

	return order, nil
}

func part1(graph *Graph) {
	order, err := topologicalSort(graph)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:")
	fmt.Println(order)
}

type Task struct {
	node     Node
	progress int
}

type Workers []*Task

func (w *Workers) FirstIdleWorker() (int, error) {
	for i, task := range *w {
		if task == nil {
			return i, nil
		}
	}
	return -1, errors.New("no idle worker available")
}

func (w *Workers) Step() {
	for i, task := range *w {
		if task != nil && task.progress == 0 {
			(*w)[i] = nil
			continue
		}
		if task != nil && task.progress > 0 {
			task.progress--
		}
	}
}

func (w *Workers) AllIdle() bool {
	for _, task := range *w {
		if task != nil {
			return false
		}
	}
	return true
}

func (w *Workers) AllOccupied() bool {
	for _, task := range *w {
		if task == nil {
			return false
		}
	}
	return true
}

func (w *Workers) Print(t int, order NodeList) {
	fmt.Printf("%5v", t)
	for _, task := range *w {
		if task != nil {
			fmt.Printf(" %5c", task.node)
		} else {
			fmt.Printf(" %5c", '.')
		}
	}
	fmt.Printf(" \"%v\"\n", order.String())
}

// Prioritized topological sort with worker assignment.
func topologicalSort2(graph *Graph, workerCount int, timeBase int) (NodeList, error) {
	inDegree := make(map[Node]int, len(graph.nodes))
	for _, n := range graph.nodes {
		inDegree[n] = 0
	}

	// Store lists of edges nodes in a map.
	adjacency := make(map[Node]NodeList, len(graph.nodes))

	for _, e := range graph.edges {
		if _, ok := adjacency[e.start]; !ok {
			adjacency[e.start] = NodeList{e.end}
		} else {
			adjacency[e.start] = append(adjacency[e.start], e.end)
		}
	}

	// Initialize in-degree for each node.
	for _, v := range adjacency {
		for _, v2 := range v {
			inDegree[v2]++
		}
	}

	// Enqueue all vertices with in-degree of 0.
	queue := make(NodeList, 0)
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
		}
	}

	// Counter of visited vertices.
	count := 0

	// Store result.
	order := make(NodeList, 0)

	t := 0
	workers := make(Workers, workerCount)

	var u Node
	for {
		if len(queue) == 0 && workers.AllIdle() {
			workers.Print(t, order)
			break
		}

		// Assign as many workers as possible.
		for {
			if len(queue) == 0 || workers.AllOccupied() {
				break
			}
			worker, err := workers.FirstIdleWorker()
			if err == nil && len(queue) > 0 {
				sort.Sort(queue)
				u, queue = queue[0], queue[1:]
				workers[worker] = &Task{u, int(u) - 'A' + timeBase}
			}
		}

		workers.Print(t, order)

		for _, task := range workers {
			if task != nil && task.progress == 0 {
				// Update state when task is finished.
				order = append(order, task.node)
				for _, n := range adjacency[task.node] {
					inDegree[n]--
					if inDegree[n] == 0 {
						queue = append(queue, n)
					}
				}
				count++
			}
		}

		workers.Step()

		t++
	}

	if count != len(graph.nodes) {
		return nil, errors.New("cycle detected")
	}

	return order, nil
}

func part2(graph *Graph) {
	order, err := topologicalSort2(graph, 5, 60)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nPart 2:")
	fmt.Println(order.String())
}

func main() {
	graph := loadData("input.txt")
	part1(graph)
	part2(graph)
}
