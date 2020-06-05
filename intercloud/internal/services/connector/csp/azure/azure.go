package azure

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ResourceSchemaAzure() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"public_prefixes": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Set:      schema.HashString,
				MinItems: 1,
				Optional: true,
			},
		},
	}
}
