package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/g"
	"github.com/moqsien/goktrl"
)

var Sock = "test"

func RunShell(kt *goktrl.Ktrl) {
	kt.CtrlShell.AddCmd(&goktrl.KtrlCmd{
		Name: "info",
		Help: "show info",
		Opts: &g.MapStrBool{
			"all,a": true,
		},
		KtrlPath:  "/ctrl/info",
		ShowTable: true,
		Func: func(k *goktrl.KtrlContext) {
			result, err := k.Client.GetResult(k.KtrlPath, k.Parser.GetOptAll(), Sock)
			if err != nil {
				fmt.Println(err)
				return
			}
			// fmt.Println("ResultString: ", result)
			k.Table.AddRowsByJsonString(result)
		},
	})
	kt.CtrlShell.Run()
}

func RunServer(kt *goktrl.Ktrl) {
	kt.CtrlServer.AddHandler("/ctrl/info", func(c *gin.Context) {
		fmt.Println("===info===")
		respStr := `{
			"headers": ["Name", "Price", "Stokes"],
			"rows": [
			  ["Apple","6", "128"],
			  ["Banana","3", "256"],
			  ["Pear","5", "121"]
			]
		  }`
		c.String(http.StatusOK, respStr)
	})
	kt.CtrlServer.Start(Sock)
}
