package main

import (
	"fmt"
	"golang.org/x/tools/container/intsets"
	"io/ioutil"
	"log"
	"strings"
)

type Units []rune

func loadFile(filename string) Units {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	trimmedLine := strings.TrimSpace(string(bytes))
	return []rune(trimmedLine)
}

func singleReduce(units Units) Units {
	// Rune difference between reactive units.
	const diff = 32
	if len(units) < 2 {
		return units
	}
	for i := 1; i < len(units); i++ {
		d := units[i] - units[i-1]
		if d == diff || d == -diff {
			var first, second Units
			if i > 1 {
				first = units[0 : i-1]
			} else {
				first = Units{}
			}
			if i < len(units)-1 {
				second = units[i+1:]
			} else {
				second = Units{}
			}
			return append(first, second...)
		}
	}
	return units
}

func part1(units Units) {
	units = fullReduce(units)
	fmt.Printf("%v : %v\n", len(units), string(units))

	unitCount := map[rune]bool{}
	for _, u := range units {
		var c rune
		if u >= 'a' {
			c = u - 32
		} else {
			c = u
		}
		unitCount[c] = true
	}
	fmt.Printf("We got %v units left.\n", len(unitCount))
}

func fullReduce(units Units) Units {
	previousLength := len(units)
	for {
		units = singleReduce(units)
		if len(units) == previousLength {
			break
		}
		previousLength = len(units)
	}
	return units
}

func part2(units Units) {
	s := string(units)
	stats := make([]int, 0)
	for c := 'A'; c <= 'Z'; c++ {
		s1 := strings.Replace(s, string(c), "", -1)
		s2 := strings.Replace(s1, string(c+32), "", -1)
		reduced := fullReduce(Units(s2))
		stats = append(stats, len(reduced))
	}
	minIndex := 0
	minValue := intsets.MaxInt
	for i, v := range stats {
		if v < minValue {
			minIndex = i
			minValue = v
		}
	}
	fmt.Println("Part2:")
	fmt.Printf("Min value = %v for %v\n", minValue, string(minIndex+'A'))
}

func main() {
	units := loadFile("input.txt")
	fmt.Println(len(units))
	part1(units)
	part2(units)
}
