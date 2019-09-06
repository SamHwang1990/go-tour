/*

Interface

	Interface：接口
		- 用于定义方法（ Method ）签名集合
		- 本身不需要有方法签名
		- Interface 中的方法名必须唯一
		- Interface 定义内部只有两种类型的表达式：
			* 方法签名
			* 嵌套 Interface 类型

	定义 Interface（ Declaring interface ）
		* spec 语法：
			```go
			InterfaceType      = "interface" "{" { MethodSpec ";" } "}" .
			MethodSpec         = MethodName Signature | InterfaceTypeName .
			MethodName         = identifier .
			InterfaceTypeName  = TypeName .
			```

		* 匿名 Interface
			仅创建了一个 interface 定义，但没有声明为一个类型，
			多用于函数参数类型声明

			```go
				interface {
					Read()  (int, error)
					Write() (int, error)
					Close() error
				}
			```

		* Interface Type
			创建 interface 定义的同时声明为一个类型
			```go
				type File interface {
					Read()  (int, error)
					Write() (int, error)
					Close() error
				}
			```

		* Embedding interfaces
			- Interface 内部允许嵌套 Interface 定义
			- 被嵌套的 Interface 中定义的方法会被提升到外部 Interface
			- Interface 自身定义的方法签名不能与 Embedding Interface 中的方法签名重名
			- Embedding Interfaces 间的方法签名也不允许重名
			- 被嵌套的 Interface 必须是具名 Interface，即 Interface Type，
				不支持嵌套匿名 Interface

			```go
				type ReadWriter interface {
					Read(b Buffer) bool
					Write(b Buffer) bool
				}

				type File interface {
					ReadWriter	// 嵌套 ReadWriter
					Close()
				}
			```

	实现 Interface（ Implementing interface ）
		- 当一个类型实现了 Interface 中声明的所有方法时，即表示该类型实现了指定的 Interface
		- 即，若类型实现了多个 Interface 方法中的方法，表示该类型同时实现了

		- 举例：
			```go
				type Lock interface {
					Lock()
					Unlock()
				}

				type Foo struct {}

				func (f Foo) Lock() {}

				func (f Foo) UnLock() {}
			```
			例子中，
				-- 创建了一个 Interface Type：Lock，该接口需要实现两个方法：Lock、Unlock
				-- 创建了一个 Struct Type：Foo，该 Struct 作为 Receiver，有 Lock、Unlock 方法
				-- 此时，Foo 实现了 Lock 接口

		`pointer` vs `value` receiver
			若方法的 Receiver 声明为指针类型，则 Type 本身并没有实现 Interface，而是 Type 类型的指针实现了 Interface
			因此，在 interface 变量赋值时，需要使用 `&` 操作符先获取对应的指针，参考下面：`pointerReceiver` 函数中的 ` lo = &f `

	Zero Value
		Interface 的 Zero Value 为 nil

	空 Interface（ Empty Interface ）
		Empty Interface，即没有声明任何方法签名的 Interface，
		- Go 中所有的类型都实现了 Empty Interface，即相当于 Typescript 中的 any
		- 多用于函数参数类型声明
		- 语法：
			```go
				interface{}
			```

	`static type`、`concrete type`、`concrete value`
		- Interface 类型自身即为 static type
		- interface 类型变量实际指向的值即为 `concrete value`，又称 `dynamic value`
		- concrete value 的实际类型即为 `concrete type`，又称 `dynamic type`

		- 在 Go Interface 中，Interface 变量与 concrete value 的数量关系是 1:N，
			-- 只要 concrete type 实现了 Interface Type 中的所有方法，该 concrete value 即可赋值给指定的 interface 变量
			-- Concrete type 可以实现多个 Interface
		- Dynamic Type 与 Interface Type 的从属关系是：is-a，
			Dynamic Type(Dynamic Value) is a Interface Type

	Interface 变量使用
		* 假设以下 Interface Type 和 DynamicType 声明：
			```go
				type Lock interface {
					Lock()
					Unlock()
				}

				type Foo struct {}

				func (f Foo) Lock() {}

				func (f Foo) UnLock() {}
			```

		- 变量声明：
			` Lock lo `
		- 变量初始化
			```go
				foo := Foo{}

				Lock lo = foo

				lo.Unlock()
				lo.Lock()
			```
		- 不支持 short-hand notation

	类型推断（ Type Assertion ）
		- 语法：
			```go
				typeInstance, isOk := x.(Type)
			```

			** `x` 为 Interface Value
			** `Type` 为指定的类型
			** 若 x 接口变量的 Dynamic Type 为指定的 Type 类型，则 assertion 返回 Type 实例，isOk 为 true
			** 若 x 接口变量 Dynamic Type 不是指定的 Type 类型，则 assertion 返回 Type 的 zero value，isOk 为 false
			** 若 x 的 Dynamic Value 为 Pointer，则 Type 也要为对应类型的 Pointer 类型，返回的 Type 实例也会是 Pointer 类型

	Type Switching：
		参见：05-flow-control-statements

	参考文章：
		- [golang spec#Interface_types](https://golang.org/ref/spec#Interface_types)
		- [golang spec#Type_assertions](https://golang.org/ref/spec#Type_assertions)
		- [Interfaces in Go](https://medium.com/rungo/interfaces-in-go-ab1601159b3a)
*/

package main

import "fmt"

type ILock interface {
	Lock()
	Unlock()
}

type Foo struct {
}

func (f Foo) Lock() {

}

// 指针类型的 Receiver
func (f *Foo) Unlock() {

}

func (f Foo) Hey() {
	fmt.Println("Foo Hey")
}

type Bar struct {
}

func (b Bar) Hello() {
	fmt.Println("Bar Hello")
}

func (b Bar) Lock() {

}

func (b Bar) Unlock() {

}

func embeddingInterface() {
	type Lock interface {
		Lock()
		Unlock()
	}

	type ReadWrite interface {
		Read(b string) error
		Write() string
	}

	type ReadWriteFile interface {
		Lock
		ReadWrite

		// Panic: 不允许嵌套匿名 Interface
		// interface {
		// 	Close() error
		// }

		Close() error
	}
}

func emptyInterfaceType(param interface{}) {
}

func pointerReceiver() {
	var lo ILock

	f := Foo{}

	// Foo 中的 Unlock 方法使用了指针作为 Receiver，所以，interface 变量赋值时需要使用对应的指针
	lo = &f

	lo.Lock()
	lo.Unlock()

	lo2 := Foo{}
	lo2.Unlock()
}

func typeAssertion() {
	var lo ILock = &Foo{}

	foo := lo.(*Foo)

	foo.Hey()
	(*foo).Hey()

	var loBar ILock = Bar{}

	bar := loBar.(Bar)
	bar.Hello()
}

func main() {
	fmt.Println("Go Interfaces")

	pointerReceiver()

	emptyInterfaceType(1)
	emptyInterfaceType(true)
	emptyInterfaceType("abc")
	emptyInterfaceType('a')
	emptyInterfaceType([]int{1, 2, 3})
	emptyInterfaceType(map[string]int{
		"foo": 1,
	})

	typeAssertion()
}
