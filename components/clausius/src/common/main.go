package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"strconv"
)

const defaultCellValue = 1

func GetGridSize() (int, int, error) {
	nbRows, err := strconv.Atoi(os.Getenv("NB_ROWS"))
	if err != nil {
		return 0, 0, err
	}
	nbCols, err := strconv.Atoi(os.Getenv("NB_COLS"))
	if err != nil {
		return 0, 0, err
	}
	return nbRows, nbCols, nil
}

func GetAllKeys(nbRows int, nbCols int) []string {
	allKeys := []string{}
	for y := 0; y < nbRows; y++ {
		for x := 0; x < nbCols; x++ {
			allKeys = append(allKeys, BuildKey(x, y))
		}
	}
	return allKeys
}

func BuildKey(x int, y int) string {
	return fmt.Sprintf("x:%v,y:%v", x, y)
}

type StoreModifier interface {
	Table() string
	Get(key string) (int, error)
	Set(key string, value int) error
	BatchGetItem(input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error)
}

type Store struct {
	StoreModifier
}

func (store Store) GetGrid(nbRows int, nbCols int) (map[string]int, error) {
	allKeys := GetAllKeys(nbRows, nbCols)
	keyChunks := SplitSSlice(allKeys, 100)
	result := make(map[string]int)
	for _, keyChunk := range keyChunks {
		m, err := store.GetValues(keyChunk)
		if err != nil {
			return nil, err
		}
		for k, v := range m {
			result[k] = v
		}
	}

	notFoundKeys := Except(allKeys, KeysOfMap(result))
	for _, k := range notFoundKeys {
		err := store.Set(k, defaultCellValue)
		if err != nil {
			return nil, err
		}
		result[k] = defaultCellValue
	}

	return result, nil
}

//TODO: test
func (store Store) GetValues(keys []string) (map[string]int, error) {
	attrValues := BuildAllAttrValues(keys)
	requestItems := BuildRequestItems(attrValues, store.Table())
	input := BuildBGIInput(requestItems)
	output, err := store.BatchGetItem(input)
	if err != nil {
		return nil, err
	}
	return FillMap(output, store.Table())
}

func (store Store) GetOrSetCell(x int, y int) (int, error) {
	got, getErr := store.GetCell(x, y)
	switch errV := getErr.(type) {
	case *ItemNotFound:
		setErr := store.SetCell(x, y, defaultCellValue)
		if setErr != nil {
			return 0, setErr
		}
		return defaultCellValue, nil
	default:
		return got, errV
	}
}

// TODO: use dynamodb atomic operation
func (store Store) RevertState(x int, y int) error {
	state, err := store.GetCell(x, y)
	if err != nil {
		return err
	}
	otherState := 1 - state
	return store.SetCell(x, y, otherState)
}

func (store Store) GetCell(x int, y int) (int, error) {
	key := BuildKey(x, y)
	return store.Get(key)
}

func (store Store) SetCell(x int, y int, v int) error {
	key := BuildKey(x, y)
	return store.Set(key, v)
}

type Cell struct {
	X int
	Y int
}

func (c Cell) IsValid(nbRows int, nbCols int) bool {
	validX := 0 <= c.X && c.X < nbCols
	validY := 0 <= c.Y && c.Y < nbRows
	return validX && validY
}
