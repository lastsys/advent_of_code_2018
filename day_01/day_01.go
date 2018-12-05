package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func loadFile(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	numbers := make([]int, 0)

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return numbers
			}
			log.Fatal(err)
		}
		if number, err := strconv.Atoi(string(line)); err != nil {
			log.Fatal(err)
		} else {
			numbers = append(numbers, number)
		}
	}
}

func part1(values []int) {
	v := 0
	for _, value := range values {
		v += value
	}
	fmt.Println("Answer part 1: ", v)
}

func part2(values []int) {
	visited := map[int]bool{0: true}
	v := 0
	for {
		for _, value := range values {
			v += value
			if _, ok := visited[v]; ok {
				fmt.Println("Answer part 2: ", v)
				return
			} else {
				visited[v] = true
			}
		}
	}
}

func main() {
	values := loadFile("input.txt")
	part1(values)
	part2(values)
}
