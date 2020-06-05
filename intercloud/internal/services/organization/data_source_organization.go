package organization

import (
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/config"
)

func dataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationRead,

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
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

func dataSourceOrganizationRead(d *schema.ResourceData, meta interface{}) (err error) {
	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	// override provider configuration, if organization_id specified in datasource
	var organization_id uuid.UUID
	if d.Get("organization_id").(string) != "" {
		organization_id = uuid.MustParse(d.Get("organization_id").(string))
	} else {
		organization_id = cfg.OrganizationID()
	}

	input := &client.ReadOrganizationInput{
		OrganisationID: organization_id,
	}

	organization, err := api.ReadOrganization(input)

	if err != nil {
		log.Println("[ERROR] erro while reading organizations")
		log.Println(err)
		return
	}

	d.SetId(organization.ID.String())
	d.Set("organization_id", organization.ID.String())
	d.Set("name", organization.Name)
	d.Set("parent_id", organization.ParentID.String())

	return nil
}
