package main

import "testing"

func TestParse(t *testing.T) {
	values := loadData("test_input.txt")
	_, root := parseData(values, 0)
	if len(root.children) != 2 {
		t.Errorf("Root should have 2 children, found %v.", len(root.children))
	}
	if len(root.metadata) != 3 {
		t.Errorf("Root should have 3 metadata, found %v.", len(root.metadata))
	}

	if len(root.children[0].children) != 0 {
		t.Errorf("First child of root should have 0 children, found %v.", len(root.children[0].children))
	}
	if len(root.children[0].metadata) != 3 {
		t.Errorf("First child of root should have 3 metadata, found %v.", len(root.children[0].metadata))
	}

	if len(root.children[1].children) != 1 {
		t.Errorf("Second child of root should have 1 children, found %v.", len(root.children[1].children))
	}
	if len(root.children[1].metadata) != 1 {
		t.Errorf("Second child of root should have 1 metadata, found %v.", len(root.children[1].children))
	}

	if len(root.children[1].children[0].children) != 0 {
		t.Errorf("Child of second child of root should have 0 children, found %v.",
			len(root.children[1].children[0].children))
	}
	if len(root.children[1].children[0].metadata) != 1 {
		t.Errorf("Child of second child of root should have 1 metadata, found %v.",
			len(root.children[1].children[0].metadata))
	}

	if root.MetadataSum() != 138 {
		t.Errorf("Metadata sum %v is not equal to 138.", root.MetadataSum())
	}
}

func TestPart2(t *testing.T) {
	values := loadData("test_input.txt")
	_, root := parseData(values, 0)
	if root.Value() != 66 {
		t.Errorf("Root value should be 66, got %v.", root.Value())
	}
}
