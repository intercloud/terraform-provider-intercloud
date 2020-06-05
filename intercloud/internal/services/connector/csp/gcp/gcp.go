package gcp

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// GcpMed maps human readable values (low, high) to int values
var GcpMed = map[string]uint64{"low": 2147483648, "high": 5000}

func ResourceSchemaGcp() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"med": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "low",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"low",
						"high",
					},
					true,
				),
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "BPS_50M",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"BPS_50M", // @FIXME: Tech preview only BPS_50M
						// "BPS_100M",
						// "BPS_200M",
						// "BPS_300M",
						// "BPS_400M",
						// "BPS_500M",
						// "BPS_1G",
						// "BPS_2G",
						// "BPS_5G",
						// "BPS_10G",
						// "BPS_20G",
						// "BPS_50G",
					},
					false,
				),
			},
			"pairing_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interconnect_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// HumanMedToInt converts a human readable med (low, high) to its int value.
// When the human readable value is unknown, the equivalent int value for "low" is returned
func HumanMedToInt(med string) uint64 {
	if val, ok := GcpMed[strings.ToLower(med)]; ok {
		return val
	}
	// default value is "low"
	return GcpMed["low"]
}

// IntMedToHuman converts an integer med value to its human readable equivalent.
// When the int value is unknown, the "low" human readable value is returned
func IntMedToHuman(med uint64) string {
	for k, v := range GcpMed {
		if v == med {
			return k
		}
	}
	// default is "low"
	return "low"
}
