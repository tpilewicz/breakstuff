package main

import (
	"clausius/common"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"math/rand"
	"time"
)

var nBCellsCDF = []float32{1.0 / 32.0, 1.0 / 16.0, 1.0 / 8.0, 1.0 / 4.0, 1.0 / 2.0, 3.0 / 4.0, 7.0 / 8.0, 15.0 / 16.0, 31.0 / 32.0, 1.0}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.CloudWatchEvent) {
	nbRows, nbCols, err := common.GetGridSize()
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	cellsToClick := randomCellsToClick(nbRows, nbCols)
	nbCellsToClick := len(cellsToClick)
	fmt.Printf("Clicking %v cells", nbCellsToClick)

	store, err := common.ConnectToFunes()
	if err != nil {
		panic(err)
	}

	errors := []error{}
	for _, c := range cellsToClick {
		err = store.RevertState(c.X, c.Y)
		if err != nil {
			errors = append(errors, err)
		}
	}
	nbErrors := len(errors)
	if nbErrors > 0 {
		for _, err = range errors {
			fmt.Println(err)
		}
		panic(fmt.Errorf("That's %v errors, out of %v clicked cells. Not cool.", nbErrors, nbCellsToClick))
	}
}

// Chooses randomly between 1 and len(nBCellsCDF) cells to click, returns their keys. Note that a key may be repeated in the returned slice.
func randomCellsToClick(nbRows int, nbCols int) []common.Cell {
	nbCellsToClick := sample(nBCellsCDF) + 1 // We want at least one cell
	cellsToClick := []common.Cell{}
	for i := 0; i < nbCellsToClick; i++ {
		cellsToClick = append(cellsToClick, randomCell(nbRows, nbCols))
	}
	return cellsToClick
}

// Returns a random int i, with i >= 0 and i < cdf.length, using cdf as the cumulative distribution function.
func sample(cdf []float32) int {
	r := rand.Float32()
	bucket := 0
	for r > cdf[bucket] {
		bucket++
	}
	return bucket
}

func randomCell(nbRows int, nbCols int) common.Cell {
	x := rand.Intn(nbCols)
	y := rand.Intn(nbRows)
	return common.Cell{X: x, Y: y}
}
