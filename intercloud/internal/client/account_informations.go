package client

import (
	"context"
	"encoding/json"
)

// ReadAccountInformations reads informations of the current logged in user
func (c *client) ReadAccountInformations() (out *ReadAccountInformations, err error) {

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		"GET", "/me/info",
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

	return
}
