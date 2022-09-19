package goktrl

import (
	"os"

	"github.com/abiosoft/ishell/v2"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
)

/*
  为shell增加参数解析功能
*/
type KtrlShell struct {
	*ishell.Shell
}

func NewShell() *KtrlShell {
	return &KtrlShell{
		Shell: ishell.New(),
	}
}

type KtrlContext struct {
	*ishell.Context
	Parser        *gcmd.Parser
	Args          []string
	Table         *KtrlTable
	KtrlPath      string
	Client        *KtrlClient
	DefaultSocket string
}

type KtrlCmd struct {
	Name          string
	Help          string
	Func          func(k *KtrlContext)
	Opts          *g.MapStrBool // 设置Options，支持别名，详见https://goframe.org/pages/viewpage.action?pageId=35357529
	KtrlPath      string
	ShowTable     bool
	DefaultSocket string
}

func (that *KtrlShell) AddCmd(cmd *KtrlCmd) {
	that.Shell.AddCmd(&ishell.Cmd{
		Name: cmd.Name,
		Help: cmd.Help,
		Func: func(c *ishell.Context) {
			os.Args = c.Args
			kc := &KtrlContext{
				Client:        NewKtrlClient(),
				Context:       c,
				KtrlPath:      cmd.KtrlPath,
				DefaultSocket: cmd.DefaultSocket,
			}
			kc.Parser, _ = gcmd.Parse(*cmd.Opts)
			kc.Args = kc.Parser.GetArgAll()
			if cmd.ShowTable {
				// 结果以table形式显示，table的数据在cmd.Func中获取
				kc.Table = NewKtrlTable()
			}
			kc.KtrlPath = cmd.KtrlPath
			cmd.Func(kc)
			if kc.Table != nil && cmd.ShowTable {
				// 打印table
				kc.Table.Render()
			}
		},
	})
}
