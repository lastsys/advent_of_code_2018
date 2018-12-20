package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Register [6]int

type Machine struct {
	register   Register
	ip         int
	ipRegister int
	program    *Program
}

func NewMachine(program *Program) *Machine {
	return &Machine{Register{}, 0, program.ipRegister, program}
}

func (m *Machine) Step() bool {
	m.register[m.ipRegister] = m.ip
	if m.ip < 0 || m.ip >= len(m.program.instructions) {
		return false
	}
	i := m.program.instructions[m.ip]
	m.register = i.f(m.register, i.a, i.b, i.c)
	m.register[m.ipRegister]++
	m.ip = m.register[m.ipRegister]
	return true
}

func (m *Machine) Print() {
	fmt.Printf("ip =%3d [%8d, %8d, %8d, %8d, %8d, %8d]\n", m.ip,
		m.register[0], m.register[1], m.register[2], m.register[3], m.register[4], m.register[5])
}

type Instruction struct {
	f Function
	a int
	b int
	c int
}

type Program struct {
	ipRegister   int
	instructions []Instruction
}

type Function func(Register, int, int, int) Register

var functionMap = map[string]Function{
	"addr": addr,
	"addi": addi,
	"mulr": mulr,
	"muli": muli,
	"banr": banr,
	"bani": bani,
	"borr": borr,
	"bori": bori,
	"setr": setr,
	"seti": seti,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
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

func loadProgram(filename string) *Program {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	program := &Program{}
	program.instructions = make([]Instruction, 0)
	_, err = fmt.Sscanf(scanner.Text(), "#ip %d", &program.ipRegister)
	if err != nil {
		log.Fatal("Failed to read IP-register.")
	}
	var fName string
	for scanner.Scan() {
		instruction := Instruction{}
		_, err = fmt.Sscanf(scanner.Text(), "%s %d %d %d",
			&fName, &instruction.a, &instruction.b, &instruction.c)
		instruction.f = functionMap[fName]
		program.instructions = append(program.instructions, instruction)
	}
	return program
}

func part1() {
	program := loadProgram("input.txt")
	machine := NewMachine(program)
	machine.Print()
	for machine.Step() {
	}
	fmt.Println("Halted.")
	machine.Print()
	fmt.Println(strings.Repeat("-", 80))
}

func part2() {
	program := loadProgram("input.txt")
	machine := NewMachine(program)
	machine.register[0] = 1
	machine.Print()
	for i := 0; i < 200; i++ {
		machine.Step()
		machine.Print()
	}
	fmt.Println("Halted.")
	machine.Print()

	v := 10551348
	sum := 0
	for i := 1; i <= v; i++ {
		// Inner loop simplifed from assembler.
		//for j := 1; j <= v; j++ {
		//	if i*j == v {
		//		sum += i
		//	}
		//}
		if v%i == 0 {
			sum += i
		}
	}
	fmt.Println(sum)
}

func main() {
	part1()
	part2()
}
