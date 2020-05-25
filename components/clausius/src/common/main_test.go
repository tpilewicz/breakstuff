package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/go-cmp/cmp"
	"os"
	"strconv"
	"testing"
)

func TestIsValid(t *testing.T) {
	nbRows := 5
	nbCols := 10

	c := Cell{X: 0, Y: 3}
	if !c.IsValid(nbRows, nbCols) {
		t.Fatal(fmt.Errorf("%v should be a valid cell.", c))
	}

	c = Cell{X: -1, Y: 3}
	if c.IsValid(nbRows, nbCols) {
		t.Fatal(fmt.Errorf("%v should NOT be a valid cell.", c))
	}

	c = Cell{X: 0, Y: -1}
	if c.IsValid(nbRows, nbCols) {
		t.Fatal(fmt.Errorf("%v should NOT be a valid cell.", c))
	}

	c = Cell{X: 10, Y: 3}
	if c.IsValid(nbRows, nbCols) {
		t.Fatal(fmt.Errorf("%v should NOT be a valid cell.", c))
	}

	c = Cell{X: 2, Y: 5}
	if c.IsValid(nbRows, nbCols) {
		t.Fatal(fmt.Errorf("%v should NOT be a valid cell.", c))
	}
}

func TestGetGridSize(t *testing.T) {
	os.Setenv("NB_ROWS", "10")
	os.Setenv("NB_COLS", "15")
	nbRows, nbCols, err := GetGridSize()

	if err != nil {
		t.Fatal(err)
	}
	wantRows := 10
	if nbRows != 10 {
		t.Fatal(fmt.Errorf("nbRows: got %v, want %v", nbRows, wantRows))
	}
	wantCols := 15
	if nbCols != 15 {
		t.Fatal(fmt.Errorf("nbCols: got %v, want %v", nbCols, wantCols))
	}
}

func TestGetGrid(t *testing.T) {
	mockModifier := &MockStoreModifier{}
	store := Store{StoreModifier: mockModifier, table: "my_tbl"}

	nbRows := 10
	nbCols := 15

	want := make(map[string]int)
	responses := []map[string]*dynamodb.AttributeValue{}
	unprocessedKeys := map[string]*dynamodb.KeysAndAttributes{
		"my_tbl": {
			Keys: []map[string]*dynamodb.AttributeValue{},
		},
	}
	for y := 0; y < nbRows; y++ {
		for x := 0; x < nbCols-1; x++ {
			v := (x + y) % 2
			want[BuildKey(x, y)] = v
			vStr := strconv.Itoa(v)
			responses = append(responses, map[string]*dynamodb.AttributeValue{
				"Key": &dynamodb.AttributeValue{
					S: aws.String(BuildKey(x, y)),
				},
				"V": &dynamodb.AttributeValue{
					N: aws.String(vStr),
				},
			})
		}
		x := nbCols - 1
		want[BuildKey(x, y)] = defaultCellValue
		unprocessedKeys["my_tbl"].Keys = append(
			unprocessedKeys["my_tbl"].Keys,
			map[string]*dynamodb.AttributeValue{
				"Key": {
					S: aws.String(BuildKey(x, y)),
				},
			},
		)
		mockModifier.On("Set", BuildKey(x, y), defaultCellValue).Return(nil).Once()
	}

	allKeys := GetAllKeys(nbRows, nbCols)
	attrValues := BuildAllAttrValues(allKeys)
	input := BuildBGIInput(attrValues, store.table)
	mockModifier.On("BatchGetItem", input).Return(
		&dynamodb.BatchGetItemOutput{
			Responses: map[string][]map[string]*dynamodb.AttributeValue{
				"my_tbl": responses,
			},
			UnprocessedKeys: unprocessedKeys,
		},
		nil,
	)

	got, err := store.GetGrid(nbRows, nbCols)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Got %#v\nWant %#v", got, want))
	}

}

func TestGetOrSetCell(t *testing.T) {
	mockModifier := &MockStoreModifier{}
	store := Store{StoreModifier: mockModifier, table: "my_tbl"}

	x := 1
	y := 15
	v := 1

	mockModifier.On("Get", BuildKey(x, y)).Return(v, nil).Once()
	got, err := store.GetOrSetCell(x, y)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, v))
	}

	expectedErr := fmt.Errorf("This is an expected error.")
	mockModifier.On("Get", BuildKey(x, y)).Return(0, expectedErr).Once()
	got, err = store.GetOrSetCell(x, y)
	if got != 0 {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, 0))
	}
	if err != expectedErr {
		t.Fatal(fmt.Errorf("Got error: %v, want error: %v", err, expectedErr))
	}

	mockModifier.On("Get", BuildKey(x, y)).Return(0, &ItemNotFound{"Nope"}).Once()
	mockModifier.On("Set", BuildKey(x, y), defaultCellValue).Return(nil).Once()
	got, err = store.GetOrSetCell(x, y)
	if got != defaultCellValue {
		t.Fatal(fmt.Errorf("Got %v, want %v (default cell value)", got, defaultCellValue))
	}
	if err != nil {
		t.Fatal(err)
	}

	mockModifier.On("Get", BuildKey(x, y)).Return(v, &ItemNotFound{"Nope"}).Once()
	mockModifier.On("Set", BuildKey(x, y), defaultCellValue).Return(expectedErr).Once()
	got, err = store.GetOrSetCell(x, y)
	if got != 0 {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, 0))
	}
	if err != expectedErr {
		t.Fatal(fmt.Errorf("Got error: %v, want error: %v", err, expectedErr))
	}

	mockModifier.AssertExpectations(t)
}

func TestGetCell(t *testing.T) {
	mockModifier := &MockStoreModifier{}
	store := Store{StoreModifier: mockModifier, table: "my_tbl"}

	x := 5
	y := 1
	v := 0
	mockModifier.On("Get", BuildKey(x, y)).Return(v, nil).Once()
	got, err := store.GetCell(x, y)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, v))
	}

	expectedErr := fmt.Errorf("This is an expected error.")
	mockModifier.On("Get", BuildKey(x, y)).Return(v, expectedErr).Once()
	got, err = store.GetCell(x, y)
	if err != expectedErr {
		t.Fatal(fmt.Errorf("Got error: %v, want error: %v", err, expectedErr))
	}
	if got != v {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, v))
	}

	mockModifier.AssertExpectations(t)
}

func TestSetCell(t *testing.T) {
	mockModifier := &MockStoreModifier{}
	store := Store{StoreModifier: mockModifier, table: "my_tbl"}

	x := 2
	y := 10
	v := 1
	mockModifier.On("Set", BuildKey(x, y), v).Return(nil).Once()
	err := store.SetCell(x, y, v)
	if err != nil {
		t.Fatal(err)
	}

	expectedErr := fmt.Errorf("This is an expected error.")
	mockModifier.On("Set", BuildKey(x, y), v).Return(expectedErr).Once()
	err = store.SetCell(x, y, v)
	if err != expectedErr {
		t.Fatal(fmt.Errorf("Got error: %v, want error: %v", err, expectedErr))
	}

	mockModifier.AssertExpectations(t)
}

func TestBuildKey(t *testing.T) {
	got := BuildKey(10, 20)
	want := "x:10,y:20"
	if got != want {
		t.Fatal(fmt.Errorf("got: %v, want: %v", got, want))
	}
}

func TestRevertState(t *testing.T) {
	mockModifier := &MockStoreModifier{}
	store := Store{StoreModifier: mockModifier, table: "my_tbl"}

	x := 2
	y := 10

	mockModifier.On("Get", BuildKey(x, y)).Return(0, nil).Once()
	mockModifier.On("Set", BuildKey(x, y), 1).Return(nil).Once()
	err := store.RevertState(x, y)
	if err != nil {
		t.Fatal(err)
	}

	mockModifier.On("Get", BuildKey(x, y)).Return(1, nil).Once()
	mockModifier.On("Set", BuildKey(x, y), 0).Return(nil).Once()
	err = store.RevertState(x, y)
	if err != nil {
		t.Fatal(err)
	}

	mockModifier.AssertExpectations(t)
}

func TestGetAllKeys(t *testing.T) {
	nbRows := 4
	nbCols := 4
	want := []string{}
	for y := 0; y < nbRows; y++ {
		for x := 0; x < nbCols; x++ {
			want = append(want, BuildKey(x, y))
		}
	}
	got := GetAllKeys(4, 4)
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Got %#v\nWant %#v", got, want))
	}
}

func TestBuildAllAttrValues(t *testing.T) {
	got := BuildAllAttrValues([]string{"key1", "key2", "key3"})
	want := []map[string]*dynamodb.AttributeValue{
		{
			"Key": &dynamodb.AttributeValue{
				S: aws.String("key1"),
			},
		},
		{
			"Key": &dynamodb.AttributeValue{
				S: aws.String("key2"),
			},
		},
		{
			"Key": &dynamodb.AttributeValue{
				S: aws.String("key3"),
			},
		},
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Got %#v\nWant %#v", got, want))
	}
}

func TestBuildAttrValue(t *testing.T) {
	got := BuildAttrValue("my_key")
	want := map[string]*dynamodb.AttributeValue{
		"Key": &dynamodb.AttributeValue{
			S: aws.String("my_key"),
		},
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Got %#v\nWant %#v", got, want))
	}
}

func TestBuildBGIInput(t *testing.T) {
	attrValues := BuildAllAttrValues([]string{"key1", "key2", "key3"})

	got := *BuildBGIInput(attrValues, "my_tbl")
	want := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			"my_tbl": {
				Keys:                 attrValues,
				ProjectionExpression: aws.String("V"),
			},
		},
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Got %#v\nWant %#v", got, want))
	}
}

func TestFillMap(t *testing.T) {
	got := make(map[string]int)
	output := dynamodb.BatchGetItemOutput{
		Responses: map[string][]map[string]*dynamodb.AttributeValue{
			"my_tbl": {
				{
					"Key": &dynamodb.AttributeValue{
						S: aws.String("key1"),
					},
					"V": &dynamodb.AttributeValue{
						N: aws.String("42"),
					},
				},
				{
					"Key": &dynamodb.AttributeValue{
						S: aws.String("key2"),
					},
					"V": &dynamodb.AttributeValue{
						N: aws.String("324"),
					},
				},
			},
		},
	}
	FillMap(got, &output, "my_tbl")

	want := map[string]int{
		"key1": 42,
		"key2": 324,
	}

	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Got %#v\nWant %#v", got, want))
	}
}
