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

func (printer *printer) PrettyPrint() {
	s, _ := json.MarshalIndent(printer.i, "", "  ")

	fmt.Println(string(s))
}

func NewPrinter(i interface{}) Printer {
	return &printer{i: i}
}
