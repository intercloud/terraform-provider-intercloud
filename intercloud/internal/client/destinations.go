package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/google/uuid"
)

type DestinationFilter int

func (d DestinationFilter) String() string {
	return [...]string{"family", "location"}[d]
}

const (
	DestinationFilterFamily DestinationFilter = iota
	DestinationFilterLocation
)

type Family struct {
	ID        uuid.UUID `json:"id"`
	ShortName string    `json:"shortName"`
}
type CloudGatewayType struct {
	ID      uuid.UUID `json:"id"`
	CspCity string    `json:"cspCity"`
	Family  Family    `json:"family"`
}

type ReadDestinationsResponse struct {
	CloudGatewayTypes []CloudGatewayType `json:"cloudGatewayTypes"`
}

func (c *client) ReadDestinations(in *ReadDestinationsInput) (out *ReadDestinationsOutput, err error) {

	log.Println(fmt.Sprintf("[WARN] ReadDestinations in %+v", in))

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		"GET", "/cloud-gateways-types",
		nil,
		"",
		nil,
	)

	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var types *ReadDestinationsResponse

	err = json.NewDecoder(resp.Body).Decode(&types)

	if err != nil {
		return
	}

	var locationsFilter *[]string
	var familiesFilter *[]string

	for _, filter := range in.Filters {
		if filter.Name == DestinationFilterLocation.String() {
			l := filter.Values
			locationsFilter = &l
		}
	}

	for _, filter := range in.Filters {
		if filter.Name == DestinationFilterFamily.String() {
			l := filter.Values
			familiesFilter = &l
		}
	}

	locations := []Destination{}

	for _, t := range types.CloudGatewayTypes {
		if contains(t.CspCity, locationsFilter) && contains(t.Family.ShortName, familiesFilter) {
			location := Destination{
				ID:              t.ID,
				FamilyID:        t.Family.ID,
				FamilyShortname: t.Family.ShortName,
				CspCity:         t.CspCity,
			}
			locations = append(locations, location)
		}
	}

	out = &ReadDestinationsOutput{
		Destinations: locations,
	}

	return

}

func contains(search string, locationsFilter *[]string) bool {
	if locationsFilter != nil {
		filters := *locationsFilter
		length := len(filters)
		sort.Strings(filters)
		found := sort.SearchStrings(filters, search)
		return found != length && filters[found] == search
	}
	return true
}
