package main

import "fmt"

type Race int

const (
	Elf Race = iota
	Goblin
)

const initialHitPoints = 200

type UnitId uint8

type Unit struct {
	id       UnitId
	race     Race
	position Vector2
	hp       int
}

func NewUnit(id UnitId, race Race, x, y int) *Unit {
	return &Unit{id, race, Vector2{x, y}, initialHitPoints}
}

func (u *Unit) Print() {
	fmt.Println(*u)
}
