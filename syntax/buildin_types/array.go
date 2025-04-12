package main

import "fmt"

func Array() {
	a1 := [6]int{3, 4, 5}
	fmt.Printf("a1: %v \n len: %d \n cap: %d \n", a1, len(a1), cap(a1))
	var a2 [6]int //默认都赋值为0
	fmt.Printf("a2: %v \n len: %d \n cap: %d", a2, len(a2), cap(a2))
}

func Slice() {
	s1 := []int{3, 4, 5}
	fmt.Printf("s1: %v \n len: %d \n cap: %d \n", s1, len(s1), cap(s1))

	s2 := make([]int, 3, 4) //切片容量为4，初始化了3个元素
	fmt.Printf("s2: %v \n len: %d \n cap: %d \n", s2, len(s2), cap(s2))
	s2 = append(s2, 7) //len=4, cap=4
	s2 = append(s2, 8) //len=5, cap=8 原容量不足，扩容，容量*2
	fmt.Printf("s2: %v \n len: %d \n cap: %d \n", s2, len(s2), cap(s2))
	//推荐写法  s3 := make([]int, 0, capacity) len=0，并预估容量capacity，因为扩容性能很差
	s4 := make([]int, 4) //切片容量为4，初始化了4个元素
	fmt.Printf("s4: %v \n len: %d \n cap: %d \n", s4, len(s4), cap(s4))
}

func SubSlice() {
	s1 := []int{2, 4, 6, 8, 10}
	fmt.Printf("s1: %v， len:%d, cap:%d \n", s1, len(s1), cap(s1))
	s2 := s1[:2] //左闭右开，右边的不取
	fmt.Printf("s2: %v， len:%d, cap:%d \n", s2, len(s2), cap(s2))
	s3 := s1[2:] //左闭右开，右边的不取
	fmt.Printf("s3: %v， len:%d, cap:%d \n", s3, len(s3), cap(s3))
	s4 := s1[1:3] //左闭右开，右边的不取
	//子切片容量：从子切片的起始取到全切片的最后
	fmt.Printf("s4: %v， len:%d, cap:%d \n", s4, len(s4), cap(s4))
}

func ShareSlice() {
	s1 := []int{2, 4, 6, 8, 10}
	fmt.Printf("s1: %v， len:%d, cap:%d \n", s1, len(s1), cap(s1))
	s2 := s1[2:]
	//修改s2元素，s1的元素被同步修改
	s2[0] = 99
	fmt.Printf("s2: %v， len:%d, cap:%d \n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v， len:%d, cap:%d \n", s1, len(s1), cap(s1))
	//给s2扩容
	s2 = append(s2, 200)
	fmt.Printf("s2: %v， len:%d, cap:%d \n", s2, len(s2), cap(s2))
	//修改s2的元素，s1元素不再被同步修改
	s2[0] = 999
	fmt.Printf("s2: %v， len:%d, cap:%d \n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v， len:%d, cap:%d \n", s1, len(s1), cap(s1))
}
