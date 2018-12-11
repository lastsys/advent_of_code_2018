package main

import "testing"

func TestTopologicalSort(t *testing.T) {
	graph := loadData("test_input.txt")
	nodes, err := topologicalSort(graph)
	if err != nil {
		t.Error(err)
	}
	if nodes.String() != "CABDFE" {
		t.Errorf("Wrong sort order, expected CABDFE, got %v.", nodes.String())
	}
}

func TestTopologicalSort2(t *testing.T) {
	graph := loadData("test_input.txt")
	nodes, err := topologicalSort2(graph, 2, 0)
	if err != nil {
		t.Error(err)
	}
	if nodes.String() != "CABFDE" {
		t.Errorf("Wrong sort order, expected CABFDE, got %v.", nodes.String())
	}
}
