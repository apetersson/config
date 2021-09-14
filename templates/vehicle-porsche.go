package templates 

import (
	"github.com/evcc-io/config/registry"
)

func init() {
	template := registry.Template{
		Class:  "vehicle",
		Type:   "porsche",
		Name:   "Porsche",
		Params: []registry.TemplateParam{
			{
				Name: "user",
				Value: "your@email.com",
				Hint: "Porsche ID (Email-Address)",
			},
			{
				Name: "Password",
				Value: "password",
				Hint: "Password of the Porsche ID",
			},
			{
				Name: "title",
				Value: "Taycan",
				Hint: "Display name of the vehicle in the UI",
			},
			{
				Name: "capacity",
				Value: "83",
				Hint: "The available battery capacity in kWh",
			},
			{
				Name: "vin",
				Value: "WP...",
				Hint: "The VIN number of your vehicle",
			},
		},
		Sample: `title: {{ .title.value }} # {{ .title.hint }}
capacity: {{ .capacity.value }} # {{ .capacity.hint }}
user: {{ .user.value }} # {{ .user.hint }}
password: {{ .password.value }} # {{ .password.hint }}
vin: {{ .vin.value }} # {{ .vin.hint }}`,
	}

	registry.Add(template)
}
