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

type KtrlContext struct {
	*ishell.Context
	Parser          *ParserPlus
	Options         interface{}
	Args            []string
	Table           *KtrlTable
	KtrlPath        string
	Client          *KtrlClient
	DefaultSocket   string
	ArgsCollectedAs string
}

func (that *KtrlContext) GetResult(sockName ...string) ([]byte, error) {
	sName := that.DefaultSocket
	if len(sockName) > 0 && len(sockName[0]) > 0 {
		sName = sockName[0]
	}
	params := that.Parser.Params
	if that.ArgsCollectedAs != "" {
		params[that.ArgsCollectedAs] = strings.Join(that.Args, ",")
	}
	return that.Client.GetResult(that.KtrlPath, params, sName)
}

func (that *KtrlShell) AddCmd(kcmd *KCommand) {
	that.Shell.AddCmd(&ishell.Cmd{
		Name:     strings.ReplaceAll(kcmd.Name, " ", ""),
		Help:     kcmd.Help,
		LongHelp: fmt.Sprintf("%s%s", kcmd.Help, kcmd.Opts.ShowHelpStr(kcmd.Opts)),
		Func: func(c *ishell.Context) {
			os.Args = c.Args
			kc := &KtrlContext{
				Client:          NewKtrlClient(),
				Context:         c,
				KtrlPath:        kcmd.GetKtrlPath(),
				DefaultSocket:   kcmd.SocketName,
				ArgsCollectedAs: kcmd.ArgsCollectedAs,
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
