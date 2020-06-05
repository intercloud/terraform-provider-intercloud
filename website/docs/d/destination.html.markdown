---
subcategory: 'CONNECTOR'
layout: 'intercloud'
page_title: 'InterCloud: intercloud_destination'
description: |-
  Provides details about a specific InterCloud destination
---

# Data Source: intercloud_destination

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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The destination ID.
