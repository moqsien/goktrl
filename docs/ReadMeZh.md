### [goktrl](https://github.com/moqsien/goktrl)

------------------
[En](https://github.com/moqsien/goktrl)

goktrl是一个用于交互式多进程管理的shell库。进程间交互使用unix套接字。

### 主要特点

------------------
- 交互式shell
- 友好的命令参数提示
- 使用结构体和标签作为参数配置
- 参数自动解析和校验
- 自动渲染表格数据，如果开启了相关选项
- 自动处理和显示数据，如果开启了相关选项
- 通过unix套接字连接到进程
- 方便的请求转发功能
- 整体非常直观，配置灵活

### 使用方法

------------------
```shell
go get -u "github.com/moqsien/goktrl@v1.3.6"
```

- 更多示例: [examples](https://github.com/moqsien/goktrl/tree/main/examples/ktrl)
- 最简示例: 

```text
file:   test.go
```

```go
package main

import (
	"fmt"

	"github.com/moqsien/goktrl"
)

/*
  一个最简单的例子: 甚至都不需要shell钩子函数。
*/

func Handler(c *goktrl.Context) {
	fmt.Printf("$$ server: args = %v\n", c.Args) // 自动解析shell传过来的位置参数到c.Args
	Result := map[string]string{
		"hello": "info",
	}
	c.Send(Result)
}

func ShowInfo() *goktrl.Ktrl {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name:            "info",          // 命令名称
		Help:            "show info",     // 命令简短介绍
		KtrlHandler:     Handler,         // shell服务端视图函数
		Auto:            true,            // 是否全自动处理和显示数据
	})
	return kt
}

func main() {
	kt := ShowInfo()
	if len(os.Args) > 1 {
		kt.RunShell()
	} else {
		kt.RunCtrl()
	}
}
```

```shell
go run test.go     # 启动服务端

go run test.go aaa # 启动shell客户端
```

- 示例效果图
![shell-0](https://github.com/moqsien/goktrl/blob/main/docs/0.png)
![shell-1](https://github.com/moqsien/goktrl/blob/main/docs/1.png)
![shell-2](https://github.com/moqsien/goktrl/blob/main/docs/2.png)

### Thanks To

------------------
[dmicro](https://github.com/osgochina/dmicro)

[goframe](https://github.com/gogf/gf)

[gin](https://github.com/gin-gonic/gin)

[ishell](https://github.com/abiosoft/ishell)

[table](https://github.com/aquasecurity/table)
