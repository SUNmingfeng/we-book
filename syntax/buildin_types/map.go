package main

func Map() {
	m1 := map[string]string{
		"key1": "value1",
	}
	println(m1)
	m2 := make(map[string]string, 4) //考虑容量
	m2["key2"] = "value2"

	val1, ok := m1["key1"]
	println(val1, ok)
	val2, ok := m1["key2"]
	println(val2, ok)
	val3 := m1["key3"]
	println(val3)
	//map的键值对是无序的

	//删除键值对
	delete(m1, "key1") //无返回值，不知道map中是否有要删除的键
}
