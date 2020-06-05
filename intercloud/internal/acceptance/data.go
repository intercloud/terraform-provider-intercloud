package acceptance

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/provider"
)

var once sync.Once

type TestData struct {
	RandomInteger int

	RandomString string

	ResourceName string

	ResourceType string

	Scope string

	resourceLabel string
}

func BuildTestData(t *testing.T, resourceType, resourceLabel string, scope uuid.UUID) (testData TestData) {

	once.Do(func() {

		intercloudProvider := provider.IntercloudProvider()
		_ = intercloudProvider.Configure(terraform.NewResourceConfigRaw(map[string]interface{}{
			"organization_id": os.Getenv("ACCEPTANCE_ORG_ID"),
		}))
		IntercloudProvider = intercloudProvider.(*schema.Provider)

		SupportedProviders = map[string]terraform.ResourceProvider{
			"intercloud": intercloudProvider,
		}
	})

	testData = TestData{
		RandomInteger: acctest.RandInt(),
		RandomString:  acctest.RandString(5),
		ResourceName:  fmt.Sprintf("%s.%s", resourceType, resourceLabel),
		ResourceType:  resourceType,
		Scope:         scope.String(),
		resourceLabel: resourceLabel,
	}

	return
}
