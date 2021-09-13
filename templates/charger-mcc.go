package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "charger",
		Type:   "mcc",
		Name:   "Mobile Charger Connect (Audi, Bentley, Porsche)",
		Params: []registry.TemplateParam{
			{
				Name: "ip",
				Value: "192.0.2.2",
				Hint: "IP address or hostname of the device",
			},
			{
				Name: "password",
				Value: "password",
				Hint: "Password of the home user",
			},
		},
		Sample: `uri: https://{{ .ip.value }} # {{ .ip.hint }}
password: {{ .password.value }} # {{ .password.hint }}`,
	}

	registry.Add(template)
}
