/*

Go Packages：

	一个文件夹即为一个 package

	声明 package 名称：
		- 使用关键字 `package [packageName]` 来声明 package 名字；
		- package name 命名规范：全小写字母或数字，不需要任何连接符
			* Good: a, fooa, foob,
			* Bad: A, foo_a, foo-a, fooA, FooA

	Packages 类型：
		- program：可执行程序
			当 packageName 为 `main` 时，表示当前 package 为 program；
			program package 同时需要一个 main；
			当运行 `go install` 后会生成 executable file 到 `bin/` 目录中；
		- library：代码库，可被 import 的库，不可执行
			当 packageName 不等于 `main` 时，表示当前 package 为 library；
			当运行 `go install` 后会生成 archived object 到 `pkg/` 目录中；

		每个 package 只能属于一种类型

	引入（import） package 依赖：
		- 使用关键字 `import` 引入依赖 package
		- 单行依赖引入：`import "[packageName]"`
		- 多行依赖引入，使用括号：
			`import (
				"[APackageName]"
				"[BPackageName]"
			)`
		- 依赖别名，给所依赖的 package 声明一个别名，**文件内** 有效：
			`import [aliasName] "[PackageName]"`，
			当 aliasName 为 "_"（下划线）时，相当于告诉编译器，当前只是把 package 引进来，但还没用到，从而避免错误提醒或被 fmt 工具格式化掉
		- Dot Import，将 package 中的 public api 以顶级变量引入到文件中，文件内访问时，不需要 packageName 作为前缀
			`import . "[PackageName]"`
			不建议使用

	package variables 可见性：
		- 顶级作用域中，首字母大写的 variables，均为 public api，package 内外均可访问
		- 顶级作用域中，非首字母大写的 variables，均为 package scope api，只允许 package 内部访问
		- 非顶级作用域中，所有 variables 均为 private api，只允许当前文件，当前作用域及子作用域允许访问

	导出（export）API：
		- 顶级作用域中，首字母大写的函数、变量、类等值，均为 public api，
		  只有 public api 可被默认导出

	Program execution order（Go 程序执行顺序）：
		当执行一个 go program 时，会进入 `main package execution` 过程：
		```go
				go run *.go
				├── Main package is executed
					├── All imported packages are initialized
					|   ├── All imported packages are initialized (recursive definition)
					|   ├── All global variables are initialized
					|   └── init functions are called in lexical file name order
					└── Main package is initialized
						├── All global variables are initialized
						└── init functions are called in lexical file name order
		```

		* 优先完成所直接依赖 package（Immediate dependency） 的 initialization 过程
			** 递归完成所有间接依赖（Transitive dependency） package 的 initialization 过程
		* 完成当前 package 的 initialization 过程
		* initialization 过程包含两个步骤：
			** 初始化 global variables
			** 按照 package 文件名顺序调用文件中包含的 init function

	Package Global Variables
		* 在 package 顶级作用域中定义的变量均为 package 内的 global variables；
		* 名字以大写字母开头的 global variables 为 public api，package 内外均可访问；
		* 名字以小写字母开头的 global variables 为 package scope api，尽在 package 内全局可访问；
		* 可以认为，global variables 是 package 中所有文件共享的内存：
			** package 中的 global variables 不可重复定义
			** 同 package 内，global variables 跨文件间可读写

	Package Init Functions：
		* package 中每个文件均可声明一个 `init` 函数，在对 global variables 进行赋值会很有用
		* package initialization 中，会按文件名顺序来提取出各文件中声明的 `init` 函数，并调用

	重新考虑 Package 的组成：
		* package 的特点提取下：
			** 一个文件夹组成一个 package，文件夹中的所有 go 文件组成 package 的逻辑；
			** 跨文件间，global variables 可读写；
			** package 的 public api 不需要显式 export，只要是大写字母开头的，都可被其他 package 访问，不管这些 api 写在哪个文件；
			** 跨文件间，所有文件顶级作用域中的函数，不管大写字母开头还是小写字母开头，均可访问；
			** 文件内的 init 函数是按文件名顺序来调用的；
		* 上面几个特点，基本可以对 go package 机制作出一个设想了：
			** package 在进行编译时，会读取目录下所有文件，至少提取以下几个：
				*** global variables 声明与赋值列表；
				*** init functions 列表；
				*** global function 列表；
			** 创建一个 package 作用域；
			** 分析所有 global variables 声明，并完成赋值过程：
				*** 解析 variables 声明列表以及赋值语句；
				*** 能优先确定值的变量优先完成声明和初始化；
				*** 当某个变量存在对其他变量的依赖时，优先完成被依赖变量的声明和初始化；
					**** 这里存在一个递归的过程，当被依赖的变量存在对其他变量的依赖时，需要递归当前过程
					**** 当递归出现死循环时，编译会出错，此时开发者应尝试解决引起死循环的赋值依赖，比如将赋值放在 init 函数中进行
			** 将完成声明及初始化的 global variables 放到 package 作用域中；
			** 在 package 作用域中依次调用 init functions；
			** 将 glboal functions 按文件名以此放到 package 作用域中；
			** 编译时，将 package 作用域中的 package scope api 简单替换下名字之类的即可将 api 设为 non public；
			** 当其他 package 要访问 public api 时，即可从 package 作用域中读取对应的 api；
		* 综上，编写 package 时，需要把文件夹内的所有 go 代码都想作一个整体，因为文件顶级作用域中的所有逻辑，
			最后都会平等地放在 package 作用域中

	参考文章：
	[Everything you need to know about Packages in Go](https://medium.com/rungo/everything-you-need-to-know-about-packages-in-go-b8bac62b74cc)
*/

package main

func main() {
	packageScopeAPI()
}
