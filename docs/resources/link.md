---
subcategory: 'Link'
layout: 'intercloud'
description: |-
  Creates and manages a link
---

# intercloud_link Resource

Manages a link between two `intercloud_connector`

## Example Usage

```hcl
resource "intercloud_link" "azure_to_aws_1" {
    name        = "Azure to AWS"
    description = "A sample link from Azure to AWS"
    from        = intercloud_connector.azure_connector_1.id
    to          = intercloud_connector.aws_connector_1.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The link name.

- `description` - (Required) The link description.

- `from` - (Required) The first connector ID to interconnect.

- `to` - (Required) The second connector ID to interconnect.

- `organization_id` - (Optional) The organization ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The link id.

- `irn` - The link IRN.

- `state` - The link state.
