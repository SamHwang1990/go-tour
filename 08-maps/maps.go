/*

Maps:

	用于保存 key-value 结构
		key、value 均可以为任意数据类型
		map 引用了一个内部数据结构来存储和管理 key、value，在 map 初始化时会初始化内部数据结构

	Map type 语法：` map[keyType]valueType `
		spec:
			```go
				MapType     = "map" "[" KeyType "]" ElementType .
				KeyType     = Type .
			```

	值类型，非引用类型：
		map 对象属于值类型
			** 赋值时，会复制新的 map，而不是使用引用
			** 将 map 作为参数传值时，会复制新的 map，而不是使用引用
			** 虽然 map 传递是值传递，
				但因为 map 复制时，不会复制内部的数据结构，而是对内部数据结构增加多了一份引用
				所以，牵一发动全身，行为上跟引用传递没啥区别

	zero value:
		map 的 zero value 为 nil，
		此时，不能对 nil 的 map 进行读写，因为内部数据结构还没初始化

	Map 声明与初始化语法：
		参考：下面的 mapInitialize 函数

		** map 字面量语法：
			```go
				map[keyType]valueType {
					key1: value1,
					key2: value2,
				}
			```
			key 值不能重复

		** 只声明不初始化，使用 zero-value
			`var mapZeroValue map[string]int 	// mapZeroValue use zero-value: nil`

		** 当变量初始值能推断出类型时，可省略变量的类型声明
		** 支持使用 short-hand notation 来声明数组变量
		** 数组字面量支持多行声明，最后一行必须添加 `,`（逗号），
			否则 go compiler 会在最后一行添加分号，导致语法错误

		** empty map 初始化，使用 make 函数
			` mapFromMake := make(map[keyType]ValueType) `

			make 函数签名：func make(map[keyType]ValueType) map[keyType]ValueType

	Map Operation:
		* len，获取 map 长度，使用内置函数：len(map)
		* getter:  ` value, isExisted := map[key] `
			getter 使用 key 来查找 value
			getter 返回两个值：
				- 第一个值为 map 中 key 对应的 value，valueType 类型
					若 map 中不存在给定的 key，则 value 会使用 valueType 的 zero value
				- 第二个值为 map 中是否存在给定的 key，boolean 类型
		* setter：` map[key] = value `
			- 若 map 中已存在对应的 key，则会使用 value 来覆盖原 value
			- 若 map 中不存在对应的 key，则会在 map 中创建新的 key，赋值为 value
		* delete，删除 map 中给定的 key，使用内置函数：delete(map, key)
		* looping，map 只支持 for-range looping
		* comparison
			map 只允许跟 nil 比较，甚至不能跟 map 自己比较：
				- success: ` map1 == nil `
				- error: ` map1 == map1 `
				- error: ` map1 == map2 `

	参考文章：
		- [golang spec#Map_types](https://golang.org/ref/spec#Map_types)

*/

package main

import "fmt"

func mapInitialize() {
	// use zero value: nil
	var mapZeroValue map[string]int

	mapEmptyWithMake := make(map[string]int)
	mapEmptyWithLiteral := map[string]int{}

	map1 := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	fmt.Println("------- mapInitialize -------")

	fmt.Println("map zero value is nil", mapZeroValue == nil)
	fmt.Println("empty map with make function", mapEmptyWithMake)
	fmt.Println("empty map with empty literal", mapEmptyWithLiteral)
	fmt.Println("map1 ", map1)

	fmt.Println("------- mapInitialize -------")
}

func mapGetterAndSetter() {
	map1 := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	fmt.Println("------- mapGetterAndSetter -------")

	fmt.Println(`map1["foo"]: `, map1["foo"])

	nonexistedZeroValue, isExisted := map1["nonexisted"]
	fmt.Printf(`nonexistedZeroValue, isExisted := map1["nonexisted"]; nonexistedZeroValue = %v, isExisted = %v
`, nonexistedZeroValue, isExisted)

	map1["foo"] = -1
	fmt.Println(`map1["foo"] = -1; map1["foo"]: `, map1["foo"])

	delete(map1, "foo")
	_, isExisted = map1["foo"]
	fmt.Println(`delete(map1, "foo") `, isExisted)

	fmt.Println("------- mapGetterAndSetter -------")
}

func mapLooping() {
	map1 := map[int]int{
		0: 1,
		1: 2,
		2: 3,
		3: 4,
		4: 5,
		5: 6,
		6: 7,
	}

	fmt.Println("------- mapLooping -------")

	fmt.Println(">> for range looping")
	for key, value := range map1 {
		fmt.Printf("key: %v, value: %v\n", key, value)
	}

	fmt.Println("------- mapLooping -------")
}

func main() {
	mapInitialize()
	mapGetterAndSetter()
	mapLooping()
}
