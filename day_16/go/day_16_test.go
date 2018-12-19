package main

import "testing"

func TestMatchingOpCodes(t *testing.T) {
	before := Register{3, 2, 1, 1}
	after := Register{3, 2, 2, 1}
	instruction := Instruction{9, 2, 1, 2}
	functions := NewFunctions()
	functionIndices := matchingFunctions(functions, before, after, instruction)
	if l := len(functionIndices); l != 3 {
		t.Errorf("Expected 3 matching OpCodes, got %v.", l)
	}
	if _, ok := functionIndices[2]; !ok {
		t.Errorf("Expected mulr to be present.")
	}
	if _, ok := functionIndices[1]; !ok {
		t.Errorf("Expected addi to be present.")
	}
	if _, ok := functionIndices[9]; !ok {
		t.Errorf("Expected seti to be present.")
	}
}
