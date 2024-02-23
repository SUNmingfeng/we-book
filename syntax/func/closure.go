package main

func Closure(name string) func() string {
	//闭包，方法+他绑定的运行上下文（用到的方法之外的参数
	//使用不当可能会引起内存泄露，一个对象被闭包引用时，是不会被垃圾回收的
	return func() string {
		return `hello ` + name
	}
}

func InvokeClosure() {
	c := Closure("qqqq")
	println(c())
}
