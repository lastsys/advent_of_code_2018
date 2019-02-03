package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func loadFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	serials := make([]string, 0)

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return serials
			}
			log.Fatal(err)
		}
		trimmedLine := strings.TrimSpace(string(line))
		serials = append(serials, trimmedLine)
	}
}

func countPairsAndTriplets(s string) (int, int) {
	count := map[rune]int{}
	runes := []rune(s)
	for _, r := range runes {
		if v, ok := count[r]; ok {
			count[r] = v + 1
		} else {
			count[r] = 1
		}
	}

	pairs := 0
	triplets := 0
	for _, v := range count {
		switch v {
		case 2:
			pairs++
		case 3:
			triplets++
		}
	}
	return pairs, triplets
}

func part1(serials []string) {
	pairs := 0
	triplets := 0
	for _, serial := range serials {
		p, t := countPairsAndTriplets(serial)
		if p >= 1 {
			pairs++
		}
		if t >= 1 {
			triplets++
		}
	}
	checksum := pairs * triplets
	fmt.Printf("Part1: %v * %v = %v\n", pairs, triplets, checksum)
}

func numberOfNonEqualCharacters(s1 string, s2 string) (int, string) {
	r1 := []rune(s1)
	r2 := []rune(s2)
	count := 0
	rest := make([]rune, 0)
	for i := 0; i < len(r1); i++ {
		if r1[i] != r2[i] {
			count++
		} else {
			rest = append(rest, r1[i])
		}
	}
	return count, string(rest)
}

func part2(serials []string) {
	for i, serial1 := range serials {
		for _, serial2 := range serials[i:] {
			if serial1 == serial2 {
				continue
			}
			count, rest := numberOfNonEqualCharacters(serial1, serial2)
			if count == 1 {
				fmt.Println(count, rest)
			}
		}
	}
}

func main() {
	serials := loadFile("input.txt")
	part1(serials)
	part2(serials)
}
