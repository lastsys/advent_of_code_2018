package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Cut struct {
	id int
	x  int
	y  int
	w  int
	h  int
}

func (c1 *Cut) Overlaps(c2 *Cut) bool {
	// Check if not overlapping.
	if c1.x+c1.w < c2.x || c1.x > c2.x+c2.w {
		return false
	}
	if c1.y+c1.h < c2.y || c1.y > c2.y+c2.h {
		return false
	}
	return true
}

var cutPattern, _ = regexp.Compile(`#(\d+)\s@\s(\d+),(\d+):\s(\d+)x(\d+)`)

func loadFile(filename string) []*Cut {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cuts := make([]*Cut, 0)

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return cuts
			}
			log.Fatal(err)
		}
		trimmedLine := strings.TrimSpace(string(line))
		cuts = append(cuts, parseLine(trimmedLine))
	}
}

func parseLine(line string) *Cut {
	matches := cutPattern.FindStringSubmatch(line)
	id, _ := strconv.Atoi(matches[1])
	x, _ := strconv.Atoi(matches[2])
	y, _ := strconv.Atoi(matches[3])
	w, _ := strconv.Atoi(matches[4])
	h, _ := strconv.Atoi(matches[5])
	return &Cut{id, x, y, w, h}
}

func part1(cuts []*Cut) {
	const size = 1000
	fabric := make([]int, size*size)

	for i := 0; i < len(cuts); i++ {
		for x := cuts[i].x; x < (cuts[i].x + cuts[i].w); x++ {
			for y := cuts[i].y; y < (cuts[i].y + cuts[i].h); y++ {
				fabric[x+y*size]++
			}
		}
	}

	sum := 0
	for i := 0; i < len(fabric); i++ {
		if fabric[i] >= 2 {
			sum++
		}
	}

	fmt.Printf("Part 1: Total overlapping fabric = %v\n", sum)
}

func part2(cuts []*Cut) {
	for i := 0; i < len(cuts); i++ {
		overlaps := false
		for j := 0; j < len(cuts); j++ {
			if i == j {
				continue
			}
			if cuts[i].Overlaps(cuts[j]) {
				overlaps = true
				break
			}
		}
		if !overlaps {
			fmt.Printf("%v does not overlap with any other cut\n", cuts[i])
		}
	}
}

func main() {
	cuts := loadFile("input.txt")
	part1(cuts)
	part2(cuts)
}
