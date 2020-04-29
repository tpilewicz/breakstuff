package main

import (
	"clausius/common"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"math/rand"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.CloudWatchEvent) {
	nbRows, nbCols, err := common.GetGridSize()
	if err != nil {
		panic(err)
	}
	//TODO: create a Cell type, that has an x and a y, and rather generate a slice of this
	cellsToClick := randomCellsToClick(nbRows, nbCols)
}

// Chooses randomly between 1 and 10 cells to click, returns their keys. Note that a key may be repeated in the returned slice.
func randomCellsToClick(nbRows int, nbCols int) []string {
	nbCellsToClick := sample([]float32{1 / 32, 1 / 16, 1 / 8, 1 / 4, 1 / 2, 3 / 4, 7 / 8, 15 / 16, 31 / 32, 1}) + 1
	cellsToClick := []string{}
	for i := 0; i < nbCellsToClick; i++ {
		cellsToClick = append(cellsToClick, randomKey(nbRows, nbCols))
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

func randomKey(nbRows int, nbCols int) string {
	x := rand.Intn(nbCols)
	y := rand.Intn(nbRows)
	return common.BuildKey(x, y)
}
