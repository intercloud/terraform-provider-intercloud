package connector

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/client"
	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/config"
)

type ErrDataSourceDestination struct {
	Msg string
	Err error
}

func (e *ErrDataSourceDestination) Error() string {
	return e.Msg
}

func (e *ErrDataSourceDestination) Unwrap() error {
	return e.Err
}

func NewDataSourceDestinationError(err error) error {
	return &ErrDataSourceDestination{err.Error(), err}
}

var (
	ErrDestinationNilResponse = errors.New("Nil response")
	ErrDestinationNone        = errors.New("No destination")
	ErrDestinationTooMany     = errors.New("Too many destinations")
)

func dataSourceDestination(allowedFamilies []string) *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDestinationRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"family": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice(allowedFamilies, false),
			},
		},
	}
}

func dataSourceDestinationRead(d *schema.ResourceData, meta interface{}) (err error) {

	api := meta.(*config.ProviderConfig).ApiClient()

	var filters []*client.Filter

	if location, ok := d.GetOk("location"); ok {
		filter := &client.Filter{
			Name:   client.DestinationFilterLocation.String(),
			Values: []string{location.(string)},
		}
		filters = append(filters, filter)
	}
	if family, ok := d.GetOk("family"); ok {
		filter := &client.Filter{
			Name:   client.DestinationFilterFamily.String(),
			Values: []string{family.(string)},
		}
		filters = append(filters, filter)
	}

	var resp *client.ReadDestinationsOutput

	resp, err = api.ReadDestinations(
		&client.ReadDestinationsInput{
			Filters: filters,
		},
	)

	if err != nil {
		return NewDataSourceDestinationError(err)
	}

	if resp == nil {
		return NewDataSourceDestinationError(ErrDestinationNilResponse)
	}

	if len(resp.Destinations) < 1 {
		return NewDataSourceDestinationError(ErrDestinationNone)
	}

	if len(resp.Destinations) > 1 {
		return NewDataSourceDestinationError(ErrDestinationNone)
	}

	d.SetId(resp.Destinations[0].ID.String())

	return
}
