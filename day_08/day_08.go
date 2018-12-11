package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	children []*Node
	metadata []int
}

func (n *Node) MetadataSum() int {
	sum := 0
	for _, child := range n.children {
		sum += child.MetadataSum()
	}
	for _, data := range n.metadata {
		sum += data
	}
	return sum
}

func (n *Node) Value() int {
	sum := 0
	if len(n.children) == 0 {
		for _, data := range n.metadata {
			sum += data
		}
		return sum
	}
	for _, childIndex := range n.metadata {
		if childIndex > 0 && childIndex <= len(n.children) {
			sum += n.children[childIndex-1].Value()
		}
	}
	return sum
}

func loadData(filename string) []int {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	fInfo, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(f)
	buffer := make([]byte, fInfo.Size())
	_, err = r.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	ts := strings.TrimSpace(string(buffer))
	parts := strings.Split(ts, " ")

	values := make([]int, 0, len(parts))
	var value int
	for _, p := range parts {
		value, err = strconv.Atoi(p)
		if err != nil {
			log.Fatal(err)
		}
		values = append(values, value)
	}

	return values
}

func parseData(values []int, i int) (int, *Node) {
	childCount := values[i]
	i++
	metadataCount := values[i]
	i++
	children := make([]*Node, 0, childCount)
	var child *Node
	for j := 0; j < childCount; j++ {
		i, child = parseData(values, i)
		children = append(children, child)
	}
	metadata := make([]int, 0, metadataCount)
	for j := 0; j < metadataCount; j++ {
		metadata = append(metadata, values[i])
		i++
	}
	return i, &Node{children, metadata}
}

func part1(root *Node) {
	fmt.Println("Part 1:")
	fmt.Printf("Metadata sum = %v\n", root.MetadataSum())
}

func part2(root *Node) {
	fmt.Println("Part 2:")
	fmt.Printf("Root node value = %v\n", root.Value())
}

func main() {
	values := loadData("input.txt")
	_, root := parseData(values, 0)
	part1(root)
	part2(root)
}
