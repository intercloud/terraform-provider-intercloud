package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var IntercloudProvider *schema.Provider
var SupportedProviders map[string]terraform.ResourceProvider

func PreCheck(t *testing.T) {
	variables := []string{
		"INTERCLOUD_ACCESS_TOKEN",
		"INTERCLOUD_API_ENDPOINT",
		"ACCEPTANCE_ORG_ID",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}
