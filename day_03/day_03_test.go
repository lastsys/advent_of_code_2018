package main

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	if c := parseLine("#1 @ 829,837: 11x22"); c.id != 1 ||
		c.x != 829 ||
		c.y != 837 ||
		c.w != 11 ||
		c.h != 22 {
		t.Fatalf("Failed to parse: %v", c)
	}
	if c := parseLine("#583 @ 110,564: 10x23"); c.id != 583 ||
		c.x != 110 ||
		c.y != 564 ||
		c.w != 10 ||
		c.h != 23 {
		t.Fatalf("Failed to parse: %v", c)
	}
}
