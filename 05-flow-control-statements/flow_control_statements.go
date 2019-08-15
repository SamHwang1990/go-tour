/*

Flow Control Statements:

	For:
		go 中只有一种循环体：for 循环

		spec 定义：
			```go
				ForStmt = "for" [ Condition | ForClause | RangeClause ] Block .
				Condition = Expression .
			```

		for 循环使用上有三种变体：
			* For statements with single condition
				( 仅包含条件判断语句 )
				举例：
				```go
					for a < b {
						...
					}
				```

			* For statements with for clause
				( 包含 init statement、post statement 的 for 从句 )
				"for clause" spec:
					```go
						ForClause = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
						InitStmt = SimpleStmt .
						PostStmt = SimpleStmt .
					```
				语法：
					** 使用分号（;）来分割从句的三个部分，且两个分号均不能忽略
					** [ InitStmt ] 为 init statement，循环开始之前执行一次
						*** init statement 可创建变量，变量在整个遍历过程均可读写
						*** init statement 可以忽略，但 Condition 之前的分号不能忽略
					** [ PostStmt ] 为 post statement，每次循环体执行后都会执行该语句
						*** post statement 可以忽略，但 Condition 之后的分号不能忽略

				举例：
					```go
						for i := 0; i < 10; i++ {
							...
						}

						// 忽略 init statemnt
						i := 0
						for ; i < 10; i++ {
							...
						}

						// 忽略 post statement
						i := 0
						for ; i < 10 ; {
							...
						}
					```


			* For statements with range clause
				( range 从句 )

				`range clause` spec:
					```go
						RangeClause = [ ExpressionList "=" | IdentifierList ":=" ] "range" Expression .
					```

				语法解析：
					** "range" 右侧的表达式称为 [range expression]
					** "range" 左侧的表达式列表或变量列表称为 [iteration variables]
						*** 当 range expression 为 Channel 时，iteration variables 长度只能等于 1
							**** 循环体中 iteration variable 为 channel 接收到的值
						*** 当 range expression 不为 Channel 时，iteration variables 长度可以是 1 或 2，
							当第二个 iteration variable 设为"_"（blank identifier，下划线）时，从句会忽略第二个 variable 的声明
							**** 第一个 iteration variable 为 索引或 key
							**** 第二个 iteration variable 为 索引或 key 对应的值，
								 当没有声明第二个 iteration variable，循环体可能就不会对该索引或 key 进行取值，
								 在只关心索引或 key 的遍历逻辑中，对性能会有提升
						*** 当 iteration variable 为表达式赋值时（ExpressionList），使用 “=”，每个循环体执行前都会对 variable 进行赋值
						*** 当 iteration variable 为变量声明时（IdentifierList），使用 “:=”，
							- 变量的作用域为 for 循环内部
							- 循环开始之前会创建对应的变量
							- 每个循环体执行前都会重新为变量赋值
					** range expression 执行时机及执行次数
						*** 当 len(range expression) 为常量，且只声明了一个 iteration variable 时，
							range expression 本身不会被取值
						*** 当不符合上面的条件时：len(range expression) 不为常量，或声明了两个 iteration variable，
							range expression 仅会在遍历过程开始之前取值一次

				四种 range expression 值：
					range clause 中的 [range expression] 支持四种类型的值：
						- array or slice
						- string
						- map
						- channel

					不同类型的值，在每次循环体中，iteration variables 的值以表格展示如下：
						```go
							Range expression                          1st value          2nd value

							array or slice  a  [n]E, *[n]E, or []E    index    i  int    a[i]       E
							string          s  string type            index    i  int    see below  rune
							map             m  map[K]V                key      k  K      m[k]       V
							channel         c  chan E, <-chan E       element  e  E
						```

						** array or slice
							- index 的值从 0 开始递增
							- 若只声明了一个 iteration variable，则循环只会遍历数组的长度，
								而不会在每个循环体中主动读取数组中对应索引的值
							- 若 [range expression] 为 nil 的 slice，则遍历会马上结束

						** string
							- 字符串的 range 循环，循环单元是字符串中完整的 utf-8 编码的 unicode code point
							- index variable 为 code point 的第一个字节在字符串 byte slice 中的索引，而不是字符串中的字符索引
							- value variable 为 code point 的 rune 值
							- 当遇到非 utf-8 编码的字符：
								-- value variable 为固定返回 OxFFFD（the Unicode replacement character）
								-- 下一次循环体会固定只读取一个 byte

						** map
							- map 的遍历是无序的
							- index variable 为 map 的 key 值
							- value variable 为 key 对应的 value
							- 若 map 为 nil，则遍历会马上结束
							- 若在循环体中删除了一个尚未遍历的 key，则该 key 不会出现在后续的循环
							- 若在循环体中新创建了一个 key，则该 key 可能出现，也可能不出现在后续的循环中

						** channel
							- 只会有一个 iteration variable，即 value variable
							- 当 channel 收到值时，会触发遍历
							- 当 channel 关闭时，遍历会结束
							- 若 channel 为 nil，则遍历会阻塞整个进程，后续代码不会得到执行
								> If the channel is nil, the range expression blocks forever.

	If-else:
		spec:
			```go
				IfStmt = "if" [ SimpleStmt ";" ] Expression Block [ "else" ( IfStmt | Block ) ] .
			```
		语法：
			* Expression 不用括号包围
			* Block 必须要大括号包围
			* Expression 之前可以添加语句
				** 若为局部变量的赋值语句（short-hand notation），则创建的局部变量在整个 if-else block内都可读写


	Switch-case:
		If-else 语句的另一种扩展模式，应对多条件分支

		spec: ` SwitchStmt = ExprSwitchStmt | TypeSwitchStmt . `

		switch-case 中，所有 case 都默认带 break 机制

		Switch Statement 分为 Expression switches 和 Type switches 两种：

			* Expression switches
				spec:
					```go
						ExprSwitchStmt = "switch" [ SimpleStmt ";" ] [ Expression ] "{" { ExprCaseClause } "}" .
						ExprCaseClause = ExprSwitchCase ":" StatementList .
						ExprSwitchCase = "case" ExpressionList | "default" .
					```

				语法：
					** Switch Expression 不需要括号包围
					** Switch Expression 之前可以添加语句
						*** 若为局部变量的赋值语句（short-hand notation），则创建的局部变量在整个 switches statement 内都可读写
					** 若 Switch Expression 为空，则意味着 switch expression 为 true
					** Switch Expression 只在进入 case 判断前取值一次
					** Switch Statement 中，case expression 的值必须要与 switch expression 的值必须同类型，若类型不一致，编译阶段就会失败
						*** switch expression 的值若没有显示声明变量类型，则会隐式类型转换，比如 `8` 会被转换为 `int(8)`
						*** case expression 的值若没有显示声明变量类型，则会隐式类型转换，比如 `8` 会被转换为 `int(8)`
					** case clause 判断逻辑：
						- case expression 支持列表，用逗号隔开
						- 遍历 case expression 列表
							expression 并取值，
								若 expression 的类型与值均与 switch expression value 相同，则进入当前 case
						- 若 case expression 列表不符合 switch expression，则跳过当前 case
					** 若无任何 case clause 符合 switch expression，则会使用 default clause 中的语句，
						default clause 可以出现在ExprCaseClause 中的任意位置

				`fallthrough statement`:
					fallthrough 表达式可以出现在 Expression Switches 非末尾 clause 中的最后一个语句，
						当前 clause 执行到 fallthrough 语句时，会跳出当前 clause，并跳到下一个 clause 的第一个语句开始执行

					划重点：
						** 只允许出现在 Expression Switches 中的 case clause 或 default clause
						** 所属的 clause 不能是 switch statement 中的最后一个 clause
						** fallthrough 语句必须位于所属 clause 的最后一个语句

				举例：
					```go
						switch tag {
						default: s3()					// default clause 位置没有强制要求
						case 0, 1, 2, 3: s1()			// case clause 支持 expression 列表
						case 4, 5, 6, 7: s2()
						}

						switch x := f(); { 				// 前置局部变量初始化语句 + missing switch expression means "true"
						case x < 0: return -x			// 若 ` x < 0 `，则 expression 的值等于 true
						default: return x
						}

						switch {						// missing switch expression means "true"
						case x < y: f1()
						case x < z: f2()
						case x == 4: f3()
						}

						switch 16 {
						case 16, int(32.0):
							fmt.Println("Expression Switch case 32 or 16")
							fallthrough					// fallthrough 语句执行后会马上执行 default clause
						default:
							fmt.Println("Expression Switch default")
						}
					```

			* Type switches
				用于根据值的类型进行分支处理

				spec:
					```go
						TypeSwitchStmt  = "switch" [ SimpleStmt ";" ] TypeSwitchGuard "{" { TypeCaseClause } "}" .
						TypeSwitchGuard = [ identifier ":=" ] PrimaryExpr "." "(" "type" ")" .
						TypeCaseClause  = TypeSwitchCase ":" StatementList .
						TypeSwitchCase  = "case" TypeList | "default" .
						TypeList        = Type { "," Type } .
					```

				语法：
					** type switch expression 取值并获取类型：`expr.(type)`，".(type)" 只允许在 TypeSwitchGuard 中使用
					** type switch expression 中允许包含局部变量的 short-hand notation
						该局部变量不在 TypeSwitchGuard 中声明，而是在 TypeSwitchCase clause 中才进行单独的声明和初始化
						*** 若 TypeSwitchCase clause 包含多个 Type，则变量会以 type switch expression 的类型来初始化
						*** 若 TypeSwitchCase clause 只包含一个 Type，且 Type 为 nil，则变量会以 type switch expression 的类型来初始化
						*** 若 TypeSwitchCase clause 只包含一个 Type，且 Type 不为 nil，则变量会以该 Type 来初始化
					** TypeSwitchCase clause 中允许一个或多个 Type，用 `,`（逗号）分隔
					** 所有 TypeSwitchCase 中，nil 类型只允许出现一次
					** Type switches 中不支持 fallthrough 语句
					** Switch Expression 之前可以添加表达式，该表达式只会在 Type Switches 中取值一次

				举例：
					```go
						switch i := x.(type) {
						case nil:
							printString("x is nil")                // type of i is type of x (interface{})
						case int:
							printInt(i)                            // type of i is int
						case float64:
							printFloat64(i)                        // type of i is float64
						case func(int) float64:
							printFunction(i)                       // type of i is func(int) float64
						case bool, string:
							printString("type is bool or string")  // type of i is type of x (interface{})
						default:
							printString("don't know the type")     // type of i is type of x (interface{})
						}
					```

					上面的 Type Switches 使用 if-else 语句来改写如下：
					```go
						v := x  								   // x is evaluated exactly once
						if v == nil {
							i := v                                 // type of i is type of x (interface{})
							printString("x is nil")
						} else if i, isInt := v.(int); isInt {
							printInt(i)                            // type of i is int
						} else if i, isFloat64 := v.(float64); isFloat64 {
							printFloat64(i)                        // type of i is float64
						} else if i, isFunc := v.(func(int) float64); isFunc {
							printFunction(i)                       // type of i is func(int) float64
						} else {
							_, isBool := v.(bool)
							_, isString := v.(string)
							if isBool || isString {
								i := v                         // type of i is type of x (interface{})
								printString("type is bool or string")
							} else {
								i := v                         // type of i is type of x (interface{})
								printString("don't know the type")
							}
						}
					```


	参考文章：
		* [golang spec#For_statements](https://golang.org/ref/spec#For_statements)
		* [golang spec#If_statements](https://golang.org/ref/spec#If_statements)
		* [golang spec#Switch_statements](https://golang.org/ref/spec#Switch_statements)
		* [golang spec#Fallthrough_statements](https://golang.org/ref/spec#Fallthrough_statements)

*/

package main

import "fmt"

func main() {
	if a := 0; a == 1 {
		fmt.Println("a == 1", a)
	} else {
		fmt.Println("a != 1", a)
	}

	switch 16 {
	case 16, int(32.0):
		fmt.Println("Expression Switch case 32 or 16")
		fallthrough
	default:
		fmt.Println("Expression Switch default")
	}

}
