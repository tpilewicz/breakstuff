package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGetParams(t *testing.T) {
	params := make(map[string]string)
	params["x"] = "x"
	params["y"] = "y"
	params["v"] = "v"

	x, y, v, ok := get_params(params)
	if x != "x" || y != "y" || v != "v" {
		t.Fatal(
			fmt.Errorf(
				"We want for x, y, v: %v, %v, %v. We got: %v, %v, %v",
				"x", "y", "v", x, y, v,
			),
		)
	}
	if !ok {
		t.Fatal(fmt.Errorf("Result should be OK"))
	}

	delete(params, "v")
	x, y, v, ok = get_params(params)
	if ok {
		t.Fatal(fmt.Errorf("Result should not be OK"))
	}
}

func TestValidateParams(t *testing.T) {
	nb_rows := 10
	nb_cols := 15

	// Negative x
	expectedErr := fmt.Errorf("weird values for x and/or y. x: %v, y: %v", -10, 5)
	_, _, _, err := convertParams("-10", "5", "0", nb_rows, nb_cols)
	if err.Error() != expectedErr.Error() {
		t.Fatal(fmt.Errorf("Wanted error: %v. Got error: %v", expectedErr, err))
	}
	// Negative y
	expectedErr = fmt.Errorf("weird values for x and/or y. x: %v, y: %v", 2, -5)
	_, _, _, err = convertParams("2", "-5", "0", nb_rows, nb_cols)
	if err.Error() != expectedErr.Error() {
		t.Fatal(fmt.Errorf("Wanted error: %v. Got error: %v", expectedErr, err))
	}

	// Large x
	expectedErr = fmt.Errorf("weird values for x and/or y. x: %v, y: %v", 15, 5)
	_, _, _, err = convertParams("15", "5", "0", nb_rows, nb_cols)
	if err.Error() != expectedErr.Error() {
		t.Fatal(fmt.Errorf("Wanted error: %v. Got error: %v", expectedErr, err))
	}
	// Large y
	expectedErr = fmt.Errorf("weird values for x and/or y. x: %v, y: %v", 2, 10)
	_, _, _, err = convertParams("2", "10", "0", nb_rows, nb_cols)
	if err.Error() != expectedErr.Error() {
		t.Fatal(fmt.Errorf("Wanted error: %v. Got error: %v", expectedErr, err))
	}

	// Bad v
	expectedErr = fmt.Errorf("v must be 0 or 1. got: %v", 42)
	_, _, _, err = convertParams("1", "5", "42", nb_rows, nb_cols)
	if err.Error() != expectedErr.Error() {
		t.Fatal(fmt.Errorf("Wanted error: %v. Got error: %v", expectedErr, err))
	}

	// Unparseable x
	_, _, _, err = convertParams("aaa", "5", "1", nb_rows, nb_cols)
	switch errType := err.(type) {
	default:
		t.Fatal(fmt.Errorf("Error should be of type *strconv.NumError, got type %v", errType))
	case *strconv.NumError:
	}
	// Unparseable y
	_, _, _, err = convertParams("0", "2zzas", "1", nb_rows, nb_cols)
	switch errType := err.(type) {
	default:
		t.Fatal(fmt.Errorf("Error should be of type *strconv.NumError, got type %v", errType))
	case *strconv.NumError:
	}
	// Unparseable v
	_, _, _, err = convertParams("2", "5", "dzdz", nb_rows, nb_cols)
	switch errType := err.(type) {
	default:
		t.Fatal(fmt.Errorf("Error should be of type *strconv.NumError, got type %v", errType))
	case *strconv.NumError:
	}

	// ok
	x, y, v, err := convertParams("0", "5", "1", nb_rows, nb_cols)
	if x != 0 || y != 5 || v != 1 {
		t.Fatal(
			fmt.Errorf(
				"We want for x, y, v: %v, %v, %v. We got: %v, %v, %v",
				0, 5, 1, x, y, v),
		)
	}
}
