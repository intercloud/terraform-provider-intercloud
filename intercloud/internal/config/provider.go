package config

import (
	"github.com/google/uuid"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/api"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/client"
	"github.com/intercloud/terraform-provider-intercloud/version"
)

type ProviderConfig struct {
	client         client.Client
	organizationID uuid.UUID
}

type ProviderConfigInput struct {
	PrivateAccessToken string
	ApiEndpoint        string
	TerraformVersion   string
	OrganizationID     uuid.UUID
}

func NewProviderConfig(input ProviderConfigInput) (config *ProviderConfig) {

	config = &ProviderConfig{
		client: client.NewClient(
			&client.Config{
				Endpoint:           input.ApiEndpoint,
				PrivateAccessToken: input.PrivateAccessToken,
				UserAgentProducts: []api.UserAgentProduct{
					{Name: "Terraform", Version: input.TerraformVersion},
					{Name: "IntercloudPluginProvider", Version: version.String()},
				},
			},
		),
		organizationID: input.OrganizationID,
	}

	return

}

func (pc *ProviderConfig) ApiClient() client.Client {
	return pc.client
}

func (pc *ProviderConfig) OrganizationID() uuid.UUID {
	return pc.organizationID
}
