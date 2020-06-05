---
subcategory: 'GROUP'
layout: 'intercloud'
page_title: 'InterCloud: intercloud_group'
description: |-
  Creates and manages a group
---

# Resource: intercloud_group

Manages a group used to regroup resources like `intercloud_connector`

## Example Usage

```hcl
resource "intercloud_group" "my_group_1" {
    name = "My group #1"
    description = "A sample group named #1"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The group name.

- `description` - (Optional) The description of the resource group.

- `organization_id` - (Optional) The organization id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The group id.

- `capacity` - The allocated capacity (currently fixed to `50Mbps`)
