package main

import (
	"fmt"

	"github.com/moqsien/goktrl"
)

func testOpts() {
	shell := goktrl.NewShell()
	shell.AddCmd(&goktrl.KtrlCmd{
		Name: "info",
		Help: "show info",
		Opts: goktrl.Opts{
			&goktrl.Option{
				Name:      "test,t",
				NeedParse: true,
			},
			&goktrl.Option{
				Name: "test2,t2",
			},
		},
		Func: func(k *goktrl.KtrlContext) {
			fmt.Println("args: ", k.Args)
			fmt.Println("t: ", k.Parser.GetOpt("t"))
			fmt.Println("test: ", k.Parser.GetOpt("test"))
			fmt.Println("t2: ", k.Parser.GetOpt("t2"))
			fmt.Println("test2: ", k.Parser.GetOpt("test2"))
		},
	})
	shell.Run()
}

func testShowTable() {
	shell := goktrl.NewShell()
	shell.AddCmd(&goktrl.KtrlCmd{
		Name: "table",
		Help: "show table",
		Opts: goktrl.Opts{
			&goktrl.Option{
				Name:      "table,t",
				NeedParse: true,
			},
		},
		ShowTable: true,
		Func: func(k *goktrl.KtrlContext) {
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
