package goktrl

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	Result        []byte
}

func (that *Context) GetResult(sockName ...string) ([]byte, error) {
	sName := that.DefaultSocket
	if len(sockName) > 0 && len(sockName[0]) > 0 {
		sName = sockName[0]
	}
	params := make(map[string]string)
	if that.Type == ContextClient {
		params = that.Parser.Params
		params[fmt.Sprintf(ArgsFormatStr, that.ShellCmdName)] = strings.Join(that.Args, ",")
	} else {
		// 服务端继续转发请求
		for k := range that.Request.URL.Query() {
			params[k] = that.Query(k)
		}
		that.KtrlPath = that.Request.URL.Path
	}
	if that.Client == nil {
		that.Client = NewKtrlClient()
	}
	return that.Client.GetResult(that.KtrlPath, params, sName)
}

func (that *Context) Send(content interface{}, code ...int) {
	if that.Context != nil {
		statusCode := http.StatusOK
		if len(code) > 0 {
			statusCode = code[0]
		}
		switch r := content.(type) {
		case string:
			that.Context.String(statusCode, r)
		case []byte:
			that.Context.String(statusCode, string(r))
		default:
			res, err := json.Marshal(content)
			if err != nil {
				fmt.Println(err)
				that.Context.String(http.StatusInternalServerError, err.Error())
				return
			}
			that.Context.String(statusCode, string(res))
		}
	}
}
