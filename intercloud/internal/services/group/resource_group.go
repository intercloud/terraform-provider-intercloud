package group

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/client"
	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/config"
)

// ===============================================
// schema
// ===============================================
func resourceGroup() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"organization_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"irn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{ //@TODO: fix import only deployed to avoid problems
			State: schema.ImportStatePassthrough,
		},
	}
}

// ===============================================
// create
// ===============================================
func resourceGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)

	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	if d.Get("organization_id").(string) == "" {
		d.Set("organization_id", organizationID.String())
	}

	input := &client.CreateGroupInput{
		CreateInput: client.CreateInput{
			OrganizationID: uuid.MustParse(d.Get("organization_id").(string)),
		},
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	group, err := api.CreateGroup(input)

	if err != nil {
		// @FIXME: Handle specific error
		// uwerr := errors.Unwrap(err)
		// switch uwerr {
		// case errors.Is(ap):
		// 	return
		// }
		return
	}

	d.SetId(group.ID.String())

	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)

	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	if d.Get("organization_id").(string) == "" {
		d.Set("organization_id", organizationID.String())
	}

	input := &client.ReadGroupInput{
		ID:             uuid.MustParse(d.Id()),
		OrganizationID: uuid.MustParse(d.Get("organization_id").(string)),
	}

	group, err := api.ReadGroup(input)

	if err != nil {
		return
	}

	d.Set("name", group.Name)
	d.Set("irn", group.Irn)
	d.Set("description", group.Description)
	d.Set("organization_id", group.OrganizationID.String())
	d.Set("capacity", ConvertCapacityWithUnits(group.AllocatedCapacity))

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	if d.Get("organization_id").(string) == "" {
		d.Set("organization_id", organizationID.String())
	}

	input := &client.DeleteGroupInput{
		DeleteInput: client.DeleteInput{
			ID:             uuid.MustParse(d.Id()),
			OrganizationID: uuid.MustParse(d.Get("organization_id").(string)),
		},
	}

	err = api.DeleteGroup(input)

	if err == nil {
		d.SetId("")
	}

	return nil
}
