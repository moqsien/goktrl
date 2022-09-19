package main

import (
	"github.com/moqsien/goktrl"
)

func renderTableFromObject() {
	type Data struct {
		Price  string `order:"2"`
		Stokes string `order:"3"`
		Name   string `order:"1"`
	}
	var dList = []*Data{
		{Name: "Apple", Price: "6", Stokes: "128"},
		{Name: "Banana", Price: "3", Stokes: "256"},
		{Name: "Pear", Price: "5", Stokes: "121"},
	}
	table := goktrl.NewKtrlTable()
	table.AddRowsByListObject(dList)
	table.Render()
}

func renderTableFromString() {
	table := goktrl.NewKtrlTable()
	table.AddRowsByJsonString(`{
		"headers": ["Name", "Price", "Stokes"],
		"rows": [
		  ["Apple","6", "128"],
		  ["Banana","3", "256"],
		  ["Pear","5", "121"]
		]
	  }`)
	table.Render()
}

func main() {
	renderTableFromObject()
}
