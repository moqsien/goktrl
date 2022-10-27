### Introduction to [goktrl](https://github.com/moqsien/goktrl) 

------------------
[Zh_CN](https://github.com/moqsien/goktrl/blob/main/docs/ReadMeZh.md)

goktrl is a powerful interactive shell designed to probe into your multi-processing go projects by using unix domain sockets.

### Characteristics
------------------
- interactive shell
- nice hints for args and options
- use struct tags for configuration of options
- automatically parse and verify args and options
- automatically render table for terminal if enabled
- automatically handle and print data if enabled(including table data)
- connect to your process using unix domain sockets
- dispatching requests easily
- very intuitive and flexible

### Usage
------------------
```shell
go get -u "github.com/moqsien/goktrl@v1.3.6"
```
- [More examples](https://github.com/moqsien/goktrl/tree/main/examples/ktrl)
- A simple One:
```go
package main

import (
	"fmt"

	"github.com/moqsien/goktrl"
)

/*
  A very simple example by implementing goktrl.
  Even a hook for the shell is not required.
*/

func Handler(c *goktrl.Context) {
	fmt.Printf("$$ server: args = %v\n", c.Args) // args are automatically parsed and stored in c.Args
	Result := map[string]string{
		"hello": "info",
	}
	c.Send(Result)
}

func ShowInfo() *goktrl.Ktrl {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name:            "info",          // name of your shell command
		Help:            "show info",     // help info for your shell command
		KtrlHandler:     Handler,         // view controller for server side
		Auto:            true,            // automatically show results or not
	})
	return kt
}

func main() {
	kt := ShowInfo()
	if len(os.Args) > 1 {
		kt.RunShell()
	} else {
		kt.RunCtrl()
	}
}
```

- Some Exibitions
![shell-0](https://github.com/moqsien/goktrl/blob/main/docs/0.png)
![shell-1](https://github.com/moqsien/goktrl/blob/main/docs/1.png)
![shell-2](https://github.com/moqsien/goktrl/blob/main/docs/2.png)

### License
[MIT](https://github.com/moqsien/goktrl/blob/main/LICENSE)

### Thanks To

------------------
[dmicro](https://github.com/osgochina/dmicro)

[goframe](https://github.com/gogf/gf)

[gin](https://github.com/gin-gonic/gin)

[ishell](https://github.com/abiosoft/ishell)

[table](https://github.com/aquasecurity/table)
