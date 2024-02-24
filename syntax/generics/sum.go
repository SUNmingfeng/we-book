package main

import (
	"fmt"
	"io"
)

func Sum[T Number](vals []T) T {
	var res T
	for _, v := range vals {
		res = res + v
	}
	return res
}

func Max[T Number](vals []T) T {
	max := vals[0]
	for i := 1; i < len(vals); i++ {
		if max < vals[i] {
			max = vals[i]
		}
	}
	return max
}

func Min[T Number](vals []T) T {
	min := vals[0]
	for i := 1; i < len(vals); i++ {
		if min > vals[i] {
			min = vals[i]
		}
	}
	return min
}

func Find[T Number](vals []T, filter func(t T) bool) T {
	for _, v := range vals {
		if filter(v) {
			return v
		}
	}
	var t T
	return t
}

func Insert[T Number](idx int, val T, vals []T) []T {
	fmt.Printf("切片原长度：%v \n", len(vals))
	if idx < 0 || idx > len(vals) {
		panic("idx不合法")
	}
	vals = append(vals, val)
	for i := len(vals) - 1; i > idx; i-- {
		if i-1 >= 0 {
			vals[i] = vals[i-1]
		}
	}
	vals[idx] = val
	return vals
}

// 实现删除切片特定下标元素的方法
func Delete[T Number](idx int, vals []T) []T {
	if idx < 0 || idx > len(vals)-1 {
		panic("idx不合法")
	}
	for i := idx; i < len(vals)-1; i++ {
		vals[i] = vals[i+1]
	}
	fmt.Printf("原始容量：%v，长度：%v\n", cap(vals), len(vals))
	vals = vals[:len(vals)-1]
	newVals := make([]T, 0, len(vals))
	newVals = append(newVals, vals...)
	fmt.Printf("缩容后容量：%v，长度：%v\n", cap(newVals), len(newVals))
	return newVals
}

// int的衍生类型
type Interger int

type Number interface {
	// ~int表示int的衍生类型
	~int | uint
}

func UseSum() {
	res := Sum[int]([]int{123, 456})
	println(res)
	resV1 := Sum[Interger]([]Interger{123, 456})
	println(resV1)
}

func Closable[T io.Closer]() {
	var t T
	t.Close()
}
