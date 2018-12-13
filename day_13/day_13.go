package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"io"
	"log"
	"os"
	"sort"
)

type Collision struct {
	x int
	y int
}

type Cart struct {
	x                    int
	y                    int
	dx                   int
	dy                   int
	nextIntersectionTurn int
	crashed              bool
}

type CartList []*Cart

func (c CartList) Len() int {
	return len(c)
}

func (c CartList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CartList) Less(i, j int) bool {
	if c[i].y < c[j].y {
		return true
	} else if c[i].y > c[j].y {
		return false
	}
	return c[i].x < c[j].x
}

type Grid [][]byte

type State struct {
	carts       CartList
	grid        Grid
	activeCarts int
}

func (s *State) Print() {
	var char aurora.Value
	for y, row := range s.grid {
		fmt.Printf("%3v", y)
		for x, tile := range row {
			foundCart := false
			for _, cart := range s.carts {
				if cart.crashed {
					continue
				}
				if cart.x == x && cart.y == y {
					if cart.dx == 1 {
						char = aurora.Blue(">")
					} else if cart.dx == -1 {
						char = aurora.Blue("<")
					} else if cart.dy == 1 {
						char = aurora.Blue("v")
					} else if cart.dy == -1 {
						char = aurora.Blue("^")
					}
					foundCart = true
					// Mark if collision.
					for _, cart2 := range s.carts {
						if cart != cart2 && cart.x == cart2.x && cart.y == cart2.y {
							char = aurora.Red("X")
							break
						}
					}
					break
				}
			}
			if !foundCart {
				char = aurora.Gray(string(tile))
			}
			fmt.Printf("%v", char)
		}
		fmt.Println()
	}
}

func (s *State) Tick() {
	sort.Sort(s.carts)
	for _, cart := range s.carts {
		if cart.crashed {
			continue
		}
		// Move cart.
		cart.x += cart.dx
		cart.y += cart.dy
		// Find out if the cart needs to turn.
		tile := s.grid[cart.y][cart.x]
		switch tile {
		case '\\':
			cart.dx, cart.dy = cart.dy, cart.dx
		case '/':
			cart.dx, cart.dy = -cart.dy, -cart.dx
		case '+':
			switch cart.nextIntersectionTurn {
			case 0: // left
				if cart.dx == -1 {
					cart.dx, cart.dy = 0, 1
				} else if cart.dx == 1 {
					cart.dx, cart.dy = 0, -1
				} else if cart.dy == -1 {
					cart.dx, cart.dy = -1, 0
				} else if cart.dy == 1 {
					cart.dx, cart.dy = 1, 0
				}
			case 1: // straight
			case 2: // right
				if cart.dx == -1 {
					cart.dx, cart.dy = 0, -1
				} else if cart.dx == 1 {
					cart.dx, cart.dy = 0, 1
				} else if cart.dy == -1 {
					cart.dx, cart.dy = 1, 0
				} else if cart.dy == 1 {
					cart.dx, cart.dy = -1, 0
				}
			}
			cart.nextIntersectionTurn++
			cart.nextIntersectionTurn %= 3
		case '|', '-':
		default:
			s.Print()
			panic(errors.New("uncaught situation"))
		}

		// Check if we collided.
		for _, cart2 := range s.carts {
			if cart2.crashed {
				continue
			}
			if cart != cart2 &&
				cart.x == cart2.x &&
				cart.y == cart2.y {

				// Stop moving carts and return.
				cart.crashed = true
				cart2.crashed = true
				s.activeCarts -= 2
				fmt.Println(&Collision{cart.x, cart.y})

				if s.activeCarts == 1 {
					for _, c := range s.carts {
						if !c.crashed {
							fmt.Println("Last cart", c)
							return
						}
					}
				}
			}
		}
	}
}

func NewState(raw [][]byte) *State {
	carts := make(CartList, 0)

	for y, row := range raw {
		for x, tile := range row {
			switch tile {
			case '<':
				raw[y][x] = '-'
				carts = append(carts, &Cart{x, y, -1, 0, 0, false})
			case '>':
				raw[y][x] = '-'
				carts = append(carts, &Cart{x, y, 1, 0, 0, false})
			case '^':
				raw[y][x] = '|'
				carts = append(carts, &Cart{x, y, 0, -1, 0, false})
			case 'v':
				raw[y][x] = '|'
				carts = append(carts, &Cart{x, y, 0, 1, 0, false})
			}
		}
	}

	return &State{carts, raw, len(carts)}
}

func loadData(filename string) *State {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Using the scanner here gives very strange results.
	reader := bufio.NewReader(f)
	raw := make([][]byte, 0)
	var bytes []byte
	for {
		bytes, err = reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		raw = append(raw, bytes[0:len(bytes)-1])
	}

	return NewState(raw)
}

func part1(state *State) {
	i := 0
	fmt.Println(i)
	state.Print()
	for {
		state.Tick()
		i++
		if state.activeCarts == 1 {
			break
		}
	}
}

func main() {
	initialState := loadData("input.txt")
	part1(initialState)
}
