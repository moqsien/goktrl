package goktrl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/util/gconv"
)

const (
	Alias       = "alias"
	NeedParse   = "needparse"
	Must        = "required"
	Description = "descr"
)

type KtrlOpt interface {
	IsKtrlOpt() bool
}

// KtrlOption 命令行参数配置基类
type KtrlOption struct{}

func (that *KtrlOption) IsKtrlOpt() bool {
	return true
}

func ShowHelpStr(o KtrlOpt) (help string) {
	if o == nil {
		return ""
	}
	if reflect.ValueOf(o).Type().Kind() != reflect.Pointer {
		fmt.Println("[Opts] should be a pointer!")
		return ""
	}
	valType := reflect.ValueOf(o).Type().Elem()
	for i := 0; i < valType.NumField(); i++ {
		name := valType.Field(i).Name
		if name == "KtrlOption" {
			continue
		}
		tag := valType.Field(i).Tag
		if valType.Field(i).Type.Kind() == reflect.Bool {
			help += fmt.Sprintf("\n  -%s; alias:{%s}; description: %s",
				strings.ToLower(name),
				tag.Get(Alias),
				tag.Get(Description))
		} else {
			help += fmt.Sprintf("\n  --%s=xxx; alias:{%s}; description: %s",
				strings.ToLower(name),
				tag.Get(Alias),
				tag.Get(Description))
		}
	}
	if len(help) > 0 {
		help = "\n options: " + help
	}
	return
}

func SetStructValue(field reflect.Value, fValue string, found bool) {
	if field.Type().Kind() == reflect.String {
		field.SetString(fValue)
	} else if field.Type().Kind() == reflect.Bool {
		field.SetBool(gconv.Bool(found))
	} else if field.CanInt() {
		field.SetInt(gconv.Int64(fValue))
	} else if field.CanUint() {
		field.SetUint(gconv.Uint64(fValue))
	} else if field.CanFloat() {
		field.SetFloat(gconv.Float64(fValue))
	}
}

type ParserPlus struct {
	*gcmd.Parser
	Params map[string]string
}

func ParseShellOptions(o KtrlOpt) (KtrlOpt, *ParserPlus) {
	if o == nil {
		parser, err := gcmd.Parse(g.MapStrBool{})
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		return nil, &ParserPlus{
			Parser: parser,
			Params: map[string]string{},
		}
	}
	val := reflect.ValueOf(o)
	if val.Type().Kind() != reflect.Pointer {
		fmt.Println("[Opts] should be a pointer!")
		return nil, nil
	}
	valType := val.Type().Elem()
	settings := g.MapStrBool{}
	for i := 0; i < valType.NumField(); i++ {
		alias := valType.Field(i).Tag.Get(Alias)
		fName := strings.ToLower(valType.Field(i).Name)
		if alias == "" {
			alias = fName
		} else if !strings.Contains(alias, fName) {
			alias += "," + fName
		}
		np := valType.Field(i).Tag.Get(NeedParse)
		if np == "" && valType.Field(i).Type.Kind() != reflect.Bool {
			np = "a"
		}
		settings[alias] = gconv.Bool(np)
	}

	parser, err := gcmd.Parse(settings)
	params := parser.GetOptAll()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	for i := 0; i < valType.NumField(); i++ {
		fName := valType.Field(i).Name
		option := strings.ToLower(fName)
		paramValue, found := parser.GetOptAll()[option]
		// 检查必传的具名参数
		if gconv.Bool(valType.Field(i).Tag.Get(Must)) && !found {
			fmt.Printf("Option: [%v] is required!\n", option)
			return nil, nil
		}
		SetStructValue(val.Elem().FieldByName(fName), paramValue, found)
		// Bool型参数在向服务端传递过程中不能为空
		if _, ok := params[option]; valType.Field(i).Type.Kind() == reflect.Bool && ok {
			params[option] = "true"
		}
	}
	return o, &ParserPlus{
		Parser: parser,
		Params: params,
	}
}

func ParseServerOptions(o KtrlOpt, c *gin.Context) KtrlOpt {
	if o == nil {
		return nil
	}
	val := reflect.ValueOf(o)
	if val.Type().Kind() != reflect.Pointer {
		fmt.Println("[Opts] should be a pointer!")
		return nil
	}
	valType := val.Type().Elem()
	for i := 0; i < valType.NumField(); i++ {
		fName := valType.Field(i).Name
		paramValue := c.Query(strings.ToLower(fName))
		SetStructValue(val.Elem().FieldByName(fName), paramValue, gconv.Bool(paramValue))
	}
	return o
}
