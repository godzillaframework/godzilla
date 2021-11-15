package godzilla

import "fmt"

func ExampleGetString() {
	b := []byte("ABC€")
	str := GetString(b)
	fmt.Println(str)
	fmt.Println(len(b) == len(str))

	b = []byte("user")
	str = GetString(b)
	fmt.Println(str)
	fmt.Println(len(b) == len(str))

	b = nil
	str = GetString(b)
	fmt.Println(str)
	fmt.Println(len(b) == len(str))
	fmt.Println(len(str))
}
