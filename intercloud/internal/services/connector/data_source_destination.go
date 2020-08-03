package connector

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/config"
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
	ErrDestinationNilResponse = errors.New("No response from the server")
	ErrDestinationNone        = errors.New("No destination has been found with the given parameters")
	ErrDestinationTooMany     = errors.New("Too many destinations found with the given parameters")
)

func dataSourceDestination() *schema.Resource {
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
				ValidateFunc: validation.StringInSlice(AllConnectionsFamiliesNames(), false),
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
