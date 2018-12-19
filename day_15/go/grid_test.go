package main

import (
	"testing"
)

func TestUnitReadOrder(t *testing.T) {
	grid := loadData("test_movement_0.txt")
	readOrder := grid.UnitReadOrder()
	for i := 0; i < len(readOrder); i++ {
		if unitId := readOrder[i].id; unitId != UnitId(i) {
			t.Errorf("Expected %v, got %v.", i, unitId)
		}
	}
}

func TestGetUnitAt(t *testing.T) {
	grid := loadData("test_movement_0.txt")
	if unit := grid.GetUnitAt(4, 4); unit != nil {
		if race := unit.race; race != Elf {
			t.Errorf("Expected Elf (%v), got %v.", Elf, race)
		}
	} else {
		t.Errorf("Expected unit, found none.")
	}
}
