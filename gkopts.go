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
	Required    = "required"
	Description = "descr"
)

type KtrlOpt interface{}

func ShowHelpStr(o KtrlOpt) (help string) {
	if o == nil {
		return ""
	}
	if reflect.ValueOf(o).Type().Kind() != reflect.Pointer {
		fmt.Println("[Opts] should be a pointer!")
		return ""
	}
	valType := reflect.ValueOf(o).Type().Elem()
	if valType.Kind() != reflect.Struct {
		fmt.Println("[Opts] should be a pointer of struct!")
		return ""
	}
	for i := 0; i < valType.NumField(); i++ {
		name := valType.Field(i).Name
		if name == "KtrlOption" {
			continue
		}
		tag := valType.Field(i).Tag
		if valType.Field(i).Type.Kind() == reflect.Bool {
			help += fmt.Sprintf("\n  -%s; alias:-{%s}; description: %s",
				strings.ToLower(name),
				tag.Get(Alias),
				tag.Get(Description))
		} else {
			help += fmt.Sprintf("\n  --%s=xxx; alias:-{%s}; description: %s",
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

type Options struct {
	OptList  g.MapStrBool // 是否解析参数的值
	Required g.MapStrBool // 是否必传
}

func ParseOptionsProperties(valType reflect.Type, k *KCommand) {
	if k.options == nil {
		optList := g.MapStrBool{}
		required := g.MapStrBool{}
		for i := 0; i < valType.NumField(); i++ {
			alias := valType.Field(i).Tag.Get(Alias)
			fName := strings.ToLower(valType.Field(i).Name)
			if alias == "" {
				alias = fName
			} else if !strings.Contains(alias, fName) {
				alias += "," + fName
			}
			np := valType.Field(i).Tag.Get(NeedParse)
			// 非布尔型参数默认需要解析其值，布尔型参数默认不需要解析其值(例如，-y出现在参数中即为ture，否则为false)
			if np == "" && valType.Field(i).Type.Kind() != reflect.Bool {
				np = "a"
			}
			optList[alias] = gconv.Bool(np)
			required[valType.Field(i).Name] = gconv.Bool(valType.Field(i).Tag.Get(Required))
		}
		k.options = &Options{
			Required: required,
			OptList:  optList,
		}
	}
}

func ParseShellOptions(o KtrlOpt, k *KCommand) (KtrlOpt, *ParserPlus) {
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
	if valType.Kind() != reflect.Struct {
		fmt.Println("[Opts] should be a pointer of struct!")
		return nil, nil
	}
	ParseOptionsProperties(valType, k)
	if k.options == nil {
		return nil, nil
	}

	parser, err := gcmd.Parse(k.options.OptList)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	params := parser.GetOptAll()
	for fieldName, isRequired := range k.options.Required {
		optName := strings.ToLower(fieldName)
		paramValue, found := params[optName]
		// 检查必传的具名参数
		if isRequired && !found {
			fmt.Printf("Option: [%v] is required!\n", optName)
			return nil, nil
		}
		SetStructValue(val.Elem().FieldByName(fieldName), paramValue, found)
		// Bool型参数在向服务端传递过程中不能为空
		if found && val.Elem().FieldByName(fieldName).Kind() == reflect.Bool {
			params[optName] = "true"
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
