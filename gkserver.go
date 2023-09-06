package goktrl

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gfile"
	logger "github.com/moqsien/processes/logger"
)

// Unix套接字
type UnixSocket struct {
	UnixSocketName string // 套接字名称
	UnixSocketPath string // 套接字存放完整路径
}

type KtrlServer struct {
	UnixSocket
	Router *gin.Engine // UnixSockHttp 服务端
}

func NewKtrlServer() *KtrlServer {
	gin.SetMode(gin.ReleaseMode)
	return &KtrlServer{
		Router: gin.New(),
	}
}

// AddHandler 为KtrlServer添加视图函数
func (that *KtrlServer) AddHandler(kcmd *KCommand) {
	if kcmd.KtrlHandler != nil {
		that.Router.GET(kcmd.GetKtrlPath(), func(c *gin.Context) {
			options := ParseServerOptions(kcmd.Opts, c) // 解析Options
			kcmd.KtrlHandler(&Context{
				Context: c,
				Type:    ContextServer,
				Options: options,
				Args:    strings.Split(c.Query(fmt.Sprintf(ArgsFormatStr, kcmd.Name)), ","),
			})
		})
	}
}

func (that *KtrlServer) SetUnixSocket(sockName string) {
	if len(sockName) > 0 {
		if !strings.HasSuffix(sockName, ".sock") {
			sockName += ".sock"
		}
		that.UnixSocketName = sockName
		that.UnixSocketPath = GetSockFilePath(sockName)
	}
}

func (that *KtrlServer) CheckUnixSocket() {
	_, err := os.Stat(that.UnixSocketPath)
	if !os.IsNotExist(err) {
		// 判断socket文件是否存在，若已存在则删除
		_ = gfile.Remove(that.UnixSocketPath)
	}
}

func (that *KtrlServer) Start(sockName ...string) error {
	if len(sockName) > 0 {
		that.SetUnixSocket(sockName[0])
	}
	if that.UnixSocketPath == "" {
		err := errors.New("unix socket not initialized")
		logger.Error(err)
		return err
	}
	that.CheckUnixSocket()
	unixAddr, err := net.ResolveUnixAddr("unix", that.UnixSocketPath)
	if err != nil {
		return err
	}
	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		logger.Error("listening error:", err)
		return err
	}
	return http.Serve(listener, that.Router)
}
