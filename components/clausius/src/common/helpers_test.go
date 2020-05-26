package common

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSplitSSlice(t *testing.T) {
	s := []string{"1", "2", "3", "4", "5", "6", "7"}
	got := SplitSSlice(s, 2)
	want := [][]string{
		{"1", "2"},
		{"3", "4"},
		{"5", "6"},
		{"7"},
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Want %v\nGot %v", want, got))
	}

	s = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	got = SplitSSlice(s, 2)
	want = [][]string{
		{"1", "2"},
		{"3", "4"},
		{"5", "6"},
		{"7", "8"},
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Want %v\nGot %v", want, got))
	}
}

func TestSplitISlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	got := SplitISlice(s, 2)
	want := [][]int{
		{1, 2},
		{3, 4},
		{5, 6},
		{7},
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Want %v\nGot %v", want, got))
	}

	s = []int{1, 2, 3, 4, 5, 6, 7, 8}
	got = SplitISlice(s, 2)
	want = [][]int{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
	}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Want %v\nGot %v", want, got))
	}
}

func TestStringInSlice(t *testing.T) {
	if !stringInSlice("s", []string{"s", "a", "d"}) {
		t.Fatal(fmt.Errorf("This should be true"))
	}
	if stringInSlice("coucou", []string{"s", "a", "d"}) {
		t.Fatal(fmt.Errorf("This should be false"))
	}
}

func TestKeysOfMap(t *testing.T) {
	got := KeysOfMap(map[string]int{"k1": 1, "k2": 42})
	want := []string{"k1", "k2"}
	gotInWant := true
	wantInGot := true
	for _, elt := range got {
		if !stringInSlice(elt, want) {
			gotInWant = false
		}
	}
	for _, elt := range want {
		if !stringInSlice(elt, got) {
			wantInGot = false
		}
	}
	if !(gotInWant && wantInGot) {
		t.Fatal(fmt.Errorf("Want %v\nGot %v", want, got))
	}
}

func TestExcept(t *testing.T) {
	got := Except([]string{"1", "2", "42", "000"}, []string{"aa", "bb", "000"})
	want := []string{"1", "2", "42"}
	if !cmp.Equal(got, want) {
		t.Fatal(fmt.Errorf("Want %v\nGot %v", want, got))
	}
}
