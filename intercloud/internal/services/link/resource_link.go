package link

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/config"
)

func resourceLink() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"organization_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"from": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"to": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"irn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceLinkCreate,
		Read:   resourceLinkRead,
		Update: resourceLinkUpdate,
		Delete: resourceLinkDelete,
		Importer: &schema.ResourceImporter{ //@TODO: fix import only deployed to avoid problems
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceLinkCreate(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	if d.Get("organization_id").(string) == "" {
		d.Set("organization_id", organizationID.String())
	}

	input := &client.CreateLinkInput{
		CreateInput: client.CreateInput{
			OrganizationID: uuid.MustParse(d.Get("organization_id").(string)),
		},
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		// GroupID:     uuid.MustParse(d.Get("group_id").(string)),
		From: uuid.MustParse(d.Get("from").(string)),
		To:   uuid.MustParse(d.Get("to").(string)),
	}

	link, err := api.CreateLink(input)

	if err != nil {
		// @FIXME: Handle specific error
		// uwerr := errors.Unwrap(err)
		// switch uwerr {
		// case errors.Is(ap):
		// 	return
		// }
		return
	}

	d.SetId(link.ID.String())

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			client.ResourceStateRegistered.String(),
			client.ResourceStateInDeployment.String(),
		},
		Target: []string{
			client.ResourceStateDeployed.String(),
			client.ResourceStateDelivered.String(),
		},
		Refresh: func() (result interface{}, state string, err error) {

			err = resourceLinkRead(d, meta)

			return 0, d.Get("state").(string), err
		},
		Timeout:    time.Duration(30 * time.Minute),
		Delay:      time.Duration(10 * time.Second),
		MinTimeout: time.Duration(15 * time.Second),
	}

	_, err = createStateConf.WaitForState()

	if err != nil {
		// Delete Resource if not created in those delay
		_ = resourceLinkDelete(d, meta)
		d.SetId("")
		return
	}

	return resourceLinkRead(d, meta)
}

func resourceLinkRead(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	if d.Get("organization_id").(string) == "" {
		d.Set("organization_id", organizationID.String())
	}

	input := &client.ReadLinkInput{
		ID:             uuid.MustParse(d.Id()),
		OrganizationID: uuid.MustParse(d.Get("organization_id").(string)),
	}

	link, err := api.ReadLink(input)

	if err != nil {
		uerr := errors.Unwrap(err)
		switch uerr {
		case client.ErrResourceNotFound:
			d.SetId("")
			return nil
		}
		return
	}

	d.Set("name", link.Name)
	d.Set("irn", link.Irn)
	d.Set("description", link.Description)

	d.Set("state", link.State.String())

	// do not update "from", "to" fields
	// check consistency between api and terraform connectors
	// connectors pair from api must be equal with the terraform one
	cpTerraform := ConnectorPair{First: uuid.MustParse(d.Get("from").(string)), Second: uuid.MustParse(d.Get("to").(string))}
	cpBackend := ConnectorPair{First: uuid.MustParse(link.From.String()), Second: uuid.MustParse(link.To.String())}
	if !cpBackend.Equal(&cpTerraform) {
		return errors.New("intercloud api return different link connectors than the expected ones")
	}

	return nil
}

func resourceLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceLinkRead(d, meta)
}

func resourceLinkDelete(d *schema.ResourceData, meta interface{}) (err error) {
	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	if d.Get("organization_id").(string) == "" {
		d.Set("organization_id", organizationID.String())
	}

	input := &client.DeleteLinkInput{
		DeleteInput: client.DeleteInput{
			ID:             uuid.MustParse(d.Id()),
			OrganizationID: uuid.MustParse(d.Get("organization_id").(string)),
		},
	}

	err = api.DeleteLink(input)

	if err == nil {
		d.SetId("")
	}

	return nil
}
