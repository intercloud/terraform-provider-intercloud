package provider

import (
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/common"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/group"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/link"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/organization"
)

func SupportedServices() []common.ServiceRegistration {
	return []common.ServiceRegistration{
		group.Registration{},
		connector.Registration{},
		link.Registration{},
		organization.Registration{},
	}
}
