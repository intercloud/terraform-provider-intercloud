package azure

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ResourceSchemaAzure() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
