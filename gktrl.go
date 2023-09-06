package goktrl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogf/gf/os/gfile"
)

type Ktrl struct {
	CtrlServer *KtrlServer // 服务端
	CtrlShell  *KtrlShell  // 客户端交互式shell
	Multiple   bool        // 是否客户端和服务端在不同进程中执行
}

func NewKtrl(ismulti ...bool) *Ktrl {
	multi := true // 默认是分别在不同的进程中运行客户端和服务端
	if len(ismulti) > 0 {
		multi = ismulti[0]
	}
	return &Ktrl{
		CtrlServer: NewKtrlServer(),
		CtrlShell:  NewShell(),
		Multiple:   multi,
	}
}

type KCommand struct {
	Name            string                  // shell 命令名称
	Help            string                  // shell 命令解释
	Func            func(c *Context)        // shell 命令钩子函数
	Opts            KtrlOpt                 // shell 命令可选参数配置
	KtrlHandler     func(c *Context)        // Ktrl服务端视图函数
	SocketName      string                  // 默认Unix套接字名称
	ArgsDescription string                  // 位置参数说明
	ArgsRequired    bool                    // 位置参数是否至少要传一个
	Auto            bool                    // 是否自动发送请求并处理结果
	TableObject     interface{}             // 空的表格对象
	ShowTable       bool                    // 结果是否在命令行中以表格显示
	options         *Options                // 缓存参数结构体reflect结果
	ArgsHook        func([]string) []string // 在shell中处理Args, 然后向Server发送处理之后的Args
}

func (that *KCommand) GetKtrlPath() string {
	return fmt.Sprintf("/ktrl/%s", that.Name)
}

func (that *Ktrl) AddKtrlCommand(kcmd *KCommand) {
	if kcmd.SocketName == "" {
		kcmd.SocketName = "ktrlDefault"
	}
	that.CtrlShell.AddCmd(kcmd)

	that.CtrlServer.AddHandler(kcmd)
	if kcmd.SocketName != "" && that.CtrlServer.UnixSocket.UnixSocketName == "" {
		// 服务端Unix套接字设置一次就好了
		that.CtrlServer.SetUnixSocket(kcmd.SocketName)
	}
}

// RunShell 运行Ktrl交互式shell
func (that *Ktrl) RunShell(sockName ...string) {
	if that.Multiple {
		that.CtrlServer = nil // 回收Server
	}
	that.CtrlShell.Run()
}

// RunCtrl 运行Ktrl服务端
func (that *Ktrl) RunCtrl(sockName ...string) {
	if that.Multiple {
		that.CtrlShell = nil // 回收Client
	}
	that.CtrlServer.Start(sockName...)
}

const (
	GoKtrlSockDirEnv string = "GOKTRL_SOCK_DIR"
)

func GetSockFilePath(sockName string) (p string) {
	sockDir := os.Getenv(GoKtrlSockDirEnv)
	if sockDir == "" {
		p = gfile.TempDir(sockName)
	} else {
		p = filepath.Join(sockDir, sockName)
	}
	return
}
