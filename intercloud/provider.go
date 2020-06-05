package intercloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/provider"
)

func Provider() terraform.ResourceProvider {
	return provider.IntercloudProvider()
}
