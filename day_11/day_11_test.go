package main

import "testing"

func checkPower(t *testing.T, x, y, serial, expected int) {
	power := powerLevel(x, y, serial)
	if power != expected {
		t.Errorf("Power level expected to be %v, got %v.", expected, power)
	}
}

func verifyPower(t *testing.T, grid Grid, x, y, serial, expected int) {
	grid.Initialize(serial)
	power := grid[y-1][x-1]
	if power != expected {
		t.Errorf("Power level expected to be %v, got %v.", expected, power)
	}
}

func TestEnergyLevel(t *testing.T) {
	checkPower(t, 3, 5, 8, 4)
	checkPower(t, 122, 79, 57, -5)
	checkPower(t, 217, 196, 39, 0)
	checkPower(t, 101, 153, 71, 4)

	grid := NewGrid(300, 300)
	verifyPower(t, grid, 3, 5, 8, 4)
	verifyPower(t, grid, 122, 79, 57, -5)
	verifyPower(t, grid, 217, 196, 39, 0)
	verifyPower(t, grid, 101, 153, 71, 4)
}

func TestFindMaxSquare(t *testing.T) {
	grid := NewGrid(300, 300)
	maxSquareTester(t, grid, 18, 33, 45, 29, [][]int{
		{-2, -4, 4, 4, 4},
		{-4, 4, 4, 4, -5},
		{4, 3, 3, 4, -4},
		{1, 1, 2, 4, -3},
		{-1, 0, 2, -5, -2},
	})
	maxSquareTester(t, grid, 42, 21, 61, 30, [][]int{
		{-3, 4, 2, 2, 2},
		{-4, 4, 3, 3, 4},
		{-5, 3, 3, 4, -4},
		{4, 3, 3, 4, -3},
		{3, 3, 3, -5, -1},
	})
}

func maxSquareTester(t *testing.T, grid Grid, serial, ex, ey, ep int, square5x5 [][]int) {
	grid.Initialize(serial)
	x, y, totalPower := grid.FindMaxSquare(3, 3)
	for yy := y - 2; yy < y; yy++ {
		for xx := x - 2; xx < x; xx++ {
			if square5x5[yy-(y-2)][xx-(x-2)] != grid[yy][xx] {
				t.Errorf("(%v,%v) expected to be %v, but is %v.", xx+1, yy+1,
					square5x5[yy-(y-2)][xx-(x-2)], grid[yy][xx])
			}
		}
	}
	if x != ex {
		t.Errorf("x = %v expected to be %v.", x, ex)
	}
	if y != ey {
		t.Errorf("y = %v expected to be %v.", y, ey)
	}
	if totalPower != ep {
		t.Errorf("Total power is %v, expected %v.", totalPower, ep)
	}
}

func TestFindMaxSquare2(t *testing.T) {
	grid := NewGrid(300, 300)
	maxSquareTester2(t, grid, 18, 90, 269, 113, 16)
	maxSquareTester2(t, grid, 42, 232, 251, 119, 12)
}

func maxSquareTester2(t *testing.T, grid Grid, serial, ex, ey, ep, esz int) {
	grid.Initialize(serial)
	x, y, power, sz := grid.FindMaxSquare2(300)
	if x != ex {
		t.Errorf("x = %v expected to be %v", x, ex)
	}
	if y != ey {
		t.Errorf("y = %v expected to be %v", y, ey)
	}
	if power != ep {
		t.Errorf("Total power is %v, expected %v.", power, ep)
	}
	if sz != esz {
		t.Errorf("sz = %v expected to be %v", sz, esz)
	}
}
