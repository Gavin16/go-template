package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorProcess(t *testing.T) {
	item, err := findItem(3)
	fmt.Printf("%T\n", ErrNotFound)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			fmt.Println("Item not found")
		} else {
			fmt.Println("An error occurred:", err)
		}
		return
	}
	fmt.Println("Found item:", item)
}
