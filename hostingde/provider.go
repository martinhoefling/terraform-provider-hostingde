package hostingde

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HOSTINGDE_TOKEN", nil),
				Description: "The API token to access the Hosting.de API.",
			},
			"owner_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HOSTINGDE_ACCOUNTID", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hostingde_zone":   resourceZone(),
			"hostingde_record": resourceRecord(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return NewClient(d.Get("auth_token").(string), d.Get("owner_account_id").(string)), nil
}
