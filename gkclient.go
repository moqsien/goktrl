package goktrl

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/gogf/gf/os/gfile"
)

type KtrlClient struct {
	UnixSocket
	Client *http.Client
	params string
}

func NewKtrlClient() *KtrlClient {
	return &KtrlClient{
		Client: &http.Client{},
	}
}

// SetConnition 设置连接
func (that *KtrlClient) SetConnection() {
	that.Client.Transport = &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", that.UnixSocketPath)
		},
	}
}

func (that *KtrlClient) SetUnixSocket(sockName string) {
	if len(sockName) > 0 {
		if !strings.HasSuffix(sockName, ".sock") {
			sockName += ".sock"
		}
		that.UnixSocketName = sockName
		that.UnixSocketPath = gfile.TempDir(sockName)
		that.SetConnection() // 连接server
	}
}

func (that *KtrlClient) ParseParams(params map[string]string) {
	that.params = ""
	for k, v := range params {
		if len(that.params) == 0 {
			that.params += fmt.Sprintf("?%s=%s", k, v)
		} else {
			that.params += fmt.Sprintf("&%s=%s", k, v)
		}
	}
}

func (that *KtrlClient) GetResult(urlPath string, params map[string]string, sockName ...string) ([]byte, error) {
	if len(sockName) > 0 {
		that.SetUnixSocket(sockName[0])
	}
	// 解析参数
	that.ParseParams(params)
	// 生成http的url
	url := fmt.Sprintf("http://%s/%s/%s",
		that.UnixSocketName, strings.Trim(urlPath, "/"),
		that.params)

	resp, err := that.Client.Get(url)
	if err != nil {
		return nil, err
	}
	if result, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
