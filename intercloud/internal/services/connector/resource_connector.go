package connector

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	apiclient "github.com/intercloud/terraform-provider-intercloud/intercloud/internal/api"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/config"
)

var cspType = []string{}
var enterpriseType = []string{}
var allType = append(cspType, enterpriseType...)

var locationType = []string{"destination_id", "cloud_appliance_id"}

// func validateCspLocation(i interface{}, k string) (warnings []string, err []error) {

// 	v, ok := i.(string)
// 	if !ok {
// 		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
// 	}

// 	sort.Strings(cspType)

// 	if idx := sort.SearchStrings(cspType, v); idx == len(cspType) {
// 		return nil, []error{fmt.Errorf("expected value of %q to be one of %s", k, strings.Join(cspType, ","))}
// 	}

// 	return nil, nil

// }

func resourceConnector() *schema.Resource {

	supportedCsps := SupportedCsps()
	supportedEnterprises := SupportedEnterprises()

	cspSchemas := make(map[string]*schema.Schema, len(supportedCsps))
	enterpriseSchemas := make(map[string]*schema.Schema, len(supportedEnterprises))

	// @FIXME: once
	for familyShortName, csp := range supportedCsps {
		cspType = append(cspType, familyShortName)
		cspSchemas[familyShortName] = &schema.Schema{
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: allType,
			Elem:         csp,
		}
	}

	resource := &schema.Resource{
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
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"destination_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.IsUUID,
				ExactlyOneOf:  locationType,
				ConflictsWith: enterpriseType,
			},
			"cloud_appliance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.IsUUID,
				ConflictsWith: cspType,
				ExactlyOneOf:  locationType,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"irn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceConnectorCreate,
		Read:   resourceConnectorRead,
		Update: resourceConnectorUpdate,
		Delete: resourceConnectorDelete,
		Importer: &schema.ResourceImporter{ //@TODO: fix import only deployed to avoid problems
			State: schema.ImportStatePassthrough,
		},
	}

	for key, csp := range cspSchemas {
		resource.Schema[key] = csp
	}
	for key, enterprise := range enterpriseSchemas {
		resource.Schema[key] = enterprise
	}

	return resource
}

func resourceConnectorCreate(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	input := &client.CreateConnectorInput{
		CreateInput: client.CreateInput{
			OrganizationID: organizationID,
		},
		Connector: client.ConnectorCreate{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			GroupID:     uuid.MustParse(d.Get("group_id").(string)),
		},
	}

	if destinationID, ok := d.GetOk("destination_id"); ok {

		familiesOut, errA := api.ReadDestinationFamilies()

		if errA != nil {
			return errA
		}

		input.Connector.Csp = &client.CspCreate{
			DestinationID: uuid.MustParse(destinationID.(string)),
		}

		for idx := range cspType {
			if v, ok := d.GetOk(cspType[idx]); ok && len(v.([]interface{})) > 0 {
				options := v.([]interface{})
				s := options[0].(map[string]interface{})

				switch cspType[idx] {
				case FamilyAws.String():
					input.Connector.Csp.AwsParams = expandFamilyAwsParams(s)
				case FamilyAzure.String():
					input.Connector.Csp.AzureParams = expandFamilyAzureParams(s)
				case FamilyGcp.String():
					input.Connector.Csp.GcpParams = expandFamilyGcpParams(s)
				}
				input.Connector.Csp.FamilyID = familiesOut.Families[cspType[idx]]
				break
			}
		}
	}

	// @TODO: Get user input organization_id

	resp, err := api.CreateConnector(input)

	if err != nil {
		return
	}

	d.SetId(resp.ID.String())

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

			err = resourceConnectorRead(d, meta)

			return 0, d.Get("state").(string), err
		},
		Timeout:    time.Duration(30 * time.Minute),
		Delay:      time.Duration(10 * time.Second),
		MinTimeout: time.Duration(15 * time.Second),
	}

	_, err = createStateConf.WaitForState()

	if err != nil {
		// Delete Resource if not created in those delay
		_ = resourceConnectorDelete(d, meta)
		d.SetId("")
		return
	}

	return resourceConnectorRead(d, meta)
}

func resourceConnectorRead(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	input := &client.ReadConnectorInput{
		ID:             uuid.MustParse(d.Id()),
		OrganizationID: organizationID,
	}

	output, err := api.ReadConnector(input)

	if err != nil {

		// @FIXME: not from apiclient
		if errors.Is(err, apiclient.ErrNotFound) {
			// Delete from state if Not Found
			d.SetId("")
			return nil
		}

		return
	}

	log.Println(fmt.Sprintf("[WARN] resourceConnectorRead %+v", output))

	d.Set("name", output.Connector.Name)
	d.Set("group_id", output.Connector.GroupID.String())
	d.Set("description", output.Connector.Description)
	d.Set("irn", output.Connector.Irn)
	d.Set("state", output.State.String())

	if output.Connector.Csp != nil {
		d.Set("destination_id", output.Connector.Csp.DestinationID)
		switch {
		case output.Connector.Csp.AwsParams != nil:
			{
				if err := d.Set(FamilyAws.String(), flattenFamilyAwsParams(output.Connector.Csp.AwsParams)); err != nil {
					return fmt.Errorf("Error setting `%q`: %+v", FamilyAws.String(), err)
				}
			}
		case output.Connector.Csp.AzureParams != nil:
			{
				if err := d.Set(FamilyAzure.String(), flattenFamilyAzureParams(output.Connector.Csp.AzureParams)); err != nil {
					return fmt.Errorf("Error setting `%q`: %+v", FamilyAzure.String(), err)
				}
			}
		case output.Connector.Csp.GcpParams != nil:
			{
				if err := d.Set(FamilyGcp.String(), flattenFamilyGcpParams(output.Connector.Csp.GcpParams)); err != nil {
					return fmt.Errorf("Error setting `%q`: %+v", FamilyGcp.String(), err)
				}
			}
		}
	}

	return nil
}

func resourceConnectorDelete(d *schema.ResourceData, meta interface{}) (err error) {

	cfg := meta.(*config.ProviderConfig)
	api := cfg.ApiClient()

	organizationID := cfg.OrganizationID()

	input := &client.DeleteConnectorInput{
		DeleteInput: client.DeleteInput{
			ID:             uuid.MustParse(d.Id()),
			OrganizationID: organizationID,
		},
	}

	return api.DeleteConnector(input)
}

func resourceConnectorUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceConnectorRead(d, meta)
}
