package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moqsien/goktrl"
)

type Data struct {
	Name     string                 `order:"1"`
	Price    float32                `order:"2"`
	Stokes   int                    `order:"3"`
	Addition []interface{}          `order:"4"`
	Sth      map[string]interface{} `order:"5"`
}

func Info(k *goktrl.KtrlContext) {
	result, err := k.GetResult()
	if err != nil {
		fmt.Println(err)
		return
	}
	content := &[]*Data{}
	err = json.Unmarshal(result, content)
	k.Table.AddRowsByListObject(*content)
}

func Handler(c *gin.Context) {
	Result := []*Data{
		{Name: "Apple", Price: 6.0, Stokes: 128, Addition: []interface{}{1, "a", "c"}},
		{Name: "Banana", Price: 3.5, Stokes: 256, Addition: []interface{}{"b", 1.2}},
		{Name: "Pear", Price: 5, Stokes: 121, Sth: map[string]interface{}{"s": 123}},
	}
	content, _ := json.Marshal(Result)
	c.String(http.StatusOK, string(content))
}

var SName = "info"

func ShowTable() {
	kt := goktrl.NewKtrl()
	kt.AddKtrlCommand(&goktrl.KCommand{
		Name: "info",
		Help: "show info",
		Func: Info,
		Opts: goktrl.Opts{&goktrl.Option{
			Name:      "all,a",
			NeedParse: true,
			Must:      true,
		}},
		ShowTable:   true,
		KtrlHandler: Handler,
		SocketName:  SName,
	})
	go kt.RunCtrl()
	kt.RunShell()
}

func main() {
	ShowTable()
}
