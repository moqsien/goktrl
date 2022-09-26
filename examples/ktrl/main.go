package main

import (
	"os"

	"github.com/moqsien/goktrl/examples/ktrl/autor"
	"github.com/moqsien/goktrl/examples/ktrl/dispatch"
	"github.com/moqsien/goktrl/examples/ktrl/manual"
	"github.com/moqsien/goktrl/examples/ktrl/simple"
	"github.com/moqsien/goktrl/examples/ktrl/single"
)

// 自动处理数据
func runAuto() {
	kt := autor.ShowTable()
	// 单进程运行会报错
	// go kt.RunCtrl()
	// kt.RunShell()
	if len(os.Args) > 1 {
		kt.RunShell()
	} else {
		kt.RunCtrl()
	}
}

// 手动处理数据
func runManual() {
	kt := manual.ShowTable()
	if len(os.Args) > 1 {
		kt.RunShell()
	} else {
		kt.RunCtrl()
	}
}

func runSimple() {
	kt := simple.ShowInfo()
	if len(os.Args) > 1 {
		kt.RunShell()
	} else {
		kt.RunCtrl()
	}
}

func runSingle() {
	// 单进程运行
	kt := single.ShowTable()
	go kt.RunCtrl()
	kt.RunShell()
}

func runDispatching() {
	if len(os.Args) > 2 {
		dispatch.RunC()
	} else if len(os.Args) > 1 {
		dispatch.RunS(dispatch.DefaultSock)
	} else {
		dispatch.RunS(dispatch.Sock1)
	}
}

func main() {
	// runAuto()
	// runManual()
	// runSimple()
	// runSingle()
	runDispatching()
}
