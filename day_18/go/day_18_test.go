package main

import (
	"fmt"
	"testing"
)

func TestSteps(t *testing.T) {
	grids := make([]*Grid, 11)

	for i := 0; i <= 10; i++ {
		grids[i] = loadData(fmt.Sprintf("%02d.txt", i))
	}

	g := grids[0]
	for i := 1; i <= 10; i++ {
		g.Step()
		if !g.Equal(grids[i]) {
			t.Errorf("Grids not equal at step %v.", i)
		}
	}

	if rv := g.ResourceValue(); rv != 1147 {
		t.Errorf("Resource value expected to be 1147, got %v.", rv)
	}
}
