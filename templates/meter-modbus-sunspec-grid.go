package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "meter",
		Type:   "modbus",
		Name:   "SunSpec compliant 3-phase meter via inverter (Grid Meter)",
		Usage: []string{
			"grid",
		},
		Params: []registry.TemplateParam{
			{
				Name: "host",
				Value: "192.0.2.2",
				Hint: "IP address or hostname",
			},
			{
				Name: "port",
				Value: "502",
				Hint: "port address",
			},
		},
		Sample: `model: sunspec
uri: {{ .host.value }}:{{ .port.value }} # {{ .host.hint }} and {{ .port.hint }}
id: 1
power: 203:W # sunspec model 203 meter`,
	}

	registry.Add(template)
}
