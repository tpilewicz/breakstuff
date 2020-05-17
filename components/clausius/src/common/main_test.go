package common

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"os"
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
	store := Store{mockModifier}

	nb_rows := 10
	nb_cols := 15

	want := make(map[string]int)
	for y := 0; y < nb_rows; y++ {
		for x := 0; x < nb_cols-1; x++ {
			v := (x + y) % 2
			want[BuildKey(x, y)] = v
			mockModifier.On("Get", BuildKey(x, y)).Return(v, nil).Once()
		}
		x := nb_cols - 1
		want[BuildKey(x, y)] = defaultCellValue
		mockModifier.On("Get", BuildKey(x, y)).Return(0, &ClausiusTestError{"Nope"}).Once()
		mockModifier.On("Set", BuildKey(x, y), defaultCellValue).Return(nil).Once()
	}

	got, err := store.GetGrid(nb_rows, nb_cols)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Got %#v, want %#v", got, want))
	}

}

func TestGetOrSetCell(t *testing.T) {
	mockModifier := &MockStoreModifier{}
	store := Store{mockModifier}

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

	mockModifier.On("Get", BuildKey(x, y)).Return(0, &ClausiusTestError{"Nope"}).Once()
	mockModifier.On("Set", BuildKey(x, y), defaultCellValue).Return(nil).Once()
	got, err = store.GetOrSetCell(x, y)
	if got != defaultCellValue {
		t.Fatal(fmt.Errorf("Got %v, want %v (default cell value)", got, defaultCellValue))
	}
	if err != nil {
		t.Fatal(err)
	}

	mockModifier.On("Get", BuildKey(x, y)).Return(v, &ClausiusTestError{"Nope"}).Once()
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
	store := Store{mockModifier}

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
	store := Store{mockModifier}

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
	store := Store{mockModifier}

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
