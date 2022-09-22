package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moqsien/goktrl"
)

var SockName string = "info"

func InfoShell(k *goktrl.KtrlContext) {
	result, err := k.Client.GetResult(k.KtrlPath,
		k.Parser.GetOptAll(),
		k.DefaultSocket)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("ResultString: ", result)
	k.Table.AddRowsByJsonString(string(result))
}

func InfoHandler(c *gin.Context) {
	fmt.Println("===", c.Query("all"))
	respStr := `{
			"headers": ["Name", "Price", "Stokes"],
			"rows": [
			  ["Apple","6", "128"],
			  ["Banana","3", "256"],
			  ["Pear","5", "121"]
			]
		  }`
	c.String(http.StatusOK, respStr)
}

func KtrlTest() {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name: "info",
		Help: "show info",
		Func: InfoShell,
		Opts: goktrl.Opts{&goktrl.Option{
			Name:      "all,a",
			NeedParse: true,
			Must:      true,
		}},
		ShowTable:   true,
		KtrlHandler: InfoHandler,
		SocketName:  SockName,
	})
	go kt.RunCtrl()
	kt.RunShell()
}
