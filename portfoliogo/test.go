package main

import (
	"fmt"
	// "strings"
	// "math"
	// "math/rand"
	// "errors"
	// "log"
)

var mySquare *[]int
var hi *myType

func test() {
	fmt.Println(mySquare)
}

func sayHi() {
	fmt.Println(hi)
}

type myType struct {
	hello string
}

func main() {
	fmt.Println(mySquare)
	mySquare = &[]int{1, 2, 3}
	fmt.Println(mySquare)
	test()
	hi = &myType{"hey!"}
	fmt.Println(hi)
}
