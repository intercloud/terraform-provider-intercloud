package connector

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/config"
)

func dataSourceDestinations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDestinationsRead,
		Schema: map[string]*schema.Schema{
			"family": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice(AllConnectionsFamiliesNames(), false),
			},
			"destinations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     destinationResource(),
			},
		},
	}
}

func destinationResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDestinationsRead(d *schema.ResourceData, meta interface{}) (err error) {

	api := meta.(*config.ProviderConfig).ApiClient()

	var filters []*client.Filter

	filter := &client.Filter{
		Name:   client.DestinationFilterFamily.String(),
		Values: []string{d.Get("family").(string)},
	}
	filters = append(filters, filter)

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

	d.SetId(d.Get("family").(string))
	d.Set("destinations", flattenDestinations(&resp.Destinations))

	return
}

func flattenDestinations(in *[]client.Destination) []map[string]string {
	result := make([]map[string]string, len(*in))
	for i, dest := range *in {
		result[i] = map[string]string{
			"family":   dest.FamilyShortname,
			"location": dest.CspCity,
			"id":       dest.ID.String(),
		}
	}
	return result
}
