package manual

import (
	"encoding/json"
	"fmt"

	"github.com/moqsien/goktrl"
)

/*
  下面是一个关于goktrl的简短示例。
  手动处理数据。
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
	result, err := c.GetResult()               // 自动根据参数向服务端发送请求，请求会到达下面的Handler路由方法
	if err != nil {
		fmt.Println(err)
		return
	}
	var content interface{} = &[]*Data{}
	err = json.Unmarshal(result, content)
	c.Table.AddRowsByListObject(content) // 如果ShowTable设置为true，此处可添加表格数据，会自动渲染和显示表格
}

func Handler(c *goktrl.Context) {
	o := c.Options.(*InfOptions)                 // 自动解析参数到结构体
	fmt.Printf("$$ server: options = %v\n", o)   // 打印结构体
	fmt.Printf("$$ server: args = %v\n", c.Args) // 自动解析shell传过来的位置参数到c.Args
	var Result interface{} = []*Data{
		{Name: "Apple", Price: 6.0, Stokes: 128, Addition: []interface{}{1, "a", "c"}},
		{Name: "Banana", Price: 3.5, Stokes: 256, Addition: []interface{}{"b", 1.2}},
		{Name: "Pear", Price: 5, Stokes: 121, Sth: map[string]interface{}{"s": 123}},
	}
	content, _ := json.Marshal(&Result)
	c.Send(content) // 发送数据给shell
}

var SName = "info" // shell客户端和服务端交互的unix套接字名称

func ShowTable() *goktrl.Ktrl {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name:            "info",          // 命令名称
		Help:            "show info",     // 命令简短介绍
		Func:            Info,            // shell命令钩子
		Opts:            &InfOptions{},   // shell命令的具名参数
		ShowTable:       true,            // 是否开启表格显示功能
		KtrlHandler:     Handler,         // shell服务端视图函数
		SocketName:      SName,           // unix套接字名称
		ArgsRequired:    true,            // 至少要传一个位置参数
		ArgsDescription: "info elements", // 位置参数功能描述
	})
	return kt
}
