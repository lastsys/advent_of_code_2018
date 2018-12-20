package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Register [4]int

type Instruction [4]int

type OpCode int

type TestCase struct {
	Before      Register
	Instruction Instruction
	After       Register
}

type Function func(Register, int, int, int) Register

func NewFunctions() []Function {
	return []Function{
		addr, addi,
		mulr, muli,
		banr, bani,
		borr, bori,
		setr, seti,
		gtir, gtri, gtrr,
		eqir, eqri, eqrr,
	}
}

func addr(r Register, a, b, c int) Register {
	r[c] = r[a] + r[b]
	return r
}

func addi(r Register, a, b, c int) Register {
	r[c] = r[a] + b
	return r
}

func mulr(r Register, a, b, c int) Register {
	r[c] = r[a] * r[b]
	return r
}

func muli(r Register, a, b, c int) Register {
	r[c] = r[a] * b
	return r
}

func banr(r Register, a, b, c int) Register {
	r[c] = r[a] & r[b]
	return r
}

func bani(r Register, a, b, c int) Register {
	r[c] = r[a] & b
	return r
}

func borr(r Register, a, b, c int) Register {
	r[c] = r[a] | r[b]
	return r
}

func bori(r Register, a, b, c int) Register {
	r[c] = r[a] | b
	return r
}

func setr(r Register, a, _, c int) Register {
	r[c] = r[a]
	return r
}

func seti(r Register, a, _, c int) Register {
	r[c] = a
	return r
}

func gtir(r Register, a, b, c int) Register {
	if a > r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

func gtri(r Register, a, b, c int) Register {
	if r[a] > b {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

func gtrr(r Register, a, b, c int) Register {
	if r[a] > r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

func eqir(r Register, a, b, c int) Register {
	if a == r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

func eqri(r Register, a, b, c int) Register {
	if r[a] == b {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

func eqrr(r Register, a, b, c int) Register {
	if r[a] == r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

func matchingFunctions(functions []Function, before Register, after Register, instruction Instruction) map[int]bool {
	matching := make(map[int]bool)
	for i, f := range functions {
		if after == f(before, instruction[1], instruction[2], instruction[3]) {
			matching[i] = true
		}
	}
	return matching
}

func loadData(filename string) ([]TestCase, []Instruction) {
	const (
		testCaseSection int = iota
		programSection
	)
	const (
		beforeState int = iota
		instructionState
		afterState
		separatorState
	)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	beforePattern, _ := regexp.Compile(`Before:\s+\[(\d+),\s*(\d+),\s*(\d+),\s*(\d+)\]`)
	instructionPattern, _ := regexp.Compile(`(\d+)\s(\d+)\s(\d+)\s(\d+)`)
	afterPattern, _ := regexp.Compile(`After:\s+\[(\d+),\s*(\d+),\s*(\d+),\s*(\d+)\]`)

	testCases := make([]TestCase, 0)
	program := make([]Instruction, 0)
	scanner := bufio.NewScanner(f)
	sectionState := testCaseSection
	readState := beforeState
	var testCase TestCase
	var line string
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if sectionState == testCaseSection {
			switch readState {
			case beforeState:
				// Happens on second separator line.
				if len(line) == 0 {
					// Read another line.
					scanner.Scan()
					sectionState = programSection
					continue
				}
				testCase = TestCase{}
				match := beforePattern.FindStringSubmatch(line)
				reg0, _ := strconv.Atoi(match[1])
				reg1, _ := strconv.Atoi(match[2])
				reg2, _ := strconv.Atoi(match[3])
				reg3, _ := strconv.Atoi(match[4])
				testCase.Before = Register{reg0, reg1, reg2, reg3}
				readState = instructionState
			case instructionState:
				match := instructionPattern.FindStringSubmatch(line)
				opCode, _ := strconv.Atoi(match[1])
				a, _ := strconv.Atoi(match[2])
				b, _ := strconv.Atoi(match[3])
				c, _ := strconv.Atoi(match[4])
				testCase.Instruction = Instruction{opCode, a, b, c}
				readState = afterState
			case afterState:
				match := afterPattern.FindStringSubmatch(line)
				reg0, _ := strconv.Atoi(match[1])
				reg1, _ := strconv.Atoi(match[2])
				reg2, _ := strconv.Atoi(match[3])
				reg3, _ := strconv.Atoi(match[4])
				testCase.After = Register{reg0, reg1, reg2, reg3}
				testCases = append(testCases, testCase)
				readState = separatorState
			case separatorState:
				readState = beforeState
			}
		} else {
			if len(line) == 0 {
				continue
			}
			match := instructionPattern.FindStringSubmatch(line)
			opCode, _ := strconv.Atoi(match[1])
			a, _ := strconv.Atoi(match[2])
			b, _ := strconv.Atoi(match[3])
			c, _ := strconv.Atoi(match[4])
			instruction := Instruction{opCode, a, b, c}
			program = append(program, instruction)
		}
	}

	return testCases, program
}

func part1(testCases []TestCase) {
	functions := NewFunctions()
	threeOrMoreMatches := 0
	for _, testCase := range testCases {
		if o := matchingFunctions(functions, testCase.Before, testCase.After, testCase.Instruction); len(o) >= 3 {
			threeOrMoreMatches++
		}
	}
	fmt.Printf("Out of %v cases, %v behave like three or more opcodes.\n",
		len(testCases), threeOrMoreMatches)
}

func part2(testCases []TestCase, program []Instruction) {
	functions := NewFunctions()
	candidates := Candidates{}
	for _, testCase := range testCases {
		matchingFunctions := matchingFunctions(functions, testCase.Before, testCase.After, testCase.Instruction)
		opCode := OpCode(testCase.Instruction[0])
		for functionIndex := range matchingFunctions {
			if _, ok := candidates[functionIndex]; !ok {
				candidates[functionIndex] = map[OpCode]bool{}
			}
			candidates[functionIndex][opCode] = true
		}
	}

	fmt.Println("Matches before (original function index vertical, opCode horizontal):")
	for functionIndex := 0; functionIndex < len(candidates); functionIndex++ {
		matches := candidates[functionIndex]
		fmt.Printf("%2d : ", functionIndex)
		v := [16]bool{}
		for m := range matches {
			v[m] = true
		}
		for i, m := range v {
			if m {
				fmt.Printf("%2d ", i)
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Println()
	}

	matchedFunctions := make([]Function, len(functions))
	for {
		functionIndex, opCode, err := candidates.findSingleMatch()
		if err != nil {
			break
		}
		fmt.Printf("%d -> %d\n", functionIndex, opCode)
		matchedFunctions[int(opCode)] = functions[functionIndex]
		delete(candidates, functionIndex)
	}
	fmt.Println(matchedFunctions)

	// Now just run the program.
	r := Register{}
	for _, instruction := range program {
		r = matchedFunctions[instruction[0]](r, instruction[1], instruction[2], instruction[3])
	}
	fmt.Println(r)
}

type Candidates map[int]map[OpCode]bool

func (c Candidates) matchCountForOpCode(opCode OpCode) int {
	count := 0
	for _, matches := range c {
		if _, ok := matches[opCode]; ok {
			count++
		}
	}
	return count
}

func (c Candidates) findSingleMatch() (int, OpCode, error) {
	count := map[OpCode]int{}
	for i := 0; i < 16; i++ {
		count[OpCode(i)] = c.matchCountForOpCode(OpCode(i))
	}
	for k, v := range count {
		if v == 1 {
			for fi, match := range c {
				if _, ok := match[k]; ok {
					return fi, k, nil
				}
			}
		}
	}
	return 0, 0, errors.New("could not find single match")
}

func main() {
	testCases, program := loadData("input.txt")
	part1(testCases)
	part2(testCases, program)
}
