package goktrl

import (
	"fmt"
	"os"
	"strings"

	"github.com/abiosoft/ishell/v2"
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

func (that *KtrlShell) AddCmd(kcmd *KCommand) {
	that.Shell.AddCmd(&ishell.Cmd{
		Name:     strings.ReplaceAll(kcmd.Name, " ", ""),
		Help:     kcmd.Help,
		LongHelp: fmt.Sprintf("%s%s", kcmd.Help, kcmd.Opts.ShowHelpStr(kcmd.Opts)),
		Func: func(c *ishell.Context) {
			os.Args = c.Args
			kc := &Context{
				Type:          ContextClient,
				Client:        NewKtrlClient(),
				ShellContext:  c,
				KtrlPath:      kcmd.GetKtrlPath(),
				DefaultSocket: kcmd.SocketName,
				ShellCmdName:  kcmd.Name,
			}
			kc.Options, kc.Parser = kcmd.Opts.ParseShellOptions(kcmd.Opts)
			if kc.Parser == nil {
				return
			}
			kc.Args = kc.Parser.GetArgAll()
			if kcmd.ShowTable {
				// 结果以table形式显示，table的数据在cmd.Func中获取
				kc.Table = NewKtrlTable()
			}
			kcmd.Func(kc)
			if kc.Table != nil && kcmd.ShowTable {
				// 打印table
				kc.Table.Render()
			}
		},
	})
}
