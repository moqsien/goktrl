package dispatch

import (
	"fmt"

	"github.com/moqsien/goktrl"
)

/*
  下面是一个关于goktrl的简短示例。
  自动处理数据。
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
	required: 是否为必传具名参数；
	descr: 具名参数描述信息；
	needparse: 一般不需要用户设置，已根据结构体字段类型进行自动处理；

  支持的字段类型有: string, bool, int, uint, float
*/
type InfOptions struct {
	All  bool   `alias:"a" required:"true" descr:"show all info or not"`
	Info string `alias:"i" descr:"infomation"`
}

func Info(c *goktrl.Context) {
	o := c.Options.(*InfOptions)               // 自动解析参数到结构体
	fmt.Printf("## client: options=%v\n", o)   // 打印结构体
	fmt.Printf("## client: args=%v\n", c.Args) // 自动收集命令行普通的位置参数
}

var (
	DefaultSock   = "info"
	Sock1         = "info1"
	IsServerSock1 = false
)

func Handler(c *goktrl.Context) {
	if !IsServerSock1 {
		fmt.Println("===dispatching request from client!")
		result, _ := c.GetResult(Sock1) // 转发请求到Sock1
		fmt.Println("===dispatching result: ", string(result))
		c.Send(result)
	} else {
		o := c.Options.(*InfOptions)                 // 自动解析参数到结构体
		fmt.Printf("$$ server: options = %v\n", o)   // 打印结构体
		fmt.Printf("$$ server: args = %v\n", c.Args) // 自动解析shell传过来的位置参数到c.Args
		Result := []*Data{
			{Name: "Apple", Price: 6.0, Stokes: 128, Addition: []interface{}{1, "a", "c"}},
			{Name: "Banana", Price: 3.5, Stokes: 256, Addition: []interface{}{"b", 1.2}},
			{Name: "Pear", Price: 5, Stokes: 121, Sth: map[string]interface{}{"s": 123}},
		}
		c.Send(Result)
	}
}

func ShowTable(sockName string) *goktrl.Ktrl {
	if sockName == "" {
		sockName = DefaultSock
	}
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name:            "info",          // 命令名称
		Help:            "show info",     // 命令简短介绍
		Func:            Info,            // shell命令钩子
		Opts:            &InfOptions{},   // shell命令的具名参数
		ShowTable:       true,            // 是否开启表格显示功能
		KtrlHandler:     Handler,         // shell服务端视图函数
		SocketName:      sockName,        // unix套接字名称
		ArgsRequired:    true,            // 至少要传一个位置参数
		ArgsDescription: "info elements", // 位置参数功能描述
		Auto:            true,            // 是否全自动处理数据
		TableObject:     &[]*Data{},      // table object，用于表格自动加载数据
	})
	return kt
}

func RunS(sockName string) {
	if sockName != DefaultSock && sockName != Sock1 {
		fmt.Println("Sock: ", sockName, "is not surpported!")
		return
	}
	if sockName == Sock1 {
		IsServerSock1 = true
	}
	kt := ShowTable(sockName)
	kt.RunCtrl()
}

func RunC() {
	kt := ShowTable(DefaultSock)
	kt.RunShell()
}
