package main

import "fmt"

func main() {
	u1 := &User{}  //取地址
	u1 = new(User) //取地址
	fmt.Println(u1)

	u2 := User{}
	fmt.Printf("%+v \n", u2)
	u2.Name = "jerry"

	var u3 User
	fmt.Printf("%+v \n", u3)
	var u4 *User //如果声明了一个指针，但没有赋值，那么他是nil
	fmt.Printf("%+v \n", u4)

	u5 := User{Name: "Jerry"}
	fmt.Printf("%+v", u5)

	u6 := User{"Jerry", 18} //难以区分 禁用
	fmt.Printf("%+v", u6)
	ChangeUser()

}

func UserList() {
	l1 := LinkedList{}
	l1Ptr := &l1    //对一个变量取地址，那么这个值是一个指针
	var l2 = *l1Ptr //对指针解引用
	fmt.Printf("%+v \n", l2)
}
