package goktrl

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell/v2"
	"github.com/gin-gonic/gin"
)

type ContextType int

const (
	ContextClient ContextType = 1 // 客户端
	ContextServer ContextType = 2 // 服务端
	ArgsFormatStr string      = "args%sargs"
)

type Context struct {
	*gin.Context
	ShellContext  *ishell.Context
	Type          ContextType
	Options       KtrlOpt
	Args          []string
	Parser        *ParserPlus
	Table         *KtrlTable
	KtrlPath      string
	Client        *KtrlClient
	DefaultSocket string
	ShellCmdName  string
}

func (that *Context) GetResult(sockName ...string) ([]byte, error) {
	sName := that.DefaultSocket
	if len(sockName) > 0 && len(sockName[0]) > 0 {
		sName = sockName[0]
	}
	params := that.Parser.Params
	params[fmt.Sprintf(ArgsFormatStr, that.ShellCmdName)] = strings.Join(that.Args, ",")
	return that.Client.GetResult(that.KtrlPath, params, sName)
}
