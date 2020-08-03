package client

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// ReadAccountInformations reads informations of the current logged in user
func (c *client) ReadAccountInformations() (out *ReadAccountInformations, err error) {

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		http.MethodGet,
		"/me/info",
		nil,
		"",
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
	log.Printf("[DEBUG] account informations (infos = %+v, err = %+v)", out, err)
	if err != nil {
		return nil, err
	}

	return out, nil
}
