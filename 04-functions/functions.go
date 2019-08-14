/*

Functions

	Function 定义语法 Spec：
		```go
			FunctionDecl = "func" FunctionName Signature [ FunctionBody ] .
			FunctionName = identifier .
			FunctionBody = Block .
		```

		关键字 `func` 来声明函数

	函数名规范约定：
		变量名使用单词或 camelCase 单词组合：
			Good: doSomething
			Bad: do_something、DoSomething、dosomething

	Function Signature (函数签名)：
		签名语法：
			* func funcName(arg1 paramType1, arg2 paramType2, ...) returnType {}
			* Return Multiple Values（多返回值）：
				** func funcName(arg1 paramType1, arg2 paramType2, ...) (returnType1, returnType2) {}
			* Named Return Value（命名返回值）：
				** func funcName(arg1 paramType1, arg2 paramType2, ...) (returnVar1 returnType1, returnVar2 returnType2) {}
				** 可归并同类型的返回值：
					func funcName(arg1 paramType1, arg2 paramType2, ...) (returnVar1, returnVar2 returnType) {}

	Return Multiple Values（多返回值）：
		* 在函数签名中声明返回值的类型列表，比如：
			```go
				func convert(a int, b int) (int, int) {	 // 声明返回值类型列表
					return b, a	// 返回多个值
				}

				x, y := convert(1, 2)		// 接受函数的多个返回值语法
			```


	Named Return Value（命名返回值）
		* 在函数签名中声明返回值的名称以及类型，同类型的返回值可合并声明，比如：
			`func convert(a int, b int) (x int, y int) {}`
			`func convert(a int, b int) (x y int) {}`

		* 声明的返回值，会在函数调用时，在函数作用域内初始化同名的变量，并在 return 时将这些变量的值作为返回值传递出去
		* 当函数逻辑完结时，必须要手动使用 `return` 语句
			** 若 return 语句显式声明了函数返回值，则会使用 return 的返回值作为函数的返回值
			** 若 return 语句没有声明函数返回值，则会将函数签名中返回声明的同名变量的值作为返回值
		* 返回值命名其实就是多返回值函数的扩展，新增一个自动创建局部变量并将局部变量自动返回的语法糖而已
		* 函数返回值的顺序以 return 中声明的返回值顺序或函数签名中的返回值列表顺序为准，
		* 接受函数返回值时，不需要与命名返回值有同样的名称：
				```go
				func convert(a int, b int) (x, y int) {	 // 声明返回值名称并合并类型
					x = b
					y = a
					return
				}

				foo, bar := convert(1, 2)		// 接受函数返回值的变量不需要跟返回值变量名称一致，
				```

	`defer` 语句
		用法：
			在函数体中使用 `defer` 关键字来声明一个 defer 函数或方法调用，具体机制：
				* 获取要被调用的函数体：函数变量、表达式等等；
				* 获取要传进函数的实参
				* 延迟该函数的调用
				* defer 函数实际调用时机
					** after: any result parameters are set by that return statement
						函数在父函数设置好返回参数之后调用：
						*** 父函数的调用到达函数底部
						*** 父函数的调用到达 return 语句
						*** the corresponding goroutine is panicking

					** before: the function returns to its caller
						函数在父函数将返回值返回给调用者之前调用
						*** 若父函数声明了命名返回值，则 defer 函数可以通过指针引用的方式修改到父函数返回值
							```go
							// f returns 42
							func f() (result int) {
								defer func() {
									// result is accessed after it was set to 6 by the return statement
									result *= 7
								}()
								return 6
							}
							```

			关键点：
				* defer 语句执行时，要被 defer 执行的函数及其实参都会马上进行取值，而不是等函数真正被执行是才取值
					```go

					// print: 1 PM
					func f() () {
						time := "1 PM"
						defer func(in: string) {
							fmt.Println(time)
						}(time)					// 函数实参会在此时取值：1 PM

						time := "2 PM"
					}
					```

			defer stack:
				函数内部类似存在一个 defer 堆栈的结构来处理多个 defer 函数调用（先进后出）
				先运行的 defer 语句，其函数实际调用时机会晚于后续运行的 defer 语句产生的函数调用：
					```go

						// prints 3 2 1 0 before surrounding function returns
						func f() {
							for i := 0; i <= 3; i++ {
								defer fmt.Print(i)
							}
						}
					```

		语法 spec：
			` DeferStmt = "defer" Expression . `
			其中的 expression 必须为函数调用，且不允许出现括号包裹，例如：
				```go

				import "fmt"

				func convert(a int, b int) (int, int) {
					// Error: defer (fmt.Println("print something"))
					// Error: defer 1+1

					defer fmt.Println("print something")
					return b, a
				}

				```

	函数类型（Function Type）：
		函数类型由函数签名来组成：形参列表（类型、数量）、返回值类型（类型、数量），与函数名无关、与形参名称无关

		当两个函数的形参列表、返回值类型均相同，则认为两个函数同类型：
				` func append(slice []Type, elms ...Type) []Type `
				` func prepend(slice []Type, elms ...Type) []Type `
			上面的 append、prepend 函数属于相同类型，尽管他们的函数名字不一样，他们属于同一个函数类型：
				` func ([]Type, ...Type) []Type `

		定义函数类型：
			` "type" TypeName "func" Signature `

			例如：
				` type CalcFunc func(int, int) int`

	匿名函数（anonymouse function）：
		创建函数是不提供函数名，常用于创建函数字面量
			```go
				var convertFn = func (x int, y int) (int, int) { ... }
			```

	快速调用函数（Immediately-invoked function）：
		创建函数的同时完成函数调用：
			```go
				var x, y = 1, 2

				x, y = func (x int, y int) (int, int) { ... }(x, y)
			```

	参考文章：
		* [golang spec#Function_declarations](https://golang.org/ref/spec#Function_declarations)
		* [golang spec#Defer_statements](https://golang.org/ref/spec#Defer_statements)
		* [golang spec#Function_types](https://golang.org/ref/spec#Function_types)
		* [golang spec#Function_literals](https://golang.org/ref/spec#Function_literals)
		* [The anatomy of Functions in Go](https://medium.com/rungo/the-anatomy-of-functions-in-go-de56c050fe11)
*/

package main

import "fmt"
import "errors"

type calcFunc func(int, int) int

func namedReturnValue(in int) (y int, x int) {
	x = in * 2
	y = in * 3
	return x, y
}

func returnMultipleValue(in int) (int, int) {
	return in * 2, in * 3
}

func myThrowAndDefer() (int, error) {
	defer fmt.Println("Defer myThrowAndDefer")

	return 0, errors.New("myThrowAndDefer Occur Error")
}

func calcNumber(x int, y int, calc calcFunc) int {
	return calc(x, y)
}

func main() {
	var d, e int
	d, e = namedReturnValue(100)

	a, b := returnMultipleValue(100)

	defer func(in int) {
		fmt.Println("Immediate Function Call", in)
	}(1)

	defer fmt.Println("fmt.Println")

	defer fmt.Println(d, e)
	fmt.Println(a, b)

	fmt.Println(calcNumber(1, 2, func(x int, y int) (result int) {
		result = x + y
		return
	}))

	fmt.Println(calcNumber(1, 2, func(x int, y int) (result int) {
		result = x * y
		return
	}))
}
