package resources

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"
)

type Printer interface {
	PrettyPrint()
}

type jsonPrinter struct {
	i interface{}
}

func (this *jsonPrinter) PrettyPrint() {
	s, _ := json.MarshalIndent(this.i, "", "  ")
	fmt.Println(string(s))
}

type yamlPrinter struct {
	i interface{}
}

func (this *yamlPrinter) PrettyPrint() {
	s, _ := yaml.Marshal(this.i)
	fmt.Println(string(s))
}

func NewJsonPrinter(i interface{}) Printer {
	return &jsonPrinter{i: i}
}

func NewYAMLPrinter(i interface{}) Printer {
	return &yamlPrinter{i: i}
}
