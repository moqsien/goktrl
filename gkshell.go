package goktrl

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/abiosoft/ishell/v2"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/moqsien/processes/logger"
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

func (that *KtrlContext) GetResult(sockName ...string) ([]byte, error) {
	sName := that.DefaultSocket
	if len(sockName) > 0 && len(sockName[0]) > 0 {
		sName = sockName[0]
	}
	return that.Client.GetResult(that.KtrlPath, that.Parser.GetOptAll(), sName)
}

type Option struct {
	Name      string // 参数名称和别名，英文逗号分隔，无空格
	NeedParse bool   // 是否需要解析值
	Must      bool   // 是否必传
}

type Opts []*Option

type KtrlCmd struct {
	Name          string
	Help          string
	Func          func(k *KtrlContext)
	Opts          Opts
	KtrlPath      string
	ShowTable     bool
	DefaultSocket string
}

func (that *KtrlShell) ParseAndCheckOpts(options Opts) (parser *gcmd.Parser, err error) {
	result := g.MapStrBool{}
	must := []string{}
	for _, o := range options {
		result[o.Name] = o.NeedParse
		if o.Must {
			must = append(must, strings.Split(o.Name, ",")[0])
		}
	}
	parser, err = gcmd.Parse(result)
	// 检查必传参数
	for _, m := range must {
		if len(parser.GetOpt(m)) == 0 {
			return nil, errors.New(fmt.Sprintf("Option:<%s> must present!", m))
		}
	}
	return
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
			var err error
			// 设置Options，支持别名，详见https://goframe.org/pages/viewpage.action?pageId=35357529
			kc.Parser, err = that.ParseAndCheckOpts(cmd.Opts)
			if err != nil {
				logger.Info(err)
				return
			}
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
