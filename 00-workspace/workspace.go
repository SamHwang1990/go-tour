/*
	Go Workspace:

	在 Go 中，所有代码、依赖均被放置到一个 workspace 中，
	建议设置 $GOPATH 路径为 workspace 所在的路径，默认是 `~/go`

	workspace 文件夹结构：
		bin/
		pkg/
		src/

	目录结构说明：
		src/
			- 所有 packages 的源代码目录
			- package 中用 import 声明的依赖均需要在 src 中存在代码

		pkg/
			- 当使用 `go install` 时，若 package 不是 `main package`，即类型为 library 而不是 program，
				则会在 `pkg/` 目录下生成当前系统平台中，对应 package 的 archived object，以 `.a` 作为文件后缀
			- 当使用 `go install` 时，若某个依赖在 `pkg/` 中有 archived object，则直接拿以来的 archived object 来编译，
				而不需要再从 src/ 中 build

		bin/
			- 当使用 `go install` 时，若 package 是 `main package`，即类型为 program 而不是 library，
				则会在 `bin/` 目录下生成当前系统平台中，对应 package 的 executable file，即可执行文件

	Package Archived Object：
		- library 类型的 package 编译后，得到的是 Archived Object，以 `.a` 作为文件后缀，
			包含 package 的二进制代码以及：debug symbols and source information
		- Archived Object 内已包含了 package 所依赖的所有 package 的 Archived Object
		- ** 好处是降低依赖 package 编译时间： **
			假设 packageA 依赖 packageB
				```
				// packageA/packageA.go
				package packagea

				import "packageb"
				```

			假设 packageB 预先进行了 `go install`，生成了 `packageb.a` 文件到 `pkg/` 目录
			则，当 packageA 要进行 `go install` 时，便不需要另行编译 packageB，直接引入 packageb.a 文件即可，
			从而，加快 packageA 的编译时间

	Package Executable File：
		- program 类型的 package 编译后，得到的是 Executable File，
			包含 package 的可执行二进制代码
		- Executable File 内已包含了 package 所依赖的所有 package 的 Archived Object

	`go build` vs `go install`：
		- 两个命令，简单理解，均用于编译 package
		- GOTMPDIR：package 编译过程中的临时文件会默认保存在一个临时文件夹
			若已设置 GOTMPDIR 变量，则临时文件会保存在变量指定的路径；
			若未设置 GOTMPDIR 变量，则临时文件会保存在系统默认的临时文件夹（比如 Unix 系统的 $TMPDIR 变量所指向的路径）
		- library package 编译缓存
			当 package 被依赖，且要触发编译时，编译期若发现仍有效的 library archived object，则会直接拿该 archived object 作为依赖编译；
			缓存有效性判断：
				> The go  build command now detects out-of-date packages purely based on the content of source files, specified build flags, and metadata stored in the compiled packages.
				> Modification times are no longer consulted or relevant.
				> -- https://pocketgophers.com/go-release-timeline/#go1.10tools
		- `build` 与 `install` 命令的差别：
			* 若编译的是 library package 时：
				** `build` 在编译后不会生成任何内容到 `$GOPATH/pkg/` 或当前目录；
				** `install` 在编译后，会在 `$GOPATH/pkg/` 目录生成 package 对应的 library archived object；
			* 若编译的是 program package 时：
				** `build` 在编译后，会在当前目录（调用命令的目录）生成 package 对应的 executable file；
				** `install` 在编译后，会在 `$GOPATH/bin/` 生成 package 对应的 executable file；

	参考文章：
		- [How to Write Go Code](https://golang.org/doc/code.html)
		- [Getting started with Go](https://medium.com/rungo/working-in-go-workspace-3b0576e0534a)
		- [go install vs go build](https://pocketgophers.com/go-install-vs-go-build/)
		- [Build & Install](https://pocketgophers.com/go-release-timeline/#go1.10tools)
*/

package main

import "github.com/SamHwang1990/go-tour/00-workspace/foo"

func main() {
	foo.Print()
}
