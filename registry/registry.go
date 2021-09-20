package registry

import "strings"

// Registry is the template registry instance
var Registry = make([]Template, 0)

const (
	UsageTypeGrid    string = "grid"
	UsageTypePV      string = "pv"
	UsageTypeBattery string = "battery"
)

var ValidUsageTypes = []string{UsageTypeGrid, UsageTypePV, UsageTypeBattery}
var UsageTypeDescriptions = map[string]string{
	UsageTypeGrid:    "Grid Meter",
	UsageTypePV:      "PV Meter",
	UsageTypeBattery: "Battery Meter",
}

const (
	ParamsChoiceTypeRS485 string = "rs485"
	ParamsChoiceTypeTCPIP string = "tcpip"
)

var ValidParamsChoiceTypes = []string{ParamsChoiceTypeRS485, ParamsChoiceTypeTCPIP}

type TemplateParam struct {
	Name     string
	Value    string
	Hint     string
	Type     string
	Choice   []string
	Optional bool
}

// Template contains the template definition
type Template struct {
	Class       string
	Type        string
	Name        string
	Usage       []string
	Params      []TemplateParam
	Sample      string
	PlainSample string
}

func Add(t Template) {
	Registry = append(Registry, t)
}

func TemplatesByClass(class string) []Template {
	templates := make([]Template, 0)
	for _, t := range Registry {
		if t.Class == class {
			templates = append(templates, t)
		}
	}
	return templates
}

type Templates []Template

func (e Templates) Len() int {
	return len(e)
}

func (e Templates) Less(i, j int) bool {
	return strings.ToLower(e[i].Class) < strings.ToLower(e[j].Class) ||
		(strings.ToLower(e[i].Class) == strings.ToLower(e[j].Class)) && strings.ToLower(e[i].Name) < strings.ToLower(e[j].Name)
}

func (e Templates) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
