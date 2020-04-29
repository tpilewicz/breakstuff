package common

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	"os"
	"strconv"
	"testing"
	"time"
)

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

type MockStoreModifier struct {
	mock.Mock
}

func (store *MockStoreModifier) Get(key string) (string, error) {
	rtns := store.Called(key)
	return rtns.String(0), rtns.Error(1)
}
func (store *MockStoreModifier) Set(key string, value interface{}, expiration time.Duration) error {
	rtns := store.Called(key, value, expiration)
	return rtns.Error(0)
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
			vStr := strconv.Itoa(v)
			want[BuildKey(x, y)] = v
			mockModifier.On("Get", BuildKey(x, y)).Return(vStr, nil).Once()
		}
		x := nb_cols - 1
		want[BuildKey(x, y)] = defaultCellValue
		mockModifier.On("Get", BuildKey(x, y)).Return("", redis.Nil).Once()
		mockModifier.On("Set", BuildKey(x, y), defaultCellValue, time.Duration(0)).Return(nil).Once()
	}

	got, err := store.GetGrid(nb_rows, nb_cols)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, want))
	}

}

func TestGetOrSetCell(t *testing.T) {
	mockModifier := &MockStoreModifier{}
	store := Store{mockModifier}

	x := 1
	y := 15
	v := 1
	vStr := "1"

	mockModifier.On("Get", BuildKey(x, y)).Return(vStr, nil).Once()
	got, err := store.GetOrSetCell(x, y)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, v))
	}

	expectedErr := fmt.Errorf("This is an expected error.")
	mockModifier.On("Get", BuildKey(x, y)).Return(vStr, expectedErr).Once()
	got, err = store.GetOrSetCell(x, y)
	if got != 0 {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, 0))
	}
	if err != expectedErr {
		t.Fatal(fmt.Errorf("Got error: %v, want error: %v", err, expectedErr))
	}

	mockModifier.On("Get", BuildKey(x, y)).Return(vStr, redis.Nil).Once()
	mockModifier.On("Set", BuildKey(x, y), defaultCellValue, time.Duration(0)).Return(nil).Once()
	got, err = store.GetOrSetCell(x, y)
	if got != defaultCellValue {
		t.Fatal(fmt.Errorf("Got %v, want %v (default cell value)", got, defaultCellValue))
	}
	if err != nil {
		t.Fatal(err)
	}

	mockModifier.On("Get", BuildKey(x, y)).Return(vStr, redis.Nil).Once()
	mockModifier.On("Set", BuildKey(x, y), defaultCellValue, time.Duration(0)).Return(expectedErr).Once()
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
	vStr := "0"
	mockModifier.On("Get", BuildKey(x, y)).Return(vStr, nil).Once()
	got, err := store.GetCell(x, y)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Fatal(fmt.Errorf("Got %v, want %v", got, v))
	}

	expectedErr := fmt.Errorf("This is an expected error.")
	mockModifier.On("Get", BuildKey(x, y)).Return(string(v), expectedErr).Once()
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
	mockModifier.On("Set", BuildKey(x, y), v, time.Duration(0)).Return(nil).Once()
	err := store.SetCell(x, y, v)
	if err != nil {
		t.Fatal(err)
	}

	expectedErr := fmt.Errorf("This is an expected error.")
	mockModifier.On("Set", BuildKey(x, y), v, time.Duration(0)).Return(expectedErr).Once()
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
