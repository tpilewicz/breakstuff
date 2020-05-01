package main

import (
	"fmt"
	"testing"
)

func TestRandomCellsToClick(t *testing.T) {
	nbRows := 10
	nbCols := 5
	nbCellsCDF := []float32{0.2, 0.5, 0.6, 0.8, 0.95, 1}
	cellsToClick := randomCellsToClick(nbRows, nbCols, nbCellsCDF)

	nbCells := len(cellsToClick)
	if nbCells < 1 || nbCells > 6 {
		t.Fatal(fmt.Errorf("We need between 1 and 10 cells to click. Got %v", nbCells))
	}
	for _, c := range cellsToClick {
		if !c.IsValid(nbRows, nbCols) {
			t.Fatal(fmt.Errorf("Woop, seems we generated a cell that isn't valid: %v. nbRows: %v, nbCols: %v", c, nbRows, nbCols))
		}
	}
}

func TestSample(t *testing.T) {
	s := sample([]float32{0.2, 0.5, 0.6, 0.8, 0.95, 1})
	if s < 0 || s > 5 {
		t.Fatal(fmt.Errorf("Sample should be between 0 and 5, got %v", s))
	}
}

func TestRandomCell(t *testing.T) {
	nbRows := 10
	nbCols := 5
	c := randomCell(nbRows, nbCols)

	if !c.IsValid(nbRows, nbCols) {
		t.Fatal(fmt.Errorf("%v is not a valid cell. nbRows: %v, nbCols: %v", c, nbRows, nbCols))
	}
}
