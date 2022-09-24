### [goktrl](https://github.com/moqsien/goktrl)

------------------
goktrl是一个用于交互式多进程管理库。可以帮助您的后端应用程序轻松实现交互式的进程内部状态管理。

### 主要特点

------------------
- 交互式shell
- 友好的命令提示，使用方法：[command] help
- 强大的参数配置、自动检测、自动解析功能，详见下面的示例
- 通过Unix Domain Socket管理正在运行的进程
- shell终端支持表格显示，表格字段支持以"order"标签的值作为排序标准，如果没有order标签，则按照字段名排序
- 清晰直观，服务端和客户端一起编写，方便集成；
- 后端集成了goktrl之后，就可以实现在不影响项目运行的情况下，查看进程内部状态，开启或停止goroutine等；

### 使用方法

------------------
```shell
go get -u "github.com/moqsien/goktrl@v1.3.1"
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
![shell-1](https://github.com/moqsien/goktrl/blob/main/docs/1.png)
![shell-2](https://github.com/moqsien/goktrl/blob/main/docs/2.png)

### 适用场景

------------------
- 在不重启进程的情况下，对调整进程中的参数、开启和关闭goroutine、显示进程内部状态等；

### Thanks To

------------------
[dmicro](https://github.com/osgochina/dmicro)

[goframe](https://github.com/gogf/gf)

[gin](https://github.com/gin-gonic/gin)

[ishell](https://github.com/abiosoft/ishell)

[table](https://github.com/aquasecurity/table)
