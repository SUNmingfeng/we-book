package main

type ListV1[T any] interface {
	Add(index int, val T)
	Append(val T) error
	Delete(index int) error
}

type LinkListV1[T any] struct {
	head *nodeV1[T]
}

type nodeV1[T any] struct {
	data T
}
