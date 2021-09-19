package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "meter",
		Type:   "modbus",
		Name:   "Eastron SDM Modbus Meter",
		Params: []registry.TemplateParam{
			{
				Name: "interface",
				Type: "modbus",
				Choice: []string{
					"serial",
					"tcprtu",
				},
			},
		},
		Sample: `model: sdm # specific non-sunspec meter
energy: Sum # only required for charge meter usage
# ::modbus-setup::`,
	}

	registry.Add(template)
}
