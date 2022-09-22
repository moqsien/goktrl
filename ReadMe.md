### goktrl

------------------
goktrl是一个用于交互式进程管理库。可以帮助您的后端应用程序轻松实现交互式的进程内部状态管理。

### 主要特点

------------------
- 交互式shell
- 通过Unix Domain Socket管理正在运行的进程
- shell终端支持表格显示，表格字段支持以"order"标签的值作为排序标准，如果没有order标签，则按照字段名排序
- shell命令支持可选参数解析，使用的是[goframe](https://goframe.org/pages/viewpage.action?pageId=35357529)参数解析组件
- 清晰直观，服务端和客户端一起编写，方便集成；后端集成了goktrl之后，就可以实现在不影响项目运行的情况下，查看进程内部状态，开启或停止goroutine等；

### 使用方法

------------------
```shell
go get "github.com/moqsien/goktrl@v1.1.0"
```
```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moqsien/goktrl"
)

type Data struct {
	Addition []interface{}          `order:"4"` // 表格会按照order进行字段排序
	Name     string                 `order:"1"`
	Price    float32                `order:"2"`
	Stokes   int                    `order:"3"`
	Sth      map[string]interface{} `order:"5"`
}

func Info(k *goktrl.KtrlContext) {
	all := k.Parser.GetOpt("all")
	fmt.Printf("$$$client: all=%s\n", all)

	result, err := k.GetResult() // 向服务端发送请求，会自动携带shell收集到的命名参数，例如all
	if err != nil {
		fmt.Println(err)
		return
	}
	content := &[]*Data{}
	err = json.Unmarshal(result, content)
	k.Table.AddRowsByListObject(*content) // 渲染表格
}

func Handler(c *gin.Context) {
	all := c.Query("all")
	fmt.Printf("$$$server: all = %v\n", all)

	Result := []*Data{
		{Name: "Apple", Price: 6.0, Stokes: 128, Addition: []interface{}{1, "a", "c"}},
		{Name: "Banana", Price: 3.5, Stokes: 256, Addition: []interface{}{"b", 1.2}},
		{Name: "Pear", Price: 5, Stokes: 121, Sth: map[string]interface{}{"s": 123}},
	}
	content, _ := json.Marshal(Result)
	c.String(http.StatusOK, string(content))
}

var SName = "info"

func ShowTable() {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name: "info",                               // 命令名称
		Help: "【show info】Usage: info -a=<sth.>", // 帮助信息
		Func: Info, // shell客户端钩子函数
		Opts: goktrl.Opts{&goktrl.Option{
			Name:      "all,a", // 参数名称和别名
			NeedParse: true,    // 是否需要解析，针对-xxx等无需传值的标记参数，详见goframe
			Must:      true,    // 是否不能为空，设置为true后会自动检测shell命令的参数是否已传
		}},
		ShowTable:   true,      // 开启表格显示
		KtrlHandler: Handler,   // 服务端视图函数
		SocketName:  SName,     // 默认unix套接字名称
	})
	go kt.RunCtrl()
	kt.RunShell()
}

func main() {
	ShowTable()
}
```
- 效果图
![shell-0](https://github.com/moqsien/goktrl/blob/main/docs/0.png)
![shell-1](https://github.com/moqsien/goktrl/blob/main/docs/1.png)
![shell-2](https://github.com/moqsien/goktrl/blob/main/docs/2.png)
- [examples/ktrl/ktrl.go](https://github.com/moqsien/goktrl/blob/main/examples/ktrl/ktrl.go)

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
