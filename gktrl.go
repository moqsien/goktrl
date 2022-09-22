package goktrl

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
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
	Opts        Opts                 // shell 命令可选参数配置
	ShowTable   bool                 // 结果是否在命令行中以表格显示
	KtrlHandler func(c *gin.Context) // Ktrl服务端视图函数
	SocketName  string               // 默认Unix套接字名称
}

func (that *KCommand) GetKtrlPath() string {
	return fmt.Sprintf("/ktrl/%s", that.Name)
}

func (that *Ktrl) AddKtrlCommand(kcmd *KCommand) {
	that.CtrlShell.AddCmd(&KtrlCmd{
		Name:          strings.ReplaceAll(kcmd.Name, " ", ""),
		Help:          kcmd.Help,
		Opts:          kcmd.Opts,
		KtrlPath:      kcmd.GetKtrlPath(),
		ShowTable:     kcmd.ShowTable,
		Func:          kcmd.Func,
		DefaultSocket: kcmd.SocketName,
	})

	that.CtrlServer.AddHandler(kcmd.GetKtrlPath(), kcmd.KtrlHandler)
	if kcmd.SocketName != "" && that.CtrlServer.UnixSocket.UnixSocketName == "" {
		// 服务端Unix套接字设置一次就好了
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
