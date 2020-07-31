package aws

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ResourceSchemaAwsHostedConnection() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"port_speed": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "50Mbps",
				ValidateFunc: validation.StringInSlice([]string{"50Mbps", "100Mbps", "200Mbps", "300Mbps", "400Mbps", "500Mbps", "1Gbps", "2Gbps", "5Gbps"}, false),
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_peer_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_peer_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bgp_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}
