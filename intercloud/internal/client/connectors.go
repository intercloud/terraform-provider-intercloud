package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrConnectorNotEnoughSpareProducts = errors.New("not enough spare products")
)

type AwsParamsApi struct {
	ASN        uint16 `json:"asNumber"`
	AwsAccount int    `json:"awsAccount"`
	Dxvif      string `json:"dxvif,omitempty"`
}

type CspResponseApi struct {
	DestinationID uuid.UUID     `json:"cloudDestinationID"`
	AwsParams     *AwsParamsApi `json:"aws,omitempty"`
}

type EnterpriseResponseApi struct {
}

type ReadConnectorResponseApi struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	GroupID     string                 `json:"groupID"`
	IRN         string                 `json:"irn"`
	Csp         *CspResponseApi        `json:"csp,omitempty"`
	Enterprise  *EnterpriseResponseApi `json:"enterprise,omitempty"`
}

type ConnectorResponseApi struct {
	Read  ReadConnectorResponseApi `json:"read"`
	State string                   `json:"state"`
}

type CreateConnectorResponseApi struct {
	ID uuid.UUID `json:"id"`
}

func (c *client) ReadConnector(input *ReadConnectorInput) (output *ReadConnectorOutput, err error) {

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		http.MethodGet,
		"/flow/connector/"+input.ID.String(),
		nil,
		input.OrganizationID.String(),
		nil,
	)

	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&output)

	if err != nil {
		return nil, err
	}

	return output, nil
}

func (c *client) CreateConnector(input *CreateConnectorInput) (output *CreateConnectorOutput, err error) {

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		http.MethodPost,
		"/flow/connector/auto-deploy",
		input.Connector,
		input.OrganizationID.String(),
		nil,
	)

	if err != nil {
		// @TODO: Handle errors
		// 402 {"error":"{\"error\":\"Not enough products available\",\"code\":\"NotEnoughProductAvailable\"}","code":"NotEnoughSpareProducts"}
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var connector *CreateConnectorResponseApi

	err = json.NewDecoder(resp.Body).Decode(&connector)

	if err != nil {
		return
	}

	output = &CreateConnectorOutput{
		CreateOutput{ID: connector.ID},
	}

	return
}

func (c *client) DeleteConnector(input *DeleteConnectorInput) (err error) {

	resp, err := c.apiClient.DoRequest(
		context.Background(),
		http.MethodDelete,
		"/flow/connector/"+input.ID.String(),
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
