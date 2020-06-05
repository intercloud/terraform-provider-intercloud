package tests

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
	"text/template"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/acceptance"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/config"
)

func TestAccGroup_basic(t *testing.T) {

	data := acceptance.BuildTestData(t, "intercloud_group", "test", uuid.MustParse(os.Getenv("ACCEPTANCE_ORG_ID")))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccCheckResourceGroupDestroy(data),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupResource(data),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupResourceExists(data),
				),
			},
		},
	})
}

func testAccGroupResourceExists(data acceptance.TestData) resource.TestCheckFunc {

	return func(s *terraform.State) (err error) {

		api := acceptance.IntercloudProvider.Meta().(*config.ProviderConfig).ApiClient()

		for _, rs := range s.RootModule().Resources {

			if rs.Type != data.ResourceType {
				continue
			}

			_, err = api.ReadGroup(
				&client.ReadGroupInput{
					ID:             uuid.MustParse(rs.Primary.ID),
					OrganizationID: uuid.MustParse(rs.Primary.Attributes["organization_id"]),
				},
			)
		}

		return
	}
}

func testAccGroupResource(data acceptance.TestData) string {

	log.Println(fmt.Sprintf("[WARN] testAccCheckResourceGroupDestroy %+v", data))

	var tpl bytes.Buffer
	t := template.Must(
		template.New("group").Parse(`
	provider "intercloud" {
		organization_id = "{{- .Scope -}}"
	}
	
	resource "{{- .ResourceType -}}" "foo" {
		name            = "{{- .ResourceName -}}"
		description     = "toto"
		organization_id = "{{- .Scope -}}"
	}
	
	`),
	)

	err := t.Execute(&tpl, data)
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("[WARN] testAccCheckResourceGroupDestroy %+s", tpl.String()))

	return tpl.String()
}

func testAccCheckResourceGroupDestroy(data acceptance.TestData) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		log.Println(fmt.Sprintf("[WARN] testAccCheckResourceGroupDestroy %+v", s))

		c := acceptance.IntercloudProvider.Meta().(*config.ProviderConfig)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != data.ResourceType {
				continue
			}

			c.ApiClient().DeleteGroup(
				&client.DeleteGroupInput{
					DeleteInput: client.DeleteInput{
						ID:             uuid.MustParse(rs.Primary.ID),
						OrganizationID: uuid.MustParse(rs.Primary.Attributes["organization_id"]),
					},
				},
			)

		}

		return nil
	}
}
