package intercloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/provider"
)

func Provider() terraform.ResourceProvider {
	return provider.IntercloudProvider()
}
