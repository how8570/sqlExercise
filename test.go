package main

import (
	"fmt"
	"strings"
)

func main() {
	a := strings.Split("a s d", " ")
	for i := range a {
		fmt.Printf("%s\n", a[i])
	}
}
