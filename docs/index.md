---
layout: 'intercloud'
description: |-
  The InterCloud provider is used to configure your InterCloud infrastructure.
---

# InterCloud Provider

The InterCloud provider is used to configure your [InterCloud](https://intercloud.com)
infrastructure. The provider needs to be configured with
the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the InterCloud Provider
provider "intercloud" {
  version         = "~> 0.0.1"
  organization_id = "6cead98e-98a4-430b-b924-499ea5638f04"
}
```

## Organization

You must provide the ID of the organization containing the managed resources.
The organization ID can be configured via an environment variable or in the provider configuration.
The reading order is:

- Static organization ID specified in provider configuration
- Environment variable

Your organization ID can be found with via the [REST API endpoint "/me/info"](https://doc.intercloud.io/apiref/#operation/getMe)

### Static organization ID

You can define the `organization_id` field into the provider configuration.

Usage:

```hcl
provider "intercloud" {
  organization_id = "6cead98e-98a4-430b-b924-499ea5638f04"
}
```

### Environment variable organization ID

You can define the `INTERCLOUD_ORGANIZATION_ID` environment variable.

Usage:

```shell
export INTERCLOUD_ORGANIZATION_ID="6cead98e-98a4-430b-b924-499ea5638f04"
terraform plan
```

## Authentication

The InterCloud provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables
- Shared credentials file

### Static credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `access_token` in-line
in the InterCloud provider block:

Usage:

```hcl
provider "intercloud" {
  access_token = "my-access-token"
}
```

### Environment variable

You can provide your credential via the `INTERCLOUD_ACCESS_TOKEN` environment
variable, representing your InterCloud Personal Access Token.

```hcl
provider "intercloud" {}
```

Usage:

```sh
export INTERCLOUD_ACCESS_TOKEN="my-access-token"
terraform plan
```

### Shared Credentials file

You can use an InterCloud credentials file to specify your credentials. The
default location is `$HOME/.intercloud/credentials.json` on Linux and OS X, or
`"%USERPROFILE%\.intercloud\credentials"` for Windows users. If we fail to
detect credentials inline, or in the environment, Terraform will check
this location.

Usage:

```hcl
provider "intercloud" {
  shared_credentials_file = "/Users/tf_user/.intercloud/my-custom-creds.json"
}
```

Shared credentials file expected content:

```json
{
  "access_token": "my-access-token"
}
```

## API Endpoint

The called API endpoint can be customized, essentially, for testing purpose on
dedicated infrastructure.
The custom value is read in this order of priority :

- Static endpoint specified in provider configuration
- Environment variable

### Static API Endpoint

Usage :

```hcl
provider "intercloud" {
  api_endpoint    = "https://api-console-lab.intercloud.io"
}
```

### Environment variable Endpoint

You can provide the custom API enpoint through the `INTERCLOUD_API_ENDPOINT`
environment variable.

```hcl
provider "intercloud" {}
```

Usage:

```sh
export INTERCLOUD_API_ENDPOINT="https://api-console-lab.intercloud.io"
terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the InterCloud
`provider` block:

- `organization_id` - (Optional) Organization where resources are managed. It
  must be provided, but it can also be sourced from the `INTERCLOUD_ORGANIZATION_ID`
  environment variable.

- `access_token` - (Optional) This is the InterCloud personal access token. It
  must be provided, but it can also be sourced from the `INTERCLOUD_ACCESS_TOKEN`
  environment variable, or via a shared credentials file.

- `api_endpoint` - (Optional) This is the customizable api endpoint. It can also
  be sourced from the `INTERCLOUD_API_ENDPOINT` environment variable.

## Rate limits strategy

A rate limiter protects the InterCloud REST API.

In order to properly handle rate limits, this plugin implements an exponential backoff strategy.

The strategy works like this:

- If the rate limit is hit, a maximum of 3 retries is done
- The elapsed time between each retry is exponential (1s -> 2s -> 4s)
- A jitter of +/- 500ms is added in order to not retry all the failed requests at the same time

This strategy is sufficient to ensure that all HTTP requests made on InterCloud REST API succeed without any error due to rate limits.
