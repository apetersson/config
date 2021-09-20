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
				Name: "password",
				Value: "password",
				Hint: "Password",
			},
			{
				Name: "title",
				Value: "Taycan",
				Hint: "display name for UI",
			},
			{
				Name: "capacity",
				Value: "83",
				Hint: "Battery capacity (kWh)",
			},
			{
				Name: "vin",
				Value: "WP...",
				Optional: true,
			},
		},
		Sample: `title: {{ .title.value }} # {{ .title.hint }}
capacity: {{ .capacity.value }} # {{ .capacity.hint }}
user: {{ .user.value }} # {{ .user.hint }}
password: {{ .password.value }} # {{ .password.hint }}
vin: {{ .vin.value }} # optional`,
	}

	registry.Add(template)
}
