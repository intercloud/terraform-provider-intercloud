package organization

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/client"
	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/config"
)

func dataSourceOrganizations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationsRead,

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     organizationResource(),
			},
		},
	}
}

func organizationResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOrganizationsRead(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	input := &client.ReadOrganizationInput{
		OrganisationID: cfg.OrganizationID(),
	}

	organizations, err := api.ReadOrganizations(input)

	if err != nil {
		log.Println("[ERROR] erro while reading organization")
		log.Println(err)
		return
	}

	children := flattenOrganizations(organizations)
	d.SetId(cfg.OrganizationID().String())
	d.Set("organization_id", cfg.OrganizationID().String())
	d.Set("children", children)

	return nil
}

func flattenOrganizations(in *[]client.ReadOrganizationOutput) []map[string]string {
	result := make([]map[string]string, len(*in))
	for i, org := range *in {
		result[i] = map[string]string{
			"organization_id": org.ID.String(),
			"name":            org.Name,
			"parent_id":       org.ParentID.String(),
		}
	}
	return result
}
