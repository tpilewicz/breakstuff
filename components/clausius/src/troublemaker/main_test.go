package main

import (
	"clausius/common"
	"fmt"
	"testing"
)

func TestRandomCellsToClick(t *testing.T) {
	nbRows := 10
	nbCols := 5
	cellsToClick := randomCellsToClick(nbRows, nbCols)

	nbCells := len(cellsToClick)
	if nbCells < 1 || nbCells > 10 {
		t.Fatal(fmt.Errorf("We need between 1 and 10 cells to click. Got %v", nbCells))
	}
	for _, k := range cellsToClick {
		if !isKey(k, nbRows, nbCols) {
			t.Fatal(fmt.Errorf("Woop, seems we generated a key that isn't one: %v. nbRows: %v, nbCols: %v", k, nbRows, nbCols))
		}
	}
}

func TestSample(t *testing.T) {
	s := sample([]float32{0.2, 0.5, 0.6, 0.8, 0.95, 1})
	if s < 0 || s > 5 {
		t.Fatal(fmt.Errorf("Sample should be between 0 and 5, got %v", s))
	}
}

func TestRandomKey(t *testing.T) {
	nbRows := 10
	nbCols := 5
	k := randomKey(nbRows, nbCols)

	if !isKey(k, nbRows, nbCols) {
		t.Fatal(fmt.Errorf("%v is not an acceptable key. nbCols: %v, nbRows: %v", k, nbCols, nbRows))
	}
}

func isKey(k string, nbRows int, nbCols int) bool {
	found_key := false
	for y := 0; y < nbRows; y++ {
		for x := 0; x < nbCols; x++ {
			if common.BuildKey(x, y) == k {
				found_key = true
			}
		}
	}
	return found_key
}
