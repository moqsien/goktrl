### goktrl

------------------
goktrl是一个用于交互式进程管理库。可以帮助您的后端应用程序轻松实现交互式的进程内部状态管理。

### 主要功能

------------------
- 交互式shell
- 通过Unix Domain Socket管理正在运行的进程
- shell终端支持表格显示，表格字段支持以"order"标签的值作为排序标准，如果没有order标签，则按照字段名排序
- shell命令支持可选参数解析，使用的是[goframe](https://goframe.org/pages/viewpage.action?pageId=35357529)参数解析组件

### 使用方法

------------------
```shell
go get -u "github.com/moqsien/goktrl"
```
```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/g"
	"github.com/moqsien/goktrl"
)

type Data struct {
	Name     string                 `order:"1"`
	Price    float32                `order:"2"`
	Stokes   int                    `order:"3"`
	Addition []interface{}          `order:"4"`
	Sth      map[string]interface{} `order:"5"`
}

func Info(k *goktrl.KtrlContext) {
	result, err := k.Client.GetResult(k.KtrlPath, k.Parser.GetOptAll(), k.DefaultSocket)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := &[]*Data{}
	err = json.Unmarshal([]byte(result), content)
	k.Table.AddRowsByListObject(*content)
}

func Handler(c *gin.Context) {
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
		Name: "info",
		Help: "show info",
		Func: Info,
		Opts: &g.MapStrBool{
			"all,a": true,
		},
		KtrlPath:    "/ctrl/info",
		ShowTable:   true,
		KtrlHandler: Handler,
		SocketName:  SName,
	})
	go kt.RunCtrl()
	kt.RunShell()
}

func main() {
	ShowTable()
}

```
- 效果图
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
