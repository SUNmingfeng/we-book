package main

func Loop1() {
	for i := 0; i < 10; i++ {
		println(i)
	}
	for i := 0; i < 10; {
		println(i)
		i++ //等价与上面
	}
}

func Loop2() {
	i := 0
	for i < 10 {
		println(i)
		i++ //等价于上面
	}
}

func Loop3() {
	for {
		println("xxx") //无限循环 ，cpu100%
	}
}

func ForSlice() {
	arr := []string{"a", "b", "c"}
	for idx, val := range arr {
		println(idx, val)
	}
}

func LoopBug() {
	users := []User{
		{
			name: "Tom",
		},
		{
			name: "Jerry",
		},
	}

	m := make(map[string]*User)
	for _, u := range users {
		m[u.name] = &u //这里的u是一个暂存位置，取&u实际上取得的就是这个暂存位的最新内容，遍历完成后，这个位置存的是最新读进来值
	}
}

type User struct {
	name string
}
