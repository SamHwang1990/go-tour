/*

Array: 数组

	特性：
		* 数组长度固定，不可修改
		* 数组中元素类型须唯一
		* 数组对象属于值类型
			** 赋值时，会复制新的数组，而不是使用引用，修改新数组不会影响原数组
			** 将数组作为参数传值时，会复制新的数组，而不是使用引用，函数中修改实参数组不会影响原数组
			** 值包含的是整个数组，而不是数组第一额元素的内存地址，
				在 C 语言中，数组的值其实是数组中第一个元素的内存地址，这里是最大的差别

	zero value：
		数组变量的声明需要指定数组长度及元素类型，可以不对数组变量进行初始化
		此时，即使不对数组变量进行初始化，go runtime 也会在内存中根据数组长度及类型申请内存地址，并创建数组对象：
			* 数组长度为声明中指定的长度
			* 数组中元素均设为对应类型的 zero value，比如 int 类型会设置为 0

		举例：
			```go
				// arr1 use zero-value: [0, 0, 0]
				var arr1 [3]int

				// arr2 use zero-value: [false,false, false]
				var arr2 [3]bool
			```

	语法：
		* 声明 & 初始化数组
			参考：下面的 arrInitialize 函数

			** 只声明不初始化，使用 zero-value
				`var arr1 [3]int 	// // arr1 use zero-value: [0, 0, 0]`

			** 当变量初始值能推断出类型时，可省略变量的类型声明
				` var arr3 = [3]int{1, 2, 3} `

			** 支持使用 short-hand notation 来声明数组变量
				` arr4 := [3]int{1, 2, 3} `

			** 数组字面量支持多行声明，最后一行必须添加 `,`（逗号），
				否则 go compiler 会在最后一行添加分号，导致语法错误
				```go
					arr5 := [3]string{
							"foo",
							"bar",
							"zoo",
						}
				```

			** 数组字面量中，初始化给出的数组元素数量可以小于指定的数组长度，
				未初始化的数组元素会使用给类型的 zero-value：` arr6 == [1, 2, 3, 0, 0, 0] `
				` arr6 := [6]int{1, 2, 3}		// arr6 == [1, 2, 3, 0, 0, 0] `

			** 数组字面量中，可以只初始化数组元素，数组长度使用 `...` 来代替，
				此时，数组会使用初始化的元素数量来作为数组长度，
				`...` 不能省略，否则会变成 slices 的初始化
				` arr7 == [1, 2, 3, 4, 5] `


		* 读取数组元素
			所用 index 访问：` arr[index] `

		* 设置数组元素
			所用 index 更新元素：` arr[index] = newValue `

		* 遍历数组（ for loop ）
			```go
				arr := [...]int{1, 2, 3, 4, 5}

				for index, value := range arr {
					...
				}

				for i := 0; i < len(arr); i = i + 1 {
					...
				}
			```

	Array Comparison：
		* Array Type
			Array Type 由两部分组成：
				- Array Length，数组长度
				- Array Item Type，数组元素的类型
			Array Length 和 Array Item Type 一样的数组属于同一个 Array Type
		* 只有 Array Type 相同的数组才可进行比较
		* 判断两个数组是否相等：
			** 两个数组的 Array Type 一样
			** 数组中每个元素都相等

		举例：
			```go
				// true：元素类型以及数组长度均相同，可进行比较，且数组元素均相等
				[...]int{1, 2, 3} == [...]int{1, 2, 3}

				// Panic error：元素类型一致，但数组长度不一致，不能进行比较
				[...]int{1, 2, 3} == [...]int{1}

				// false：元素类型以及数组长度均相同，可进行比较，但数组元素不完全相等
				[...]int{1, 2, 3} == [...]int{3, 2, 1}
			```

	值类型，非引用类型：
		数组对象属于值类型
			** 赋值时，会复制新的数组，而不是使用引用，修改新数组不会影响原数组
			** 将数组作为参数传值时，会复制新的数组，而不是使用引用，函数中修改实参数组不会影响原数组
			** 值包含的是整个数组，而不是数组第一额元素的内存地址，
				在 C 语言中，数组的值其实是数组中第一个元素的内存地址，这里是最大的差别

		举例：
			```go
				arr1 := [...]int{1, 2, 3, 4, 5}

				// 在内存中重新申请空间，并将 arr1 的值全部复制到新的内存空间中
				arr2 := arr1

				// true，因为两个数组的 Array Type 以及数组元素均相等
				arr1 == arr2

				// false，两个数组使用了不同的内存地址，不是同一个引用指向
				&arr1 == &arr2
			```

	参考文章：
		* [the-anatomy-of-arrays-in-go](https://medium.com/rungo/the-anatomy-of-arrays-in-go-24429e4491b7)

*/

package main

import "fmt"

func arrInitialize() {
	// arr1 use zero-value: [0, 0, 0]
	var arr1 [3]int

	// arr2 use zero-value: [false,false, false]
	var arr2 [3]bool

	// 当变量初始值能推断出类型时，可省略变量的类型声明
	var arr3 = [3]int{1, 2, 3}

	// 支持使用 short-hand notation 来声明数组变量
	arr4 := [3]int{1, 2, 3}

	// 数组字面量支持多行声明，最后一行必须添加 `,`（逗号），否则 go compiler 会在最后一行添加分号，导致语法错误
	arr5 := [3]string{
		"foo",
		"bar",
		"zoo",
	}

	// 数组字面量中，初始化给出的数组元素数量可以小于指定的数组长度，
	// 未初始化的数组元素会使用给类型的 zero-value
	// arr6 == [1, 2, 3, 0, 0, 0]
	arr6 := [6]int{1, 2, 3}

	// 数组字面量中，可以只初始化数组元素，数组长度使用 `...` 来代替，
	// 此时，数组会使用初始化的元素数量来作为数组长度，
	// `...` 不能省略，否则会变成 slices 的初始化
	// arr7 == [1, 2, 3, 4, 5]
	arr7 := [...]int{1, 2, 3, 4, 5}

	fmt.Println("------- arrInitialize -------")
	fmt.Println(arr1)
	fmt.Println(arr2)
	fmt.Println(arr3)
	fmt.Println(arr4)
	fmt.Println(arr5)
	fmt.Println(arr6)
	fmt.Println(arr7)
	fmt.Println("------- arrInitialize -------")
}

func arrAsValue() {
	arr1 := [...]int{1, 2, 3, 4, 5}

	// 在内存中重新申请空间，并将 arr1 的值全部复制到新的内存空间中
	arr2 := arr1

	fmt.Println("------- arrAsValue -------")

	fmt.Println(arr1)
	fmt.Println(arr2)
	fmt.Println(arr1 == arr2)
	fmt.Println(&arr1 == &arr2)

	for index, value := range arr1 {
		arr1[index] = value * 3
	}

	for i := 0; i < len(arr2); i = i + 1 {
		arr2[i] = arr2[i] * 2
	}

	fmt.Println(arr1)
	fmt.Println(arr2)
	fmt.Println(arr1 == arr2)

	fmt.Println("------- arrAsValue -------")
}

func arrComparison() {
	fmt.Println("------- arrComparison -------")

	fmt.Println("[...]int{1, 2, 3} == [...]int{1, 2, 3}", [...]int{1, 2, 3} == [...]int{1, 2, 3})

	// Mismatch types，不同的 Array Type 不能进行比较
	// fmt.Println("[...]int{1, 2, 3} == [...]int{1}", [...]int{1, 2, 3} == [...]int{1})

	fmt.Println("[...]int{3, 2, 1} == [...]int{1, 2, 3}", [...]int{1, 2, 3} == [...]int{3, 2, 1})

	fmt.Println("------- arrComparison -------")
}

func main() {
	arrInitialize()
	arrAsValue()
	arrComparison()
}
