package main

import (
	"container/list"
	"testing"
)

func listToIntArray(l *list.List) []int {
	a := make([]int, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		a = append(a, e.Value.(int))
	}
	return a
}

func intArrayToList(a []int) *list.List {
	l := list.New()
	for _, i := range a {
		l.PushBack(i)
	}
	return l
}

func elementOfIndex(l *list.List, idx int) *list.Element {
	e := l.Front()
	for i := 0; i < idx; i++ {
		e = e.Next()
	}
	return e
}

func TestInsertMarbleIntoCircle(t *testing.T) {
	circle := list.New()
	circle.PushBack(0)
	i0 := circle.Back()
	i1 := insertMarbleIntoCircle(circle, i0, 1)
	if i1 != elementOfIndex(circle, 1) {
		t.Errorf("Index of list %v is not equal to %v.", i1, elementOfIndex(circle, 1))
		t.Errorf("%v", listToIntArray(circle))
	}
	if circle.Len() != 2 {
		t.Errorf("Circle should contain 2 elements, but contains %v.", circle.Len())
	}
	a := listToIntArray(circle)
	if a[0] != 0 || a[1] != 1 {
		t.Errorf("Circle should be 0, 1 but is %v.", a)
	}

	i2 := insertMarbleIntoCircle(circle, i1, 2)
	if i2 != elementOfIndex(circle, 1) {
		t.Error("Index of list is not correct.")
	}
	if circle.Len() != 3 {
		t.Errorf("Circle should contain 3 elements, but contains %v.", circle.Len())
	}
	a = listToIntArray(circle)
	if a[0] != 0 || a[1] != 2 || a[2] != 1 {
		t.Errorf("Circle should be 0, 2, 1 but is %v.", a)
	}

	i3 := insertMarbleIntoCircle(circle, i2, 3)
	if i3 != elementOfIndex(circle, 3) {
		t.Errorf("Index after insert should be 3 but is %v.", i3.Value)
	}
	if circle.Len() != 4 {
		t.Errorf("Circle should contain 4 elements, but contains %v.", circle.Len())
	}
	a = listToIntArray(circle)
	if a[0] != 0 || a[1] != 2 || a[2] != 1 || a[3] != 3 {
		t.Errorf("Circle should be 0, 2, 1, 3 but is %v.", a)
	}
}

func TestHandleMultipleOf23(t *testing.T) {
	a1 := []int{0, 16, 8, 17, 4, 18, 9, 19, 2, 20, 10, 21, 5, 22, 11, 1, 12, 6, 13, 3, 14, 7, 15}
	circle := intArrayToList(a1)
	i1 := circle.Front()
	for i := 0; i < 13; i++ {
		i1 = i1.Next()
	}
	score, i2 := handleMultipleOf23(circle, i1, 23)
	if score != 23+9 {
		t.Errorf("Expected score 32 but got %v.", score)
	}
	element6 := elementOfIndex(circle, 6)
	if i2 != element6 {
		t.Errorf("Expected index 6 but got %v.", i2)
	}
	a2 := []int{0, 16, 8, 17, 4, 18, 19, 2, 20, 10, 21, 5, 22, 11, 1, 12, 6, 13, 3, 14, 7, 15}
	expectedCircle := intArrayToList(a2)
	if circle.Len() != expectedCircle.Len() {
		t.Errorf("Expected new circle length %v but got %v", expectedCircle.Len(), circle.Len())
	}
	for i := 0; i < expectedCircle.Len(); i++ {
		if a2[i] != elementOfIndex(expectedCircle, i).Value.(int) {
			t.Errorf("Circle value %v was %v but expected %v.", i, a2[i],
				elementOfIndex(expectedCircle, 6).Value)
		}
	}
}

func TestPlayGames(t *testing.T) {
	games := loadData("test_input.txt", true)
	var score int
	for _, game := range games {
		score = play(&game)
		if game.highScore != score {
			t.Errorf("Game %v expected highscore = %v, but got %v.", game, game.highScore, score)
		}
	}
}
