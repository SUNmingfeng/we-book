package main

func Defer() {
	defer func() {}()
}

func DeferClosure() {
	i := 0
	defer func() {
		println(i)
	}()
	i = 1
}

func DeferClosureV1() {
	i := 111
	j := 333
	defer func(val int) {
		println("i:", val) //111
		println("j", j)    //444
	}(i)
	i = 222
	j = 444
}

func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}

// 如果是带名字的返回值，那么可以修改这个返回值
func DeferReturnV1() (a int) {
	a = 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturnV2() *MyStruct {
	res := &MyStruct{
		name: "Tom",
	}
	defer func() {
		//修改的不是res，是res指向的结构体
		//尽量不要在defer中修改返回值
		res.name = "Jerry"
	}()
	return res
}

type MyStruct struct {
	name string
}

// 10,10,10...
func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

func DeferClosureLoop2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

func DeferClosureLoop3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}
