/*

String:

	字符串编码方案：UTF-8

	创建字符串字面量：
		* 双引号："foo bar\n zoo"
			支持字符转义
		* 反引号 (Raw String)：`foo\n "bar" 'zoo'`
			不支持字符转义
			支持多行
		* 不支持单引号创建字符串字面量
		* 更多例子：
			```go
			"中国香港"                                 			// UTF-8 input text
			`中国香港`                                 			// UTF-8 input text as a raw literal
			"\u4e2d\u56fd\u9999\u6e2f"                    		// the explicit Unicode code points
			"\U00004e2d\U000056fd\U00009999\U6e2f"        		// the explicit Unicode code points
			"\xe4\xb8\xad\xe5\x9b\xbd\xe9\xa6\x99\xe6\xb8\xaf"  // the explicit UTF-8 bytes
			```


	实际上，字符串是将字符 code point 进行 utf-8 编码后得到的只读字节 slice
		> In Go, ** a string is in effect a read-only slice of bytes. **
		关键点：
			* 字符串本质上是 bytes slice
			* slice 中元素是字符串 utf-8 编码的字节序列
			* 字符串是只读的
		因为字符串是 utf-8 的字节序列，鉴于 utf-8 的编码方式，有几个关键点需要注意：
			* len() 得到的长度并不是字符串中包含的字符数量，而是这些字符被 utf-8 编码后占据的字节总数
			* 位置索引（index）反应的是utf-8 编码后字节序列的索引，而不是字符在字符串中的位置
			* 只有当字符串内容是 ASCII 字符时，len() 才等同于字符数量，index 才等同于字符串中指定位置的字符
		可以说，utf-8 编码给数据的存储、传输带来很大的好处，但对程序编写来说，一定程度上是带来负担的


Rune:
	类似于 char，本质上是 int32 类型的数值
	使用 4 字节的整数来表示 utf-8 编码的 unicode code point
		4 字节足够表示基础平面的 unicode 字符，
		与 utf-16 类似的长度
	使用 `[]rune` 类型的数组可以一定程度上表达字符串中的字符数组

	创建 Rune 字面量：
		* 单引号：'中国香港'

	String 转换为 []rune：
		` []rune("中国香港") `

用法：
	* 正确获取字符串中字符数量：
		```go

		import "unicode/utf8"
		import "strings"

		const str = "中国香港"

		utf8.RuneCountInString("世界") 	// 方法一

		len([]rune(str))				// 方法二

		strings.Count(str, "") - 1		// 方法三
		```

	* 正确获取字符串中指定索引处的字符：
		```go

		const str = "中国香港"

		[]rune(str)[1]
		[]rune(str)[3]
		```

	* 遍历字符串中的字符
		```go

		import "fmt"

		const str = "中国香港"

		// 方式一，获取字符串长度后，从 0 到 length-1 循环

		// 方式二，使用 range，配合 iterate 协议，获取 rune 单元以及 rune 单元所在的字符串 byte slice 位置索引
		// 注意，range 的 index 是 byte slice 的索引，而不是字符在字符串中位置
		for index, char := range str {
			fmt.Printf("charactor at index %d is %c\n", index, char)
		}
		```

	* 更多对 string 的操作参照内置的 string package：
		[Package strings](https://golang.org/pkg/strings/)


参考文章：
	* [golang spec#String_literals](https://golang.org/ref/spec#String_literals)
	* [String Data Type in Go](https://medium.com/rungo/string-data-type-in-go-8af2b639478)

*/

package main

import (
	"fmt"
	_ "unicode/utf8"
)

func main() {
	fmt.Printf("%c", []rune("世界")[1])
	fmt.Println("")

	const str = "中国香港"

	for index, char := range str {
		fmt.Printf("charactor at index %d is %c\n", index, char)
	}
}
