package main

import (
	"fmt"

	"github.com/moqsien/goktrl"
)

type InfoOptions struct {
	*goktrl.KtrlOption
	Test  string `alias:"t"`
	Test2 string `alias:"t2" needparse:"false"`
}

func testOpts() {
	shell := goktrl.NewShell()
	shell.AddCmd(&goktrl.KCommand{
		Name: "info",
		Help: "show info",
		Opts: &InfoOptions{},
		Func: func(k *goktrl.Context) {
			fmt.Println("args: ", k.Args)
			fmt.Println("t: ", k.Parser.GetOpt("t"))
			fmt.Println("test: ", k.Parser.GetOpt("test"))
			fmt.Println("t2: ", k.Parser.GetOpt("t2"))
			fmt.Println("test2: ", k.Parser.GetOpt("test2"))
		},
	})
	shell.Run()
}

type TableOptions struct {
	*goktrl.KtrlOption
	Table string `alias:"t"`
}

func testShowTable() {
	shell := goktrl.NewShell()
	shell.AddCmd(&goktrl.KCommand{
		Name:      "table",
		Help:      "show table",
		Opts:      &TableOptions{},
		ShowTable: true,
		Func: func(k *goktrl.Context) {
			//命令： table -t abc
			if table := k.Parser.GetOpt("table"); len(table) > 0 {
				k.Table.AddRowsByJsonString(`{
					"headers": ["Name", "Price", "Stokes"],
					"rows": [
					  ["Apple","6", "128"],
					  ["Banana","3", "256"],
					  ["Pear","5", "121"]
					]
				  }`)
			}
		},
	})
	shell.Run()
}

func main() {
	testOpts()
	// testShowTable()
}
