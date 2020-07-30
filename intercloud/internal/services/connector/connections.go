package connector

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector/csp/gcp"
)

type Connection int

const (
	ConnectionAws Connection = iota
	ConnectionAwsHosted
	ConnectionAzure
	ConnectionGcp
)

type connectionData struct {
	family         string
	connectionType string
}

var (
	sliceConnections = []connectionData{
		{family: FamilyAws.String()},
		{family: FamilyAws.String(), connectionType: "hosted_connection"},
		{family: FamilyAzure.String()},
		{family: FamilyGcp.String()},
	}
)

func (c Connection) Family() string {
	return sliceConnections[c].family
}

func (c Connection) ConnectionType() string {
	return sliceConnections[c].connectionType
}

func AllConnections() []Connection {
	return []Connection{ConnectionAws, ConnectionAwsHosted, ConnectionAzure, ConnectionGcp}
}

func expandConnectionAwsParams(m map[string]interface{}) *client.AwsParams {
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

func expandConnectionAwsHostedParams(m map[string]interface{}) *client.AwsHostedParams {
	params := client.AwsHostedParams{}

	if v, ok := m["aws_account_id"]; ok {
		params.AwsAccount = v.(int)
	}
	if v, ok := m["aws_bgp_asn"]; ok {
		params.ASN = uint16(v.(int))
	}
	if v, ok := m["hosted_connection"]; ok {
		hosted := v.(*schema.Set).List()[0].(map[string]interface{})
		log.Printf("[DEBUG] hosted params to expand : %+v", hosted)
		if v, ok := hosted["port_speed"]; ok {
			params.PortSpeed = v.(string)
		}
		if v, ok := hosted["vlan_id"]; ok {
			params.VlanID = v.(int)
		}
		if v, ok := hosted["connection_id"]; ok {
			params.ConnectionID = v.(string)
		}
		if v, ok := hosted["customer_peer_ip"]; ok {
			params.CustomerPeerIP = v.(string)
		}
		if v, ok := hosted["aws_peer_ip"]; ok {
			params.AwsPeerIP = v.(string)
		}
		if v, ok := hosted["bgp_key"]; ok {
			params.BgpKey = v.(string)
		}
	}

	return &params
}

func flattenConnectionAwsParams(params *client.AwsParams) []interface{} {
	result := make(map[string]interface{})

	result["aws_bgp_asn"] = int(params.ASN)
	result["aws_account_id"] = params.AwsAccount
	result["dxvif"] = params.Dxvif

	return []interface{}{result}
}

func flattenConnectionAwsHostedParams(params *client.AwsHostedParams) []interface{} {
	result := make(map[string]interface{})

	result["aws_bgp_asn"] = int(params.ASN)
	result["aws_account_id"] = params.AwsAccount

	hostedConnection := map[string]interface{}{
		"port_speed":       params.PortSpeed,
		"vlan_id":          params.VlanID,
		"connection_id":    params.ConnectionID,
		"customer_peer_ip": params.CustomerPeerIP,
		"aws_peer_ip":      params.AwsPeerIP,
		"bgp_key":          params.BgpKey,
	}
	result["hosted_connection"] = []interface{}{hostedConnection}

	return []interface{}{result}
}

func expandConnectionAzureParams(m map[string]interface{}) *client.AzureParams {
	params := client.AzureParams{}

	if v, ok := m["service_key"]; ok {
		params.ServiceKey = v.(string)
	}

	return &params
}

func flattenConnectionAzureParams(params *client.AzureParams) []interface{} {
	result := make(map[string]interface{})

	result["service_key"] = params.ServiceKey

	return []interface{}{result}
}

func expandConnectionGcpParams(m map[string]interface{}) *client.GcpParams {
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

func flattenConnectionGcpParams(params *client.GcpParams) []interface{} {
	result := make(map[string]interface{})

	result["bandwidth"] = params.Bandwidth
	result["pairing_key"] = params.PairingKey
	result["med"] = gcp.IntMedToHuman(params.Med)
	result["interconnect_id"] = params.InterconnectID

	return []interface{}{result}
}
