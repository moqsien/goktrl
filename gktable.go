package goktrl

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/aquasecurity/table"
	"github.com/gogf/gf/container/gtree"
	"github.com/gogf/gf/util/gutil"
	"github.com/tidwall/gjson"
)

/*
  命令行打印表格
*/

type KtrlTable struct {
	*table.Table
	Headers  []string
	RowCache []string
}

func NewKtrlTable() *KtrlTable {
	return &KtrlTable{
		Table:   table.New(os.Stdout),
		Headers: []string{},
	}
}

// RestTable 重置KtrlTable
func (that *KtrlTable) RestTable() {
	that.Table = table.New(os.Stdout)
	that.Headers = that.Headers[:0] // 无需新建
}

const (
	TableFieldOrderFlag = "order"
)

func (that *KtrlTable) ParseHeadersFromObject(vType reflect.Type) {
	Headers := gtree.NewAVLTree(gutil.ComparatorString, true)
	for j := 0; j < vType.NumField(); j++ {
		order := vType.Field(j).Tag.Get(TableFieldOrderFlag)
		name := vType.Field(j).Name
		if order == "" {
			order = name
		}
		Headers.Set(order, name)
	}
	// Headers按order标签进行排序
	for _, s := range Headers.Values() {
		v, _ := s.(string)
		that.Headers = append(that.Headers, v)
	}
	if len(that.Headers) > 0 {
		that.AddHeaders(that.Headers...) // 添加表头
	}
}

/*
  AddRowsByListObject 为表格添加数据行；
  dataList 数据类型格式: []Data 或 []*Data；
  Data为struct, 结构示例如下：
    type Data struct {
		FieldOne string `order:"1"`
		FieldTwo string `order:"2"`
	}
*/
func (that *KtrlTable) AddRowsByListObject(dataList interface{}) {
	dList := reflect.ValueOf(dataList)
	if dList.Kind() == reflect.Pointer {
		dList = dList.Elem()
	}
	if dList.Kind() != reflect.Slice {
		fmt.Printf("Unsurpported table object: [%s]", dList.Kind())
		return
	}
	dLength := dList.Len()
	for i := 0; i < dLength; i++ {
		v := dList.Index(i)
		vType := v.Type()
		// 如果List中的对象是结构体指针
		if vType.Kind() == reflect.Pointer {
			vType = vType.Elem()
		}

		if vType.Kind() == reflect.Struct {
			if len(that.Headers) == 0 && i == 0 {
				// 解析表头Headers
				that.ParseHeadersFromObject(vType)
			}
			row := []string{}
			for _, name := range that.Headers {
				// field := v.Elem().FieldByName(name).String()
				fStr, _ := json.Marshal(v.Elem().FieldByName(name).Interface())
				row = append(row, string(fStr))
			}
			that.AddRow(row...) // 添加当前行
		}
	}
}

const (
	HeadersInJson = "headers"
	RowsInJson    = "rows"
)

func (that *KtrlTable) ParseHeadersFromString(jsonStr string) {
	hValue := gjson.Get(jsonStr, HeadersInJson).Array()
	if len(hValue) > 0 {
		for _, v := range hValue {
			that.Headers = append(that.Headers, v.String())
		}
		that.AddHeaders(that.Headers...) // 添加表头
	}
}

/*
  AddRowsByJsonString 为表格添加数据行；
  jsonString 格式: {headers: ["", "", ""], rows: [["", "", ""], ["", "", ""], ["", "", ""]]}
*/
func (that *KtrlTable) AddRowsByJsonString(jsonString string) {
	if len(that.Headers) == 0 {
		that.ParseHeadersFromString(jsonString)
	}
	rValue := gjson.Get(jsonString, "rows").Array()
	for _, rv := range rValue {
		l := []string{}
		row := rv.Array()
		for _, field := range row {
			l = append(l, field.String())
		}
		if len(l) > 0 {
			that.AddRow(l...)
		}
	}
}
