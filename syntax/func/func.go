package main

func main() {
	//Invoke()
	//Func5()
	//Func6Invoke()
	//InvokeClosure()
	//YourNameInvoke()
	//DeferClosureV1()
	//DeferReturn()
	//fmt.Println(DeferReturnV2().name)
	//DeferClosureLoopV1()
	//DeferClosureLoop2()
	DeferClosureLoop3()
}

func Func0(name string) string {
	return "hello, " + name
}

func Func2(a int, b int) (str string, err error) {
	str = "hello"
	return "", err //可以给返回值命名但不使用
}

func Func3(a int, b int) (str string, err error) {
	return "abc", nil
}

func Invoke() {
	str := Func0("ssss")
	println(str)
	str1, err := Func2(1, 2)
	println(str1, err)
	_, err = Func3(1, 3)
	_, _ = Func3(4, 5) //可以全部返回值不接受
}

func Recursive() {
	Recursive()
}
