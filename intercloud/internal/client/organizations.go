package client

import (
	"context"
	"encoding/json"
)

// ReadOrganizations finds children organizations of the specified organization
// Inputs:
//	- in: contains the organization parent criteria for finding children
// Outputs:
//	- out: list of children organizations
func (c *client) ReadOrganizations(in *ReadOrganizationInput) (out *[]ReadOrganizationOutput, err error) {

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		"GET", "/organisations",
		nil,
		in.OrganisationID.String(),
		nil,
	)
	if err != nil {
		// @TODO: Handle specific error
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&out)

	return
}

// ReadOrganization finds an organization by parameters specified in the input
// Inputs:
//	- in: contains the organization criteria for finding the organization
// Outputs:
//	- out: the found organization
func (c *client) ReadOrganization(in *ReadOrganizationInput) (out *ReadOrganizationOutput, err error) {

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		"GET", "/organisations/"+in.OrganisationID.String(),
		nil,
		in.OrganisationID.String(),
		nil,
	)
	if err != nil {
		// @TODO: Handle specific error
		// like : { "status": "", "error": "record not found" }
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&out)

	return
}
