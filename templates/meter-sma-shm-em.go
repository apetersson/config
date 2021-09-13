package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "meter",
		Type:   "sma",
		Name:   "SMA Sunny Home Manager 2.0 / Energy Meter (Grid, PV or Battery Meter)",
		Params: []registry.TemplateParam{
			{
				Name: "ip",
				Value: "192.0.2.2",
				Hint: "IP address or hostname of the device",
			},
		},
		Sample: `uri: {{ .ip.value }}`,
	}

	registry.Add(template)
}
