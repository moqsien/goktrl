### go开发QQ群推荐
------------------
Golang/Go 语言开发群: 6848027

### [goktrl](https://github.com/moqsien/goktrl)

------------------
goktrl是一个用于交互式进程管理库。可以帮助您的后端应用程序轻松实现交互式的进程内部状态管理。

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
go get -u "github.com/moqsien/goktrl@v1.2.3"
```
```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/moqsien/goktrl"
)

/*
  下面是一个关于goktrl的简短示例。
*/

/*
  表格字段，如果需要显示表格，可以自行定义；
  order标签用于显示时的字段排序，若不设置order标签则按字段名排序；
*/
type Data struct {
	Name     string                 `order:"1"`
	Price    float32                `order:"2"`
	Stokes   int                    `order:"3"`
	Addition []interface{}          `order:"4"`
	Sth      map[string]interface{} `order:"5"`
}

/*
  命令的具名参数(options)的配置；
  必须继承*goktrl.KtrlOption；
  结构体字段名即为参数名；
  标签功能解释：

    alias: 设置别名；
	must: 是否为必传具名参数；
	descr: 具名参数描述信息；
	needparse: 一般不需要用户设置，已根据结构体字段类型进行自动处理；

  支持的字段类型有: string, bool, int, uint, float
*/
type InfOptions struct {
	*goktrl.KtrlOption
	All  bool   `alias:"a" must:"true" descr:"show all info or not"`
	Info string `alias:"i" descr:"infomation"`
}

func Info(c *goktrl.Context) {
	o := c.Options.(*InfOptions)               // 自动解析参数到结构体
	fmt.Printf("## client: options=%v\n", o)   // 打印结构体
	fmt.Printf("## client: args=%v\n", c.Args) // 自动收集命令行普通的位置参数
	result, err := c.GetResult()               // 自动根据参数向服务端发送请求，请求会到达下面的Handler路由方法
	if err != nil {
		fmt.Println(err)
		return
	}
	content := &[]*Data{}
	err = json.Unmarshal(result, content)
	c.Table.AddRowsByListObject(*content) // 如果ShowTable设置为true，此处可添加表格数据，会自动渲染和显示表格
}

func Handler(c *goktrl.Context) {
	o := c.Options.(*InfOptions)                 // 自动解析参数到结构体
	fmt.Printf("$$ server: options = %v\n", o)   // 打印结构体
	fmt.Printf("$$ server: args = %v\n", c.Args) // 自动解析shell传过来的位置参数到c.Args
	Result := []*Data{
		{Name: "Apple", Price: 6.0, Stokes: 128, Addition: []interface{}{1, "a", "c"}},
		{Name: "Banana", Price: 3.5, Stokes: 256, Addition: []interface{}{"b", 1.2}},
		{Name: "Pear", Price: 5, Stokes: 121, Sth: map[string]interface{}{"s": 123}},
	}
	content, _ := json.Marshal(Result)
	c.String(http.StatusOK, string(content)) // 发送数据给shell
}

var SName = "info" // shell客户端和服务端交互的unix套接字名称

func ShowTable() {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name:            "info",          // 命令名称
		Help:            "show info",     // 命令简短介绍
		Func:            Info,            // shell命令钩子
		Opts:            &InfOptions{},   // shell命令的具名参数
		ShowTable:       true,            // 是否开启表格显示功能
		KtrlHandler:     Handler,         // shell服务端视图函数
		SocketName:      SName,           // unix套接字名称
		ArgsMust:        true,            // 至少要传一个位置参数
		ArgsDescription: "info elements", // 位置参数功能描述
	})
	go kt.RunCtrl() // 开启服务端
	kt.RunShell()   // 开启shell客户端
}

func main() {
	ShowTable()
}
```
- 示例效果图
![shell-0](https://github.com/moqsien/goktrl/blob/main/docs/1.png)
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
