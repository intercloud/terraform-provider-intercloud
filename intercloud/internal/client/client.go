package client

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/api"
)

var (
	ErrResourceNotFound   = errors.New("not found")
	ErrResourceDuplicated = errors.New("already exists")
	ErrNil                = errors.New("is nil")
)

type Config struct {
	PrivateAccessToken string
	Endpoint           string
	UserAgentProducts  api.UserAgentProducts
}

type Filter struct {
	Name   string
	Values []string
}

type CreateGroupInput struct {
	CreateInput
	Name        string `json:"name"`
	Description string `json:"description"`
}
type CreateGroupOutput struct {
	ID uuid.UUID `json:"id"`
}
type ReadGroupInput struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}
type ReadGroupOutput struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Irn               string    `json:"irn"`
	OrganizationID    uuid.UUID `json:"organisationId"`
	AllocatedCapacity int       `json:"allocatedCapacity"`
}
type DeleteGroupInput struct {
	DeleteInput
}

type ReadDestinationsInput struct {
	Filters []*Filter
}

type Destination struct {
	ID              uuid.UUID
	FamilyID        uuid.UUID
	FamilyShortname string
	CspCity         string
}

type ReadDestinationsOutput struct {
	Destinations []Destination
}

type ReadConnectorInput struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}

type CspRead struct {
	DestinationID   uuid.UUID        `json:"cloudDestinationID"`
	AwsParams       *AwsParams       `json:"aws,omitempty"`
	AwsHostedParams *AwsHostedParams `json:"awshostedconnection,omitempty"`
	AzureParams     *AzureParams     `json:"azure,omitempty"`
	GcpParams       *GcpParams       `json:"gcp,omitempty"`
}

type ConnectorRead struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupID     uuid.UUID `json:"groupID"`
	Irn         string    `json:"irn"`
	Csp         *CspRead  `json:"csp,omitempty"`
}

type ReadConnectorOutput struct {
	Connector *ConnectorRead `json:"read"`
	State     ResourceState  `json:"state"`
}

type AwsParams struct {
	ASN        uint16 `json:"asNumber"`
	AwsAccount int    `json:"awsAccount"`
	Dxvif      string `json:"dxvif,omitempty"`
}

type AwsHostedParams struct {
	ASN          uint16 `json:"asNumber"`
	AwsAccount   int    `json:"awsAccount"`
	PortSpeed    string `json:"portSpeed"`
	VlanID       int    `json:"vlanId,omitempty"`
	ConnectionID string `json:"connectionId,omitempty"`
}

type AzureParams struct {
	ServiceKey string `json:"azureServiceKey"`
}

type GcpParams struct {
	Med            uint64 `json:"med"`
	Bandwidth      string `json:"bandwidth"`
	PairingKey     string `json:"pairingKey"`
	InterconnectID string `json:"interconnectId,omitempty"`
}

type CspCreate struct {
	DestinationID   uuid.UUID        `json:"cloudDestinationID"`
	AwsParams       *AwsParams       `json:"aws,omitempty"`
	AwsHostedParams *AwsHostedParams `json:"awshostedconnection,omitempty"`
	AzureParams     *AzureParams     `json:"azure,omitempty"`
	GcpParams       *GcpParams       `json:"gcp,omitempty"`
}

type EnterpriseCreate struct {
	// @TODO:
}

type ConnectorCreate struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	GroupID     uuid.UUID         `json:"groupID"`
	Csp         *CspCreate        `json:"csp,omitempty"`
	Enterprise  *EnterpriseCreate `json:"enterprise,omitempty"`
}

func (c ConnectorCreate) String() string {
	b, _ := json.MarshalIndent(c, "ConnectorCreate", "\t")
	return string(b)
}

type CreateConnectorInput struct {
	CreateInput
	Connector ConnectorCreate
}
type CreateConnectorOutput struct {
	CreateOutput
}
type DeleteConnectorInput struct {
	DeleteInput
}

type ReadLinkInput struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}
type ReadLinkOutput struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Irn            string        `json:"irn"`
	From           uuid.UUID     `json:"from"`
	To             uuid.UUID     `json:"to"`
	State          ResourceState `json:"state"`
	OrganizationID uuid.UUID
}
type CreateLinkInput struct {
	CreateInput
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupID     uuid.UUID `json:"-"` // `json:"groupID"`
	From        uuid.UUID `json:"from"`
	To          uuid.UUID `json:"to"`
}
type CreateLinkOutput struct {
	CreateOutput
}
type DeleteLinkInput struct {
	DeleteInput
}

type DeleteInput struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}
type CreateInput struct {
	OrganizationID uuid.UUID
}
type CreateOutput struct {
	ID uuid.UUID
}
type ReadOrganizationInput struct {
	OrganisationID uuid.UUID `json:"organisationId"`
}

type ReadOrganizationOutput struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ParentID uuid.UUID `json:"parentId"`
}

type ReadAccountInformations struct {
	ID             uuid.UUID `json:"id"`
	Firstname      string    `json:"firstname"`
	Lastname       string    `json:"lastname"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	OrganizationID uuid.UUID `json:"organisationId"`
}

type Client interface {
	ReadGroup(*ReadGroupInput) (*ReadGroupOutput, error)
	CreateGroup(*CreateGroupInput) (*CreateGroupOutput, error)
	DeleteGroup(*DeleteGroupInput) error

	ReadConnector(*ReadConnectorInput) (*ReadConnectorOutput, error)
	CreateConnector(*CreateConnectorInput) (*CreateConnectorOutput, error)
	DeleteConnector(*DeleteConnectorInput) error

	ReadLink(*ReadLinkInput) (*ReadLinkOutput, error)
	CreateLink(*CreateLinkInput) (*CreateLinkOutput, error)
	DeleteLink(*DeleteLinkInput) error

	ReadDestinations(*ReadDestinationsInput) (*ReadDestinationsOutput, error)

	ReadOrganization(in *ReadOrganizationInput) (out *ReadOrganizationOutput, err error)
	ReadOrganizations(in *ReadOrganizationInput) (out *[]ReadOrganizationOutput, err error)

	ReadAccountInformations() (*ReadAccountInformations, error)
}

type client struct {
	apiClient *api.Client
}

func getHttpClient() *http.Client {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(5 * time.Second),
		}).Dial,
		TLSHandshakeTimeout: time.Duration(5 * time.Second),
	}

	return &http.Client{
		Timeout:   time.Duration(10 * time.Second),
		Transport: netTransport,
	}
}

func NewClient(config *Config) (c Client) {
	h := getHttpClient()
	c = &client{
		apiClient: api.NewClient(
			config.Endpoint,
			config.PrivateAccessToken,
			config.UserAgentProducts,
			h,
		),
	}

	return
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(config *Config, fn RoundTripFunc) (c Client) {
	c = &client{
		apiClient: api.NewClient(
			config.Endpoint,
			config.PrivateAccessToken,
			config.UserAgentProducts,
			&http.Client{
				Transport: RoundTripFunc(fn),
			},
		),
	}
	return
}
