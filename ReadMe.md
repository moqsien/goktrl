### goktrl

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
go get -u "github.com/moqsien/goktrl@v1.2.0"
```
```go
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/moqsien/goktrl"
)

// 表格字段
type Data struct {
	Addition []interface{}          `order:"4"` // 可以按order标签排序，没有order标签，默认按字段名排序
	Name     string                 `order:"1"`
	Price    float32                `order:"2"`
	Stokes   int                    `order:"3"`
	Sth      map[string]interface{} `order:"5"`
}

/*
  定义参数。
  必须继承goktrl.KtrlOption基类。
  支持的标签有：
    alias: 设置参数别名；
	must: 是否必传；
	descr: 参数的描述；
	needparse: 参考goframe命令解析设置为true；默认值为true；
	字段支持Bool Int UInt Float String类型；
*/
type InfOptions struct {
	*goktrl.KtrlOption
	All  bool   `alias:"a" must:"true" descr:"show all info or not"`
	Info string `alias:"i" descr:"infomation"`
}

// 客户端钩子函数
func Info(k *goktrl.KtrlContext) {
	o := k.Options.(*InfOptions) // 可以自动解析参数
	fmt.Printf("## client: options=%v\n", o)
	fmt.Printf("## client: args=%v\n", k.Args)
	result, err := k.GetResult() // 向服务端发送请求：会请求到下面的Handler方法
	if err != nil {
		fmt.Println(err)
		return
	}
	content := &[]*Data{}
	err = json.Unmarshal(result, content)
	k.Table.AddRowsByListObject(*content) // 向表格中添加数据
}

func Handler(c *goktrl.ServerContext) {
	o := c.Options.(*InfOptions) // 可以自动解析参数
	fmt.Printf("$$ server: options = %v\n", o)
	fmt.Printf("$$ server: args = %v\n", c.Args) // 如果设置了ArgsCollectedAs，则可以获取到这些普通参数
	Result := []*Data{
		{Name: "Apple", Price: 6.0, Stokes: 128, Addition: []interface{}{1, "a", "c"}},
		{Name: "Banana", Price: 3.5, Stokes: 256, Addition: []interface{}{"b", 1.2}},
		{Name: "Pear", Price: 5, Stokes: 121, Sth: map[string]interface{}{"s": 123}},
	}
	content, _ := json.Marshal(Result)
	c.String(http.StatusOK, string(content)) // 向客户端发送内容
}

var SName = "info"

func ShowTable() {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name:            "info", // 命令名称
		Help:            "show info", // 命令简单描述
		Func:            Info,
		Opts:            &InfOptions{}, // 命令的命名的参数配置
		ShowTable:       true,
		KtrlHandler:     Handler,
		SocketName:      SName,
		ArgsCollectedAs: "in", // 如果设置了，则收集普通参数
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
