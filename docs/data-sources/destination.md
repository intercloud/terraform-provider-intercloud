---
subcategory: 'Connector'
layout: 'intercloud'
description: |-
  Provides details about a specific InterCloud destination
---

# intercloud_destination Data Source

Provide details about a specific InterCloud CSP (Cloud Solution Provider) destination.

This datasource differs from `intercloud_destinations` which list destinations
matching criteria.

## Example Usage

Get details about an AWS destination.

```hcl
data "intercloud_destination" "aws_ireland" {
  family= "aws"
  location = "Ireland"
}
```

Get details about an AWS destination for an hosted connection.

```hcl
data "intercloud_destination" "aws_hosted_connection_ireland" {
  family= "awshostedconnection"
  location = "Ireland"
}
```

Get details about an Azure destination.

```hcl
data "intercloud_destination" "azure_ireland" {
  family= "azure"
  location = "Ireland"
}
```

Get details about an Google Cloud destination.

```hcl
data "intercloud_destination" "gcp_london" {
  family= "gcp"
  location = "London"
}
```

## Argument Reference

The following arguments are supported:

- `family` - (Required) The cloud provider family ( `aws`, `azure` or `gcp`).

- `location` - (Required) The on-ramp location.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The destination ID.
