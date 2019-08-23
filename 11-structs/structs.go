/*

Struct:

	自定义结构体类型（在 Golang 中，没有 class，没有 oop

	Struct 中，需要学习下面几点：
		* 定义 Struct Type
			** 字段定义
			** 字段类型
			** 字段 meta-data，又称 StructTag
		* 初始化 struct
			** struct zero value
			** struct 字段值初始化
		* struct 字段读写
		* struct comparison
		* 字段访问权限控制

	在 Struct Type 中，字段间没有任何联系，也没有所谓的上下文

	zero value:
		Struct 的 zero value 为各字段的 zero value 集合

	值类型，非引用类型：
		struct 对象属于值类型
			** 赋值时，会复制新的 struct，而不是使用引用
			** 将 struct 作为参数传值时，会复制新的 struct，而不是使用引用

	Struct Type 定义：
		* spec 参考：https://golang.org/ref/spec#Struct_types
			```go
				StructType    = "struct" "{" { FieldDecl ";" } "}" .
				FieldDecl     = (IdentifierList Type | EmbeddedField) [ Tag ] .
				EmbeddedField = [ "*" ] TypeName .
				Tag           = string_lit .
			```

			举例：
				```go
					struct {
						fieldName1 TypeName1
						fieldName2 TypeName2

						fieldName3, fieldName4 TypeName3
					}
				```

			** 使用关键字 `struct` 来声明 Struct Type 命名
			** 字段及类型声明：` fieldName TypeName `，
				每个字段单独一行
			** 同类型的字段可归到同一行声明，字段名之间以 `,`（逗号）分隔
				` fieldName3, fieldName4 TypeName `

		* 匿名 Struct Type
			（即没有给 Struct Type 进行类型命名（Type alias））
			匿名 Struct Type 没有重用行，一般仅用于 struct 字面量创建：
				```go
					structVar := struct {
						name string
					} {
						name: "foo"
					}
				```

		* Type Alias Struct Type
			相当于具名 Struct Type，基本 Struct Type 都要 Type Alias
			```go
				type Person struct {
					name string
					age int
				}
			```

			上面创建了一个 Struct Type，并命名为 Person

		* 字段类型
			** 支持几乎所有类型，包括自定义类型，比如 Slice 和 Struct Type
			** 举例
				```go
					type Person struct {
						// 基础类型
						name string
						age int
						male boolean

						// 指针类型
						weight *int
						nickName *string

						// 函数类型
						sayHi func(name string, age int) string

						// Nested Struct Field，嵌套 Struct Type
						mother Person
						father Person
					}
				```
			** Nested Struct Field，嵌套 Struct Type
			** Anonymous fields，或称 Embedded field（匿名字段）
				不声明字段名，只声明字段类型，则会以字段类型的名字作为字段名
				```go
					type Person struct {
						// 声明了一个匿名字段，该字段名等于类型名，
						// 比如这里的 string
						string
					}
				```

				Struct Type 中，匿名字段的类型名不能冲突：
					- Success：
						```go
							// A struct with four embedded fields of types T1, *T2, P.T3 and *P.T4
							struct {
								T1        // field name is T1
								*T2       // field name is T2
								P.T3      // field name is T3
								*P.T4     // field name is T4
								x, y int  // field names are x and y
							}
						```
					- Error:
						```go
							struct {
								T     // conflicts with embedded field *T and *P.T
								*T    // conflicts with embedded field T and *P.T
								*P.T  // conflicts with embedded field T and *T
							}
						```

			** Promoted Fields（ 字段提升 ）
				若匿名字段是 struct type，struct 中的字段名读写会被提升，
				（前提是被提升的字段名没有冲突）
				字段提升的层级无限制，只要不出现字段名冲突即可
					即，任意层级的匿名 struct type 字段均可被提升到顶级 struct 变量中
				达到的效果接近于 Mixins 和 Inherits，非常好用

		* 字段 meta-data（ StructTag ）
			- 在字段声明中，可以在类型后面，以字符串的形式声明相关的元信息，若要容纳多个元信息，以空格相隔
			- 使用文档参考：https://golang.org/pkg/reflect/#StructTag
			- 举例：https://play.golang.org/p/o4SanceyFoI
				```go
					package main

					import (
						"fmt"
						"reflect"
					)

					func main() {
						type S struct {
							F string `species:"gopher" color:"blue"`
						}

						s := S{}
						st := reflect.TypeOf(s)
						field := st.Field(0)
						fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))

					}

				```

	struct 声明及初始化
		* 声明 struct，不初始化，此时使用 zero value
			```go
				var sam Person
			```
		* struct 字面量：
			```go
				Person{
					name: "foo",
					age: 18,

					// Anonymous Field 初始化，需要声明类型名
					string: "Bar"
				}
			```

			若字段初始化顺序于 Struct Type 字段声明顺序一致，则初始化时可不声明字段名：
			```go
				Person{
					"foo",
					18,
					"Bar",
				}
			```
		* 字段初始化时均为赋值行为，到底是值复制赋值还是引用复制赋值，取决于字段类型：
			比如，Struct Type 类型字段的赋值行为是复制值，所以赋值后的字段值与原 struct 是独立的两个 struct 对象
			```go
				type Person struct {
					name string
				}

				type Employee struct {
					Person
				}

				foo := Person{
					name: "foo",
				}

				fooEmployee := Employee{
					Person: foo,
				}

				fooEmployee.name = 'bar'

				// true
				fooEmployee.name == 'bar'

				// true
				foo.name == 'foo'
			```

	struct operation，读写操作：
		- 使用 `.`（句点）操作符来对字段值进行读写
		- 若字段是 Pointer 类型时，golang 提供了语法糖来快速读写 Pointer 类型的字段，跳过繁琐的 `*`、`&` 操作符
			举例：
			```go
				type A struct {
					name   string
					parent *A
				}

				a := A{
					name: "foo",
				}

				b := A{
					name:   "bar",
					parent: &a,
				}
			```

			-- 正常的 pointer getter
				` &(b.parent).name `
			-- 语法糖的 pointer getter
				` b.parent.name `

			-- 正常的 pointer setter
				` &(b.parent).name = "aslkdfj" `
			-- 语法糖的 pointer setter
				` b.parent.name = "alkfjakls" `

	sturct comparison：
		- Struct Comparison: Struct Type 一样，field 的值相等，则 stuct 相等
		- 若 Struct Type 中含有不可比较的字段类型，则 struct 之间不能进行比较
		- 虽然结构体内部是一样，甚至是值都是一样的，但只要是不同的 type alias，struct 就不能进行比较

	Struct Type 访问性、exported fields
		* Struct Type 定义在 package-scope 中，该类型的访问性遵循包变量的访问规则
		* 当字段名以大写字母开头时，该字段可被其他 package 访问

	参考文章：
		- [structures-in-go](https://medium.com/rungo/structures-in-go-76377cc106a2)
		- [golang spec#Struct_types](https://golang.org/ref/spec#Struct_types)
		- [golang reflect#StructTag](https://golang.org/pkg/reflect/#StructTag)

*/

package main

import "fmt"

func sugerGetterAndSetterOfPointerField() {
	type A struct {
		name   string
		parent *A
	}

	a := A{
		name: "foo",
	}

	b := A{
		name:   "bar",
		parent: &a,
	}

	fmt.Println(b.parent.name)

	b.parent.name = "ajdfklajds"
	fmt.Println(a.name)
	fmt.Println(b.parent.name)
}

func main() {
	sugerGetterAndSetterOfPointerField()

	// 定义 struct type 结构体，并设置类型名为 Person
	type Person struct {
		// 字段定义：` fieldName FieldType `
		firstName string
		lastName  string

		// 字段类型可以为任意类型，除了 ` *InterfaceType `
		name func(person *Person) string

		// 同类型的字段可以合并类型声明
		weight, height int
	}

	type Employee struct {
		// Nested Struct，即字段类型为 Struct Type
		// Nested Struct Field 的赋值是 Struct 值复制的，即副本
		person Person

		// Pointer Field
		personPointer *Person

		// Anonymous Struct，匿名 Struct Type，重用性极低
		// 初始化该字段时，需要调用者手动重新写一个一样的 Struct Type 声明
		// 感觉 Anonymous Struct 只能用在一次性的 struct 字面量初始化中
		salary struct {
			basic     int
			insurance int
			allowance int
		}
	}

	type Human struct {
		gender int

		// Anonymous fields，或称 Embedded field
		// 匿名字段，若匿名字段是 struct type，struct 中的字段读写会被提升
		// 类似与 Mixin，感觉非常好用
		Person
	}

	var personWithZeroValue Person
	fmt.Printf(
		"Person struct with zero value: %v, (firstName == \"\"): %v, (lastName == \"\"): %v, (name func == nil): %v, (weight and height == 0): %v\n",
		personWithZeroValue,
		personWithZeroValue.firstName == "",
		personWithZeroValue.lastName == "",
		personWithZeroValue.name == nil,
		personWithZeroValue.weight == 0 && personWithZeroValue.height == 0,
	)

	foo := Person{
		firstName: "Foo",
		lastName:  "Lueng",
		name: func(person *Person) string {
			return person.firstName + " " + person.lastName
		},
		weight: 140,
		height: 175,
	}
	fmt.Println("initializeing struct with literal", foo, foo.name(&foo))
	foo.firstName = "Bar"
	fmt.Println("Struct Field Setter", foo.firstName)

	fooEmployee := Employee{
		person:        foo,
		personPointer: &foo,

		// Anonymous Struct 真鸡肋，字段声明和初始化都要单独声明 Struct Type
		salary: struct {
			basic     int
			insurance int
			allowance int
		}{
			basic:     0,
			insurance: 0,
			allowance: 0,
		},
	}

	fmt.Println("Nested  Struct Getter: Employee.person", fooEmployee.person)

	fooEmployee.person.lastName = "Wong"
	fmt.Println(
		"Nested Struct Field Setter,",
		"fooEmployer.person got change: [",
		fooEmployee.person.lastName,
		"], foo.person won't change: [",
		foo.lastName,
		"]",
	)

	fooEmployee.personPointer.lastName = "Hwang"
	fmt.Println("Pointer Field Setter, foo.person will get change:", foo.lastName)

	fooHuman := Human{
		gender: 1,
		Person: foo,
	}

	fmt.Println("---------- Anonymous fields getter and setter ----------")
	fmt.Println("normal getter: fooHuman.Person.name()", fooHuman.Person.name(&fooHuman.Person))
	fooHuman.Person.firstName = "ping"
	fmt.Println("normal seeter: fooHuman.Person.firstName == 'ping'", fooHuman.Person.firstName, foo.firstName)

	fmt.Println("promoted getter: fooHuman.name()", fooHuman.name(&fooHuman.Person))
	fooHuman.firstName = "tao"
	fmt.Println("promoted seeter: fooHuman.firstName == 'tao'", fooHuman.firstName, foo.firstName)

	fmt.Println("---------- Anonymous fields getter and setter ----------")

	fmt.Println("\n---------- Struct Comparison ----------")

	type S1 struct {
		name string
	}

	type S2 struct {
		name string
	}

	type S3 struct {
		age int
		// map1 map[string]int
	}

	// Struct Comparison: Struct Type 一样，field 的值相等
	// 若 Struct Type 中含有不可比较的字段类型，则 struct 之间不能进行比较
	// panic error: 虽然结构体内部是一样，甚至是值都是一样的，但只要是不同的 type alias，struct 就不能进行比较
	// fmt.Println(s1 == s2)
	fmt.Println(S3{1} == S3{2})
	fmt.Println(S3{3} == S3{3})

	fmt.Println("---------- Struct Comparison ----------")

	promotedFields()
}

func promotedFields() {
	type S1 struct {
		name string
	}

	type S2 struct {
		age int
	}

	type S3 struct {
		name string
		S1
		S2
	}

	type S4 struct {
		S3
	}

	type S5 struct {
		S4
		// S1
	}

	s5 := S5{
		S4: S4{
			S3: S3{
				name: "bar",
				S1: S1{
					name: "foo",
				},
				S2: S2{
					age: 18,
				},
			},
		},
	}

	fmt.Println("---------- promotedFields ----------")

	fmt.Println(s5.name, s5.S3.name, s5.S1.name)
	fmt.Println(s5.age)

	s5.age = 20

	fmt.Println("---------- promotedFields ----------")
}
