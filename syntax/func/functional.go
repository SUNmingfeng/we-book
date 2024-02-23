package main

func Func4() {
	myfunc3 := Func3 //myfunc3本质是一个变量，只是这个变量是一个方法
	s, err := myfunc3(1, 5)
	println(s, err)
}

func Func5() {
	fn := func(name string) string {
		return "hello, " + name
	}
	str := fn("ssss")
	println(str)
}

// 返回值是方法
func Func6() func(name string) string {
	return func(name string) string {
		return "hello, " + name
	}
}

func Func6Invoke() {
	fn := Func6()
	str := fn("dddd")
	println(str)
}

// 匿名方法发起调用
func Func7() {
	fn := func(name string) string {
		return "hello, " + name
	}("aaaa")
	println(fn)

}
