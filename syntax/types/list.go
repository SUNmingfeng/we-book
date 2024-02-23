package main

// 接口是一组行为，没有属性，只有方法
type List interface {
	Add(index int, val any)
	Append(val any) error
	Delete(index int) error
}

type LinkedList struct {
	head node
}

func (l *LinkedList) Add(index int, val any) {
	//
}

type node struct {
}
