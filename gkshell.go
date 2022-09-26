package goktrl

import (
	"encoding/json"
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
	if kcmd.ArgsDescription == "" {
		kcmd.ArgsDescription = "no descriptions for args."
	}
	that.Shell.AddCmd(&ishell.Cmd{
		Name:     strings.ReplaceAll(kcmd.Name, " ", ""),
		Help:     kcmd.Help,
		LongHelp: fmt.Sprintf("%s%s\n args: \n  %s", kcmd.Help, ShowHelpStr(kcmd.Opts), kcmd.ArgsDescription),
		Func: func(c *ishell.Context) {
			os.Args = c.Args
			kc := &Context{
				Type:          ContextClient,
				ShellContext:  c,
				KtrlPath:      kcmd.GetKtrlPath(),
				DefaultSocket: kcmd.SocketName,
				ShellCmdName:  kcmd.Name,
			}
			kc.Options, kc.Parser = ParseShellOptions(kcmd.Opts, kcmd)
			if kc.Parser == nil {
				return
			}
			kc.Args = kc.Parser.GetArgAll()
			if kcmd.ArgsRequired && len(kc.Args) == 0 {
				fmt.Println("At least one argument must be provided!")
				return
			}
			if kcmd.ShowTable {
				// 结果以table形式显示，table的数据在cmd.Func中获取
				kc.Table = NewKtrlTable()
			}
			// 全自动获取result并显示
			if kcmd.Auto {
				var err error
				kc.Result, err = kc.GetResult()
				if err != nil {
					fmt.Println(err)
					return
				}
				if !kcmd.ShowTable {
					// 普通结果显示
					fmt.Println(string(kc.Result))
				} else if kcmd.TableObject != nil {
					// 自动显示表格
					err = json.Unmarshal(kc.Result, kcmd.TableObject)
					kc.Table.AddRowsByListObject(kcmd.TableObject)
				} else {
					fmt.Println("Table object is required!")
				}
			}
			if kcmd.Func != nil {
				kcmd.Func(kc)
			}
			if kc.Table != nil && kcmd.ShowTable {
				// 渲染表格
				kc.Table.Render()
			}
		},
	})
}
