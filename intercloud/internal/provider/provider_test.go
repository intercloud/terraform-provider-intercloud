package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestProvider(t *testing.T) {
	if err := IntercloudProvider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = IntercloudProvider()
}

func TestAccDefault(t *testing.T) {
	err := IntercloudProvider().Configure(terraform.NewResourceConfigRaw(map[string]interface{}{
		"organization_id": os.Getenv("ACCEPTANCE_ORG_ID"),
	}))
	if err != nil {
		t.Fatal(err)
	}
}
