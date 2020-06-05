package aws

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ResourceSchemaAws() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"aws_account_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"aws_bgp_asn": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtMost(int(^uint16(0))), // 65535
				Default:      64512,
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
				Optional: true,
			},
			"dxvif": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
