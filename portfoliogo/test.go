package main

import (
	"fmt"
)

type child struct {
	myNum int
}

func (c *child) add(n int) {
	c.myNum += 1
}

type parentif interface {
	add(n int)
}

type parent struct {
	*child
}

func main() {
	//c := &child{1}
	x := &parent{&child{1}}
	x.add(1)
	//fmt.Println(c.myNum)
	fmt.Println(x.myNum)
}
