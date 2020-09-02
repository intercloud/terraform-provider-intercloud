---
subcategory: 'Connector'
layout: 'intercloud'
description: |-
  Creates and manages a connector
---

# intercloud_connector Resource

Manages a connector

## Example Usage

AWS connector

```hcl
resource "intercloud_connector" "aws_connector_1" {
  name = "My AWS connector #1"
  description = "A sample AWS connector named #1"
  destination_id = data.intercloud_destination.aws_ireland.id
  group_id = intercloud_group.my_group_1.id
  aws {
    aws_account_id = data.aws_caller_identity.current.account_id
    aws_bgp_asn = 65536
  }
}
```

AWS hosted connection connector

```hcl
resource "intercloud_connector" "aws_connector_1" {
  name = "My AWS connector #1"
  description = "A sample AWS connector named #1"
  destination_id = data.intercloud_destination.aws_ireland.id
  group_id = intercloud_group.my_group_1.id
  aws {
    aws_account_id = data.aws_caller_identity.current.account_id
    aws_bgp_asn = 65536
    hosted_connection {
      port_speed = "50Mbps"
    }
  }
}
```

Azure connector

```hcl
resource "intercloud_connector" "azure_connector_1" {
  name = "My AZURE connector #1"
  description = "A sample AZURE connector named #1"
  destination_id = data.intercloud_destination.azure_1.id
  group_id = intercloud_resource_group.my_group_1.id
  azure {
    service_key = azurerm_express_route_circuit.my_route.service_key
  }
}
```

Google Cloud connector

```hcl
resource "intercloud_connector" "gcp_connector_1" {
  name = "My AZURE connector #1"
  description = "A sample AZURE connector named #1"
  destination_id = data.intercloud_destination.azure_1.id
  group_id = intercloud_resource_group.my_group_1.id
  gcp {
    med = 1
    bandwidth = "BPS_50M"
    pairing_key = "5fc2f532-4483-4767-a7b7-ea7355052156/europe-west4/1"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The connector name.

- `description` - (Required) The connector description.

- `destination_id` - (Required) The destination ID.

- `group_id` - (Required) The group id.

~> **NOTE:** Only one cloud section (`aws`, `azure`, `gcp`) can be set at the
same time.

When managing AWS connector the following arguments are supported :

- `aws` - (Required) Key/value pairs for passing parameters :

  - `aws_account_id` - (Required) The AWS account ID.

  - `aws_bgp_asn` - (Optional) The AWS BGP ASN.

  - `hosted_connection` (Optional) Hosted connection parameters :
  
    - `port_speed` (Optional) The port speed among allowed values `50Mbps`,
    `100Mbps`, `200Mbps`, `300Mbps`, `400Mbps`, `500Mbps`, `1Gbps`, `2Gbps`,
    `5Gbps` (default: `50Mbps`)

When managing Azure connector the following arguments are supported :

- `azure` - (Required) Key/value pairs for passing parameters :

  - `azure_service_key` - (Required) The Azure express route service key.

When managing Google Cloud connector the following arguments are supported :

- `gcp` - (Required) Key/value pairs for passing parameters :

  - `pairing_key` - (Required) The Google Cloud pairing key.

  - `med` - (Optional) The priority value among allowed values `low`, `high`
  (default: `low`)

  - `bandwith` - (Optional) The allocated bandwith (only one available value
  for now: `BPS_50M`)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The connector id.

- `irn` - The connector IRN.

- `state` - The connector state.

When managing a AWS connector the following attributes are exposed:

- `aws` - Key/value pairs for exported attributes :

  - `dxvif` - (Direct connection only) The virtual interface created
  on AWS side.

  - `hosted_connection` - (Hosted connection only) Key/value pairs for
  hosted connection exported attributes :

    - `vlan_id` - VLAN ID.

    - `connection_id` - Hosted connection ID (e.g. `dxcon-fg31dyv6`).

    - `aws_peer_ip` - Amazon router peer IP to use when [creating a virtual interface ("Your router peer ip" field)](https://docs.aws.amazon.com/directconnect/latest/UserGuide/create-vif.html) on the hosted connection (e.g. `169.254.0.1`)

    - `customer_peer_ip` - Customer router peer IP to use when [creating a virtual interface ("Amazon router peer ip" field)](https://docs.aws.amazon.com/directconnect/latest/UserGuide/create-vif.html) on the hosted connection (e.g. `169.254.0.2`).

    - `bgp_key` - BGP authentication key to use when [creating a virtual interface ("BGP authentication key" field)](https://docs.aws.amazon.com/directconnect/latest/UserGuide/create-vif.html) on the hosted connection.

When managing a GCP connector the following attributes are exposed:

- `gcp` - Key/value pairs for exported attrbiutes :

  - `interconnect_id` - The virtual interface created on GCP side.
