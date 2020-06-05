---
subcategory: 'ORGANIZATION'
layout: 'intercloud'
page_title: 'InterCloud: intercloud_organizations'
description: |-
  Provides a list of children organizations
---

# Data Source: intercloud_organizations

Provide access to children organizations.

This datasource is different from `intercloud_organization` which get details
about a specific organization.

## Example Usage

List all available organizations.

```hcl
data "intercloud_organizations" "children" {
}
```

## Attributes Reference

The following attributes are exported:

- `children` - A list of organizations records. Structure is documented below.

The organizations block supports :

- `organization_id` - The organization ID

- `name` - The organization name.

- `parent_id` - The organization parent ID.
