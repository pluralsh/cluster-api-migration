package resources

import (
	"encoding/json"
	"fmt"
)

type Printer interface {
	PrettyPrint()
}

type printer struct {
	i interface{}
}

func (this *printer) PrettyPrint() {
	s, _ := json.MarshalIndent(this.i, "", "\t")

	fmt.Println(string(s))
}

func NewPrinter(i interface{}) Printer {
	return &printer{i: i}
}
