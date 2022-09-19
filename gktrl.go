package goktrl

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/g"
)

type Ktrl struct {
	CtrlServer *KtrlServer
	CtrlShell  *KtrlShell
}

func NewKtrl() *Ktrl {
	return &Ktrl{
		CtrlServer: NewKtrlServer(), // 服务端
		CtrlShell:  NewShell(),      // 交互式shell
	}
}

type KCommand struct {
	Name        string               // shell 命令名称
	Help        string               // shell 命令解释
	Func        func(k *KtrlContext) // shell 命令钩子函数
	Opts        *g.MapStrBool        // shell 命令可选参数配置
	KtrlPath    string               // 路由
	ShowTable   bool                 // 结果是否在命令行中以表格显示
	KtrlHandler func(c *gin.Context) // Ktrl服务端视图函数
	SocketName  string               // 默认Unix套接字名称
}

func (that *Ktrl) AddKtrlCommand(kcmd *KCommand) {
	that.CtrlShell.AddCmd(&KtrlCmd{
		Name:          kcmd.Name,
		Help:          kcmd.Help,
		Opts:          kcmd.Opts,
		KtrlPath:      kcmd.KtrlPath,
		ShowTable:     kcmd.ShowTable,
		Func:          kcmd.Func,
		DefaultSocket: kcmd.SocketName,
	})

	that.CtrlServer.AddHandler(kcmd.KtrlPath, kcmd.KtrlHandler)
	if kcmd.SocketName != "" {
		that.CtrlServer.SetUnixSocket(kcmd.SocketName)
	}
}

// RunShell 运行Ktrl交互式shell
func (that *Ktrl) RunShell(sockName ...string) {
	that.CtrlShell.Run()
}

// RunCtrl 运行Ktrl服务端
func (that *Ktrl) RunCtrl(sockName ...string) {
	that.CtrlServer.Start(sockName...)
}