package main

import "fmt"

func Array() {
	a1 := [5]int{3, 4, 5}
	fmt.Printf("a1: %v \n len: %d \n cap: %d", a1, len(a1), cap(a1))
}

func Slice() {
	s1 := []int{3, 4, 5}
	fmt.Printf("s1: %v \n len: %d \n cap: %d \n", s1, len(s1), cap(s1))

	s2 := make([]int, 3, 4)
	fmt.Printf("s2: %v \n len: %d \n cap: %d \n", s2, len(s2), cap(s2))
	s2 = append(s2, 7) //len=4, cap=4
	s2 = append(s2, 8) //len=5, cap=9
	//推荐写法  s3 := make([]int, 0, capacity) len=0，并预估容量capacity，因为扩容性能很差
}

func SubSlice() {
	s1 := []int{2, 4, 6, 8}
	fmt.Printf("s1: %v， len:%d, cap:%d", s1, len(s1), cap(s1))
	s2 := s1[:1] //左闭右开，右边的不取
	fmt.Printf("s2: %v， len:%d, cap:%d", s2, len(s2), cap(s2))
}
