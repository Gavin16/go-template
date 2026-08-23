package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid")
)

func findItem(id int) (string, error) {
	if id == 1 {
		return "Item 1", nil
	} else if id == 2 {
		return "Item 2", nil
	} else {
		return "", fmt.Errorf("findItem: %w", ErrNotFound)
	}
}
