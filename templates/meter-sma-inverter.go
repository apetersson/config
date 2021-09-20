package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "meter",
		Type:   "sma",
		Name:   "SMA Speedwire Inverter (Battery Meter, PV Meter)",
		Usage: []string{
			"pv",
			"battery",
		},
		Params: []registry.TemplateParam{
			{
				Name: "host",
				Value: "192.0.2.2",
				Hint: "IP address or hostname",
			},
			{
				Name: "password",
				Value: "",
				Optional: true,
			},
		},
		Sample: `uri: {{ .host.value }} # {{ .host.hint }}
password: {{ .password.value }} # optional`,
	}

	registry.Add(template)
}
