/*

Variables：

	变量类型：
		参考：https://gist.github.com/thatisuday/c17e05de591c2e2021ab402e4c2d4bdc#file-medium-go-variables-data-types-csv

	变量名规范约定：
		变量名使用单词或 camelCase 单词组合：
			Good: fooBar
			Bad: foo_bar、FooBar、foobar

	变量声明与初始化：
		* 单变量声明：` var variableName dataType = initialValue `，变量的初始化可选
		* 多同类型变量声明，可以拆开单独声明，亦可合并同类型变量声明：
			` var var1, var2, var3 dataType = value1, value2, value3 `，变量的初始化可选
		* 声明变量时如果同时进行变量初始化，同时初始化值的类型是确定的，则 `dataType` 可省略，go 可推断：
			` var variableName = initialValue `

	zero value：
		若只对变量进行声明不进行初始化，则变量的值会使用对应类型的 `zero value`，例如：
		* boolean 类型：false
		* int、float 类型：0
		* string 类型：empty string
		* rune 类型（其实也是 int32）：0
		* 复杂类型：nil

	Short-hand notation：
		在函数定义中声明并初始化变量时，若可以忽略变量的类型（即初始值的类型可被推断），则可以使用 short-hand notation：
			* 单变量声明：` variableName := initialValue `
			* 多变量声明：` var1, var2, var3 := value1, value2, value3 `
		注意：
			* 只允许在函数体中才允许使用，package 顶级作用域不允许使用
			* 不能用于 constant 变量的声明

	类型转换：
		go 中不存在类型隐式转换，需要手动转换类型：
			* 大部分原声的二元操作符都要求两个操作数属于同一类型：
				```go
				var1 := 10 // int
				var2 := 10.5 // float

				// illegal
				// var3 := var1 + var2

				// legal
				var3 := var1 + int(var2) // var3 == 20
				```

			* 数值操作表达式中，表达式的结果跟操作数的类型有关：
				** 当操作数都是整型时，结果也会是整型：
					`var a = 11/2 // a == 5`
				** 当操作数不全是整型时，结果就一定不是整型：
					`fmt.Printf("%T", 10.0/2) // float64`

	类型别名：
		* 语法：` type aliasName aliasTo `
		* 为类型 `aliasTo` 创建  `aliasName` 的别名类型
		* 别名类型其实可以看作一个变量：
			* 若在文件顶级作用域中声明别名类型，则该类型就类似于 Package Global Variables
				- 该别名类型可被整个 package 使用
				- 若名称首字母大写，则可被其他 package 使用
			* 若在函数体内声明，则该类型仅在该函数体作用域内使用
		* 别名类型与原始类型不是同个类型，因此，类型判断上不会出现相等的情况，需要类型转换

	Constant 变量：
		变量的值仅能在声明语句中初始化，且值不可变的变量：

		单行声明语法：
			* ` const var1, var2 = value1, value2 `
			* 若声明变量时没有初始化，则该变量值会使用对应类型的 `zero value`；
			* 若声明变量时初始化了值，就使用该值

		括号声明语法（ Parenthesized const declaration list ）:
			* 每个变量单独赋值
			```go
			const (
				var1 = value1
				var2, var3 = value2, value3
			)
			```
			* 变量的赋值表达式可以为空（除了第一行的赋值语句），
				此时，变量的赋值语句会使用在当前行之前最近的非空赋值表达式，举例：
				```go
				const (
					var1 = 1 + 1	// 第一行赋值语句不能空
					var2			// 等于 var1 的赋值语句：1 + 1
					var3			// 等于 var1 的赋值语句：1 + 1
					var4 = var2
					var5			// 等于 var4 的赋值语句：var2
				)
				```

	`iota`：
		在 Constant 变量的括号声明语法中，go 预定义了一个按行递增的变量标识符：` iota `，
			括号声明中，在每一行中，iota 的值都不一样，第一行时，值为0，并固定按行递增 1。

		配合括号声明语法中空的赋值表达式会使用当前行之前最近的非空赋值表达式的特性，可以有非常简洁的 Constant 变量初始化写法。

		注意：
			** iota 的值是按行递增的，也就是，同一行的变量赋值表达式中，iota 变量的值是一样的
			** 当要跳过某一行的变量声明时，可以使用 `_`（下划线），通常用于忽略某个 iota 值

		```go
		const (
			a = 1 << iota  // a == 1  (iota == 0)
			b			   // b == 2  (iota == 1)，赋值语句等于：1 << iota
			c = 3          // c == 3  (iota == 2, unused)
			d = 1 << iota  // d == 8  (iota == 3)
		)

		const (
			bit0, mask0 = 1 << iota, 1<<iota - 1  // bit0 == 1, mask0 == 0  (iota == 0)
			bit1, mask1                           // bit1 == 2, mask1 == 1  (iota == 1)
			_, _                                  //                        (iota == 2, unused)
			bit3, mask3                           // bit3 == 8, mask3 == 7  (iota == 3)
		)
		```

	参考文章：
	- [Variables and Constants in Go Programming](https://medium.com/rungo/variables-and-constants-in-go-programming-c715443fa788)
	- [golang spec#Constant_declarations](https://golang.org/ref/spec#Constant_declarations)
	- [golang spec#Iota](https://golang.org/ref/spec#Iota)

*/

package main

import "fmt"

var foo, bar = 1, 2

type float float64

var (
	cat    = 1
	dog    = 2
	monkey = 3
)

func main() {
	fmt.Println(foo)

	var1, var2 := 1, 2

	fmt.Println(var1 + var2)
	fmt.Printf("%T", 10.0/2)
	fmt.Println("")

	fmt.Printf("%T vs %T", float(32), float64(32))
	fmt.Println("")

	const const1, const2 = 1, 2

	const (
		const3 = iota
		_
		const5
	)

	const (
		const6 = 1
		const7
		const8 = 2
		const9
	)

	// 每个 iota 都独立
	const const10, const11 = iota, iota

	fmt.Print(const1, const2, const3, const5, const6, const7, const8, const9, const10, const11)
	fmt.Println("")
}
