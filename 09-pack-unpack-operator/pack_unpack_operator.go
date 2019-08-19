/*

Variadic Function：可变参数函数
	函数可使用 pack operator（...Type）来接受相同类型、可变数量的参数列表，并将这些同类型的参数列表打包为一个 slice
		函数声明中，variadic params 必须为最后形参的最后一位

	语法：` func variadicFunc(elms ...Type) {}`
		可以往 variadicFunc 传入多个 Type 类型的实参，数量不限制，
		函数内部可使用 elms 来获取实参列表，elms 类型为 []Type

	举例：
		* append(slice []Type, elms ...Type)
			append 第一个参数为源 slice，另外还可以接受不定数量的元素，将这些元素添加到新 slice 中，
			append 签名中，使用 pack operator 操作符来接受这些不定数量的实参

Pack Operator:
	Pack Operator 用于 Variadic Function 中（可变参数函数）声明并接受可变参数，并打包为 slice

Unpack Operator:
	Unpack Operator 用于在 Variadic Function 函数调用时，将 slice 中的元素快速拆分，并作为实参传到 Variadic Function 中
	unpack operator 解构只能在调用 Variadic Function 时用来解构 slice，不支持解构其他类型的对象

	语法：` variadicFunc(slice...) `

	举例：
		以 slice append operation 举例：
		- 构建一个原始 sliceOri，拥有两个元素，[1, 2]，
			` sliceOri := []int{1, 2} `
		- 假设存在另一个 sliceAnother，若想把 sliceAnother 中的数量全部添加到 sliceOri：
			-- 不使用 unpack operator，就要遍历 sliceAnother，并逐一执行 append 添加到 sliceOri 中
			-- 使用 unpack operator，一行语句就完成了：
				` slice := append(sliceOri, sliceAnother...) `

总结：
	* 使用 Pack Operator 才能在函数中声明可变参数，使其成为 Variadic Function
	* 只有往 Variadic Function 传参才能使用 Unpack Operator 来解构 slice
	* unpack operator 只支持解构 slice

*/

package main

import "fmt"

func func1(elms ...int) []int {
	return elms[:]
}

func main() {
	fmt.Println(func1([]int{1, 2}...))
}
