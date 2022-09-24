package main

import (
	"fmt"
	"net/http"

	"github.com/moqsien/goktrl"
)

var SockName string = "info"

type InfoOptions struct {
	*goktrl.KtrlOption
	All string `alias:"a" must:"true"`
}

func InfoShell(k *goktrl.Context) {
	all := k.Options.(*InfoOptions)
	fmt.Println("##client all: ", all)
	result, err := k.GetResult()
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("ResultString: ", result)
	k.Table.AddRowsByJsonString(string(result))
}

func InfoHandler(sc *goktrl.Context) {
	all := sc.Options.(*InfoOptions)
	fmt.Println("$$sever all: ", all)
	respStr := `{
			"headers": ["Name", "Price", "Stokes"],
			"rows": [
			  ["Apple","6", "128"],
			  ["Banana","3", "256"],
			  ["Pear","5", "121"]
			]
		  }`
	sc.String(http.StatusOK, respStr)
}

func KtrlTest() {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name:        "info",
		Help:        "show info",
		Func:        InfoShell,
		Opts:        &InfoOptions{},
		ShowTable:   true,
		KtrlHandler: InfoHandler,
		SocketName:  SockName,
	})
	go kt.RunCtrl()
	kt.RunShell()
}
