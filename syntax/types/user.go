package main

import "fmt"

func ChangeUser() {
	u1 := User{
		Name: "Tom",
		Age:  12,
	}
	fmt.Printf("%+v \n", u1)
	fmt.Printf("u1 address: %p \n", &u1)
	u1.ChangeName("Jerry")
	u1.ChanfgeAge(23)
	fmt.Printf("%+v \n", u1)

	//这里用不用指针不重要，重要的是接收器中用指针
	u2 := &User{
		Name: "Tom2",
		Age:  122,
	}
	fmt.Printf("%+v \n", u2)
	fmt.Printf("u2 address: %p \n", &u2)
	u2.ChangeName("Jerry2")
	u2.ChanfgeAge(23)
	fmt.Printf("%+v \n", u2)
}

type User struct {
	Name string
	Age  int
}

// 结构体接收器
// 下面这个等于 funcChangeName(u User, name string) {}
func (u User) ChangeName(name string) {
	fmt.Printf("ChangeName address: %p \n", &u)

	u.Name = name
}

// 指针接收器
func (u *User) ChanfgeAge(age int) {
	u.Age = age
}

// 结构体内部引用自己，只能用指针
type node1 struct {
	next *node1
}

func (n node1) Add(index int, val any) {
	//TODO implement me
	panic("implement me")
}

func (n node1) Append(val any) error {
	//TODO implement me
	panic("implement me")
}

func (n node1) Delete(index int) error {
	//TODO implement me
	panic("implement me")
}
