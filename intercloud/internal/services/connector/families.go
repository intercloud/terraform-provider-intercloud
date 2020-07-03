package connector

import (
	"errors"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector/csp/gcp"
)

var (
	ErrInvalidFamily = errors.New("Invalid family")
)

type Family int

const (
	FamilyAws Family = iota + 1
	FamilyAzure
	FamilyGcp
)

var (
	sliceFamilies = []string{
		"aws",
		"azure",
		"gcp",
	}
)

func (s Family) String() string {
	return sliceFamilies[s-1]
}

func GetFamily(s string) (*Family, error) {
	for idx := range sliceFamilies {
		if s == sliceFamilies[idx] {
			f := Family(idx + 1)
			return &f, nil
		}
	}
	return nil, ErrInvalidFamily
}

func expandFamilyAwsParams(m map[string]interface{}) *client.AwsParams {
	params := client.AwsParams{}

	if v, ok := m["aws_account_id"]; ok {
		params.AwsAccount = v.(int)
	}
	if v, ok := m["aws_bgp_asn"]; ok {
		params.ASN = uint16(v.(int))
	}
	if v, ok := m["dxvif"]; ok {
		params.Dxvif = v.(string)
	}

	return &params
}

func flattenFamilyAwsParams(params *client.AwsParams) []interface{} {
	result := make(map[string]interface{})

	result["aws_bgp_asn"] = int(params.ASN)
	result["aws_account_id"] = params.AwsAccount
	result["dxvif"] = params.Dxvif

	return []interface{}{result}
}

func expandFamilyAzureParams(m map[string]interface{}) *client.AzureParams {
	params := client.AzureParams{}

	if v, ok := m["service_key"]; ok {
		params.ServiceKey = v.(string)
	}

	return &params
}

func flattenFamilyAzureParams(params *client.AzureParams) []interface{} {
	result := make(map[string]interface{})

	result["service_key"] = params.ServiceKey

	return []interface{}{result}
}

func expandFamilyGcpParams(m map[string]interface{}) *client.GcpParams {
	params := client.GcpParams{}

	if v, ok := m["bandwidth"]; ok {
		params.Bandwidth = v.(string)
	}
	if v, ok := m["pairing_key"]; ok {
		params.PairingKey = v.(string)
	}
	if v, ok := m["med"]; ok {
		params.Med = gcp.HumanMedToInt(v.(string))
	}
	if v, ok := m["interconnect_id"]; ok {
		params.InterconnectID = v.(string)
	}

	return &params
}

func flattenFamilyGcpParams(params *client.GcpParams) []interface{} {
	result := make(map[string]interface{})

	result["bandwidth"] = params.Bandwidth
	result["pairing_key"] = params.PairingKey
	result["med"] = gcp.IntMedToHuman(params.Med)
	result["interconnect_id"] = params.InterconnectID

	return []interface{}{result}
}
