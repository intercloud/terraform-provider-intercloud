---
subcategory: 'Organization'
layout: 'intercloud'
description: |-
  Provides details about a specific organization
---

# intercloud_organization Data Source

Provide detail about a specific organization.

This datasource is different from `intercloud_organizations` which provides a
way to list available organizations.

## Example Usage

Get details about current organization.

```hcl
data "intercloud_organization" "current" {
}
```

Get details about an organization.

```hcl
data "intercloud_organization" "my_organization" {
  organization_id = "my-organization-id"
}
```

## Argument Reference

The following arguments are supported:

- `organization_id` - (Required) The organization ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `name` - The organization name.

- `parent_id` - The organization parent ID.
