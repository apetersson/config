package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "meter",
		Type:   "custom",
		Name:   "Fronius Symo GEN24 Plus (Battery Meter)",
		Params: []registry.TemplateParam{
			{
				Name: "interface",
				Type: "modbus",
				Choice: []string{
					"tcpip",
				},
			},
		},
		Sample: `power:
  source: calc
  add:
  - source: modbus
    model: sunspec
    # ::modbus-setup::
    value: 160:3:DCW # mppt 3 charge
    scale: -1
  - source: modbus
    model: sunspec
    # ::modbus-setup::
    value: 160:4:DCW # mppt 4 discharge
soc:
  source: modbus
  model: sunspec
  # ::modbus-setup::
  value: ChargeState`,
	}

	registry.Add(template)
}
