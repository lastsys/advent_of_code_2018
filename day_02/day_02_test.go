package main

import "testing"

func TestNumberOfNonEqualCharacters(t *testing.T) {
	//serials := []string{
	//	"abcde", "fghij", "klmno", "pqrst", "fguij", "axcye", "wcxyz",
	//}

	if count, rest := numberOfNonEqualCharacters("abcde", "axcye"); count != 2 {
		t.Fatalf("Expected 2, got %v", count)
	} else if rest != "ace" {
		t.Fatalf("Expected ace, got %v", rest)
	}

	if count, rest := numberOfNonEqualCharacters("fghij", "fguij"); count != 1 {
		t.Fatalf("Expected 1, got %v", count)
	} else if rest != "fgij" {
		t.Fatalf("Expected fgij, got %v", rest)
	}
}
