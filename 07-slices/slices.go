/*

Slices:
	> 中文翻译是：分片，建议交流时使用英文

	slice 个人理解：
		golang 中，slice 是为了满足可变数组需求，而从数组演变出来的结构体，注意，slice 是一个结构体而已（struct），
		数组不可变，那 slice 可变的原理，简单理解：
			* slice 底层使用的仍然是数组来存储数据，slice 中存在一个指向数组元素的引用
			* 当对 slice 进行修改时，比如增加元素，
				若 slice 底层的数组仍有空间放置新元素，则会所谓的 slice appending，其实只是在数组的适当位置存储新的元素值，
				若 slice 底层的数组不够空间放置新元素，则会创建一个新的，足够大的数组来放置slice 原来的元素以及新的元素

	Slice is a struct：
		> A slice is a descriptor of an array segment.
		* Slice 内部是一个 Struct 结构，描述了一个数组片段：
			- ptr（ *Elem ），指向了数组中的某一个元素的内存地址，
				该元素是 slice 片段的第一个元素，
				或者说，slice 片段从该元素开始
				该指针只读
			- length（ int ），slice 片段的长度，只读
			- cap（int），slice 的容量，该值等于 slice 片段第一个元素到数组最后一个元素的数量，只读

		* Slice 信息本身其实是只读，不可改的
			- slice 指向的数组不会变化
			- slice ptr 指向的数组元素位置不会变化
			- slice 片段的长度不会变化
			- slice 片段的容量不会变化

		* Slice 是如何满足可变数组长度的需求的：
			- 数组本身不可变
			- Slice 本身也不可变
			- 重点：当 slice 长度发生变化时，肯定都是通过生成新的 slice 甚至生成新的数组了
				slice 长度变化场景中，拿最基本的两个场景举例：
				-- slice 长度缩减：` sliceShorter := sliceLonger[0:len(sliceLonger) - 1] `
					代码通过 slicing expression 将 slice 长度减一，结果是：
						--- sliceLonger 本身其实没有变化
						--- 返回了新的 slice 并赋值给 sliceShorter
				-- slice 长度增加：` sliceLonger := append(sliceShorter, element1) `
					代码通过 append 一个新元素的方法，将 slice 长度加一，结果是：
						--- sliceShorter 本身其实没有变化
						--- 返回了新的 slice 并赋值给 sliceLonger

		* Slice 元素值的调整会反映到内部的数组中

		* 在 Golang 中，我们其实不会经常使用数组本身，而是直接使用 Slice 来给数据存储带来足够的灵活度，
			但鉴于 slice 其实是数组的引用，所以需要非常小心 slice 的内存占用问题

	值类型，非引用类型：
		slice 对象属于值类型
			** 赋值时，会复制新的 slice，而不是使用引用
			** 将 slice作为参数传值时，会复制新的 slice，而不是使用引用
			** 虽然 slice 传递是值传递，但因为 slice 复制后引用的数组还是同一个数组，
				所以，牵一发动全身，行为上跟引用传递没啥区别

	zero value：
		slice 的 zero value 是 `nil`，即任何只声明但没有初始化的 slice 均等于 `nil`

	Slice 声明与初始化语法：
		参考：下面的 sliceInitialize 函数

		** slice 字面量的语法与 array 字面量语法类似，差别在于不需要声明数量，golang 会自动生成对应的数组
			` slice1 := []int{1,2,3} `

		** 只声明不初始化，使用 zero-value
			`var slice1 []int 	// slice1 use zero-value: nil`

		** 当变量初始值能推断出类型时，可省略变量的类型声明
		** 支持使用 short-hand notation 来声明数组变量
		** 数组字面量支持多行声明，最后一行必须添加 `,`（逗号），
			否则 go compiler 会在最后一行添加分号，导致语法错误

		** empty slice 初始化，使用 make 函数
			` slice1 := make([]Type, len, cap) `

			make 函数签名：func make([]Type, len int, cap int) []Type
				cap 可选

	Slicing Expression：
		除了使用 slice 字面量来初始化 slice 的值，
			golang 还支持 slicing expression 来创建 slice

		参考：下面的 operationSlicing 函数

		* Simpleslice expressions
			** spec: ` a[low : high] `
			** `a` 可以为 字符串、数组、数组指针、slice
			** low、high 均为 `a` 中的索引值，low、high 均可省略，但 ` : `（冒号）不可省略
			** low、high 有效值范围：
				- 若 `a` 是数组或字符串：0 <= low <= high <= len(a)
				- 若 `a` 是指针：
					-- 0 <= low <= len(a)
					-- 若 high 被忽略，high 等于 len(a)
					-- 若 high 没被忽略，显式声明的话：low <= high <= cap(a)
						即，极值等于 cap(a)，
						这里跟 high 被忽略时的极值不一致

			** 从 a 中生成子 slice 片段，片段取值区间为：[low, high)，即包含 low，不包含 high
			** 子 slice 片段的元素包含：a[low]、a[low + 1]、...、a[high - 1]
			** 子 slice 片段长度为 high - low

		* Full slice expressions
			** spec：` a[low : high : max] `
			** `a` 可以为数组、数组指针、slice，不能是字符串
			** 可指定子 slice 的 cap：max - low

	Slice Operation：
		* len，获取 slice 长度，使用内置函数：len(slice)
		* cap，获取 slice 容量，使用内置函数：cap(slice)
		* getter，使用索引即可访问 slice 中指定位置的值，` slice1[0] `
			index 最大值为 len(slice) - 1
		* setter，使用索引来设置 slice 中制定位置的值，` slice1[0] = ... `
			注意，setter 的操作会反映到 slice 引用的数组中
			index 最大值为 len(slice) - 1
		* copy，使用预设方法：` copy(dst []Type, src []Type) int `
			- 参考下面的 `operationCopy` 用例
			- copy 方法会将 src 中元素赋值到 dst 中，从索引 0 开始，并返回复制的元素数量
			- 当 dst 长度大于 src 时，则会将 src 的所有元素覆盖到 dst 中，从索引 0 开始，
				被复制的元素数量为 src 的长度
			- 当 dst 长度小于 src 时，则会读取 src 中的一个子片段 ` src[0:len(dst)] `，
				并将子片段的值覆盖到 dst 中，从索引 0 开始，
				被复制的元素数量为 dst 的长度
		* append，使用预设方法：` append(src []Type, elem1 Type, ...) []Type `
			- 参考下面的 ` operationAppend ` 用例
			- 假设要往 src 中 append n 个元素：
			- 若 cap(src) >= len(src) + n，则意味着 slice 引用的数组还有足够空间放置新元素，
				此时，更新这 n 个元素到原 slice 引用的数组中，
					创建一个新的 slice 指向新的 index 区间：
				举例：
					-- `arr := [5]int{0, 1, 2, 3, 4}`
						创建了一个长度为 5 的数组
					-- `slice := arr[0:1]`
						slice 指向数组 arr，ptr 指向 arr[0]，len 为 1， cap 为 5，
						即：slice 目前的内容为：[0]
					-- `slice2 := append(slice, -1, -2)`
						往 slice 中添加了两个元素 -1、-2，
						arr 目前的内容为：[0, -1, -2, 3, 4]，
						slice 目前的内容为：[0]，
						slice2 目前的内容为：[0, -1, -2]
					-- 可以看到，此时，append 会修改原数组内容，原 slice 不会被改变，新返回的 slice2 则包含了新 append 的内容

			- 若 cap(src) < len(src) + n，则意味着 slice 引用的数组没有足够的空间容纳新的 n 个元素
				此时，创建新的数组，来包含原 slice 中的元素以及新的 n 个元素，
					注意，只会把原 slice 的内容复制，而不是把原 slice 指向的数组内容全部复制。
					创建一个新的 slice 指向新的 index 区间：
				举例：
					-- `arr := [2]int{0, 1}`
						创建了一个长度为 2 的数组
					-- `slice := arr[0:1]`
						slice 指向数组 arr，ptr 指向 arr[0]，len 为 1， cap 为 2，
						即：slice 目前的内容为：[0]
					-- `slice2 := append(slice, -1, -2)`
						往 slice 中添加了两个元素 -1、-2，
						此时，会创建新的数组来容纳原 slice 元素和新的元素
						arr 目前的内容为：[0, 1]，
						slice 目前的内容为：[0]，
						slice2 目前的内容为：[0, -1, -2]
					-- 可以看到，此时，append 不会修改原数组，原 slice 不会被改变，新返回的 slice2 则仅包含了原 slice 元素和新 append 的内容
		* comparison
			slice 只允许跟 nil 比较，甚至不能跟 slice 自己比较：
				- success: ` slice1 == nil `
				- error: ` slice1 == slice1 `
				- error: ` slice1 == slice2 `

		* unpack
			spec：` []Type... `，类似于 js 中的 spread operator，将 array 或 slice 的元素结构到函数实参列表中
			当调用 append 时，可以 unpack slice 实现快速将 slice 元素拆开并添加

			```go
				slice1 := []int{1, 2}
				slice2 := []int{3, 4, 5}

				// 将 slice2 的元素 unpack，添加到 slice1 中
				slice3 := append(slice1, slice2...)
			```

		* looping：
			参考下面的 sliceLoopping 用例

	Memory Optimization：
		因为 slice 会一直保留对原数组的引用依赖，导致数组本身可能不会被 GC 回收，当数组体量大的时候，可能会引起内存性能问题

		所以，小心 slice 的内存问题，比如
			- 若函数打算返回 slice，则尽量对原 slice 进行定长复制，避免原 slice 中的大数组被持久分发引用
				Below is a bad program：
					```go
						func getCountries() []string {
							countries := []string{
								"United states", "United kingdom", "Austrilia",
								"India", "China", "Russia",
								"France", "Germany", "Spain",
							} // can be much more
							return countries[:3]
						}
					```

				The following program is a good program:
					```go
						func getCountries() []string {
							countries := []string{
								"United states", "United kingdom", "Austrilia",
								"India", "China", "Russia",
								"France", "Germany", "Spain",
							} // can be much more

							c = make([]string, 3) // made empty of length and capacity 3
							copy(c, countries[:3]) // copied to `c`
							return c
						}
					```
	参考文章：
		* [The anatomy of Slices in Go](https://medium.com/rungo/the-anatomy-of-slices-in-go-6450e3bb2b94)
		* [Go Slices: usage and internals](https://blog.golang.org/go-slices-usage-and-internals)
		* [golang spec#Slice_expressions](https://golang.org/ref/spec#Slice_expressions)
*/

package main

import "fmt"

func sliceInitialize() {
	// slice1 use zero value, slice1 == nil
	var slice1 []int

	slice2 := []int{1, 2, 3}
	slice3 := []int{
		1,
		2,
		3,
	}

	// make 创建 empty slice，slice 中每个元素都为 zero value
	// slice4 == [0, 0, 0]
	slice4 := make([]int, 3, 4)

	fmt.Println("------- sliceInitialize -------")

	fmt.Println("slice zero value == nil", slice1 == nil)
	fmt.Println(slice2)
	fmt.Println(slice3)
	fmt.Println(slice4)

	fmt.Println("------- sliceInitialize -------")
}

func relationOfArrayAndSlice() {
	arr := [...]int{0, 1, 2, 3, 4, 5}

	slice1 := arr[:]
	slice2 := arr[2:]
	slice3 := arr[:4]
	slice4 := arr[2:4]

	fmt.Println("------- relationOfArrayAndSlice -------")
	fmt.Println("slice1 := arr[:]", slice1, len(slice1), cap(slice1))
	fmt.Println("slice2 := arr[2:]", slice2, len(slice2), cap(slice2))
	fmt.Println("slice3 := arr[:4]", slice3, len(slice3), cap(slice3))
	fmt.Println("slice4 := arr[2:4]", slice4, len(slice4), cap(slice4))

	arr[3] = 33

	fmt.Printf(
		"after arr[3] = 33\narr[3]=%v, slice1[3]=%v, slice2[1]=%v, slice3[3]=%v, slice4[1]=%v \n",
		arr[3],
		slice1[3],
		slice2[1],
		slice3[3],
		slice4[1],
	)

	fmt.Println("------- relationOfArrayAndSlice -------")

}

// https://golang.org/ref/spec#Slice_expressions
func operationSlicing() {
	arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	slice1 := arr[:]
	slice2 := slice1[2:]
	slice3 := slice1[:7]
	slice4 := slice1[2:7]
	slice5 := slice4[1:3]
	slice6 := slice4[4:]
	slice7 := slice4[len(slice4):] // lower 取值范围：[0, len(slice)]，upper 取值范围：[lower,cap(slice)]
	slice8 := slice3[6:10]

	fmt.Println("------- operationSlicing -------")
	fmt.Println("slice1 := arr[:]", slice1, len(slice1), cap(slice1))
	fmt.Println("slice2 := slice1[2:]", slice2, len(slice2), cap(slice2))
	fmt.Println("slice3 := slice1[:7]", slice3, len(slice3), cap(slice3))
	fmt.Println("slice4 := slice1[2:7]", slice4, len(slice4), cap(slice4))
	fmt.Println("slice5 := slice4[1:3]", slice5, len(slice5), cap(slice5))
	fmt.Println("slice6 := slice4[4:]", slice6, len(slice6), cap(slice6))
	fmt.Println("slice7 := slice4[5:]", slice7, len(slice7), cap(slice7))
	fmt.Println("slice8 := slice3[6:]", slice8, len(slice8), cap(slice8))
	fmt.Println("------- operationSlicing -------")
}

func operationGetterAndSetter() {
	arr := [...]int{0, 1, 2, 3, 4, 5}

	slice1 := arr[:]

	fmt.Println("------- operationGetterAndSetter -------")

	fmt.Println("getter slice1[3]", slice1[3])

	slice1[3] = 33
	fmt.Println("setter slice1[3]=33", slice1[3])

	// panic error，index 最大值为 len(slice) - 1
	// slice1[len(slice1)]
	// slice1[len(slice1)] = 6

	fmt.Println("------- operationGetterAndSetter -------")
}

func operationCopy() {
	var slice1, slice2 []int

	fmt.Println("------- operationCopy -------")

	slice1 = []int{1, 2, 3}
	slice2 = []int{4, 5}

	copy(slice1, slice2)
	fmt.Println("copy elements from shorten slice to longer slice", slice1, len(slice1), cap(slice1))
	slice1[0] = -1 * slice1[0]
	fmt.Println("slice1[0] = -1 * slice1[0]", slice1, slice2)

	slice1 = []int{1, 2, 3}
	slice2 = []int{4, 5}

	copy(slice2, slice1)
	fmt.Println("copy elements from longer slice to shorten slice", slice2, len(slice2), cap(slice2))

	fmt.Println("------- operationCopy -------")
}

func operationAppend() {
	fmt.Println("------- operationAppend -------")

	fmt.Println(">> origin arr and slice1")
	arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println("arr: ", arr)

	slice1 := arr[2:4]
	fmt.Println("slice1 := arr[2:4]", slice1, len(slice1), cap(slice1))

	fmt.Println("")
	fmt.Println(">> slice2 := append(slice1, -4, -5, -6)")
	fmt.Println("slice1 has enough capacity, use the original array, thus append operation changed the original array")
	slice2 := append(slice1, -4, -5, -6)
	fmt.Println("arr: ", arr)
	fmt.Println("slice1: ", slice1, len(slice1), cap(slice1))
	fmt.Println("slice2: ", slice2, len(slice2), cap(slice2))

	fmt.Println("")
	fmt.Println(">> origin arr2 and slice3")

	arr2 := [...]int{1, 2, 3}
	fmt.Println("arr2: ", arr2)

	slice3 := arr2[:]
	fmt.Println("slice3 := arr2[:] ", slice3, len(slice3), cap(slice3))

	fmt.Println("")
	fmt.Println(">> slice4 := append(slice3, 4, 5, 6)")
	fmt.Println("slice3 didn't have enough capacity, append operation will create new array, the original array did not change")
	slice4 := append(slice3, 4, 5, 6)
	fmt.Println("arr2: ", arr2)
	fmt.Println("slice3: ", slice3, len(slice3), cap(slice3))
	fmt.Println("slice4: ", slice4, len(slice4), cap(slice4))

	fmt.Println("------- operationAppend -------")
}

func operationDelete() {
	fmt.Println("------- operationDelete -------")

	arr := [...]int{0, 1, 2, 3, 4, 5, 6}

	// delete element at index 3
	slice1 := append(arr[:3], arr[4:]...)

	fmt.Println(arr, slice1)

	fmt.Println("------- operationDelete -------")
}

func operationSliceComparison() {
	var slice1 []int
	fmt.Println("------- operationSliceComparison -------")
	fmt.Println("nil slice == nil", slice1 == nil)
	fmt.Println("------- operationSliceComparison -------")
}

func unpackOperator() {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{}

	fmt.Println("------- unpackOperator -------")
	fmt.Println("slice1 = ", slice1)
	fmt.Println("slice2 = ", slice2)

	slice2 = append(slice2, slice1...)
	fmt.Println("slice1 = ", slice1)
	fmt.Println("slice2 = ", slice2)

	slice2[1] = -1 * slice2[1]
	fmt.Println("slice1 = ", slice1)
	fmt.Println("slice2 = ", slice2)

	fmt.Println("------- unpackOperator -------")
}

func sliceLoopping() {
	slice1 := []int{0, 1, 2, 3, 4, 5}

	fmt.Println("------- sliceLoopping -------")

	fmt.Println(">> Looping with for clause")
	for index := 0; index < len(slice1); index = index + 1 {
		fmt.Printf("index: %v, value: %v\n", index, slice1[index])
	}

	fmt.Print("\n>> Looping with for-range \n")
	for index, value := range slice1 {
		fmt.Printf("index: %v, value: %v\n", index, value)
	}

	fmt.Println("------- sliceLoopping -------")
}

func memoryOptimization() {
	fmt.Println("------- memoryOptimization -------")

	fmt.Println("------- memoryOptimization -------")
}

func main() {
	sliceInitialize()

	// relationOfArrayAndSlice()

	operationSlicing()
	// operationGetterAndSetter()
	// operationCopy()
	// operationAppend()
	// operationDelete()
	// operationSliceComparison()
	// unpackOperator()
	// sliceLoopping()
	// memoryOptimization()
}
