/*

Methods:

	方法：
		- 类似于 js 中的原型方法，同类型的实例均可使用
		- 通过一个 Receiver 来作为方法的载体，方法作用域内可访问该 Receiver 中的属性或字段
		- Method 与 Function Field 的差别在：
			-- field 函数只是一个函数字面量，与其 struct 本身几乎没有任何联系，更没有上下文的概念
			-- method 绑定了指定的类型，在调用时可访问该类型实例调用者的属性或字段

	Method 只能在 package scope 中定义，不能在其他作用域中定义

	只有 package scope 定义的类型（Type）才允许在 package 内进行 Method 定义：
		举例：
			- 要对 string 类型进行扩展，比如添加一个 length 函数来获取字符串长度，
			- 因为 string 时 golang 预设的类型，其他 package 无法扩展该类型：
				```go
					// panic error: 只允许对 package 内的类型进行 Method 定义
					func (str string) length() int {
						...
					}
				```
			- fix: 通过在 package 内对预设的 string 类型进行 type alias，即创建一个新的类型：
					` type MyString string `

				此时，就可以给 MyString 类型扩展方法了：
					```go
						func (str MyString) length() int {
							...
						}
					```

	语法：
		```go
			func (r ReceiverType) functionName(...ArgType) ReturnType {
				// do something with r
				// get receiver field: r.fieldName
				// update receiver field: r.fieldName = val...
			}
		```
		可以看出，Method 与 Function 的定义很类似，差别在于 func 关键字后面要声明 Receiver：`(r ReceiverType)`，
			Method 调用时会通过变量 r 来接收指定类型的实例对象

	Receiver 可以为任意类型，比如指针，但不允许使用 interface

	指针类型的 Receiver
		使用指针类型的 Receiver，可以进行引用传递，对原触发者进行修改
		举例：
			```go
				func (p *Person) grow() {
					// 获取指针指向的 struct 对象，并更改字段值
					(*p).age++
				}

				person := Person{}

				// 获取 person 变量的指针对象并触发指针对象上的 grow 方法
				(&p).grow()
			```

		一般来说，Method 的 Receiver 都是使用指针类型的，
			配合 struct 对 pointer 字段的语法糖读写语法，使用起来几乎不会感受到指针 reference 与 dereference 繁琐的地方

	调用方法时，receiver 为值传递
		若 method receiver 不是指针类型，则会复制一个 receiver 对象，并传到 method 中，
		所以，比较好的做法是，将 Method 的 Receiver 类型设置为对应的指针类型

		```go
			// 非指针 Receiver
			func (p Person) growToAnother() {
				p.age++
			}

			// 指针 Receiver
			func (p *Person) grow() {
				p.age++
			}

			person := Person{
				age: 1,
			}

			person.growToAnother()
			// true，调用 growToAnother 时，会先拷贝一个 person 对象，并以副本作为 receiver 来调用方法，而不是 person 本身
			// 所以，growToAnother 没有修改到 person 对象本身
			person.age == 1

			person.grow()
			// true，调用 grow 方法时，会先获取 person 对象的指针对象并以此作为 receiver 来调用方法
			// 此时 grow 方法内部使用指针来修改 person 对象本身
			person.age == 2
		```


	方法定义时 Receiver 会区分是否是指针类型，但调用触发则不区分
		只是用法上不区分，实际作用方式则只会按照方法定义的 Receiver 类型进行

		* 若 Method 声明的 Receiever 为具体的值类型时：
			- 调用时，若调用者为具体类型的对象，则拷贝一个调用者来作为 receiver
			- 调用时，若调用者为具体类型对象的指针对象，则会先使用 `*` 操作符来获取实际的对象作为调用者，然后拷贝一个对象来作为 receiver

		* 若 Method 声明的 Receiver 为类型指针时：
			- 调用时，若调用者为具体类型的对象，则会先使用 `&` 操作符来获取对象的指针，并将指针作为 receiver
			- 调用时，若调用者为指针对象，则以该指针为 receiver

	Struct Type 中，Anonymous Field 的 Method 也可以得到提升
		Struct 匿名字段类型若存在 Method 声明，Method 均符合字段提升的规则和特性，无论 Method 的 Receiver 是否指针类型


	参考文章：
		- [anatomy-of-methods-in-go](https://medium.com/rungo/anatomy-of-methods-in-go-f552aaa8ac4a)

*/

package main

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	age       int
}

type Employee struct {
	*Person
}

// Method 只能在 package scope 中定义，不能在其他作用域中定义
func (p Person) Name() string {
	return p.FirstName + " " + p.LastName
}

func (p Person) growToAnother() {
	p.age++
}

func (p *Person) grow() {
	p.age++
}

// MyString local type of predeclared type
type MyString string

func (str MyString) length() int {
	return len(str)
}

func main() {
	fmt.Println("---------- Go Methods ----------")

	ping := Person{
		FirstName: "ping",
		LastName:  "D xiao",
	}

	fmt.Println(ping.Name())

	fmt.Println(ping.age)

	(&ping).growToAnother()
	fmt.Println(ping.age)
	ping.growToAnother()
	fmt.Println(ping.age)

	ping.grow()
	fmt.Println(ping.age)
	(&ping).grow()
	fmt.Println(ping.age)

	// create string literal: "ping deng", and convert to local type  MyString
	myName := MyString("ping deng")
	fmt.Println(myName, myName.length())

	pingP := &ping

	pingP.FirstName = "ping pointer"
	fmt.Println(pingP.Name())

	(*pingP).LastName = "deng pointer"
	fmt.Println((*pingP).Name())

	employee := Employee{
		Person: &ping,
	}

	fmt.Println(employee.LastName)
}
