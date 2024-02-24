package main

import (
	"fmt"
)

func main() {
	//UseSum()
	//vals := Insert[int](2, 33, []int{11, 22})
	//fmt.Printf("%v", vals)
	vals := Delete(2, []int{1, 2, 3, 4, 5})
	fmt.Printf("%v", vals)
}
