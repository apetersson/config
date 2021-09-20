package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/evcc-io/config/registry"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

const (
	ext     = ".yaml"
	summary = "template.md"
)

var (
	confYaml       string
	confGo         bool
	confOutGo      string
	confSummary    bool
	confOutSummary string
	confHelp       bool
	tmpl           *template.Template
)

func init() {
	flag.StringVarP(&confYaml, "yaml", "y", "yaml", "yaml path")
	flag.StringVarP(&confOutGo, "output-go", "o", "", "output go files path")
	flag.StringVarP(&confOutSummary, "output-summary", "f", "", "output summary file")
	flag.BoolVarP(&confGo, "go", "g", false, "generate go files")
	flag.BoolVarP(&confSummary, "summary", "s", false, "generate summary")
	flag.BoolVarP(&confHelp, "help", "h", false, "help")
	flag.Parse()
}

const (
	magicCommentModbusSetup = "# ::modbus-setup::"
	modbusChoiceRS485       = "rs485"
	modbusChoiceTCPIP       = "tcpip"
	keyModbusId             = "id"
	keyModbusRS485Serial    = "rs485serial"
	keyModbusRS485TCPIP     = "rs485tcpip"
	keyModbusTCPIP          = "tcpip"
)

var modbusID = `id: 1`

var modbusRS485Serial = `
# locally attached:
device: /dev/ttyUSB0 # serial port
baudrate: 9600
comset: 8N1`

var modbusRS485TCPIP = `
# RS485 via TCP:
uri: 192.0.2.2:502
rtu: true # serial modbus rtu (rs485) device connected using simple ethernet adapter`

var modbusTCPIP = `
# via TCP:
uri: 192.0.2.2:502`

var modbusTemplate = `{{.` + keyModbusId + ` | indent 0}}{{.` + keyModbusRS485Serial + ` | indent 0}}{{.` + keyModbusRS485TCPIP + ` | indent 0}}{{.` + keyModbusTCPIP + ` | indent 0}}`

var sourceTemplate = `package templates {{/* Define backtick variable */}}{{$tick := "` + "`" + `"}}

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "{{.Class}}",
		Type:   "{{.Type}}",
		Name:   "{{.Name}}",
{{- if eq .Class "meter" }}
{{- if .Usage }}
		Usage: []string{
{{- range .Usage }}
			"{{.}}",
{{- end }}
		},
{{- end }}
{{- end }}
{{- if .Params }}
		Params: []registry.TemplateParam{
{{- range .Params }}
			{
				Name: "{{.Name}}",
{{- if not (eq .Type "") }}
				Type: "{{.Type}}",
				Choice: []string{
{{- range .Choice }}
					"{{.}}",
{{- end }}
				},
{{- else }}
				Value: "{{.Value}}",
{{- if .Optional }}
				Optional: {{.Optional}},
{{- end }}
			{{- if .Hint }}
				Hint: "{{.Hint}}",
			{{- end }}
{{- end }}
			},
{{- end }}
		},
{{- end }}
		Sample: {{$tick}}{{escape .PlainSample}}{{$tick}},
	}

	registry.Add(template)
}
`

func scanFolder(root string) (files []string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(info.Name()) == ext {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func parseSample(file, typ string) registry.Template {
	src, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var sample registry.Template
	if err := yaml.Unmarshal(src, &sample); err != nil {
		panic(err)
	}

	sample.Class = typ

	// trim trailing linebreaks
	sample.PlainSample = strings.TrimRight(sample.Sample, "\r\n")
	sample.Sample = strings.TrimRight(sample.Sample, "\r\n")

	sample = renderSample(sample)

	return sample
}

func renderSample(sample registry.Template) registry.Template {
	if len(sample.Params) == 0 && len(sample.Usage) == 0 {
		return sample
	}

	sampleTmpl, err := template.New("sample").Option("missingkey=zero").Parse(sample.Sample)
	if err != nil {
		panic(err)
	}

	if sample.Class == "meter" {
		usageItems := []string{}
		for _, item := range sample.Usage {
			if !contains(registry.ValidUsageTypes, item) {
				panic("usage " + item + " is invalid")
			}
			usageItems = append(usageItems, registry.UsageTypeDescriptions[item])
		}
		if len(usageItems) > 0 {
			sort.Slice(usageItems, func(i, j int) bool {
				return usageItems[i] < usageItems[j]
			})
			sample.Name = sample.Name + " (" + strings.Join(usageItems, ", ") + ")"
		}
	}

	var modbusChoices []string
	paramItems := make(map[string]interface{})

	for _, item := range sample.Params {
		paramItem := make(map[string]string)

		if item.Name == "" {
			panic("params name is required")
		}
		if item.Value == "" && !item.Optional && item.Type == "" {
			panic("params value or type is required")
		}
		if item.Type != "" && len(item.Choice) == 0 {
			panic("params choice is required with type")
		}

		if len(item.Choice) > 0 {
			for _, choice := range item.Choice {
				if !contains(registry.ValidParamsChoiceTypes, choice) {
					panic("param choice " + choice + " is invalid")
				}
			}
			modbusChoices = item.Choice
		}

		if item.Value != "" {
			paramItem["value"] = item.Value
		}
		if item.Hint != "" {
			paramItem["hint"] = item.Hint
		}
		paramItems[item.Name] = paramItem
	}

	var tpl bytes.Buffer
	if err = sampleTmpl.Execute(&tpl, paramItems); err != nil {
		panic(err)
	}

	sample.Sample = tpl.String()

	if len(modbusChoices) > 0 {
		var choices = map[string]string{
			keyModbusRS485Serial: "",
			keyModbusRS485TCPIP:  "",
			keyModbusTCPIP:       "",
		}
		if contains(modbusChoices, modbusChoiceRS485) {
			choices[keyModbusRS485Serial] = modbusRS485Serial
			choices[keyModbusRS485TCPIP] = modbusRS485TCPIP
		}
		if contains(modbusChoices, modbusChoiceTCPIP) {
			choices[keyModbusTCPIP] = modbusTCPIP
		}
		choices[keyModbusId] = modbusID

		// search for magicCommentModbusSetup and replace it with the correct indentation
		r := regexp.MustCompile(`.*` + magicCommentModbusSetup + `.*`)
		matches := r.FindAllString(sample.Sample, -1)
		for _, match := range matches {
			indentation := strings.Repeat(" ", strings.Index(match, magicCommentModbusSetup))

			result := renderModbus(modbusTemplate, len(indentation), choices)

			sample.Sample = strings.ReplaceAll(sample.Sample, match, result)
		}
	}

	return sample
}

func contains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}

func renderModbus(tmpl string, indentlength int, modbusChoices map[string]string) string {
	tmpl = strings.ReplaceAll(tmpl, " | indent 0", " | indent "+strconv.Itoa(indentlength))

	var tpl bytes.Buffer
	if err := template.Must(template.New("modbus").Funcs(template.FuncMap(sprig.FuncMap())).Parse(tmpl)).Execute(&tpl, modbusChoices); err != nil {
		panic(err)
	}

	return tpl.String()
}

func render(wr io.Writer, sample registry.Template) {
	if tmpl == nil {
		var err error
		tmpl, err = template.New("test").Funcs(template.FuncMap{
			// escape backticks in raw strings
			"escape": func(s string) string {
				return strings.ReplaceAll(s, "`", "`+\"`\"+`")
			},
		}).Parse(sourceTemplate)

		if err != nil {
			panic(err)
		}
	}

	tmpl.Execute(wr, sample)
}

func renderSummary(wr io.Writer, samples []registry.Template) {
	summaryTemplate, err := os.ReadFile(summary)
	if err != nil {
		panic(err)
	}

	// prepare outside of loop
	re, err := regexp.Compile("[^a-zA-ZäöüÄÖÜ0-9]")
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("test").Funcs(template.FuncMap{
		// filter samples by class
		"filter": func(class string, samples []registry.Template) (reg []registry.Template) {
			for _, sample := range samples {
				if sample.Class == class {
					reg = append(reg, sample)
				}
			}
			return
		},
		// https://github.com/Masterminds/sprig/blob/48e6b77026913419ba1a4694dde186dc9c4ad74d/strings.go#L109
		"indent": func(spaces int, v string) string {
			pad := strings.Repeat(" ", spaces)
			return pad + strings.Replace(v, "\n", "\n"+pad, -1)
		},
		// unique link target
		"href": func(class, name string) string {
			link := strings.ReplaceAll(re.ReplaceAllString(strings.ToLower(name), "-"), "--", "-")
			return strings.Trim(strings.ToLower(class)+"-"+link, "-")
		},
	}).Parse(string(summaryTemplate))

	if err != nil {
		panic(err)
	}

	tmpl.Execute(wr, samples)
}

func output(file string, fun func(io.Writer)) {
	wr := os.Stdout
	if file != "" {
		var err error
		wr, err = os.Create(file)
		if err != nil {
			panic(err)
		}
	}

	fun(wr)
	wr.Close()
}

func main() {
	if confHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var samples []registry.Template

	files := scanFolder(confYaml)
	for _, file := range files {
		// example type
		dir := filepath.Dir(file)
		typ := filepath.Base(dir)
		typ = strings.TrimRight(typ, "s") // de-pluralize

		sample := parseSample(file, typ)

		samples = append(samples, sample)

		if confGo {
			var out string
			if confOutGo != "" {
				name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
				out = fmt.Sprintf("%s/%s-%s.go", confOutGo, typ, name)
			}

			println(out)

			output(out, func(wr io.Writer) {
				render(wr, sample)
			})
		}
	}

	if confSummary {
		sort.Sort(registry.Templates(samples))
		output(confOutSummary, func(wr io.Writer) {
			renderSummary(wr, samples)
		})
	}
}
