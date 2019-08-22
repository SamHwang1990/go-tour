/*

Pointers:

	指针类型：保存内存地址，并能对内存地址中的值进行读写
		* 变量的内存地址可以用 int 类型的 16 进制整数表达
		* 内存地址从本质上只是一个整数
		* 要想通过内存地址来读写内存中存储的数据，就需要通过 Pointer 类型
		* 除了内存地址，Pointer 还需要直到内存存储的数据类型

	spec：
		```go
			PointerType = "*" BaseType .
			BaseType    = Type .
		```

	zero value：
		pointer 类型变量的 zero value 是 nil

	Pointer Operation：
		* 声明 Pointer 变量
			** Pointer Type 的声明语法：星号 + 值类型：
				- ` var pointer *int `
				- ` var pointer *string `
			** Point 的 zero value 是 nil
		* 初始化 Pointer 对象
			** 使用 "&" 操作符，操作数只能是变量，不能使用任何类型的字面量
				- Error:
					-- ` pointer := &1 `
					-- ` pointer := &"raw_string" `
				- Success:
					-- ` a := 1; pointer := &a`
					-- ` a := "raw_string"; pointer := &a `
			** 使用 new 函数创建指定类型的 Pointer 对象，
				Pointer 指向的内存数据为指定类型的 zero value，
				该过程又称 Allocation
				- spec: ` new(T) `，T 为 Type 对象
				- 用例：
					```go
						// allocation an nil slice in memory
						// and create a pointer which referencing that memory
						slicePointer := new([]int)

						// true
						*slicePointer == nil

						// allocation an int with zero value in memory
						// and create a pointer which referencing that memory
						intPointer := new(int)

						// true
						*intPointer == 0
					```
		* Dereferencing the Pointer（指针反引用）：
			指的是，根据 pointer 获取内存中的数据

			使用 "*" 操作符来操作 Pointer 变量

			** 读取 Pointer 中的数据：
				```go
					function pointerValueGetter(point *int) {

						// 直接使用 `*` 操作符即可完成读取操作
						fmt.Println(*point)
					}
				```

			** 更改 Pointer 中的数据，此更改操作会影响到所有引用到该内存地址的变量的值：
			```go
				function pointerValueSetter(point *int) {
					// 直接使用 `*` 操作符获取引用并进行赋值操作
					*point = (*pointer) * 2
				}

				a := 3
				pointerValueSetter(&a)

				// true
				a == 6
			```


	参考文章：
		- [pointers-in-go](https://medium.com/rungo/pointers-in-go-a789eafccd53)
		- [golang spec#Allocation](https://golang.org/ref/spec#Allocation)
*/

package main

import "fmt"

func main() {
	a := 1
	pa := &a

	fmt.Printf("pointer type: %T, pointer value: %v, memory value: %v\n", pa, pa, *pa)

	fmt.Println(*pa == a, a)

	*pa = 2
	fmt.Println(*pa == a, a)

	pSlice := new([]int)
	fmt.Printf("new slice pointer type: %T, new slice pointer value: %v, new slice is nil: %v\n", pSlice, pSlice, *pSlice == nil)

	*pSlice = append(*pSlice, 1, 2)
	fmt.Printf("appended slice pointer type: %T, appended slice pointer value: %v, appended slice value: %v\n", pSlice, pSlice, *pSlice)
}
