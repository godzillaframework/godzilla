package godzilla_test

import (
	"fmt"
	"testing"
)

func maintest(t *testing.T) {
	fmt.Println("hello world")

	i := 1

	if i == 2 {
		fmt.Println("one")
	} else {
		fmt.Println("two")
	}

	msg := "MSG"

	fmt.Println(msg)
}
