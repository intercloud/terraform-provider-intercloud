package connector

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector/csp/aws"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector/csp/azure"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector/csp/gcp"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Connector"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"intercloud_destination":  dataSourceDestination(AllFamilies()),
		"intercloud_destinations": dataSourceDestinations(AllFamilies()),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"intercloud_connector": resourceConnector(),
	}
}

func SupportedCspConnections() map[string]*schema.Resource {
	sup := make(map[string]*schema.Resource, 2)
	sup[FamilyAws.String()] = aws.ResourceSchemaAws()
	sup[FamilyAzure.String()] = azure.ResourceSchemaAzure()
	sup[FamilyGcp.String()] = gcp.ResourceSchemaGcp()
	return sup
}

func SupportedEnterprises() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}
