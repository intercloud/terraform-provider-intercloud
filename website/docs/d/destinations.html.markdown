---
subcategory: 'CONNECTOR'
layout: 'intercloud'
page_title: 'InterCloud: intercloud_destinations'
description: |-
  List InterCloud destinations
---

# Data Source: intercloud_destinations

List InterCloud CSP (Cloud Solution Provider) destinations matching criteria.

This datasource differs from `intercloud_destination` which provides details
about a specific destination.

## Example Usage

List AWS destinations.

```hcl
data "intercloud_destinations" "aws_destinations" {
  family = "aws"
}
```

List AWS destinations for an hosted connection.

```hcl
data "intercloud_destinations" "aws_hosted_connection_destinations" {
  family = "awshostedconnection"
}
```

List Azure destinations.

```hcl
data "intercloud_destinations" "azure_destinations" {
  family = "azure"
}
```

List Google destinations.

```hcl
data "intercloud_destinations" "gcp_destinations" {
  family = "gcp"
}
```

## Argument Reference

The following arguments are supported:

- `family` - (Required) The cloud provider family ( `aws`, `awshostedconnection`, `azure` or `gcp`).

## Attributes Reference

The following attributes are exported:

- `destinations` - A list of destinations records. Structure is documented below.

The organizations block supports :

- `id` - The destination ID.

- `family` - The destination family.

- `location` - The destination location.
