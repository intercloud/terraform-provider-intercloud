package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/api"
)

func (c *client) ReadGroup(in *ReadGroupInput) (out *ReadGroupOutput, err error) {
	if in == nil {
		return nil, fmt.Errorf("input: %w", ErrNil)
	}

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		"GET", "/groups/"+in.ID.String(),
		nil,
		in.OrganizationID.String(),
		nil,
	)
	if err != nil {
		// @TODO: Handle specific error
		uerr := errors.Unwrap(err)
		switch uerr {
		case api.ErrNotFound:
			return nil, fmt.Errorf("%q: %w", in.ID, ErrResourceNotFound)
		}
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&out)

	return
}

func (c *client) CreateGroup(in *CreateGroupInput) (out *CreateGroupOutput, err error) {

	resp, err := c.apiClient.DoRequest(context.Background(), "POST", "/groups", in, in.OrganizationID.String(), nil)
	if err != nil {
		// @TODO: Handle specific error
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&out)

	log.Println(fmt.Sprintf("[WARN] Create out %+s", out))
	log.Println(fmt.Sprintf("[WARN] Create err %+s", err))

	return
}

func (c *client) DeleteGroup(in *DeleteGroupInput) (err error) {

	log.Println(fmt.Sprintf("[WARN] DESTROY in %+s", in))

	_, err = c.apiClient.DoRequest(
		context.Background(),
		http.MethodDelete,
		"/groups/"+in.ID.String(),
		nil,
		in.OrganizationID.String(),
		nil,
	)
	return err
}
