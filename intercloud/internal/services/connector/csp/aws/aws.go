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
			"hosted_connection": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_speed": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "50Mbps",
							ValidateFunc: validation.StringInSlice([]string{"50Mbps", "100Mbps", "200Mbps", "300Mbps", "400Mbps", "500Mbps", "1Gbps", "2Gbps", "5Gbps"}, false),
						},
						"vlan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dxvif": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
