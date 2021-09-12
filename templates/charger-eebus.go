package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "charger",
		Type:   "eebus",
		Name:   "EEBUS compatible wallbox (e.g. Mobile Charger Connect)",
		Params: []registry.TemplateParam{
			{
				Name: "ski",
				Value: "1234-5678-90ab-cdef-1234-5678-90ab-cdef-1234-5678",
				Hint: "Enter the SKI of the wallbox which can usually be found in the wallbox web interface",
			},
		},
		Sample: `ski: {{ .values.ski }} # {{ .hints.ski }}
forcePVLimits: true # use Overload Protection to limit PV charging, if false PV surplus is sent as recommended charging levels to the EV, but this is currently unreliable`,
	}

	registry.Add(template)
}
