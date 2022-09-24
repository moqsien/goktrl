package simple

import (
	"fmt"

	"github.com/moqsien/goktrl"
)

/*
  下面是一个关于goktrl的简短示例。
  一个最最简单的例子: 只要求至少有一个位置参数。
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
		ArgsRequired:    true,            // 至少要传一个位置参数
		ArgsDescription: "info elements", // 位置参数功能描述
		Auto:            true,            // 是否全自动处理和显示数据
	})
	return kt
}
