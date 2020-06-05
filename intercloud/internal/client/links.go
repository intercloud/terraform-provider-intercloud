package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/api"
)

func (c *client) ReadLink(input *ReadLinkInput) (output *ReadLinkOutput, err error) {
	resp, err := c.apiClient.DoRequest(
		context.Background(),
		"GET",
		"/flow/link/"+input.ID.String(),
		input,
		input.OrganizationID.String(),
		nil,
	)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		uerr := errors.Unwrap(err)
		switch uerr {
		case api.ErrNotFound:
			return nil, fmt.Errorf("%q: %w", input.ID, ErrResourceNotFound)
		}
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&output)

	return
}
func (c *client) CreateLink(input *CreateLinkInput) (output *CreateLinkOutput, err error) {

	// @FIXME: Change to auto-deploy link when done
	var respSet, respDep *http.Response

	defer func() {
		if respSet != nil {
			_ = respSet.Body.Close()
		}
		if respDep != nil {
			_ = respDep.Body.Close()
		}
	}()
	if respSet, err = c.apiClient.DoRequest(
		context.Background(),
		"POST",
		"/flow/link/set-registered",
		input,
		input.OrganizationID.String(),
		nil,
	); err != nil {
		return nil, err
	}

	err = json.NewDecoder(respSet.Body).Decode(&output)

	if respDep, err = c.apiClient.DoRequest(
		context.Background(),
		"POST",
		"/flow/link/"+output.ID.String()+"/set-in-deployment",
		nil,
		input.OrganizationID.String(),
		nil,
	); err != nil {
		// @FIXME: Delete Here ?
		_ = c.DeleteLink(&DeleteLinkInput{DeleteInput{ID: output.ID}})

		return nil, err
	}

	return
}
func (c *client) DeleteLink(input *DeleteLinkInput) (err error) {

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		"DELETE",
		"/flow/link/"+input.ID.String(),
		nil,
		input.OrganizationID.String(),
		nil,
	)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	return
}
