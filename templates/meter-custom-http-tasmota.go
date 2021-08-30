package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "meter",
		Type:   "custom",
		Name:   "Tasmota (Grid Meter, PV Meter)",
		Sample: `power: # power reading (W)
  source: http
  uri: http://192.0.2.2/cm?cmnd=Status%208
  jq: .StatusSNS.ENERGY.Power`,
	}

	registry.Add(template)
}
