package main

import "testing"

func testUnits(t *testing.T, s string, e string) {
	if r := string(singleReduce(Units(s))); r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
}

func TestSingleReduce(t *testing.T) {
	testUnits(t, "aA", "")
	testUnits(t, "abBA", "aA")
	testUnits(t, "abAB", "abAB")
	testUnits(t, "aabAAB", "aabAAB")
	testUnits(t, "dabAcCaCBAcCcaDA", "dabAaCBAcCcaDA")
	testUnits(t, "dabAaCBAcCcaDA", "dabCBAcCcaDA")
	testUnits(t, "dabCBAcCcaDA", "dabCBAcaDA")
	testUnits(t, "dabCBAcaDA", "dabCBAcaDA")
}
