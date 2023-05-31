package resources

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"
)

type PrinterMode string

const (
	PrinterModeJSON = PrinterMode("json")
	PrinterModeYAML = PrinterMode("yaml")
)

type Printer interface {
	PrettyPrint()
}

type printer struct {
	mode PrinterMode
	i    interface{}
}

func (this *printer) prettyPrint(mode PrinterMode) {
	switch mode {
	case PrinterModeJSON:
		this.prettyPrintJSON()
		return
	case PrinterModeYAML:
		this.prettyPrintYAML()
	}
}

func (this *printer) prettyPrintJSON() {
	s, _ := json.MarshalIndent(this.i, "", "  ")
	fmt.Println(string(s))
}

func (this *printer) prettyPrintYAML() {
	s, _ := yaml.Marshal(this.i)
	fmt.Println(string(s))
}

func (this *printer) PrettyPrint() {
	this.prettyPrint(this.mode)
}

func NewPrinter(i interface{}, mode PrinterMode) Printer {
	return &printer{i: i, mode: mode}
}

func NewJsonPrinter(i interface{}) Printer {
	return &printer{i: i, mode: PrinterModeJSON}
}

func NewYAMLPrinter(i interface{}) Printer {
	return &printer{i: i, mode: PrinterModeYAML}
}
